// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCategory is the golang structure for table demo_category.
type DemoCategory struct {
	Id         uint64      `orm:"id"          description:"分类ID（Snowflake）"`                         // 分类ID（Snowflake）
	ParentId   uint64      `orm:"parent_id"   description:"父分类"`                                     // 父分类
	Name       string      `orm:"name"        description:"分类名称|search:like|keyword:on|priority:95"` // 分类名称|search:like|keyword:on|priority:95
	Icon       string      `orm:"icon"        description:"图标"`                                      // 图标
	Sort       int         `orm:"sort"        description:"排序（升序）"`                                  // 排序（升序）
	Status     int         `orm:"status"      description:"状态:0=禁用,1=启用"`                            // 状态:0=禁用,1=启用
	TenantId   uint64      `orm:"tenant_id"   description:"租户"`                                      // 租户
	MerchantId uint64      `orm:"merchant_id" description:"商户"`                                      // 商户
	CreatedBy  uint64      `orm:"created_by"  description:"创建人ID"`                                   // 创建人ID
	DeptId     uint64      `orm:"dept_id"     description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt  *gtime.Time `orm:"created_at"  description:"创建时间"`                                    // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"  description:"更新时间"`                                    // 更新时间
	DeletedAt  *gtime.Time `orm:"deleted_at"  description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
