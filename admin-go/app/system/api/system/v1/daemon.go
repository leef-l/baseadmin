package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// DaemonCreateReq 创建守护进程请求
type DaemonCreateReq struct {
	g.Meta       `path:"/daemon/create" method:"post" tags:"守护进程" summary:"创建守护进程"`
	Name         string `json:"name" v:"required#显示名称不能为空" dc:"显示名称"`
	Program      string `json:"program" v:"required#进程名不能为空" dc:"Supervisor进程名"`
	Command      string `json:"command" v:"required#启动命令不能为空" dc:"启动命令"`
	Directory    string `json:"directory" v:"required#运行目录不能为空" dc:"运行目录"`
	RunUser      string `json:"runUser" dc:"运行用户"`
	Numprocs     int    `json:"numprocs" dc:"进程数量"`
	Priority     int    `json:"priority" dc:"启动优先级"`
	Autostart    int    `json:"autostart" dc:"是否随Supervisor启动"`
	Autorestart  int    `json:"autorestart" dc:"异常退出是否自动重启"`
	Startsecs    int    `json:"startsecs" dc:"启动稳定秒数"`
	Startretries int    `json:"startretries" dc:"启动重试次数"`
	StopSignal   string `json:"stopSignal" dc:"停止信号"`
	Environment  string `json:"environment" dc:"环境变量"`
	Remark       string `json:"remark" dc:"备注"`
}

type DaemonCreateRes struct {
	g.Meta `mime:"application/json"`
}

// DaemonUpdateReq 更新守护进程请求
type DaemonUpdateReq struct {
	g.Meta       `path:"/daemon/update" method:"put" tags:"守护进程" summary:"更新守护进程"`
	ID           snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
	Name         string              `json:"name" v:"required#显示名称不能为空" dc:"显示名称"`
	Program      string              `json:"program" dc:"Supervisor进程名"`
	Command      string              `json:"command" v:"required#启动命令不能为空" dc:"启动命令"`
	Directory    string              `json:"directory" v:"required#运行目录不能为空" dc:"运行目录"`
	RunUser      string              `json:"runUser" dc:"运行用户"`
	Numprocs     int                 `json:"numprocs" dc:"进程数量"`
	Priority     int                 `json:"priority" dc:"启动优先级"`
	Autostart    int                 `json:"autostart" dc:"是否随Supervisor启动"`
	Autorestart  int                 `json:"autorestart" dc:"异常退出是否自动重启"`
	Startsecs    int                 `json:"startsecs" dc:"启动稳定秒数"`
	Startretries int                 `json:"startretries" dc:"启动重试次数"`
	StopSignal   string              `json:"stopSignal" dc:"停止信号"`
	Environment  string              `json:"environment" dc:"环境变量"`
	Remark       string              `json:"remark" dc:"备注"`
}

type DaemonUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// DaemonDeleteReq 删除守护进程请求
type DaemonDeleteReq struct {
	g.Meta `path:"/daemon/delete" method:"delete" tags:"守护进程" summary:"删除守护进程"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
}

type DaemonDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// DaemonBatchDeleteReq 批量删除守护进程请求
type DaemonBatchDeleteReq struct {
	g.Meta `path:"/daemon/batch-delete" method:"delete" tags:"守护进程" summary:"批量删除守护进程"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"守护进程ID列表"`
}

type DaemonBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonBatchOperationOutput
}

// DaemonDetailReq 获取守护进程详情请求
type DaemonDetailReq struct {
	g.Meta `path:"/daemon/detail" method:"get" tags:"守护进程" summary:"获取守护进程详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
}

type DaemonDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonDetailOutput
}

// DaemonListReq 获取守护进程列表请求
type DaemonListReq struct {
	g.Meta   `path:"/daemon/list" method:"get" tags:"守护进程" summary:"获取守护进程列表"`
	PageNum  int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"10" dc:"每页数量"`
	Keyword  string `json:"keyword" dc:"关键词"`
	Program  string `json:"program" dc:"进程名"`
}

type DaemonListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.DaemonListOutput `json:"list" dc:"列表数据"`
	Total  int                       `json:"total" dc:"总数"`
}

// DaemonRestartReq 重启守护进程请求
type DaemonRestartReq struct {
	g.Meta `path:"/daemon/restart" method:"post" tags:"守护进程" summary:"重启守护进程"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
}

type DaemonRestartRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonOperationOutput
}

// DaemonBatchRestartReq 批量重启守护进程请求
type DaemonBatchRestartReq struct {
	g.Meta `path:"/daemon/batch-restart" method:"post" tags:"守护进程" summary:"批量重启守护进程"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"守护进程ID列表"`
}

type DaemonBatchRestartRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonBatchOperationOutput
}

// DaemonStopReq 暂停守护进程请求
type DaemonStopReq struct {
	g.Meta `path:"/daemon/stop" method:"post" tags:"守护进程" summary:"暂停守护进程"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
}

type DaemonStopRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonOperationOutput
}

// DaemonBatchStopReq 批量暂停守护进程请求
type DaemonBatchStopReq struct {
	g.Meta `path:"/daemon/batch-stop" method:"post" tags:"守护进程" summary:"批量暂停守护进程"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"守护进程ID列表"`
}

type DaemonBatchStopRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonBatchOperationOutput
}

// DaemonLogReq 获取守护进程日志请求
type DaemonLogReq struct {
	g.Meta  `path:"/daemon/log" method:"get" tags:"守护进程" summary:"获取守护进程日志"`
	ID      snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"守护进程ID"`
	LogType string              `json:"logType" dc:"日志类型:normal/error"`
	Lines   int                 `json:"lines" dc:"行数"`
}

type DaemonLogRes struct {
	g.Meta `mime:"application/json"`
	*model.DaemonLogOutput
}
