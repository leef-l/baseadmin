// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemMerchant is the golang structure for table system_merchant.
type SystemMerchant struct {
	Id           uint64      `json:"id"           orm:"id"            description:"商户ID（Snowflake）"`    // 商户ID（Snowflake）
	TenantId     uint64      `json:"tenantId"     orm:"tenant_id"     description:"租户"`                 // 租户
	MerchantId   uint64      `json:"merchantId"   orm:"merchant_id"   description:"商户"`                 // 商户
	Name         string      `json:"name"         orm:"name"          description:"商户名称"`               // 商户名称
	Code         string      `json:"code"         orm:"code"          description:"商户编码"`               // 商户编码
	ContactName  string      `json:"contactName"  orm:"contact_name"  description:"联系人"`                // 联系人
	ContactPhone string      `json:"contactPhone" orm:"contact_phone" description:"联系电话"`               // 联系电话
	Address      string      `json:"address"      orm:"address"       description:"商户地址"`               // 商户地址
	Status       int         `json:"status"       orm:"status"        description:"状态:0=关闭,1=开启"`       // 状态:0=关闭,1=开启
	Remark       string      `json:"remark"       orm:"remark"        description:"备注"`                 // 备注
	CreatedBy    uint64      `json:"createdBy"    orm:"created_by"    description:"创建人ID"`              // 创建人ID
	DeptId       uint64      `json:"deptId"       orm:"dept_id"       description:"所属部门ID"`             // 所属部门ID
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`               // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`               // 更新时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"` // 软删除时间，非 NULL 表示已删除
}
