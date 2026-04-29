// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberShopCategory is the golang structure for table member_shop_category.
type MemberShopCategory struct {
	Id         uint64      `orm:"id"          description:"ID（Snowflake）"`                            // ID（Snowflake）
	ParentId   uint64      `orm:"parent_id"   description:"上级分类"`                                     // 上级分类
	Name       string      `orm:"name"        description:"分类名称|search:like|keyword:on|priority:100"` // 分类名称|search:like|keyword:on|priority:100
	Icon       string      `orm:"icon"        description:"分类图标"`                                     // 分类图标
	Sort       int         `orm:"sort"        description:"排序（升序）"`                                   // 排序（升序）
	Status     int         `orm:"status"      description:"状态:0=关闭,1=开启|search:select"`               // 状态:0=关闭,1=开启|search:select
	TenantId   uint64      `orm:"tenant_id"   description:"租户"`                                       // 租户
	MerchantId uint64      `orm:"merchant_id" description:"商户"`                                       // 商户
	CreatedBy  uint64      `orm:"created_by"  description:"创建人ID"`                                    // 创建人ID
	DeptId     uint64      `orm:"dept_id"     description:"所属部门ID"`                                   // 所属部门ID
	CreatedAt  *gtime.Time `orm:"created_at"  description:"创建时间"`                                     // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"  description:"更新时间"`                                     // 更新时间
	DeletedAt  *gtime.Time `orm:"deleted_at"  description:"软删除时间，非 NULL 表示已删除"`                       // 软删除时间，非 NULL 表示已删除
}
