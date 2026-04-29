// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DaemonDao is the data access object for the table system_daemon.
type DaemonDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DaemonColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DaemonColumns defines and stores column names for the table system_daemon.
type DaemonColumns struct {
	Id           string // 守护进程ID（Snowflake）
	Name         string // 显示名称
	Program      string // Supervisor进程名
	Command      string // 启动命令
	Directory    string // 运行目录
	RunUser      string // 运行用户
	Numprocs     string // 进程数量
	Priority     string // 启动优先级
	Autostart    string // 是否随Supervisor启动
	Autorestart  string // 异常退出是否自动重启
	Startsecs    string // 启动稳定秒数
	Startretries string // 启动重试次数
	StopSignal   string // 停止信号
	Environment  string // 环境变量，Supervisor environment格式
	Remark       string // 备注
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// daemonColumns holds the columns for the table system_daemon.
var daemonColumns = DaemonColumns{
	Id:           "id",
	Name:         "name",
	Program:      "program",
	Command:      "command",
	Directory:    "directory",
	RunUser:      "run_user",
	Numprocs:     "numprocs",
	Priority:     "priority",
	Autostart:    "autostart",
	Autorestart:  "autorestart",
	Startsecs:    "startsecs",
	Startretries: "startretries",
	StopSignal:   "stop_signal",
	Environment:  "environment",
	Remark:       "remark",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDaemonDao creates and returns a new DAO object for table data access.
func NewDaemonDao(handlers ...gdb.ModelHandler) *DaemonDao {
	return &DaemonDao{
		group:    "default",
		table:    "system_daemon",
		columns:  daemonColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DaemonDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DaemonDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DaemonDao) Columns() DaemonColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DaemonDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DaemonDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DaemonDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
