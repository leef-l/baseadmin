// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemDept is the golang structure for table system_dept.
type SystemDept struct {
	Id         uint64      `json:"id"         orm:"id"          description:"部门ID（Snowflake）"`    // 部门ID（Snowflake）
	ParentId   uint64      `json:"parentId"   orm:"parent_id"   description:"上级部门ID，0 表示顶级部门"`    // 上级部门ID，0 表示顶级部门
	Title      string      `json:"title"      orm:"title"       description:"部门名称"`               // 部门名称
	Username   string      `json:"username"   orm:"username"    description:"部门负责人姓名"`            // 部门负责人姓名
	Email      string      `json:"email"      orm:"email"       description:"负责人邮箱"`              // 负责人邮箱
	Sort       int         `json:"sort"       orm:"sort"        description:"排序（升序）"`             // 排序（升序）
	Status     int         `json:"status"     orm:"status"      description:"状态:0=关闭,1=开启"`       // 状态:0=关闭,1=开启
	CreatedBy  uint64      `json:"createdBy"  orm:"created_by"  description:"创建人ID"`              // 创建人ID
	DeptId     uint64      `json:"deptId"     orm:"dept_id"     description:"所属部门ID"`             // 所属部门ID
	TenantId   uint64      `json:"tenantId"   orm:"tenant_id"   description:"租户"`                 // 租户
	MerchantId uint64      `json:"merchantId" orm:"merchant_id" description:"商户"`                 // 商户
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`               // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`               // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"软删除时间，非 NULL 表示已删除"` // 软删除时间，非 NULL 表示已删除
}
