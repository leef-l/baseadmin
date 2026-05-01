package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 注册 -----

// RegisterReq C 端会员注册请求。
//
// 业务规则：
//   - phone + smsCode 必须先经过 sms.Default().VerifyCode 校验通过（scene=register）
//   - inviteCode 必填且必须能定位到一个 status=1 的会员，否则拒绝
//   - 同手机号已注册时拒绝
//   - 注册成功后同事务创建：member_user + 三个 member_wallet（type=1/2/3，余额 0）
//   - 自动生成本人 invite_code（base32 8 位，全局唯一）
//   - 自动登录并返回 token
type RegisterReq struct {
	g.Meta     `path:"/auth/register" method:"post" tags:"会员-认证" summary:"注册"`
	Phone      string `json:"phone" v:"required|phone-loose#手机号不能为空|手机号格式不正确" dc:"手机号"`
	SmsCode    string `json:"smsCode" v:"required|length:4,6#验证码不能为空|验证码长度不正确" dc:"短信验证码"`
	Password   string `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32位" dc:"登录密码"`
	InviteCode string `json:"inviteCode" v:"required|length:1,32#邀请码不能为空|邀请码长度不正确" dc:"邀请码（必填）"`
	Nickname   string `json:"nickname" v:"max-length:50" dc:"昵称（可选，默认手机号脱敏）"`
}

// RegisterRes 注册响应（自动登录）。
type RegisterRes struct {
	g.Meta `mime:"application/json"`
	*LoginResult
}

// ----- 账号密码登录 -----

// LoginReq 登录请求。
//
// account 支持以下三种值，按顺序匹配：
//  1. 11 位手机号 → member_user.phone
//  2. 邀请码（区分大小写）→ member_user.invite_code
//  3. 自定义 username → member_user.username
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"会员-认证" summary:"账号密码登录"`
	Account  string `json:"account" v:"required|max-length:64#账号不能为空|账号长度不能超过64位" dc:"账号（手机号 / 邀请码 / 用户名）"`
	Password string `json:"password" v:"required|length:6,32#密码不能为空|密码长度6-32位" dc:"密码"`
}

// LoginRes 登录响应。
type LoginRes struct {
	g.Meta `mime:"application/json"`
	*LoginResult
}

// LoginResult 登录公共返回体（注册/登录复用）。
type LoginResult struct {
	Token       string `json:"token" dc:"会员 JWT，前端需放到 Authorization Bearer"`
	MemberID    string `json:"memberId" dc:"会员 ID（雪花字符串）"`
	Phone       string `json:"phone" dc:"手机号"`
	Nickname    string `json:"nickname" dc:"昵称"`
	Avatar      string `json:"avatar" dc:"头像"`
	InviteCode  string `json:"inviteCode" dc:"我的邀请码"`
	LevelID     string `json:"levelId" dc:"当前等级 ID"`
	IsQualified int    `json:"isQualified" dc:"仓库资格:0=已失效,1=有效"`
}

// ----- 找回密码 -----

// ForgetPasswordReq 找回密码请求。
//
// 业务规则：
//   - phone + smsCode 必须经过 sms.Default().VerifyCode 通过（scene=forget_password）
//   - 手机号未注册时直接报错（注册场景不会泄露注册状态，找回场景告知用户更友好）
//   - 重置后旧 token 不主动作废（依赖 jwt 自然过期）
type ForgetPasswordReq struct {
	g.Meta      `path:"/auth/forget-password" method:"post" tags:"会员-认证" summary:"找回密码"`
	Phone       string `json:"phone" v:"required|phone-loose#手机号不能为空|手机号格式不正确" dc:"手机号"`
	SmsCode     string `json:"smsCode" v:"required|length:4,6#验证码不能为空" dc:"短信验证码"`
	NewPassword string `json:"newPassword" v:"required|length:6,32#新密码不能为空|密码长度6-32位" dc:"新密码"`
}

// ForgetPasswordRes 找回密码响应。
type ForgetPasswordRes struct {
	g.Meta `mime:"application/json"`
}

// ----- 邀请码反查 -----

// InvitePreviewReq 通过邀请码预览上级信息（注册前 H5 调用展示"邀请人：xxx"）。
type InvitePreviewReq struct {
	g.Meta     `path:"/auth/invite-preview" method:"get" tags:"会员-认证" summary:"邀请码预览"`
	InviteCode string `json:"inviteCode" v:"required|length:1,32#邀请码不能为空" dc:"邀请码"`
}

// InvitePreviewRes 邀请人简要信息。
type InvitePreviewRes struct {
	g.Meta   `mime:"application/json"`
	Found    bool   `json:"found" dc:"邀请码是否有效"`
	Nickname string `json:"nickname,omitempty" dc:"邀请人昵称（脱敏可选）"`
	Avatar   string `json:"avatar,omitempty" dc:"邀请人头像"`
}
