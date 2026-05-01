// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberContractTemplate is the golang structure for table member_contract_template.
type MemberContractTemplate struct {
	Id           uint64      `orm:"id"            description:"模板ID（Snowflake）"`                                               // 模板ID（Snowflake）
	TemplateName string      `orm:"template_name" description:"模板名称|search:like|keyword:on|priority:100"`                      // 模板名称|search:like|keyword:on|priority:100
	TemplateType string      `orm:"template_type" description:"模板类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义"` // 模板类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	Content      string      `orm:"content"       description:"模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）|search:off"`     // 模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）|search:off
	IsDefault    int         `orm:"is_default"    description:"是否默认模板:0=否,1=是|search:select"`                                  // 是否默认模板:0=否,1=是|search:select
	Remark       string      `orm:"remark"        description:"备注|search:off"`                                                 // 备注|search:off
	Sort         int         `orm:"sort"          description:"排序（升序）"`                                                        // 排序（升序）
	Status       int         `orm:"status"        description:"状态:0=关闭,1=开启|search:select"`                                    // 状态:0=关闭,1=开启|search:select
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                                            // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                                            // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                                         // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                                        // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                                          // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                                          // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间"`                                                         // 软删除时间
}
