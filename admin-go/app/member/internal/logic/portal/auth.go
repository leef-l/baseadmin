// Package portal 是 C 端会员业务逻辑层。
//
// 与 logic/user (codegen 生成的后台 CRUD) 的关键差异：
//   - 写库不带租户/部门数据权限校验：C 端会员归属是 token 维度，不需要 tenant scope
//   - 注册时强制邀请码定位上级；后台 user.Create 不强制
//   - 登录返回 jwt.GenerateMemberToken（独立 secret），不混用后台 token
package portal

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/logic/teamops"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/cache"
	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/sms"
	"gbaseadmin/utility/snowflake"
)

// portalLogic 同包内 me/wallet/team/order/mall/warehouse 等扩展方法的共享单例。
// 通过 Logic() / AuthLogic() 暴露给 controller 调用，避免每加一个方法都改 service 接口。
var portalLogic = &sPortalAuth{}

// Logic 返回 portal 业务实现单例（仅 internal 包内使用）。
func Logic() *sPortalAuth { return portalLogic }

// AuthLogic 等价于 Logic()，命名更贴近 controller 调用语境。
func AuthLogic() *sPortalAuth { return portalLogic }

func init() {
	service.RegisterPortalAuth(portalLogic)
}

// 限流参数。
const (
	loginFailLimit       = 5
	loginFailWindow      = 10 * time.Minute
	inviteCodeMaxRetries = 8 // 同一注册流程内最多生成 8 次邀请码（雪花碰撞极小概率）
)

type sPortalAuth struct{}

// ---------- Register ----------

func (s *sPortalAuth) Register(ctx context.Context, in *service.PortalRegisterInput) (*service.PortalLoginOutput, error) {
	if in == nil {
		return nil, gerror.New("注册参数不能为空")
	}
	phone := strings.TrimSpace(in.Phone)
	smsCode := strings.TrimSpace(in.SmsCode)
	pwd := strings.TrimSpace(in.Password)
	inviteCode := strings.TrimSpace(in.InviteCode)
	nickname := strings.TrimSpace(in.Nickname)

	if phone == "" || smsCode == "" || pwd == "" || inviteCode == "" {
		return nil, gerror.New("手机号 / 验证码 / 密码 / 邀请码不能为空")
	}
	if err := password.ValidatePolicy(pwd); err != nil {
		return nil, gerror.New(err.Error())
	}

	// 1. 验证短信验证码（注册场景一次性消费）
	if _, err := sms.Default().VerifyCode(ctx, &sms.VerifyCodeInput{
		Phone:   phone,
		Scene:   "register",
		Code:    smsCode,
		Consume: true,
	}); err != nil {
		return nil, err
	}

	// 2. 校验手机号未注册
	exist, err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Phone, phone).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return nil, err
	}
	if exist > 0 {
		return nil, gerror.New("该手机号已注册，请直接登录")
	}

	// 3. 解析邀请码 → 上级
	var parent entity.MemberUser
	err = dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().InviteCode, inviteCode).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Where(dao.MemberUser.Columns().Status, 1).
		Scan(&parent)
	if err != nil {
		return nil, err
	}
	if parent.Id == 0 {
		return nil, gerror.New("邀请码无效或邀请人已被禁用")
	}

	// 4. 加密密码
	hashed, err := password.Hash(pwd)
	if err != nil {
		return nil, err
	}

	// 5. 生成本人邀请码（重试避免冲突）
	myInviteCode, err := s.generateUniqueInviteCode(ctx)
	if err != nil {
		return nil, err
	}

	// 6. 默认昵称
	if nickname == "" {
		nickname = maskPhoneAsNickname(phone)
	}

	// 7. 生成 ID
	memberID := snowflake.Generate()
	if int64(memberID) <= 0 {
		return nil, gerror.New("生成会员 ID 失败")
	}

	// 8. 同事务建 user + 三个 wallet
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.MemberUser.Table()).Ctx(ctx).Data(do.MemberUser{
			Id:           memberID,
			ParentId:     parent.Id,
			Username:     phone, // 默认 username = phone，用户后续可改
			Password:     hashed,
			Nickname:     nickname,
			Phone:        phone,
			LevelId:      0,
			TeamCount:    0,
			DirectCount:  0,
			ActiveCount:  0,
			TeamTurnover: 0,
			IsActive:     1,
			IsQualified:  1,
			InviteCode:   myInviteCode,
			RegisterIp:   in.RegisterIP,
			Sort:         0,
			Status:       1,
			TenantId:     parent.TenantId,
			MerchantId:   parent.MerchantId,
			CreatedBy:    0,
			DeptId:       0,
		}).Insert(); err != nil {
			return err
		}

		now := gtime.Now()
		walletData := make([]do.MemberWallet, 0, 3)
		for _, walletType := range []int{1, 2, 3} {
			walletData = append(walletData, do.MemberWallet{
				Id:           snowflake.Generate(),
				UserId:       memberID,
				WalletType:   walletType,
				Balance:      0,
				TotalIncome:  0,
				TotalExpense: 0,
				FrozenAmount: 0,
				Status:       1,
				TenantId:     parent.TenantId,
				MerchantId:   parent.MerchantId,
				CreatedBy:    memberID,
				DeptId:       0,
				CreatedAt:    now,
				UpdatedAt:    now,
			})
		}
		if _, err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).Data(walletData).Insert(); err != nil {
			return err
		}

		// 链式更新祖先 team_count / direct_count；注册即激活，所以同时累加 active_count，
		// 并触发各祖先 TryUpgrade 检查（teamops.IncrAncestorActiveCount 内部已自带 TryUpgrade）。
		if err := teamops.IncrAncestorTeamCount(ctx, tx, int64(memberID)); err != nil {
			return err
		}
		if err := teamops.IncrAncestorActiveCount(ctx, tx, int64(memberID)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 9. 签发 token
	token, err := jwt.GenerateMemberToken(int64(memberID), phone, 0, 0, "member")
	if err != nil {
		return nil, gerror.New("生成 token 失败")
	}

	return &service.PortalLoginOutput{
		Token:       token,
		MemberID:    fmt.Sprintf("%d", int64(memberID)),
		Phone:       phone,
		Nickname:    nickname,
		Avatar:      "",
		InviteCode:  myInviteCode,
		LevelID:     "0",
		IsQualified: 1,
	}, nil
}

// ---------- Login ----------

func (s *sPortalAuth) Login(ctx context.Context, in *service.PortalLoginInput) (*service.PortalLoginOutput, error) {
	if in == nil {
		return nil, gerror.New("登录参数不能为空")
	}
	account := strings.TrimSpace(in.Account)
	pwd := strings.TrimSpace(in.Password)
	if account == "" || pwd == "" {
		return nil, gerror.New("账号 / 密码不能为空")
	}

	// 限流
	limitKey := fmt.Sprintf("member:portal:login_fail:%s", strings.ToLower(account))
	count, _ := cache.IncrWithTTL(ctx, limitKey, loginFailWindow)
	if count > loginFailLimit {
		return nil, gerror.New("登录失败次数过多，请 10 分钟后再试")
	}

	user, err := s.findLoginUser(ctx, account)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Id == 0 {
		password.DummyVerify(pwd) // 防时序攻击：始终消耗一次 bcrypt
		return nil, gerror.New("账号或密码错误")
	}
	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用，请联系客服")
	}
	if !password.Verify(user.Password, pwd) {
		return nil, gerror.New("账号或密码错误")
	}
	// 登录成功 → 清空限流
	_ = cache.Delete(ctx, limitKey)

	// 异步升级 hash（如果是历史 SHA-256）
	if password.NeedsRehash(user.Password) {
		if upgraded, hashErr := password.Hash(pwd); hashErr == nil {
			_, _ = dao.MemberUser.Ctx(ctx).
				Where(dao.MemberUser.Columns().Id, user.Id).
				Data(do.MemberUser{Password: upgraded}).
				Update()
		}
	}

	// 更新最后登录时间
	_, _ = dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, user.Id).
		Data(do.MemberUser{LastLoginAt: gtime.Now()}).
		Update()

	token, err := jwt.GenerateMemberToken(int64(user.Id), user.Phone, 0, 0, "member")
	if err != nil {
		return nil, gerror.New("生成 token 失败")
	}

	return &service.PortalLoginOutput{
		Token:       token,
		MemberID:    fmt.Sprintf("%d", user.Id),
		Phone:       user.Phone,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		InviteCode:  user.InviteCode,
		LevelID:     fmt.Sprintf("%d", user.LevelId),
		IsQualified: user.IsQualified,
	}, nil
}

// findLoginUser 按手机号 / 邀请码 / 用户名顺序匹配会员。
func (s *sPortalAuth) findLoginUser(ctx context.Context, account string) (*entity.MemberUser, error) {
	cols := dao.MemberUser.Columns()
	// 11 位纯数字优先按手机号
	if isPhone(account) {
		var u entity.MemberUser
		if err := dao.MemberUser.Ctx(ctx).
			Where(cols.Phone, account).
			Where(cols.DeletedAt, nil).
			Scan(&u); err != nil {
			return nil, err
		}
		if u.Id != 0 {
			return &u, nil
		}
	}
	// 邀请码（一般是大写 base32）
	{
		var u entity.MemberUser
		if err := dao.MemberUser.Ctx(ctx).
			Where(cols.InviteCode, account).
			Where(cols.DeletedAt, nil).
			Scan(&u); err != nil {
			return nil, err
		}
		if u.Id != 0 {
			return &u, nil
		}
	}
	// username
	{
		var u entity.MemberUser
		if err := dao.MemberUser.Ctx(ctx).
			Where(cols.Username, account).
			Where(cols.DeletedAt, nil).
			Scan(&u); err != nil {
			return nil, err
		}
		if u.Id != 0 {
			return &u, nil
		}
	}
	return nil, nil
}

// ---------- Forget Password ----------

func (s *sPortalAuth) ForgetPassword(ctx context.Context, in *service.PortalForgetPasswordInput) error {
	if in == nil {
		return gerror.New("参数不能为空")
	}
	phone := strings.TrimSpace(in.Phone)
	smsCode := strings.TrimSpace(in.SmsCode)
	newPwd := strings.TrimSpace(in.NewPassword)
	if phone == "" || smsCode == "" || newPwd == "" {
		return gerror.New("手机号 / 验证码 / 新密码不能为空")
	}
	if err := password.ValidatePolicy(newPwd); err != nil {
		return gerror.New(err.Error())
	}
	if _, err := sms.Default().VerifyCode(ctx, &sms.VerifyCodeInput{
		Phone:   phone,
		Scene:   "forget_password",
		Code:    smsCode,
		Consume: true,
	}); err != nil {
		return err
	}

	var u entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Phone, phone).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Scan(&u); err != nil {
		return err
	}
	if u.Id == 0 {
		return gerror.New("该手机号未注册")
	}

	hashed, err := password.Hash(newPwd)
	if err != nil {
		return err
	}
	_, err = dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, u.Id).
		Data(do.MemberUser{Password: hashed}).
		Update()
	return err
}

// ---------- Invite Preview ----------

func (s *sPortalAuth) InvitePreview(ctx context.Context, inviteCode string) (*service.PortalInvitePreviewOutput, error) {
	inviteCode = strings.TrimSpace(inviteCode)
	if inviteCode == "" {
		return &service.PortalInvitePreviewOutput{Found: false}, nil
	}
	var u entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().InviteCode, inviteCode).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Where(dao.MemberUser.Columns().Status, 1).
		Scan(&u); err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return &service.PortalInvitePreviewOutput{Found: false}, nil
	}
	return &service.PortalInvitePreviewOutput{
		Found:    true,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}, nil
}

// ---------- helpers ----------

// generateUniqueInviteCode 生成不冲突的邀请码（base32 8 位，去掉容易混淆的 0/1/I/O 留 32 字符集）。
func (s *sPortalAuth) generateUniqueInviteCode(ctx context.Context) (string, error) {
	length := g.Cfg().MustGet(ctx, "member.inviteCodeLength", 8).Int()
	if length <= 0 {
		length = 8
	}
	for i := 0; i < inviteCodeMaxRetries; i++ {
		code := randInviteCode(length)
		exist, err := dao.MemberUser.Ctx(ctx).
			Where(dao.MemberUser.Columns().InviteCode, code).
			Count()
		if err != nil {
			return "", err
		}
		if exist == 0 {
			return code, nil
		}
	}
	return "", gerror.New("生成邀请码失败，请重试")
}

// randInviteCode 用 base32 编码随机字节得到大写邀请码。
func randInviteCode(length int) string {
	if length <= 0 {
		length = 8
	}
	// 取 length*5/8 + 1 字节足够 base32 编出 length 字符
	byteLen := (length*5)/8 + 1
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		// 兜底：用时间戳
		return strings.ToUpper(fmt.Sprintf("FX%010d", time.Now().UnixNano()%9999999999))[:length]
	}
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)
	if len(encoded) < length {
		encoded = encoded + strings.Repeat("F", length-len(encoded))
	}
	return encoded[:length]
}

// isPhone 简易判断 11 位手机号。
func isPhone(s string) bool {
	if len(s) != 11 {
		return false
	}
	for i, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
		if i == 0 && ch != '1' {
			return false
		}
	}
	return true
}

// maskPhoneAsNickname 把手机号脱敏作为默认昵称：138****8000。
func maskPhoneAsNickname(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
