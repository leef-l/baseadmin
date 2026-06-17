// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemTenant is the golang structure for table system_tenant.
type SystemTenant struct {
	Id           uint64      `json:"id"           orm:"id"            description:"租户ID（Snowflake）"`    // 租户ID（Snowflake）
	Name         string      `json:"name"         orm:"name"          description:"租户名称"`               // 租户名称
	Code         string      `json:"code"         orm:"code"          description:"租户编码"`               // 租户编码
	ContactName  string      `json:"contactName"  orm:"contact_name"  description:"联系人"`                // 联系人
	ContactPhone string      `json:"contactPhone" orm:"contact_phone" description:"联系电话"`               // 联系电话
	Domain       string      `json:"domain"       orm:"domain"        description:"租户域名"`               // 租户域名
	ExpireAt     *gtime.Time `json:"expireAt"     orm:"expire_at"     description:"到期时间"`               // 到期时间
	Status       int         `json:"status"       orm:"status"        description:"状态:0=关闭,1=开启"`       // 状态:0=关闭,1=开启
	Remark       string      `json:"remark"       orm:"remark"        description:"备注"`                 // 备注
	CreatedBy    uint64      `json:"createdBy"    orm:"created_by"    description:"创建人ID"`              // 创建人ID
	DeptId       uint64      `json:"deptId"       orm:"dept_id"       description:"所属部门ID"`             // 所属部门ID
	TenantId     uint64      `json:"tenantId"     orm:"tenant_id"     description:"租户"`                 // 租户
	MerchantId   uint64      `json:"merchantId"   orm:"merchant_id"   description:"商户"`                 // 商户
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`               // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`               // 更新时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"` // 软删除时间，非 NULL 表示已删除
}
