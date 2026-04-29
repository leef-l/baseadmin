// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Daemon is the golang structure for table daemon.
type Daemon struct {
	Id           uint64      `orm:"id"           description:"守护进程ID（Snowflake）"`             // 守护进程ID（Snowflake）
	Name         string      `orm:"name"         description:"显示名称"`                          // 显示名称
	Program      string      `orm:"program"      description:"Supervisor进程名"`                 // Supervisor进程名
	Command      string      `orm:"command"      description:"启动命令"`                          // 启动命令
	Directory    string      `orm:"directory"    description:"运行目录"`                          // 运行目录
	RunUser      string      `orm:"run_user"     description:"运行用户"`                          // 运行用户
	Numprocs     uint        `orm:"numprocs"     description:"进程数量"`                          // 进程数量
	Priority     uint        `orm:"priority"     description:"启动优先级"`                         // 启动优先级
	Autostart    int         `orm:"autostart"    description:"是否随Supervisor启动"`               // 是否随Supervisor启动
	Autorestart  int         `orm:"autorestart"  description:"异常退出是否自动重启"`                    // 异常退出是否自动重启
	Startsecs    uint        `orm:"startsecs"    description:"启动稳定秒数"`                        // 启动稳定秒数
	Startretries uint        `orm:"startretries" description:"启动重试次数"`                        // 启动重试次数
	StopSignal   string      `orm:"stop_signal"  description:"停止信号"`                          // 停止信号
	Environment  string      `orm:"environment"  description:"环境变量，Supervisor environment格式"` // 环境变量，Supervisor environment格式
	Remark       string      `orm:"remark"       description:"备注"`                            // 备注
	CreatedBy    uint64      `orm:"created_by"   description:"创建人ID"`                         // 创建人ID
	DeptId       uint64      `orm:"dept_id"      description:"所属部门ID"`                        // 所属部门ID
	TenantId     uint64      `orm:"tenant_id"    description:"租户"`                            // 租户
	MerchantId   uint64      `orm:"merchant_id"  description:"商户"`                            // 商户
	CreatedAt    *gtime.Time `orm:"created_at"   description:"创建时间"`                          // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"   description:"更新时间"`                          // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"   description:"软删除时间，非 NULL 表示已删除"`            // 软删除时间，非 NULL 表示已删除
}
