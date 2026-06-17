package cron

import (
	"context"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/service"
)

// Trigger 手动执行定时任务
func (c *cCron) Trigger(ctx context.Context, req *v1.CronTriggerReq) (res *v1.CronTriggerRes, err error) {
	success, msg := service.Cron().Trigger(ctx, req.Name, "")
	res = &v1.CronTriggerRes{Success: success, Message: msg}
	return
}

// LogList 查询执行日志
func (c *cCron) LogList(ctx context.Context, req *v1.CronLogListReq) (res *v1.CronLogListRes, err error) {
	res = &v1.CronLogListRes{}
	res.List, res.Total, err = service.Cron().LogList(ctx, req.CronID, req.PageNum, req.PageSize)
	return
}
