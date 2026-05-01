package portal

import (
	"context"
	"time"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/logic/bizconfig"
)

type cBizConfig struct{}

// BizConfig 公开的 C 端业务配置接口。
var BizConfig = cBizConfig{}

// Get 获取业务配置（仅时间窗与服务器时间戳，用于前端倒计时）。
func (c cBizConfig) Get(ctx context.Context, req *v1.PortalBizConfigReq) (res *v1.PortalBizConfigRes, err error) {
	cfg, err := bizconfig.GetCachedConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PortalBizConfigRes{
		PurchaseStart:   cfg.Purchase.StartTime,
		PurchaseEnd:     cfg.Purchase.EndTime,
		PurchaseDays:    cfg.Purchase.AllowedWeekdays,
		ConsignStart:    cfg.Consign.StartTime,
		ConsignEnd:      cfg.Consign.EndTime,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
