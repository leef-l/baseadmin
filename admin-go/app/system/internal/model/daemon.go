package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// DaemonCreateInput 创建守护进程输入
type DaemonCreateInput struct {
	Name         string `json:"name"`
	Program      string `json:"program"`
	Command      string `json:"command"`
	Directory    string `json:"directory"`
	RunUser      string `json:"runUser"`
	Numprocs     int    `json:"numprocs"`
	Priority     int    `json:"priority"`
	Autostart    int    `json:"autostart"`
	Autorestart  int    `json:"autorestart"`
	Startsecs    int    `json:"startsecs"`
	Startretries int    `json:"startretries"`
	StopSignal   string `json:"stopSignal"`
	Environment  string `json:"environment"`
	Remark       string `json:"remark"`
}

// DaemonUpdateInput 更新守护进程输入
type DaemonUpdateInput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Program      string              `json:"program"`
	Command      string              `json:"command"`
	Directory    string              `json:"directory"`
	RunUser      string              `json:"runUser"`
	Numprocs     int                 `json:"numprocs"`
	Priority     int                 `json:"priority"`
	Autostart    int                 `json:"autostart"`
	Autorestart  int                 `json:"autorestart"`
	Startsecs    int                 `json:"startsecs"`
	Startretries int                 `json:"startretries"`
	StopSignal   string              `json:"stopSignal"`
	Environment  string              `json:"environment"`
	Remark       string              `json:"remark"`
}

// DaemonDetailOutput 守护进程详情输出
type DaemonDetailOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Program      string              `json:"program"`
	Command      string              `json:"command"`
	Directory    string              `json:"directory"`
	RunUser      string              `json:"runUser"`
	Numprocs     int                 `json:"numprocs"`
	Priority     int                 `json:"priority"`
	Autostart    int                 `json:"autostart"`
	Autorestart  int                 `json:"autorestart"`
	Startsecs    int                 `json:"startsecs"`
	Startretries int                 `json:"startretries"`
	StopSignal   string              `json:"stopSignal"`
	Environment  string              `json:"environment"`
	Remark       string              `json:"remark"`
	ConfigPath   string              `json:"configPath"`
	OutLogPath   string              `json:"outLogPath"`
	ErrLogPath   string              `json:"errLogPath"`
	RunStatus    string              `json:"runStatus"`
	Pid          string              `json:"pid"`
	Uptime       string              `json:"uptime"`
	StatusText   string              `json:"statusText"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// DaemonListOutput 守护进程列表输出
type DaemonListOutput = DaemonDetailOutput

// DaemonListInput 守护进程列表查询输入
type DaemonListInput struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
	Program  string `json:"program"`
}

// DaemonOperationOutput 守护进程操作输出
type DaemonOperationOutput struct {
	Program   string `json:"program"`
	RunStatus string `json:"runStatus"`
	Message   string `json:"message"`
}

// DaemonBatchOperationOutput 守护进程批量操作输出
type DaemonBatchOperationOutput struct {
	Results []*DaemonOperationOutput `json:"results"`
}

// DaemonLogOutput 守护进程日志输出
type DaemonLogOutput struct {
	Program string `json:"program"`
	LogType string `json:"logType"`
	Content string `json:"content"`
}
