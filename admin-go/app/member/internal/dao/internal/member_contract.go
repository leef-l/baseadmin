// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberContractDao is the data access object for the table member_contract.
type MemberContractDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  MemberContractColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// MemberContractColumns defines and stores column names for the table member_contract.
type MemberContractColumns struct {
	Id              string // 合同ID（Snowflake）
	UserId          string // 会员|ref:member_user.nickname|search:select
	ContractNo      string // 合同编号|search:eq|keyword:on|priority:100
	ContractType    string // 合同类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	TemplateId      string // 模板|ref:member_contract_template.template_name
	RelatedId       string // 关联业务ID（订单/升级记录等）
	SignedHtml      string // 签署时实际渲染的 HTML（已替换占位符，含签名图）|search:off
	SignatureImage  string // 手写签名 base64 PNG（data:image/png;base64,...）|search:off
	SignedAt        string // 签署时间|search:date
	SignedIp        string // 签署IP
	SignedUserAgent string // UA
	PdfPath         string // PDF存储路径（OSS或本地）
	PdfStatus       string // PDF生成状态:0=未生成,1=生成中,2=已生成,3=失败|search:select
	PdfError        string // PDF生成错误信息
	Remark          string // 备注|search:off
	Sort            string // 排序
	Status          string // 状态:0=作废,1=正常|search:select
	TenantId        string // 租户
	MerchantId      string // 商户
	CreatedBy       string // 创建人ID
	DeptId          string // 所属部门ID
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	DeletedAt       string // 软删除时间
}

// memberContractColumns holds the columns for the table member_contract.
var memberContractColumns = MemberContractColumns{
	Id:              "id",
	UserId:          "user_id",
	ContractNo:      "contract_no",
	ContractType:    "contract_type",
	TemplateId:      "template_id",
	RelatedId:       "related_id",
	SignedHtml:      "signed_html",
	SignatureImage:  "signature_image",
	SignedAt:        "signed_at",
	SignedIp:        "signed_ip",
	SignedUserAgent: "signed_user_agent",
	PdfPath:         "pdf_path",
	PdfStatus:       "pdf_status",
	PdfError:        "pdf_error",
	Remark:          "remark",
	Sort:            "sort",
	Status:          "status",
	TenantId:        "tenant_id",
	MerchantId:      "merchant_id",
	CreatedBy:       "created_by",
	DeptId:          "dept_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewMemberContractDao creates and returns a new DAO object for table data access.
func NewMemberContractDao(handlers ...gdb.ModelHandler) *MemberContractDao {
	return &MemberContractDao{
		group:    "default",
		table:    "member_contract",
		columns:  memberContractColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberContractDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberContractDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberContractDao) Columns() MemberContractColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberContractDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberContractDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MemberContractDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
