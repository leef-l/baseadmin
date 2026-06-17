package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/model"
)

// CronTriggerReq 手动触发定时任务
type CronTriggerReq struct {
	g.Meta `path:"/cron/trigger" method:"post" tags:"定时任务表" summary:"手动执行定时任务"`
	Name   string `json:"name" v:"required#任务名称不能为空" dc:"任务名称"`
}

type CronTriggerRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CronLogListReq 查询执行日志
type CronLogListReq struct {
	g.Meta   `path:"/cron/log-list" method:"get" tags:"定时任务表" summary:"查询执行日志"`
	CronID   int64 `json:"cronId" v:"required#任务ID不能为空" dc:"任务ID"`
	PageNum  int   `json:"pageNum" d:"1" dc:"页码"`
	PageSize int   `json:"pageSize" d:"20" dc:"每页数量"`
}

type CronLogListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CronLogItem `json:"list"`
	Total  int                  `json:"total"`
}
