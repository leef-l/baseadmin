// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberContractTemplateDao is the data access object for the table member_contract_template.
type MemberContractTemplateDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  MemberContractTemplateColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// MemberContractTemplateColumns defines and stores column names for the table member_contract_template.
type MemberContractTemplateColumns struct {
	Id           string // 模板ID（Snowflake）
	TemplateName string // 模板名称|search:like|keyword:on|priority:100
	TemplateType string // 模板类型|search:select|enum:register=注册协议,upgrade=升级协议,custom=自定义
	Content      string // 模板正文（HTML，支持{{nickname}}{{phone}}{{date}}等占位符）|search:off
	IsDefault    string // 是否默认模板:0=否,1=是|search:select
	Remark       string // 备注|search:off
	Sort         string // 排序（升序）
	Status       string // 状态:0=关闭,1=开启|search:select
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间
}

// memberContractTemplateColumns holds the columns for the table member_contract_template.
var memberContractTemplateColumns = MemberContractTemplateColumns{
	Id:           "id",
	TemplateName: "template_name",
	TemplateType: "template_type",
	Content:      "content",
	IsDefault:    "is_default",
	Remark:       "remark",
	Sort:         "sort",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewMemberContractTemplateDao creates and returns a new DAO object for table data access.
func NewMemberContractTemplateDao(handlers ...gdb.ModelHandler) *MemberContractTemplateDao {
	return &MemberContractTemplateDao{
		group:    "default",
		table:    "member_contract_template",
		columns:  memberContractTemplateColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberContractTemplateDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberContractTemplateDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberContractTemplateDao) Columns() MemberContractTemplateColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberContractTemplateDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberContractTemplateDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberContractTemplateDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
