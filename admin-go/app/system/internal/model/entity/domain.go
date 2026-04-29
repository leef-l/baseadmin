// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Domain is the golang structure for table domain.
type Domain struct {
	Id           uint64      `orm:"id"            description:"域名ID（Snowflake）"`        // 域名ID（Snowflake）
	Domain       string      `orm:"domain"        description:"绑定域名"`                   // 绑定域名
	OwnerType    int         `orm:"owner_type"    description:"主体类型:1=租户,2=商户"`         // 主体类型:1=租户,2=商户
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                     // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                     // 商户
	AppCode      string      `orm:"app_code"      description:"应用编码：admin/upload/shop"` // 应用编码：admin/upload/shop
	VerifyToken  string      `orm:"verify_token"  description:"域名校验令牌"`                 // 域名校验令牌
	VerifyStatus int         `orm:"verify_status" description:"校验状态:0=未校验,1=已校验"`       // 校验状态:0=未校验,1=已校验
	SslStatus    int         `orm:"ssl_status"    description:"SSL状态:0=未配置,1=已配置"`      // SSL状态:0=未配置,1=已配置
	NginxStatus  int         `orm:"nginx_status"  description:"Nginx配置状态:0=未应用,1=已应用"`  // Nginx配置状态:0=未应用,1=已应用
	Status       int         `orm:"status"        description:"状态:0=关闭,1=开启"`           // 状态:0=关闭,1=开启
	Remark       string      `orm:"remark"        description:"备注"`                     // 备注
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                  // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                 // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                   // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                   // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`     // 软删除时间，非 NULL 表示已删除
}
