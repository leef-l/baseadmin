// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TenantPlanDao is the data access object for the table system_tenant_plan.
type TenantPlanDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  TenantPlanColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// TenantPlanColumns defines and stores column names for the table system_tenant_plan.
type TenantPlanColumns struct {
	Id         string // 订阅ID（Snowflake）
	TenantId   string // 租户ID
	PlanId     string // 套餐ID
	StartAt    string // 生效时间
	ExpireAt   string // 到期时间
	Status     string // 状态:0=关闭,1=开启
	Remark     string // 备注
	CreatedBy  string // 创建人ID
	DeptId     string // 所属部门ID
	MerchantId string // 商户
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间
}

// systemTenantPlanColumns holds the columns for the table system_tenant_plan.
var systemTenantPlanColumns = TenantPlanColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	PlanId:     "plan_id",
	StartAt:    "start_at",
	ExpireAt:   "expire_at",
	Status:     "status",
	Remark:     "remark",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	MerchantId: "merchant_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewTenantPlanDao creates and returns a new DAO object for table data access.
func NewTenantPlanDao(handlers ...gdb.ModelHandler) *TenantPlanDao {
	return &TenantPlanDao{
		group:    "default",
		table:    "system_tenant_plan",
		columns:  systemTenantPlanColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TenantPlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TenantPlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TenantPlanDao) Columns() TenantPlanColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TenantPlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TenantPlanDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TenantPlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
