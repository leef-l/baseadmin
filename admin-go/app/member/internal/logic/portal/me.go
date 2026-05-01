package portal

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/logic/walletops"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/sms"
	"gbaseadmin/utility/snowflake"
)

// ----- Profile -----

// GetMyProfile 返回当前会员的个人资料 + 等级 + 团队统计。
func (s *sPortalAuth) GetMyProfile(ctx context.Context, userID int64) (*MyProfile, error) {
	if userID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	var u entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Scan(&u); err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return nil, gerror.New("会员不存在")
	}
	levelName := ""
	if u.LevelId > 0 {
		val, _ := dao.MemberLevel.Ctx(ctx).
			Where(dao.MemberLevel.Columns().Id, u.LevelId).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			Value(dao.MemberLevel.Columns().Name)
		if val != nil {
			levelName = val.String()
		}
	}
	expireAt := ""
	if u.LevelExpireAt != nil && !u.LevelExpireAt.IsZero() {
		expireAt = u.LevelExpireAt.String()
	}

	return &MyProfile{
		MemberID:      fmt.Sprintf("%d", u.Id),
		Phone:         u.Phone,
		Username:      u.Username,
		Nickname:      u.Nickname,
		Avatar:        u.Avatar,
		RealName:      u.RealName,
		InviteCode:    u.InviteCode,
		ParentID:      fmt.Sprintf("%d", u.ParentId),
		LevelID:       fmt.Sprintf("%d", u.LevelId),
		LevelName:     levelName,
		LevelExpireAt: expireAt,
		IsActive:      u.IsActive,
		IsQualified:   u.IsQualified,
		TeamCount:     int(u.TeamCount),
		DirectCount:   int(u.DirectCount),
		ActiveCount:   int(u.ActiveCount),
		TeamTurnover:  int64(u.TeamTurnover),
		InviteURL:     buildInviteURL(ctx, u.InviteCode),
	}, nil
}

// MyProfile 返回值（不与 api v1 耦合，由 controller 转换）。
type MyProfile struct {
	MemberID      string
	Phone         string
	Username      string
	Nickname      string
	Avatar        string
	RealName      string
	InviteCode    string
	ParentID      string
	LevelID       string
	LevelName     string
	LevelExpireAt string
	IsActive      int
	IsQualified   int
	TeamCount     int
	DirectCount   int
	ActiveCount   int
	TeamTurnover  int64
	InviteURL     string
}

// ----- Update profile -----

// UpdateMyProfileInput 修改个人资料入参。
type UpdateMyProfileInput struct {
	UserID   int64
	Nickname string
	Avatar   string
	RealName string
}

// UpdateMyProfile 修改昵称 / 头像 / 真实姓名。
// 实名（real_name）一旦填写，本次实现不限制再次修改；后续可加"已实名后不可改"。
func (s *sPortalAuth) UpdateMyProfile(ctx context.Context, in *UpdateMyProfileInput) error {
	if in == nil || in.UserID <= 0 {
		return gerror.New("会员未登录")
	}
	data := do.MemberUser{}
	hasChange := false
	if v := strings.TrimSpace(in.Nickname); v != "" {
		data.Nickname = v
		hasChange = true
	}
	if v := strings.TrimSpace(in.Avatar); v != "" {
		data.Avatar = v
		hasChange = true
	}
	if v := strings.TrimSpace(in.RealName); v != "" {
		data.RealName = v
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	_, err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, in.UserID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Data(data).
		Update()
	return err
}

// ----- Change password -----

// ChangeMyPasswordInput 已登录态修改密码。
type ChangeMyPasswordInput struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

// ChangeMyPassword 已登录态修改密码：校验旧密码 + 设置新密码。
func (s *sPortalAuth) ChangeMyPassword(ctx context.Context, in *ChangeMyPasswordInput) error {
	if in == nil || in.UserID <= 0 {
		return gerror.New("会员未登录")
	}
	oldPwd := strings.TrimSpace(in.OldPassword)
	newPwd := strings.TrimSpace(in.NewPassword)
	if oldPwd == "" || newPwd == "" {
		return gerror.New("旧密码 / 新密码不能为空")
	}
	if oldPwd == newPwd {
		return gerror.New("新密码不能与旧密码相同")
	}
	if err := password.ValidatePolicy(newPwd); err != nil {
		return gerror.New(err.Error())
	}
	value, err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, in.UserID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Value(dao.MemberUser.Columns().Password)
	if err != nil {
		return err
	}
	if value.IsEmpty() {
		return gerror.New("会员不存在")
	}
	if !password.Verify(value.String(), oldPwd) {
		return gerror.New("旧密码错误")
	}
	hashed, err := password.Hash(newPwd)
	if err != nil {
		return err
	}
	_, err = dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, in.UserID).
		Data(do.MemberUser{Password: hashed}).
		Update()
	return err
}

// ----- Change phone -----

// ChangeMyPhoneInput 修改手机号入参。
type ChangeMyPhoneInput struct {
	UserID   int64
	NewPhone string
	SmsCode  string
}

// ChangeMyPhone 已登录态换绑手机号。新号必须未被注册，且收到 scene=change_phone 的验证码。
func (s *sPortalAuth) ChangeMyPhone(ctx context.Context, in *ChangeMyPhoneInput) error {
	if in == nil || in.UserID <= 0 {
		return gerror.New("会员未登录")
	}
	newPhone := strings.TrimSpace(in.NewPhone)
	smsCode := strings.TrimSpace(in.SmsCode)
	if newPhone == "" || smsCode == "" {
		return gerror.New("新手机号 / 验证码不能为空")
	}
	if !isPhone(newPhone) {
		return gerror.New("新手机号格式不正确")
	}
	if _, err := sms.Default().VerifyCode(ctx, &sms.VerifyCodeInput{
		Phone:   newPhone,
		Scene:   "change_phone",
		Code:    smsCode,
		Consume: true,
	}); err != nil {
		return err
	}
	// 新号未被占用
	exist, err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Phone, newPhone).
		WhereNot(dao.MemberUser.Columns().Id, in.UserID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return gerror.New("该手机号已被其他会员绑定")
	}
	_, err = dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, in.UserID).
		Data(do.MemberUser{Phone: newPhone}).
		Update()
	return err
}

// ----- Wallets -----

// GetMyWallets 返回三钱包余额（如果某钱包不存在自动创建一个空记录，避免历史数据缺失）。
func (s *sPortalAuth) GetMyWallets(ctx context.Context, userID int64) (*MyWallets, error) {
	if userID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	out := &MyWallets{}
	for _, walletType := range []int{walletops.WalletTypeCoupon, walletops.WalletTypeReward, walletops.WalletTypePromote} {
		w, err := s.loadOrInitWallet(ctx, userID, walletType)
		if err != nil {
			return nil, err
		}
		info := WalletInfoData{
			Balance:      formatCent(w.Balance),
			BalanceCent:  w.Balance,
			TotalIncome:  formatCent(int64(w.TotalIncome)),
			TotalExpense: formatCent(int64(w.TotalExpense)),
			FrozenAmount: formatCent(int64(w.FrozenAmount)),
		}
		switch walletType {
		case walletops.WalletTypeCoupon:
			out.Coupon = info
		case walletops.WalletTypeReward:
			out.Reward = info
		case walletops.WalletTypePromote:
			out.Promote = info
		}
	}
	return out, nil
}

// MyWallets 三钱包返回。
type MyWallets struct {
	Coupon  WalletInfoData
	Reward  WalletInfoData
	Promote WalletInfoData
}

// WalletInfoData 单钱包数据。
type WalletInfoData struct {
	Balance      string
	BalanceCent  int64
	TotalIncome  string
	TotalExpense string
	FrozenAmount string
}

// loadOrInitWallet 历史数据兜底：钱包记录缺失时自动建一行。
func (s *sPortalAuth) loadOrInitWallet(ctx context.Context, userID int64, walletType int) (*entity.MemberWallet, error) {
	var w entity.MemberWallet
	err := dao.MemberWallet.Ctx(ctx).
		Where(dao.MemberWallet.Columns().UserId, userID).
		Where(dao.MemberWallet.Columns().WalletType, walletType).
		Where(dao.MemberWallet.Columns().DeletedAt, nil).
		Scan(&w)
	if err != nil {
		return nil, err
	}
	if w.Id != 0 {
		return &w, nil
	}
	// 缺失：用一次事务初始化
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := gtime.Now()
		_, err := tx.Model(dao.MemberWallet.Table()).Ctx(ctx).Data(do.MemberWallet{
			Id:           genWalletID(),
			UserId:       userID,
			WalletType:   walletType,
			Balance:      0,
			TotalIncome:  0,
			TotalExpense: 0,
			FrozenAmount: 0,
			Status:       1,
			CreatedBy:    userID,
			CreatedAt:    now,
			UpdatedAt:    now,
		}).Insert()
		return err
	})
	if err != nil {
		return nil, err
	}
	if err := dao.MemberWallet.Ctx(ctx).
		Where(dao.MemberWallet.Columns().UserId, userID).
		Where(dao.MemberWallet.Columns().WalletType, walletType).
		Scan(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// ----- Wallet logs -----

// MyWalletLogsInput 钱包流水分页入参。
type MyWalletLogsInput struct {
	UserID     int64
	WalletType int // 0=全部
	PageNum    int
	PageSize   int
}

// MyWalletLogsOutput 钱包流水列表。
type MyWalletLogsOutput struct {
	Total int
	List  []*WalletLogItem
}

// WalletLogItem 单条流水（金额已转字符串元）。
type WalletLogItem struct {
	ID             string
	WalletType     int
	WalletTypeText string
	ChangeType     int
	ChangeTypeText string
	ChangeAmount   string
	BeforeBalance  string
	AfterBalance   string
	RelatedOrderNo string
	Remark         string
	CreatedAt      string
}

// ListMyWalletLogs 钱包流水分页。
func (s *sPortalAuth) ListMyWalletLogs(ctx context.Context, in *MyWalletLogsInput) (*MyWalletLogsOutput, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	pageNum := in.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	m := dao.MemberWalletLog.Ctx(ctx).
		Where(dao.MemberWalletLog.Columns().UserId, in.UserID).
		Where(dao.MemberWalletLog.Columns().DeletedAt, nil)
	if in.WalletType > 0 {
		m = m.Where(dao.MemberWalletLog.Columns().WalletType, in.WalletType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberWalletLog
	if err := m.OrderDesc(dao.MemberWalletLog.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&rows); err != nil {
		return nil, err
	}
	out := &MyWalletLogsOutput{Total: total, List: make([]*WalletLogItem, 0, len(rows))}
	for _, row := range rows {
		out.List = append(out.List, &WalletLogItem{
			ID:             fmt.Sprintf("%d", row.Id),
			WalletType:     row.WalletType,
			WalletTypeText: walletTypeText(row.WalletType),
			ChangeType:     row.ChangeType,
			ChangeTypeText: changeTypeText(row.ChangeType),
			ChangeAmount:   formatSignedCent(row.ChangeAmount),
			BeforeBalance:  formatCent(row.BeforeBalance),
			AfterBalance:   formatCent(row.AfterBalance),
			RelatedOrderNo: row.RelatedOrderNo,
			Remark:         row.Remark,
			CreatedAt:      timeStr(row.CreatedAt),
		})
	}
	return out, nil
}

// ----- Team -----

// MyTeamInput 团队列表入参。
type MyTeamInput struct {
	UserID   int64
	Scope    string // direct=仅直推 all=全部团队（递归）
	PageNum  int
	PageSize int
}

// MyTeamOutput 团队列表。
type MyTeamOutput struct {
	Total int
	List  []*TeamMember
}

// TeamMember 团队成员简要。
type TeamMember struct {
	MemberID    string
	Nickname    string
	Avatar      string
	Phone       string
	LevelName   string
	IsQualified int
	JoinedAt    string
}

// ListMyTeam 列出团队成员。
//   - direct 模式：直接 parent_id = userID
//   - all 模式：递归所有下级 ID 集合，再按集合查
func (s *sPortalAuth) ListMyTeam(ctx context.Context, in *MyTeamInput) (*MyTeamOutput, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	scope := strings.TrimSpace(in.Scope)
	if scope != "all" {
		scope = "direct"
	}
	pageNum := in.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var ids []int64
	if scope == "direct" {
		ids = []int64{in.UserID}
	} else {
		ids = []int64{in.UserID}
		// BFS 递归收集下级
		seen := map[int64]struct{}{in.UserID: {}}
		queue := []int64{in.UserID}
		for len(queue) > 0 && len(seen) < 5000 {
			batch := queue
			queue = nil
			var rows []struct {
				Id int64 `json:"id"`
			}
			if err := dao.MemberUser.Ctx(ctx).
				Fields(dao.MemberUser.Columns().Id).
				WhereIn(dao.MemberUser.Columns().ParentId, batch).
				Where(dao.MemberUser.Columns().DeletedAt, nil).
				Scan(&rows); err != nil {
				return nil, err
			}
			for _, row := range rows {
				if _, ok := seen[row.Id]; ok {
					continue
				}
				seen[row.Id] = struct{}{}
				ids = append(ids, row.Id)
				queue = append(queue, row.Id)
			}
		}
	}

	// scope=direct: parent_id=userID；scope=all: parent_id IN ids 且不含 userID 本人
	m := dao.MemberUser.Ctx(ctx).Where(dao.MemberUser.Columns().DeletedAt, nil)
	if scope == "direct" {
		m = m.Where(dao.MemberUser.Columns().ParentId, in.UserID)
	} else {
		m = m.WhereIn(dao.MemberUser.Columns().ParentId, ids)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var users []entity.MemberUser
	if err := m.OrderDesc(dao.MemberUser.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&users); err != nil {
		return nil, err
	}

	// 批量查等级名
	levelIDs := make([]uint64, 0, len(users))
	for _, u := range users {
		if u.LevelId > 0 {
			levelIDs = append(levelIDs, u.LevelId)
		}
	}
	levelMap := make(map[uint64]string, len(levelIDs))
	if len(levelIDs) > 0 {
		var levels []entity.MemberLevel
		if err := dao.MemberLevel.Ctx(ctx).
			WhereIn(dao.MemberLevel.Columns().Id, levelIDs).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			Scan(&levels); err == nil {
			for _, lv := range levels {
				levelMap[lv.Id] = lv.Name
			}
		}
	}

	out := &MyTeamOutput{Total: total, List: make([]*TeamMember, 0, len(users))}
	for _, u := range users {
		out.List = append(out.List, &TeamMember{
			MemberID:    fmt.Sprintf("%d", u.Id),
			Nickname:    u.Nickname,
			Avatar:      u.Avatar,
			Phone:       maskPhone(u.Phone),
			LevelName:   levelMap[u.LevelId],
			IsQualified: u.IsQualified,
			JoinedAt:    timeStr(u.CreatedAt),
		})
	}
	return out, nil
}

// ----- helpers -----

// formatCent 把分换算为元字符串（保留两位小数）。负数保留负号。
func formatCent(cent int64) string {
	negative := cent < 0
	if negative {
		cent = -cent
	}
	yuan := cent / 100
	fen := cent % 100
	sign := ""
	if negative {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%02d", sign, yuan, fen)
}

// formatSignedCent 流水变动金额：保留正负号。
func formatSignedCent(cent int64) string {
	if cent > 0 {
		return "+" + formatCent(cent)
	}
	return formatCent(cent)
}

func walletTypeText(t int) string {
	switch t {
	case walletops.WalletTypeCoupon:
		return "优惠券余额"
	case walletops.WalletTypeReward:
		return "奖金余额"
	case walletops.WalletTypePromote:
		return "推广奖余额"
	}
	return "未知"
}

func changeTypeText(t int) string {
	switch t {
	case walletops.ChangeTypeRecharge:
		return "充值"
	case walletops.ChangeTypeConsume:
		return "消费"
	case walletops.ChangeTypePromote:
		return "推广奖"
	case walletops.ChangeTypeWHIncome:
		return "仓库卖出"
	case walletops.ChangeTypeFee:
		return "平台扣除"
	case walletops.ChangeTypeAdjust:
		return "后台调整"
	}
	return "其它"
}

func timeStr(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("Y-m-d H:i:s")
}

// maskPhone 把手机号脱敏：138****8000。
func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

// buildInviteURL 拼邀请链接。
func buildInviteURL(ctx context.Context, inviteCode string) string {
	if strings.TrimSpace(inviteCode) == "" {
		return ""
	}
	base := strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(ctx, "member.h5BaseURL").String()), "/")
	if base == "" {
		base = "https://funddisk.easytestdev.online/h5"
	}
	q := url.Values{}
	q.Set("inviteCode", inviteCode)
	return fmt.Sprintf("%s/#/auth/register?%s", base, q.Encode())
}

// genWalletID 生成钱包行 ID（雪花）。
func genWalletID() snowflake.JsonInt64 { return snowflake.Generate() }
