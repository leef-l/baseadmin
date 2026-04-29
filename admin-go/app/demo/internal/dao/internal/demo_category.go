// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoCategoryDao is the data access object for the table demo_category.
type DemoCategoryDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DemoCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DemoCategoryColumns defines and stores column names for the table demo_category.
type DemoCategoryColumns struct {
	Id         string // 分类ID（Snowflake）
	ParentId   string // 父分类
	Name       string // 分类名称|search:like|keyword:on|priority:95
	Icon       string // 图标
	Sort       string // 排序（升序）
	Status     string // 状态:0=禁用,1=启用
	TenantId   string // 租户
	MerchantId string // 商户
	CreatedBy  string // 创建人ID
	DeptId     string // 所属部门ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间，非 NULL 表示已删除
}

// demoCategoryColumns holds the columns for the table demo_category.
var demoCategoryColumns = DemoCategoryColumns{
	Id:         "id",
	ParentId:   "parent_id",
	Name:       "name",
	Icon:       "icon",
	Sort:       "sort",
	Status:     "status",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewDemoCategoryDao creates and returns a new DAO object for table data access.
func NewDemoCategoryDao(handlers ...gdb.ModelHandler) *DemoCategoryDao {
	return &DemoCategoryDao{
		group:    "default",
		table:    "demo_category",
		columns:  demoCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoCategoryDao) Columns() DemoCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
