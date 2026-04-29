// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Daemon is the golang structure of table system_daemon for DAO operations like Where/Data.
type Daemon struct {
	g.Meta       `orm:"table:system_daemon, do:true"`
	Id           any         // 守护进程ID（Snowflake）
	Name         any         // 显示名称
	Program      any         // Supervisor进程名
	Command      any         // 启动命令
	Directory    any         // 运行目录
	RunUser      any         // 运行用户
	Numprocs     any         // 进程数量
	Priority     any         // 启动优先级
	Autostart    any         // 是否随Supervisor启动
	Autorestart  any         // 异常退出是否自动重启
	Startsecs    any         // 启动稳定秒数
	Startretries any         // 启动重试次数
	StopSignal   any         // 停止信号
	Environment  any         // 环境变量，Supervisor environment格式
	Remark       any         // 备注
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
