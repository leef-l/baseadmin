// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberContractTemplate is the golang structure of table member_contract_template for DAO operations like Where/Data.
type MemberContractTemplate struct {
	g.Meta       `orm:"table:member_contract_template, do:true"`
	Id           any         // 模板ID（Snowflake）
	TemplateName any         // 模板名称|search:like|keyword:on|priority:100
	TemplateType any         // 模板类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	Content      any         // 模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）|search:off
	IsDefault    any         // 是否默认模板:0=否,1=是|search:select
	Remark       any         // 备注|search:off
	Sort         any         // 排序（升序）
	Status       any         // 状态:0=关闭,1=开启|search:select
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间
}
