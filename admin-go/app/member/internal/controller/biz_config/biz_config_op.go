package biz_config

import (
	"context"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/logic/bizconfig"
	"gbaseadmin/app/member/internal/middleware"
)

// Get 后台获取业务配置（单例）。
func (c cBizConfig) Get(ctx context.Context, req *v1.MemberBizConfigGetReq) (res *v1.MemberBizConfigGetRes, err error) {
	cfg, err := bizconfig.GetRaw(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MemberBizConfigGetRes{BizConfigData: toAPI(cfg)}, nil
}

// Save 后台保存业务配置。
func (c cBizConfig) Save(ctx context.Context, req *v1.MemberBizConfigSaveReq) (res *v1.MemberBizConfigSaveRes, err error) {
	if req.BizConfigData == nil {
		req.BizConfigData = &v1.BizConfigData{}
	}
	operatorID := int64(middleware.GetUserID(ctx))
	if err := bizconfig.Save(ctx, fromAPI(req.BizConfigData), operatorID); err != nil {
		return nil, err
	}
	return &v1.MemberBizConfigSaveRes{}, nil
}

func toAPI(cfg *bizconfig.Config) *v1.BizConfigData {
	tiers := make([]v1.BizConfigRebate, 0, len(cfg.SelfRebateTiers))
	for _, t := range cfg.SelfRebateTiers {
		tiers = append(tiers, v1.BizConfigRebate{NthOrder: t.NthOrder, RewardYuan: t.RewardYuan})
	}
	return &v1.BizConfigData{
		Purchase: v1.BizConfigPurchase{
			StartTime:       cfg.Purchase.StartTime,
			EndTime:         cfg.Purchase.EndTime,
			AllowedWeekdays: cfg.Purchase.AllowedWeekdays,
		},
		Consign: v1.BizConfigConsign{
			StartTime: cfg.Consign.StartTime,
			EndTime:   cfg.Consign.EndTime,
		},
		SelfRebateTiers:        tiers,
		SelfTurnoverRewardRate: cfg.SelfTurnoverRewardRate,
		DirectPromoteRate:      cfg.DirectPromoteRate,
	}
}

func fromAPI(in *v1.BizConfigData) *bizconfig.Config {
	tiers := make([]bizconfig.SelfRebateTier, 0, len(in.SelfRebateTiers))
	for _, t := range in.SelfRebateTiers {
		tiers = append(tiers, bizconfig.SelfRebateTier{NthOrder: t.NthOrder, RewardYuan: t.RewardYuan})
	}
	return &bizconfig.Config{
		Purchase: bizconfig.PurchaseWindow{
			StartTime:       in.Purchase.StartTime,
			EndTime:         in.Purchase.EndTime,
			AllowedWeekdays: in.Purchase.AllowedWeekdays,
		},
		Consign: bizconfig.ConsignWindow{
			StartTime: in.Consign.StartTime,
			EndTime:   in.Consign.EndTime,
		},
		SelfRebateTiers:        tiers,
		SelfTurnoverRewardRate: in.SelfTurnoverRewardRate,
		DirectPromoteRate:      in.DirectPromoteRate,
	}
}
