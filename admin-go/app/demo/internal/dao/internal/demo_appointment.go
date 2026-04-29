// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoAppointmentDao is the data access object for the table demo_appointment.
type DemoAppointmentDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  DemoAppointmentColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// DemoAppointmentColumns defines and stores column names for the table demo_appointment.
type DemoAppointmentColumns struct {
	Id            string // 预约ID（Snowflake）
	AppointmentNo string // 预约编号|search:eq|priority:100
	CustomerId    string // 客户
	Subject       string // 预约主题|search:like|keyword:on|priority:95
	AppointmentAt string // 预约时间
	ContactPhone  string // 联系电话|search:like|keyword:on|priority:90
	Address       string // 预约地址|keyword:only
	Remark        string // 备注|keyword:only
	Status        string // 状态:0=待确认,1=已确认,2=已完成,3=已取消
	TenantId      string // 租户
	MerchantId    string // 商户
	CreatedBy     string // 创建人ID
	DeptId        string // 所属部门ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 软删除时间，非 NULL 表示已删除
}

// demoAppointmentColumns holds the columns for the table demo_appointment.
var demoAppointmentColumns = DemoAppointmentColumns{
	Id:            "id",
	AppointmentNo: "appointment_no",
	CustomerId:    "customer_id",
	Subject:       "subject",
	AppointmentAt: "appointment_at",
	ContactPhone:  "contact_phone",
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

// NewDemoAppointmentDao creates and returns a new DAO object for table data access.
func NewDemoAppointmentDao(handlers ...gdb.ModelHandler) *DemoAppointmentDao {
	return &DemoAppointmentDao{
		group:    "default",
		table:    "demo_appointment",
		columns:  demoAppointmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoAppointmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoAppointmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoAppointmentDao) Columns() DemoAppointmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoAppointmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoAppointmentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoAppointmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
