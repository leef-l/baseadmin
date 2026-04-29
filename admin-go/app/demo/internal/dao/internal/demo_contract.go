// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoContractDao is the data access object for the table demo_contract.
type DemoContractDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DemoContractColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DemoContractColumns defines and stores column names for the table demo_contract.
type DemoContractColumns struct {
	Id             string // 合同ID（Snowflake）
	ContractNo     string // 合同编号|search:eq|priority:100
	CustomerId     string // 客户
	OrderId        string // 订单
	Title          string // 合同标题|search:like|keyword:on|priority:95
	ContractFile   string // 合同文件
	SignImage      string // 签章图片
	ContractAmount string // 合同金额（分）
	SignPassword   string // 签署密码
	SignedAt       string // 签署时间
	ExpiresAt      string // 到期时间
	Status         string // 状态:0=待审核,1=已通过,2=已拒绝,3=已取消
	TenantId       string // 租户
	MerchantId     string // 商户
	CreatedBy      string // 创建人ID
	DeptId         string // 所属部门ID
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 软删除时间，非 NULL 表示已删除
}

// demoContractColumns holds the columns for the table demo_contract.
var demoContractColumns = DemoContractColumns{
	Id:             "id",
	ContractNo:     "contract_no",
	CustomerId:     "customer_id",
	OrderId:        "order_id",
	Title:          "title",
	ContractFile:   "contract_file",
	SignImage:      "sign_image",
	ContractAmount: "contract_amount",
	SignPassword:   "sign_password",
	SignedAt:       "signed_at",
	ExpiresAt:      "expires_at",
	Status:         "status",
	TenantId:       "tenant_id",
	MerchantId:     "merchant_id",
	CreatedBy:      "created_by",
	DeptId:         "dept_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewDemoContractDao creates and returns a new DAO object for table data access.
func NewDemoContractDao(handlers ...gdb.ModelHandler) *DemoContractDao {
	return &DemoContractDao{
		group:    "default",
		table:    "demo_contract",
		columns:  demoContractColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoContractDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoContractDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoContractDao) Columns() DemoContractColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoContractDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoContractDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoContractDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
