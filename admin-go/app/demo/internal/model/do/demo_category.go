// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCategory is the golang structure of table demo_category for DAO operations like Where/Data.
type DemoCategory struct {
	g.Meta     `orm:"table:demo_category, do:true"`
	Id         any         // 分类ID（Snowflake）
	ParentId   any         // 父分类
	Name       any         // 分类名称|search:like|keyword:on|priority:95
	Icon       any         // 图标
	Sort       any         // 排序（升序）
	Status     any         // 状态:0=禁用,1=启用
	TenantId   any         // 租户
	MerchantId any         // 商户
	CreatedBy  any         // 创建人ID
	DeptId     any         // 所属部门ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间，非 NULL 表示已删除
}
