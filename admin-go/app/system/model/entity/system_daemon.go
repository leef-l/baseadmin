// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SystemDaemon is the golang structure for table system_daemon.
type SystemDaemon struct {
	Id           uint64      `json:"id"           orm:"id"           description:"守护进程ID（Snowflake）"`             // 守护进程ID（Snowflake）
	Name         string      `json:"name"         orm:"name"         description:"显示名称"`                          // 显示名称
	Program      string      `json:"program"      orm:"program"      description:"Supervisor进程名"`                 // Supervisor进程名
	Command      string      `json:"command"      orm:"command"      description:"启动命令"`                          // 启动命令
	Directory    string      `json:"directory"    orm:"directory"    description:"运行目录"`                          // 运行目录
	RunUser      string      `json:"runUser"      orm:"run_user"     description:"运行用户"`                          // 运行用户
	Numprocs     uint        `json:"numprocs"     orm:"numprocs"     description:"进程数量"`                          // 进程数量
	Priority     uint        `json:"priority"     orm:"priority"     description:"启动优先级"`                         // 启动优先级
	Autostart    int         `json:"autostart"    orm:"autostart"    description:"是否随Supervisor启动"`               // 是否随Supervisor启动
	Autorestart  int         `json:"autorestart"  orm:"autorestart"  description:"异常退出是否自动重启"`                    // 异常退出是否自动重启
	Startsecs    uint        `json:"startsecs"    orm:"startsecs"    description:"启动稳定秒数"`                        // 启动稳定秒数
	Startretries uint        `json:"startretries" orm:"startretries" description:"启动重试次数"`                        // 启动重试次数
	StopSignal   string      `json:"stopSignal"   orm:"stop_signal"  description:"停止信号"`                          // 停止信号
	Environment  string      `json:"environment"  orm:"environment"  description:"环境变量，Supervisor environment格式"` // 环境变量，Supervisor environment格式
	Remark       string      `json:"remark"       orm:"remark"       description:"备注"`                            // 备注
	CreatedBy    uint64      `json:"createdBy"    orm:"created_by"   description:"创建人ID"`                         // 创建人ID
	DeptId       uint64      `json:"deptId"       orm:"dept_id"      description:"所属部门ID"`                        // 所属部门ID
	TenantId     uint64      `json:"tenantId"     orm:"tenant_id"    description:"租户"`                            // 租户
	MerchantId   uint64      `json:"merchantId"   orm:"merchant_id"  description:"商户"`                            // 商户
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"   description:"创建时间"`                          // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"   description:"更新时间"`                          // 更新时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"   description:"软删除时间，非 NULL 表示已删除"`            // 软删除时间，非 NULL 表示已删除
}
