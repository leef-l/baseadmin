// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoWorkOrderDao is the data access object for the table demo_work_order.
type DemoWorkOrderDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  DemoWorkOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// DemoWorkOrderColumns defines and stores column names for the table demo_work_order.
type DemoWorkOrderColumns struct {
	Id             string // 工单ID（Snowflake）
	TicketNo       string // 工单号|search:eq|priority:100
	CustomerId     string // 客户
	ProductId      string // 商品
	OrderId        string // 订单
	Title          string // 工单标题|search:like|keyword:on|priority:95
	Priority       string // 优先级:1=低,2=普通,3=高,4=紧急
	SourceType     string // 来源:1=官网,2=电话,3=微信,4=后台
	Description    string // 问题描述|search:like|keyword:only
	AttachmentFile string // 附件
	DueAt          string // 截止时间
	Status         string // 状态:0=待处理,1=进行中,2=已完成,3=已取消
	TenantId       string // 租户
	MerchantId     string // 商户
	CreatedBy      string // 创建人ID
	DeptId         string // 所属部门ID
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 软删除时间，非 NULL 表示已删除
}

// demoWorkOrderColumns holds the columns for the table demo_work_order.
var demoWorkOrderColumns = DemoWorkOrderColumns{
	Id:             "id",
	TicketNo:       "ticket_no",
	CustomerId:     "customer_id",
	ProductId:      "product_id",
	OrderId:        "order_id",
	Title:          "title",
	Priority:       "priority",
	SourceType:     "source_type",
	Description:    "description",
	AttachmentFile: "attachment_file",
	DueAt:          "due_at",
	Status:         "status",
	TenantId:       "tenant_id",
	MerchantId:     "merchant_id",
	CreatedBy:      "created_by",
	DeptId:         "dept_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewDemoWorkOrderDao creates and returns a new DAO object for table data access.
func NewDemoWorkOrderDao(handlers ...gdb.ModelHandler) *DemoWorkOrderDao {
	return &DemoWorkOrderDao{
		group:    "default",
		table:    "demo_work_order",
		columns:  demoWorkOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoWorkOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoWorkOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoWorkOrderDao) Columns() DemoWorkOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoWorkOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoWorkOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoWorkOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
