// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoAuditLogDao is the data access object for the table demo_audit_log.
type DemoAuditLogDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  DemoAuditLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// DemoAuditLogColumns defines and stores column names for the table demo_audit_log.
type DemoAuditLogColumns struct {
	Id          string // 审计日志ID（Snowflake）
	LogNo       string // 日志编号|search:eq|priority:100
	OperatorId  string // 操作人|ref:system_users.username
	Action      string // 动作:1=创建,2=修改,3=删除,4=导出,5=导入
	TargetType  string // 对象类型:1=客户,2=商品,3=订单,4=工单
	TargetCode  string // 对象编号|search:eq|priority:88
	RequestJson string // 请求JSON
	Result      string // 结果:0=失败,1=成功
	ClientIp    string // 客户端IP|search:eq|priority:80
	OccurredAt  string // 发生时间
	Remark      string // 备注|keyword:only
	TenantId    string // 租户
	MerchantId  string // 商户
	CreatedBy   string // 创建人ID
	DeptId      string // 所属部门ID
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 软删除时间，非 NULL 表示已删除
}

// demoAuditLogColumns holds the columns for the table demo_audit_log.
var demoAuditLogColumns = DemoAuditLogColumns{
	Id:          "id",
	LogNo:       "log_no",
	OperatorId:  "operator_id",
	Action:      "action",
	TargetType:  "target_type",
	TargetCode:  "target_code",
	RequestJson: "request_json",
	Result:      "result",
	ClientIp:    "client_ip",
	OccurredAt:  "occurred_at",
	Remark:      "remark",
	TenantId:    "tenant_id",
	MerchantId:  "merchant_id",
	CreatedBy:   "created_by",
	DeptId:      "dept_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewDemoAuditLogDao creates and returns a new DAO object for table data access.
func NewDemoAuditLogDao(handlers ...gdb.ModelHandler) *DemoAuditLogDao {
	return &DemoAuditLogDao{
		group:    "default",
		table:    "demo_audit_log",
		columns:  demoAuditLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoAuditLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoAuditLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoAuditLogDao) Columns() DemoAuditLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoAuditLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoAuditLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoAuditLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
