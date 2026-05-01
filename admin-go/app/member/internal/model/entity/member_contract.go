// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberContract is the golang structure for table member_contract.
type MemberContract struct {
	Id              uint64      `orm:"id"                description:"合同ID（Snowflake）"`                                               // 合同ID（Snowflake）
	UserId          uint64      `orm:"user_id"           description:"会员|ref:member_user.nickname|search:select"`                     // 会员|ref:member_user.nickname|search:select
	ContractNo      string      `orm:"contract_no"       description:"合同编号|search:eq|keyword:on|priority:100"`                        // 合同编号|search:eq|keyword:on|priority:100
	ContractType    string      `orm:"contract_type"     description:"合同类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义"` // 合同类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	TemplateId      uint64      `orm:"template_id"       description:"模板|ref:member_contract_template.template_name"`                 // 模板|ref:member_contract_template.template_name
	RelatedId       uint64      `orm:"related_id"        description:"关联业务ID（订单/升级记录等）"`                                              // 关联业务ID（订单/升级记录等）
	SignedHtml      string      `orm:"signed_html"       description:"签署时实际渲染的 HTML（已替换占位符，含签名图）|search:off"`                         // 签署时实际渲染的 HTML（已替换占位符，含签名图）|search:off
	SignatureImage  string      `orm:"signature_image"   description:"手写签名 base64 PNG（data:image/png;base64,...）|search:off"`         // 手写签名 base64 PNG（data:image/png;base64,...）|search:off
	SignedAt        *gtime.Time `orm:"signed_at"         description:"签署时间|search:date"`                                              // 签署时间|search:date
	SignedIp        string      `orm:"signed_ip"         description:"签署IP"`                                                          // 签署IP
	SignedUserAgent string      `orm:"signed_user_agent" description:"UA"`                                                            // UA
	PdfPath         string      `orm:"pdf_path"          description:"PDF存储路径（OSS或本地）"`                                               // PDF存储路径（OSS或本地）
	PdfStatus       int         `orm:"pdf_status"        description:"PDF生成状态:0=未生成,1=生成中,2=已生成,3=失败|search:select"`                  // PDF生成状态:0=未生成,1=生成中,2=已生成,3=失败|search:select
	PdfError        string      `orm:"pdf_error"         description:"PDF生成错误信息"`                                                     // PDF生成错误信息
	Remark          string      `orm:"remark"            description:"备注|search:off"`                                                 // 备注|search:off
	Sort            int         `orm:"sort"              description:"排序"`                                                            // 排序
	Status          int         `orm:"status"            description:"状态:0=作废,1=正常|search:select"`                                    // 状态:0=作废,1=正常|search:select
	TenantId        uint64      `orm:"tenant_id"         description:"租户"`                                                            // 租户
	MerchantId      uint64      `orm:"merchant_id"       description:"商户"`                                                            // 商户
	CreatedBy       uint64      `orm:"created_by"        description:"创建人ID"`                                                         // 创建人ID
	DeptId          uint64      `orm:"dept_id"           description:"所属部门ID"`                                                        // 所属部门ID
	CreatedAt       *gtime.Time `orm:"created_at"        description:"创建时间"`                                                          // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"        description:"更新时间"`                                                          // 更新时间
	DeletedAt       *gtime.Time `orm:"deleted_at"        description:"软删除时间"`                                                         // 软删除时间
}
