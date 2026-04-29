// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopCategory is the golang structure of table member_shop_category for DAO operations like Where/Data.
type MemberShopCategory struct {
	g.Meta     `orm:"table:member_shop_category, do:true"`
	Id         any         // ID（Snowflake）
	ParentId   any         // 上级分类
	Name       any         // 分类名称|search:like|keyword:on|priority:100
	Icon       any         // 分类图标
	Sort       any         // 排序（升序）
	Status     any         // 状态:0=关闭,1=开启|search:select
	TenantId   any         // 租户
	MerchantId any         // 商户
	CreatedBy  any         // 创建人ID
	DeptId     any         // 所属部门ID
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 软删除时间，非 NULL 表示已删除
}
