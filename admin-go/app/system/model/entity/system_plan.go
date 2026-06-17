// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemPlan is the golang structure for table system_plan.
type SystemPlan struct {
	Id          uint64      `json:"id"          orm:"id"          description:"套餐ID（Snowflake）"` // 套餐ID（Snowflake）
	Name        string      `json:"name"        orm:"name"        description:"套餐名称"`            // 套餐名称
	Code        string      `json:"code"        orm:"code"        description:"套餐编码"`            // 套餐编码
	Description string      `json:"description" orm:"description" description:"套餐描述"`            // 套餐描述
	UserLimit   int         `json:"userLimit"   orm:"user_limit"  description:"用户数量上限（-1=无限）"`   // 用户数量上限（-1=无限）
	StorageMb   int         `json:"storageMb"   orm:"storage_mb"  description:"存储空间上限MB（-1=无限）"` // 存储空间上限MB（-1=无限）
	Features    string      `json:"features"    orm:"features"    description:"功能开关列表（JSON数组）"`  // 功能开关列表（JSON数组）
	Price       int         `json:"price"       orm:"price"       description:"价格（分）"`           // 价格（分）
	Sort        int         `json:"sort"        orm:"sort"        description:"排序（升序）"`          // 排序（升序）
	Status      int         `json:"status"      orm:"status"      description:"状态:0=关闭,1=开启"`    // 状态:0=关闭,1=开启
	CreatedBy   uint64      `json:"createdBy"   orm:"created_by"  description:"创建人ID"`           // 创建人ID
	DeptId      uint64      `json:"deptId"      orm:"dept_id"     description:"所属部门ID"`          // 所属部门ID
	TenantId    uint64      `json:"tenantId"    orm:"tenant_id"   description:"租户"`              // 租户
	MerchantId  uint64      `json:"merchantId"  orm:"merchant_id" description:"商户"`              // 商户
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`            // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`            // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"  description:"软删除时间"`           // 软删除时间
}
