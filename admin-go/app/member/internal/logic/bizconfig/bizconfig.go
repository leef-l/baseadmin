// Package bizconfig 提供会员业务配置的读写与缓存。
//
// 单 row 模式：member_business_config 表只有 config_key='global' 一行，整个 payload 是 JSON。
// 商城下单 / 寄售 / 钱包返佣等核心流程都从这里读取参数（GetCachedConfig 带 5 秒缓存）。
package bizconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
	"gbaseadmin/utility/snowflake"
)

const configKeyGlobal = "global"

// PurchaseWindow 进货时间窗。
type PurchaseWindow struct {
	StartTime       string `json:"startTime"`        // "10:00"
	EndTime         string `json:"endTime"`          // "10:30"
	AllowedWeekdays []int  `json:"allowedWeekdays"`  // 1=Mon ... 7=Sun
}

// ConsignWindow 寄售时间窗（endTime 为空表示无截止）。
type ConsignWindow struct {
	StartTime string  `json:"startTime"`
	EndTime   *string `json:"endTime,omitempty"`
}

// SelfRebateTier 自购阶梯返佣档位。
type SelfRebateTier struct {
	NthOrder   int     `json:"nthOrder"`   // 当用户的 total_purchase_count 达到该值时奖励
	RewardYuan float64 `json:"rewardYuan"` // 奖励元
}

// Config 业务配置完整结构。
type Config struct {
	Purchase               PurchaseWindow   `json:"purchase"`
	Consign                ConsignWindow    `json:"consign"`
	SelfRebateTiers        []SelfRebateTier `json:"selfRebateTiers"`
	SelfTurnoverRewardRate float64          `json:"selfTurnoverRewardRate"` // 自购按金额比例返奖励钱包，单位百分比 (1.0 = 1%)
	DirectPromoteRate      float64          `json:"directPromoteRate"`      // 直推按金额比例返推广钱包，单位百分比
}

// 默认配置（迁移失败时兜底，不会进库）。
var defaultConfig = &Config{
	Purchase: PurchaseWindow{
		StartTime:       "10:00",
		EndTime:         "10:30",
		AllowedWeekdays: []int{1, 2, 3, 4, 5},
	},
	Consign: ConsignWindow{StartTime: "14:30"},
	SelfRebateTiers: []SelfRebateTier{
		{NthOrder: 2, RewardYuan: 88},
		{NthOrder: 3, RewardYuan: 188},
		{NthOrder: 4, RewardYuan: 288},
	},
	SelfTurnoverRewardRate: 1.0,
	DirectPromoteRate:      0.4,
}

// 内存缓存（5 秒），避免每个下单都查表。
var (
	cachedAt     atomic.Int64
	cachedValue  atomic.Pointer[Config]
	cacheTTLNano = int64(5 * time.Second)
)

// GetCachedConfig 返回最近 5 秒内的配置；过期会重新查 DB。
// 任意调用方修改配置后应该 InvalidateCache()。
func GetCachedConfig(ctx context.Context) (*Config, error) {
	now := time.Now().UnixNano()
	if last := cachedAt.Load(); last > 0 && now-last < cacheTTLNano {
		if v := cachedValue.Load(); v != nil {
			return v, nil
		}
	}
	cfg, err := loadFromDB(ctx)
	if err != nil {
		return nil, err
	}
	cachedValue.Store(cfg)
	cachedAt.Store(now)
	return cfg, nil
}

// InvalidateCache 显式失效缓存（后台保存配置后调用）。
func InvalidateCache() {
	cachedAt.Store(0)
}

func loadFromDB(ctx context.Context) (*Config, error) {
	var row entity.MemberBusinessConfig
	if err := dao.MemberBusinessConfig.Ctx(ctx).
		Where(dao.MemberBusinessConfig.Columns().ConfigKey, configKeyGlobal).
		Where(dao.MemberBusinessConfig.Columns().DeletedAt, nil).
		Scan(&row); err != nil {
		return nil, err
	}
	if row.Id == 0 || row.Payload == "" {
		return defaultConfig, nil
	}
	cfg := &Config{}
	if err := json.Unmarshal([]byte(row.Payload), cfg); err != nil {
		return nil, gerror.Wrap(err, "解析业务配置失败")
	}
	if cfg.Purchase.StartTime == "" || cfg.Purchase.EndTime == "" {
		return nil, gerror.New("业务配置中进货时间窗为空")
	}
	if cfg.Consign.StartTime == "" {
		return nil, gerror.New("业务配置中寄售开始时间为空")
	}
	return cfg, nil
}

// GetRaw 后台 GET 接口用，返回完整 payload。
func GetRaw(ctx context.Context) (*Config, error) {
	cfg, err := loadFromDB(ctx)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save 后台 PUT 接口用，整体覆盖 payload。
func Save(ctx context.Context, cfg *Config, operatorID int64) error {
	if cfg == nil {
		return gerror.New("配置不能为空")
	}
	if err := validate(cfg); err != nil {
		return err
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return gerror.Wrap(err, "序列化配置失败")
	}
	var existing entity.MemberBusinessConfig
	if err := dao.MemberBusinessConfig.Ctx(ctx).
		Where(dao.MemberBusinessConfig.Columns().ConfigKey, configKeyGlobal).
		Where(dao.MemberBusinessConfig.Columns().DeletedAt, nil).
		Scan(&existing); err != nil {
		return err
	}
	cols := dao.MemberBusinessConfig.Columns()
	if existing.Id == 0 {
		newID := snowflake.Generate()
		if _, err := dao.MemberBusinessConfig.Ctx(ctx).Insert(map[string]any{
			cols.Id:        int64(newID),
			cols.ConfigKey: configKeyGlobal,
			cols.Payload:   string(raw),
			cols.CreatedBy: operatorID,
		}); err != nil {
			return err
		}
	} else {
		if _, err := dao.MemberBusinessConfig.Ctx(ctx).
			Where(cols.Id, existing.Id).
			Update(map[string]any{cols.Payload: string(raw)}); err != nil {
			return err
		}
	}
	InvalidateCache()
	return nil
}

func validate(cfg *Config) error {
	if _, err := parseTimeOfDay(cfg.Purchase.StartTime); err != nil {
		return gerror.Newf("进货开始时间格式错误：%s", cfg.Purchase.StartTime)
	}
	if _, err := parseTimeOfDay(cfg.Purchase.EndTime); err != nil {
		return gerror.Newf("进货结束时间格式错误：%s", cfg.Purchase.EndTime)
	}
	if _, err := parseTimeOfDay(cfg.Consign.StartTime); err != nil {
		return gerror.Newf("寄售开始时间格式错误：%s", cfg.Consign.StartTime)
	}
	if cfg.Consign.EndTime != nil {
		if _, err := parseTimeOfDay(*cfg.Consign.EndTime); err != nil {
			return gerror.Newf("寄售结束时间格式错误：%s", *cfg.Consign.EndTime)
		}
	}
	for _, w := range cfg.Purchase.AllowedWeekdays {
		if w < 1 || w > 7 {
			return gerror.Newf("工作日定义非法：%d", w)
		}
	}
	if cfg.SelfTurnoverRewardRate < 0 || cfg.SelfTurnoverRewardRate > 100 {
		return gerror.New("自购返奖比例必须在 0-100 之间")
	}
	if cfg.DirectPromoteRate < 0 || cfg.DirectPromoteRate > 100 {
		return gerror.New("直推返奖比例必须在 0-100 之间")
	}
	return nil
}

// parseTimeOfDay 解析 "HH:MM" → time.Duration。
func parseTimeOfDay(s string) (time.Duration, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
}

// IsPurchaseAllowed 判断当前时刻是否允许下单（时间窗 + 工作日）。
func (c *Config) IsPurchaseAllowed(now time.Time) (bool, string) {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday → 7
	}
	allowed := false
	for _, w := range c.Purchase.AllowedWeekdays {
		if w == weekday {
			allowed = true
			break
		}
	}
	if !allowed {
		return false, fmt.Sprintf("今天不在允许进货的工作日内（仅 %v）", c.Purchase.AllowedWeekdays)
	}
	startD, _ := parseTimeOfDay(c.Purchase.StartTime)
	endD, _ := parseTimeOfDay(c.Purchase.EndTime)
	now0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(now0.Add(startD)) || !now.Before(now0.Add(endD)) {
		return false, fmt.Sprintf("当前不在进货时间内（%s ~ %s）", c.Purchase.StartTime, c.Purchase.EndTime)
	}
	return true, ""
}

// IsConsignAllowed 判断当前是否允许寄售操作（挂卖/购买）。
func (c *Config) IsConsignAllowed(now time.Time) (bool, string) {
	startD, _ := parseTimeOfDay(c.Consign.StartTime)
	now0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(now0.Add(startD)) {
		return false, fmt.Sprintf("当前不在寄售开放时间内（%s 起）", c.Consign.StartTime)
	}
	if c.Consign.EndTime != nil {
		endD, _ := parseTimeOfDay(*c.Consign.EndTime)
		if !now.Before(now0.Add(endD)) {
			return false, fmt.Sprintf("当前不在寄售开放时间内（%s ~ %s）", c.Consign.StartTime, *c.Consign.EndTime)
		}
	}
	return true, ""
}

// FindRebateTier 查找购买第 nth 单时的奖励金额（元），找不到返回 0。
func (c *Config) FindRebateTier(nth int) float64 {
	for _, t := range c.SelfRebateTiers {
		if t.NthOrder == nth {
			return t.RewardYuan
		}
	}
	return 0
}
