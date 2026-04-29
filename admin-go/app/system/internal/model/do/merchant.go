// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Merchant is the golang structure of table system_merchant for DAO operations like Where/Data.
type Merchant struct {
	g.Meta       `orm:"table:system_merchant, do:true"`
	Id           any         // 商户ID（Snowflake）
	TenantId     any         // 租户
	MerchantId   any         // 商户
	Name         any         // 商户名称
	Code         any         // 商户编码
	ContactName  any         // 联系人
	ContactPhone any         // 联系电话
	Address      any         // 商户地址
	Status       any         // 状态:0=关闭,1=开启
	Remark       any         // 备注
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
