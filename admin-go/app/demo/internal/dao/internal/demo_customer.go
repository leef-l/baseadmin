// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoCustomerDao is the data access object for the table demo_customer.
type DemoCustomerDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DemoCustomerColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DemoCustomerColumns defines and stores column names for the table demo_customer.
type DemoCustomerColumns struct {
	Id           string // 客户ID（Snowflake）
	Avatar       string // 头像
	Name         string // 客户名称|search:like|keyword:on|priority:95
	CustomerNo   string // 客户编号|search:eq|priority:100
	Phone        string // 联系电话|search:like|keyword:on|priority:90
	Email        string // 邮箱|search:like|keyword:on|priority:90
	Gender       string // 性别:0=未知,1=男,2=女
	Level        string // 等级:1=普通,2=VIP,3=付费,4=冻结
	SourceType   string // 来源:1=官网,2=小程序,3=线下,4=导入
	IsVip        string // 是否VIP:0=否,1=是
	RegisteredAt string // 注册时间
	Remark       string // 备注|search:like|keyword:only
	Status       string // 状态:0=禁用,1=启用
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// demoCustomerColumns holds the columns for the table demo_customer.
var demoCustomerColumns = DemoCustomerColumns{
	Id:           "id",
	Avatar:       "avatar",
	Name:         "name",
	CustomerNo:   "customer_no",
	Phone:        "phone",
	Email:        "email",
	Gender:       "gender",
	Level:        "level",
	SourceType:   "source_type",
	IsVip:        "is_vip",
	RegisteredAt: "registered_at",
	Remark:       "remark",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDemoCustomerDao creates and returns a new DAO object for table data access.
func NewDemoCustomerDao(handlers ...gdb.ModelHandler) *DemoCustomerDao {
	return &DemoCustomerDao{
		group:    "default",
		table:    "demo_customer",
		columns:  demoCustomerColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoCustomerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoCustomerDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoCustomerDao) Columns() DemoCustomerColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoCustomerDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoCustomerDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoCustomerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
