// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberUser is the golang structure of table member_user for DAO operations like Where/Data.
type MemberUser struct {
	g.Meta        `orm:"table:member_user, do:true"`
	Id            any         // 会员ID（Snowflake）
	ParentId      any         // 上级会员
	Username      any         // 用户名（登录账号）|search:eq|keyword:on|priority:100
	Password      any         // 密码（bcrypt加密）
	Nickname      any         // 昵称|search:like|keyword:on|priority:95
	Phone         any         // 手机号|search:eq|keyword:on|priority:90
	Avatar        any         // 头像
	RealName      any         // 真实姓名|search:like|keyword:on
	LevelId       any         // 当前等级|ref:member_level.name|search:select
	LevelExpireAt *gtime.Time // 等级到期时间
	TeamCount     any         // 团队总人数
	DirectCount   any         // 直推人数
	ActiveCount   any         // 有效用户数
	TeamTurnover  any         // 团队总营业额（分）
	IsActive      any         // 是否激活:0=未激活,1=已激活|search:select
	IsQualified   any         // 仓库资格:0=已失效,1=有效|search:select
	InviteCode    any         // 邀请码|search:eq
	RegisterIp    any         // 注册IP
	LastLoginAt   *gtime.Time // 最后登录时间
	Remark        any         // 备注|search:off
	Sort          any         // 排序（升序）
	Status        any         // 状态:0=冻结,1=正常|search:select
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
