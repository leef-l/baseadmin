// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tenant is the golang structure of table system_tenant for DAO operations like Where/Data.
type Tenant struct {
	g.Meta       `orm:"table:system_tenant, do:true"`
	Id           any         // 租户ID（Snowflake）
	Name         any         // 租户名称
	Code         any         // 租户编码
	ContactName  any         // 联系人
	ContactPhone any         // 联系电话
	Domain       any         // 租户域名
	ExpireAt     *gtime.Time // 到期时间
	Status       any         // 状态:0=关闭,1=开启
	Remark       any         // 备注
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
