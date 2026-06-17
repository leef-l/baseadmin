// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SystemTenantDao is the data access object for the table system_tenant.
type SystemTenantDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  SystemTenantColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// SystemTenantColumns defines and stores column names for the table system_tenant.
type SystemTenantColumns struct {
	Id           string // 租户ID（Snowflake）
	Name         string // 租户名称
	Code         string // 租户编码
	ContactName  string // 联系人
	ContactPhone string // 联系电话
	Domain       string // 租户域名
	ExpireAt     string // 到期时间
	Status       string // 状态:0=关闭,1=开启
	Remark       string // 备注
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// systemTenantColumns holds the columns for the table system_tenant.
var systemTenantColumns = SystemTenantColumns{
	Id:           "id",
	Name:         "name",
	Code:         "code",
	ContactName:  "contact_name",
	ContactPhone: "contact_phone",
	Domain:       "domain",
	ExpireAt:     "expire_at",
	Status:       "status",
	Remark:       "remark",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewSystemTenantDao creates and returns a new DAO object for table data access.
func NewSystemTenantDao(handlers ...gdb.ModelHandler) *SystemTenantDao {
	return &SystemTenantDao{
		group:    "default",
		table:    "system_tenant",
		columns:  systemTenantColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SystemTenantDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SystemTenantDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SystemTenantDao) Columns() SystemTenantColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SystemTenantDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SystemTenantDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SystemTenantDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
