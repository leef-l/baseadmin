package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

var (
	managerOnce sync.Once
	managerInst *Manager
	watcherMu   sync.Mutex
)

// Manager 定时任务管理器（基于 gcron + DB 热更新）
type Manager struct {
	cron     *gcron.Cron
	handlers map[string]func(ctx context.Context) error
}

// GetManager 获取单例管理器
func GetManager() *Manager {
	managerOnce.Do(func() {
		managerInst = &Manager{
			cron: gcron.New(),
			handlers: map[string]func(ctx context.Context) error{
				"tenant_expire_check": handleTenantExpireCheck,
			},
		}
	})
	return managerInst
}

// Start 启动定时任务管理器和热更新监听
func (m *Manager) Start(ctx context.Context) {
	m.reload(ctx)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				m.cron.Stop()
				return
			case <-ticker.C:
				m.reload(ctx)
			}
		}
	}()
}

func (m *Manager) reload(ctx context.Context) {
	watcherMu.Lock()
	defer watcherMu.Unlock()

	rows, err := loadCronRows(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "cron reload: load failed: %v", err)
		return
	}

	desired := make(map[string]cronRow)
	for _, row := range rows {
		if row.Status == 0 {
			continue
		}
		if _, ok := m.handlers[row.Handler]; !ok {
			continue
		}
		desired[row.Name] = row
	}

	for _, entry := range m.cron.Entries() {
		name := entry.Name
		if _, ok := desired[name]; !ok {
			m.cron.Remove(name)
			g.Log().Infof(ctx, "cron: removed task [%s]", name)
		}
	}

	for name, row := range desired {
		existing := m.cron.Search(name)
		if existing != nil {
			continue
		}
		handler := m.handlers[row.Handler]
		m.addJob(ctx, name, row.Expression, row.ID, handler)
	}
}

func (m *Manager) addJob(ctx context.Context, name, expression string, cronID int64, handler func(ctx context.Context) error) {
	_, err := m.cron.Add(ctx, expression, func(ctx context.Context) {
		logID := recordCronLogStart(ctx, cronID, name)
		startTime := time.Now()
		err := handler(ctx)
		duration := time.Since(startTime).Milliseconds()
		if err != nil {
			recordCronLogEnd(ctx, logID, "fail", err.Error(), int(duration))
			g.Log().Errorf(ctx, "cron: task [%s] failed: %v", name, err)
		} else {
			recordCronLogEnd(ctx, logID, "success", "", int(duration))
		}
		updateCronRunStatus(ctx, cronID, err == nil)
	}, name)
		if err != nil {
			g.Log().Errorf(ctx, "cron: add task [%s] failed: %v (expression=%s, expected 6-field format: second minute hour day month week)", name, err, expression)
		} else {
		g.Log().Infof(ctx, "cron: added task [%s] expr=%s", name, expression)
	}
}

// Trigger 手动触发一个任务，返回执行结果
func (m *Manager) Trigger(ctx context.Context, name string) (success bool, message string) {
	handler, ok := m.handlers[name]
	if !ok {
		return false, fmt.Sprintf("未知的处理器: %s", name)
	}

	// 查任务ID和显示名称
	var row struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
		_ = g.DB().Ctx(ctx).Model("system_cron").Fields("id", "name").Where("handler", name).Where("deleted_at", nil).Scan(&row)

	// 优先使用数据库中的显示名称，回退到 handler 名
	logName := row.Name
	if logName == "" {
		logName = name
	}

	logID := recordCronLogStart(ctx, row.ID, logName)
	startTime := time.Now()
	err := handler(ctx)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		recordCronLogEnd(ctx, logID, "fail", err.Error(), int(duration))
		updateCronRunStatus(ctx, row.ID, false)
		return false, err.Error()
	}
	recordCronLogEnd(ctx, logID, "success", "", int(duration))
	updateCronRunStatus(ctx, row.ID, true)
	return true, "执行成功"
}

func updateCronRunStatus(ctx context.Context, id int64, success bool) {
	result := "success"
	if !success {
		result = "fail"
	}
	_, _ = g.DB().Ctx(ctx).Model("system_cron").
		Where("id", id).
		Data(g.Map{
			"last_run_at": gtime.Now(),
			"last_result": result,
		}).
		Update()
}

func recordCronLogStart(ctx context.Context, cronID int64, cronName string) int64 {
	id := snowflake.Generate()
	_, _ = g.DB().Ctx(ctx).Model("system_cron_log").Insert(g.Map{
		"id":        id,
		"cron_id":   cronID,
		"cron_name": cronName,
		"start_at":  gtime.Now(),
		"result":    "pending",
	})
	return int64(id)
}

func recordCronLogEnd(ctx context.Context, logID int64, result, message string, durationMs int) {
	_, _ = g.DB().Ctx(ctx).Model("system_cron_log").
		Where("id", logID).
		Update(g.Map{
			"end_at":      gtime.Now(),
			"duration_ms": durationMs,
			"result":      result,
			"message":     message,
		})
}

type cronRow struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Handler    string `json:"handler"`
	Status     int    `json:"status"`
}

func loadCronRows(ctx context.Context) ([]cronRow, error) {
	var rows []cronRow
	err := g.DB().Ctx(ctx).Model("system_cron").
		Fields("id", "name", "expression", "handler", "status").
		Where("deleted_at", nil).
		Scan(&rows)
	return rows, err
}
