// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoOrderDao is the data access object for the table demo_order.
type DemoOrderDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DemoOrderColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DemoOrderColumns defines and stores column names for the table demo_order.
type DemoOrderColumns struct {
	Id            string // 订单ID（Snowflake）
	OrderNo       string // 订单号|search:eq|priority:100
	CustomerId    string // 客户
	ProductId     string // 商品
	Quantity      string // 购买数量
	Amount        string // 订单金额（分）
	PayStatus     string // 支付状态:0=待支付,1=已支付,2=已退款
	DeliverStatus string // 发货状态:0=待发货,1=已发货,2=已签收
	PaidAt        string // 支付时间
	DeliverAt     string // 发货时间
	ReceiverPhone string // 收货电话
	Address       string // 收货地址|keyword:only
	Remark        string // 备注|keyword:only
	Status        string // 状态:0=待确认,1=已确认,2=已取消
	TenantId      string // 租户
	MerchantId    string // 商户
	CreatedBy     string // 创建人ID
	DeptId        string // 所属部门ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 软删除时间，非 NULL 表示已删除
}

// demoOrderColumns holds the columns for the table demo_order.
var demoOrderColumns = DemoOrderColumns{
	Id:            "id",
	OrderNo:       "order_no",
	CustomerId:    "customer_id",
	ProductId:     "product_id",
	Quantity:      "quantity",
	Amount:        "amount",
	PayStatus:     "pay_status",
	DeliverStatus: "deliver_status",
	PaidAt:        "paid_at",
	DeliverAt:     "deliver_at",
	ReceiverPhone: "receiver_phone",
	Address:       "address",
	Remark:        "remark",
	Status:        "status",
	TenantId:      "tenant_id",
	MerchantId:    "merchant_id",
	CreatedBy:     "created_by",
	DeptId:        "dept_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewDemoOrderDao creates and returns a new DAO object for table data access.
func NewDemoOrderDao(handlers ...gdb.ModelHandler) *DemoOrderDao {
	return &DemoOrderDao{
		group:    "default",
		table:    "demo_order",
		columns:  demoOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoOrderDao) Columns() DemoOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
