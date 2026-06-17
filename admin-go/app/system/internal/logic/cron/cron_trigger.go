package cron

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/model"
)

// Trigger 手动执行定时任务
func (s *sCron) Trigger(ctx context.Context, name string, _ string) (bool, string) {
	// 先通过 name 查 handler
	var row struct {
		Name    string `json:"name"`
		Handler string `json:"handler"`
	}
	err := g.DB().Ctx(ctx).Model("system_cron").
		Fields("name", "handler").
		Where("name", name).
		Where("deleted_at", nil).
		Scan(&row)
	if err != nil || row.Handler == "" {
		return false, "未找到任务: " + name
	}
	return GetManager().Trigger(ctx, row.Handler)
}

// LogList 查询执行日志
func (s *sCron) LogList(ctx context.Context, cronID int64, pageNum, pageSize int) ([]*model.CronLogItem, int, error) {
	if cronID <= 0 {
		return nil, 0, nil
	}
	m := g.DB().Ctx(ctx).Model("system_cron_log").Where("cron_id", cronID)
	total, _ := m.Count()
	if total == 0 {
		return nil, 0, nil
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var list []*model.CronLogItem
	err := m.OrderDesc("id").Page(pageNum, pageSize).Scan(&list)
	return list, total, err
}
