// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberBusinessConfig is the golang structure of table member_business_config for DAO operations like Where/Data.
type MemberBusinessConfig struct {
	g.Meta     `orm:"table:member_business_config, do:true"`
	Id         any         // 配置ID（Snowflake）
	ConfigKey  any         // 配置键|search:eq
	Payload    any         // 业务配置JSON（进货时间窗/寄售时间窗/工作日/返佣比例等）|search:off
	Remark     any         // 备注|search:off
	TenantId   any         // 租户
	MerchantId any         // 商户
	CreatedBy  any         // 创建人ID
	DeptId     any         // 所属部门ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间，非 NULL 表示已删除
}
