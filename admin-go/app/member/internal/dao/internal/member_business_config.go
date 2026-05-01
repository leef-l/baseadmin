// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberBusinessConfigDao is the data access object for the table member_business_config.
type MemberBusinessConfigDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  MemberBusinessConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// MemberBusinessConfigColumns defines and stores column names for the table member_business_config.
type MemberBusinessConfigColumns struct {
	Id         string // 配置ID（Snowflake）
	ConfigKey  string // 配置键|search:eq
	Payload    string // 业务配置JSON（进货时间窗/寄售时间窗/工作日/返佣比例等）|search:off
	Remark     string // 备注|search:off
	TenantId   string // 租户
	MerchantId string // 商户
	CreatedBy  string // 创建人ID
	DeptId     string // 所属部门ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间，非 NULL 表示已删除
}

// memberBusinessConfigColumns holds the columns for the table member_business_config.
var memberBusinessConfigColumns = MemberBusinessConfigColumns{
	Id:         "id",
	ConfigKey:  "config_key",
	Payload:    "payload",
	Remark:     "remark",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewMemberBusinessConfigDao creates and returns a new DAO object for table data access.
func NewMemberBusinessConfigDao(handlers ...gdb.ModelHandler) *MemberBusinessConfigDao {
	return &MemberBusinessConfigDao{
		group:    "default",
		table:    "member_business_config",
		columns:  memberBusinessConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberBusinessConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberBusinessConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberBusinessConfigDao) Columns() MemberBusinessConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberBusinessConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberBusinessConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberBusinessConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
