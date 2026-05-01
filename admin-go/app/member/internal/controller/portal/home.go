package portal

import (
	"context"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/logic/portal"
	"gbaseadmin/app/member/internal/middleware"
)

// Home 控制器（首页聚合）。
var Home = cHome{}

type cHome struct{}

// Index 首页聚合。
func (c *cHome) Index(ctx context.Context, req *v1.HomeReq) (res *v1.HomeRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().GetHome(ctx, memberID)
	if err != nil {
		return nil, err
	}

	res = &v1.HomeRes{
		Banners:           make([]*v1.HomeBanner, 0, len(out.Banners)),
		LevelProgress:     &v1.HomeLevelProgress{},
		WalletBriefs:      &v1.HomeWalletBriefs{},
		RecommendedGoods:  convertGoodsItems(out.RecommendedGoods),
		WarehouseListings: make([]*v1.WarehouseMarketListing, 0, len(out.WarehouseListings)),
	}
	for _, b := range out.Banners {
		res.Banners = append(res.Banners, &v1.HomeBanner{Image: b.Image, Link: b.Link, Title: b.Title})
	}
	if out.LevelProgress != nil {
		res.LevelProgress = &v1.HomeLevelProgress{
			CurrentLevelName: out.LevelProgress.CurrentLevelName,
			NextLevelName:    out.LevelProgress.NextLevelName,
			NeedActiveCount:  out.LevelProgress.NeedActiveCount,
			NeedTurnover:     out.LevelProgress.NeedTurnover,
			IsTopLevel:       out.LevelProgress.IsTopLevel,
		}
	}
	if out.WalletBriefs != nil {
		res.WalletBriefs = &v1.HomeWalletBriefs{
			Coupon:  out.WalletBriefs.Coupon,
			Reward:  out.WalletBriefs.Reward,
			Promote: out.WalletBriefs.Promote,
		}
	}
	for _, item := range out.WarehouseListings {
		res.WarehouseListings = append(res.WarehouseListings, &v1.WarehouseMarketListing{
			ListingID:    item.ListingID,
			GoodsID:      item.GoodsID,
			GoodsNo:      item.GoodsNo,
			Title:        item.Title,
			Cover:        item.Cover,
			ListingPrice: item.ListingPrice,
			SellerName:   item.SellerName,
			TradeCount:   item.TradeCount,
			ListedAt:     item.ListedAt,
		})
	}
	return res, nil
}
