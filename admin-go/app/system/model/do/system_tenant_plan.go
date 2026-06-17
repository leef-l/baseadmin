// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemTenantPlan is the golang structure of table system_tenant_plan for DAO operations like Where/Data.
type SystemTenantPlan struct {
	g.Meta     `orm:"table:system_tenant_plan, do:true"`
	Id         any         // 订阅ID（Snowflake）
	TenantId   any         // 租户ID
	PlanId     any         // 套餐ID
	StartAt    *gtime.Time // 生效时间
	ExpireAt   *gtime.Time // 到期时间
	Status     any         // 状态:0=关闭,1=开启
	Remark     any         // 备注
	CreatedBy  any         // 创建人ID
	DeptId     any         // 所属部门ID
	MerchantId any         // 商户
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间
}
