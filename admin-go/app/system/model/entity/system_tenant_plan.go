// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemTenantPlan is the golang structure for table system_tenant_plan.
type SystemTenantPlan struct {
	Id         uint64      `json:"id"         orm:"id"          description:"订阅ID（Snowflake）"` // 订阅ID（Snowflake）
	TenantId   uint64      `json:"tenantId"   orm:"tenant_id"   description:"租户ID"`            // 租户ID
	PlanId     uint64      `json:"planId"     orm:"plan_id"     description:"套餐ID"`            // 套餐ID
	StartAt    *gtime.Time `json:"startAt"    orm:"start_at"    description:"生效时间"`            // 生效时间
	ExpireAt   *gtime.Time `json:"expireAt"   orm:"expire_at"   description:"到期时间"`            // 到期时间
	Status     int         `json:"status"     orm:"status"      description:"状态:0=关闭,1=开启"`    // 状态:0=关闭,1=开启
	Remark     string      `json:"remark"     orm:"remark"      description:"备注"`              // 备注
	CreatedBy  uint64      `json:"createdBy"  orm:"created_by"  description:"创建人ID"`           // 创建人ID
	DeptId     uint64      `json:"deptId"     orm:"dept_id"     description:"所属部门ID"`          // 所属部门ID
	MerchantId uint64      `json:"merchantId" orm:"merchant_id" description:"商户"`              // 商户
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`            // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`            // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"软删除时间"`           // 软删除时间
}
