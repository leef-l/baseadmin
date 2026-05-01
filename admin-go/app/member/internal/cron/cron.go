// Package cron 注册 member 应用的定时任务。
//
// 入口在 cmd/cmd.go 启动 server 之前调用 Setup()，gcron 在后台 goroutine 内执行。
//
// 当前注册的任务：
//   - 每天 01:30 扫描会员等级过期，把 is_qualified 置 0
//   - 每小时一次轻量补偿：补充再次校验过期（防止跨日未扫到）
//   - 每天 00:05 把 today_purchase_count 全部清零（兜底；下单时也会按日期懒重置）
package cron

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/logic/teamops"
)

// Setup 注册所有定时任务。多次调用幂等。
func Setup(ctx context.Context) error {
	// 每天 01:30 跑一次过期扫描
	if _, err := gcron.AddSingleton(ctx, "0 30 1 * * *", func(ctx context.Context) {
		count, err := teamops.ScanExpiredLevels(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "[cron] ScanExpiredLevels err=%v", err)
			return
		}
		g.Log().Infof(ctx, "[cron] ScanExpiredLevels processed=%d", count)
	}, "member.scan_expired_levels"); err != nil {
		return err
	}

	// 每小时一次轻量补扫（防止 01:30 任务漏跑或重启错过）
	if _, err := gcron.AddSingleton(ctx, "0 0 * * * *", func(ctx context.Context) {
		count, err := teamops.ScanExpiredLevels(ctx)
		if err != nil {
			g.Log().Warningf(ctx, "[cron] hourly compensate ScanExpiredLevels err=%v", err)
			return
		}
		if count > 0 {
			g.Log().Infof(ctx, "[cron] hourly compensate processed=%d", count)
		}
	}, "member.scan_expired_levels.hourly"); err != nil {
		return err
	}

	// 每天 00:05 重置 today_purchase_count；下单时也会按日期懒重置，此处是兜底。
	if _, err := gcron.AddSingleton(ctx, "0 5 0 * * *", func(ctx context.Context) {
		cols := dao.MemberUser.Columns()
		r, err := dao.MemberUser.Ctx(ctx).
			Where(cols.DeletedAt, nil).
			Where(cols.TodayPurchaseCount+" > 0").
			Data(g.Map{cols.TodayPurchaseCount: 0}).
			Update()
		if err != nil {
			g.Log().Errorf(ctx, "[cron] reset today_purchase_count err=%v", err)
			return
		}
		affected, _ := r.RowsAffected()
		g.Log().Infof(ctx, "[cron] reset today_purchase_count affected=%d", affected)
	}, "member.reset_today_purchase_count"); err != nil {
		return err
	}

	return nil
}
