package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 个人信息 -----

// MeProfileReq 获取当前会员资料。
type MeProfileReq struct {
	g.Meta `path:"/me/profile" method:"get" tags:"会员-个人中心" summary:"获取个人资料"`
}

// MeProfileRes 个人资料。
type MeProfileRes struct {
	g.Meta        `mime:"application/json"`
	MemberID      string `json:"memberId"`
	Phone         string `json:"phone"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	RealName      string `json:"realName"`
	InviteCode    string `json:"inviteCode" dc:"我的邀请码"`
	ParentID      string `json:"parentId" dc:"上级会员 ID"`
	LevelID       string `json:"levelId"`
	LevelName     string `json:"levelName"`
	LevelExpireAt string `json:"levelExpireAt"`
	IsActive      int    `json:"isActive"`
	IsQualified   int    `json:"isQualified"`
	TeamCount     int    `json:"teamCount"`
	DirectCount   int    `json:"directCount"`
	ActiveCount   int    `json:"activeCount"`
	TeamTurnover       int64  `json:"teamTurnover" dc:"团队总营业额（分）"`
	InviteURL          string `json:"inviteUrl" dc:"邀请链接（含邀请码参数）"`
	DailyPurchaseLimit int    `json:"dailyPurchaseLimit" dc:"每日限购单数"`
	TodayPurchaseCount int    `json:"todayPurchaseCount" dc:"今日已购单数"`
	TotalPurchaseCount int    `json:"totalPurchaseCount" dc:"历史累计购单数"`
}

// ----- 修改资料 -----

// MeUpdateReq 修改昵称、头像、真实姓名。
type MeUpdateReq struct {
	g.Meta   `path:"/me/update" method:"put" tags:"会员-个人中心" summary:"修改个人资料"`
	Nickname string `json:"nickname" v:"max-length:50" dc:"昵称"`
	Avatar   string `json:"avatar" v:"max-length:500" dc:"头像 URL"`
	RealName string `json:"realName" v:"max-length:50" dc:"真实姓名（实名）"`
}

// MeUpdateRes 修改响应。
type MeUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ----- 修改密码 -----

// MeChangePasswordReq 已登录态修改密码（旧密码 + 新密码）。
type MeChangePasswordReq struct {
	g.Meta      `path:"/me/change-password" method:"post" tags:"会员-个人中心" summary:"修改密码"`
	OldPassword string `json:"oldPassword" v:"required#旧密码不能为空" dc:"旧密码"`
	NewPassword string `json:"newPassword" v:"required|length:6,32#新密码不能为空|密码长度6-32位" dc:"新密码"`
}

// MeChangePasswordRes 修改密码响应。
type MeChangePasswordRes struct {
	g.Meta `mime:"application/json"`
}

// ----- 修改手机号 -----

// MeChangePhoneReq 已登录态换绑手机号（新手机号 + 验证码）。
type MeChangePhoneReq struct {
	g.Meta   `path:"/me/change-phone" method:"post" tags:"会员-个人中心" summary:"修改手机号"`
	NewPhone string `json:"newPhone" v:"required|phone-loose#新手机号不能为空|新手机号格式不正确" dc:"新手机号"`
	SmsCode  string `json:"smsCode" v:"required|length:4,6#验证码不能为空" dc:"新手机号收到的验证码（scene=change_phone）"`
}

// MeChangePhoneRes 修改手机号响应。
type MeChangePhoneRes struct {
	g.Meta `mime:"application/json"`
}

// ----- 钱包 -----

// MeWalletsReq 获取三钱包余额。
type MeWalletsReq struct {
	g.Meta `path:"/me/wallets" method:"get" tags:"会员-个人中心" summary:"获取我的三钱包余额"`
}

// MeWalletsRes 三钱包余额。
type MeWalletsRes struct {
	g.Meta  `mime:"application/json"`
	Coupon  WalletInfo `json:"coupon" dc:"优惠券余额（仅可购买商城商品）"`
	Reward  WalletInfo `json:"reward" dc:"奖金余额（寄售卖出收入）"`
	Promote WalletInfo `json:"promote" dc:"推广奖余额"`
}

// WalletInfo 单个钱包余额（单位元，前端直接展示）。
type WalletInfo struct {
	Balance      string `json:"balance" dc:"当前余额（元，保留两位小数字符串）"`
	BalanceCent  int64  `json:"balanceCent" dc:"当前余额（分）"`
	TotalIncome  string `json:"totalIncome" dc:"累计收入（元）"`
	TotalExpense string `json:"totalExpense" dc:"累计支出（元）"`
	FrozenAmount string `json:"frozenAmount" dc:"冻结金额（元）"`
}

// MeWalletLogsReq 钱包流水分页查询。
type MeWalletLogsReq struct {
	g.Meta     `path:"/me/wallet-logs" method:"get" tags:"会员-个人中心" summary:"我的钱包流水"`
	WalletType int `json:"walletType" v:"in:0,1,2,3" dc:"钱包类型 0=全部 1=优惠券 2=奖金 3=推广奖"`
	PageNum    int `json:"pageNum" d:"1"`
	PageSize   int `json:"pageSize" d:"20"`
}

// MeWalletLogsRes 钱包流水列表。
type MeWalletLogsRes struct {
	g.Meta `mime:"application/json"`
	Total  int                `json:"total"`
	List   []*WalletLogRecord `json:"list"`
}

// WalletLogRecord 单条流水。
type WalletLogRecord struct {
	ID             string `json:"id"`
	WalletType     int    `json:"walletType"`
	WalletTypeText string `json:"walletTypeText"`
	ChangeType     int    `json:"changeType"`
	ChangeTypeText string `json:"changeTypeText"`
	ChangeAmount   string `json:"changeAmount" dc:"变动金额（元，带正负号）"`
	BeforeBalance  string `json:"beforeBalance"`
	AfterBalance   string `json:"afterBalance"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	Remark         string `json:"remark"`
	CreatedAt      string `json:"createdAt"`
}

// ----- 团队 -----

// MeTeamReq 获取我的团队（直推 / 全部）。
type MeTeamReq struct {
	g.Meta   `path:"/me/team" method:"get" tags:"会员-个人中心" summary:"我的团队"`
	Scope    string `json:"scope" v:"in:direct,all" d:"direct" dc:"direct=仅直推 all=全部团队"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"20"`
}

// MeTeamRes 团队列表。
type MeTeamRes struct {
	g.Meta `mime:"application/json"`
	Total  int               `json:"total"`
	List   []*TeamMemberItem `json:"list"`
}

// TeamMemberItem 团队成员简要。
type TeamMemberItem struct {
	MemberID    string `json:"memberId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone" dc:"手机号脱敏"`
	LevelName   string `json:"levelName"`
	IsQualified int    `json:"isQualified"`
	JoinedAt    string `json:"joinedAt"`
}
