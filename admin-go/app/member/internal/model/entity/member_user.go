// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberUser is the golang structure for table member_user.
type MemberUser struct {
	Id            uint64      `orm:"id"              description:"会员ID（Snowflake）"`                             // 会员ID（Snowflake）
	ParentId      uint64      `orm:"parent_id"       description:"上级会员"`                                        // 上级会员
	Username      string      `orm:"username"        description:"用户名（登录账号）|search:eq|keyword:on|priority:100"` // 用户名（登录账号）|search:eq|keyword:on|priority:100
	Password      string      `orm:"password"        description:"密码（bcrypt加密）"`                                // 密码（bcrypt加密）
	Nickname      string      `orm:"nickname"        description:"昵称|search:like|keyword:on|priority:95"`       // 昵称|search:like|keyword:on|priority:95
	Phone         string      `orm:"phone"           description:"手机号|search:eq|keyword:on|priority:90"`        // 手机号|search:eq|keyword:on|priority:90
	Avatar        string      `orm:"avatar"          description:"头像"`                                          // 头像
	RealName      string      `orm:"real_name"       description:"真实姓名|search:like|keyword:on"`                 // 真实姓名|search:like|keyword:on
	LevelId       uint64      `orm:"level_id"        description:"当前等级|ref:member_level.name|search:select"`    // 当前等级|ref:member_level.name|search:select
	LevelExpireAt *gtime.Time `orm:"level_expire_at" description:"等级到期时间"`                                      // 等级到期时间
	TeamCount     uint        `orm:"team_count"      description:"团队总人数"`                                       // 团队总人数
	DirectCount   uint        `orm:"direct_count"    description:"直推人数"`                                        // 直推人数
	ActiveCount   uint        `orm:"active_count"    description:"有效用户数"`                                       // 有效用户数
	TeamTurnover  uint64      `orm:"team_turnover"   description:"团队总营业额（分）"`                                   // 团队总营业额（分）
	IsActive      int         `orm:"is_active"       description:"是否激活:0=未激活,1=已激活|search:select"`              // 是否激活:0=未激活,1=已激活|search:select
	IsQualified   int         `orm:"is_qualified"    description:"仓库资格:0=已失效,1=有效|search:select"`               // 仓库资格:0=已失效,1=有效|search:select
	InviteCode    string      `orm:"invite_code"     description:"邀请码|search:eq"`                               // 邀请码|search:eq
	RegisterIp    string      `orm:"register_ip"     description:"注册IP"`                                        // 注册IP
	LastLoginAt   *gtime.Time `orm:"last_login_at"   description:"最后登录时间"`                                      // 最后登录时间
	Remark        string      `orm:"remark"          description:"备注|search:off"`                               // 备注|search:off
	Sort          int         `orm:"sort"            description:"排序（升序）"`                                      // 排序（升序）
	Status        int         `orm:"status"          description:"状态:0=冻结,1=正常|search:select"`                  // 状态:0=冻结,1=正常|search:select
	TenantId      uint64      `orm:"tenant_id"       description:"租户"`                                          // 租户
	MerchantId    uint64      `orm:"merchant_id"     description:"商户"`                                          // 商户
	CreatedBy     uint64      `orm:"created_by"      description:"创建人ID"`                                       // 创建人ID
	DeptId        uint64      `orm:"dept_id"         description:"所属部门ID"`                                      // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"      description:"创建时间"`                                        // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"      description:"更新时间"`                                        // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"      description:"软删除时间，非 NULL 表示已删除"`                          // 软删除时间，非 NULL 表示已删除
}
