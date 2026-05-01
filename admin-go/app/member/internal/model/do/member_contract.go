// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberContract is the golang structure of table member_contract for DAO operations like Where/Data.
type MemberContract struct {
	g.Meta          `orm:"table:member_contract, do:true"`
	Id              any         // 合同ID（Snowflake）
	UserId          any         // 会员|ref:member_user.nickname|search:select
	ContractNo      any         // 合同编号|search:eq|keyword:on|priority:100
	ContractType    any         // 合同类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	TemplateId      any         // 模板|ref:member_contract_template.template_name
	RelatedId       any         // 关联业务ID（订单/升级记录等）
	SignedHtml      any         // 签署时实际渲染的 HTML（已替换占位符，含签名图）|search:off
	SignatureImage  any         // 手写签名 base64 PNG（data:image/png;base64,...）|search:off
	SignedAt        *gtime.Time // 签署时间|search:date
	SignedIp        any         // 签署IP
	SignedUserAgent any         // UA
	PdfPath         any         // PDF存储路径（OSS或本地）
	PdfStatus       any         // PDF生成状态:0=未生成,1=生成中,2=已生成,3=失败|search:select
	PdfError        any         // PDF生成错误信息
	Remark          any         // 备注|search:off
	Sort            any         // 排序
	Status          any         // 状态:0=作废,1=正常|search:select
	TenantId        any         // 租户
	MerchantId      any         // 商户
	CreatedBy       any         // 创建人ID
	DeptId          any         // 所属部门ID
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
	DeletedAt       *gtime.Time // 软删除时间
}
