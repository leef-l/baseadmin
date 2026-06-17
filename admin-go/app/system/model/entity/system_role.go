// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemRole is the golang structure for table system_role.
type SystemRole struct {
	Id         uint64      `json:"id"         orm:"id"          description:"角色ID（Snowflake）"`                      // 角色ID（Snowflake）
	ParentId   uint64      `json:"parentId"   orm:"parent_id"   description:"上级角色ID，0 表示顶级角色"`                      // 上级角色ID，0 表示顶级角色
	Title      string      `json:"title"      orm:"title"       description:"角色名称"`                                 // 角色名称
	DataScope  int         `json:"dataScope"  orm:"data_scope"  description:"数据范围:1=全部,2=本部门及以下,3=本部门,4=仅本人,5=自定义"` // 数据范围:1=全部,2=本部门及以下,3=本部门,4=仅本人,5=自定义
	IsAdmin    int         `json:"isAdmin"    orm:"is_admin"    description:"是否超级管理员:0=否,1=是"`                      // 是否超级管理员:0=否,1=是
	Sort       int         `json:"sort"       orm:"sort"        description:"排序（升序）"`                               // 排序（升序）
	Status     int         `json:"status"     orm:"status"      description:"状态:0=关闭,1=开启"`                         // 状态:0=关闭,1=开启
	CreatedBy  uint64      `json:"createdBy"  orm:"created_by"  description:"创建人ID"`                                // 创建人ID
	DeptId     uint64      `json:"deptId"     orm:"dept_id"     description:"所属部门ID"`                               // 所属部门ID
	TenantId   uint64      `json:"tenantId"   orm:"tenant_id"   description:"租户"`                                   // 租户
	MerchantId uint64      `json:"merchantId" orm:"merchant_id" description:"商户"`                                   // 商户
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`                                 // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`                                 // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"软删除时间，非 NULL 表示已删除"`                   // 软删除时间，非 NULL 表示已删除
}
