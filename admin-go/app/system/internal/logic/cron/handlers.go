package cron

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	authlogic "gbaseadmin/app/system/internal/logic/auth"
)

// handleTenantExpireCheck 检查过期租户并自动禁用
func handleTenantExpireCheck(ctx context.Context) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	var tenants []struct {
		Id     int64  `json:"id"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	}
	err := g.DB().Ctx(ctx).Model("system_tenant").
		Fields("id", "name", "status").
		Where("status", 1).
		Where("expire_at IS NOT NULL AND expire_at < ?", now).
		Where("deleted_at", nil).
		Scan(&tenants)
	if err != nil {
		return fmt.Errorf("查询过期租户失败: %w", err)
	}

	count := 0
	for _, t := range tenants {
		_, err := g.DB().Ctx(ctx).Model("system_tenant").
			Where("id", t.Id).
			Data(g.Map{"status": 0}).
			Update()
		if err != nil {
			g.Log().Errorf(ctx, "cron: disable tenant [%s] failed: %v", t.Name, err)
			continue
		}
		authlogic.ClearTenantTokens(ctx, t.Id)
		g.Log().Infof(ctx, "cron: auto-disabled expired tenant [%s] (id=%d)", t.Name, t.Id)
		count++
	}
	return nil
}
