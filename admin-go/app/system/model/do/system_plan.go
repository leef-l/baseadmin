// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemPlan is the golang structure of table system_plan for DAO operations like Where/Data.
type SystemPlan struct {
	g.Meta      `orm:"table:system_plan, do:true"`
	Id          any         // 套餐ID（Snowflake）
	Name        any         // 套餐名称
	Code        any         // 套餐编码
	Description any         // 套餐描述
	UserLimit   any         // 用户数量上限（-1=无限）
	StorageMb   any         // 存储空间上限MB（-1=无限）
	Features    any         // 功能开关列表（JSON数组）
	Price       any         // 价格（分）
	Sort        any         // 排序（升序）
	Status      any         // 状态:0=关闭,1=开启
	CreatedBy   any         // 创建人ID
	DeptId      any         // 所属部门ID
	TenantId    any         // 租户
	MerchantId  any         // 商户
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 软删除时间
}
