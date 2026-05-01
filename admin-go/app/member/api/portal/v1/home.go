package v1

import "github.com/gogf/gf/v2/frame/g"

// HomeReq 首页聚合数据。
//
// 一次返回：
//   - banner（暂用配置或商城推荐，预留扩展）
//   - 等级进度（当前等级、下一级、缺多少 active_count / team_turnover）
//   - 三钱包余额简版
//   - 推荐商城商品（最多 8 个 IsRecommend=1）
//   - 仓库市场动态（最近 4 个挂卖）
type HomeReq struct {
	g.Meta `path:"/home" method:"get" tags:"会员-首页" summary:"首页聚合"`
}

// HomeRes 首页聚合返回。
type HomeRes struct {
	g.Meta            `mime:"application/json"`
	Banners           []*HomeBanner            `json:"banners"`
	LevelProgress     *HomeLevelProgress       `json:"levelProgress"`
	WalletBriefs      *HomeWalletBriefs        `json:"walletBriefs"`
	RecommendedGoods  []*MallGoodsListItem     `json:"recommendedGoods"`
	WarehouseListings []*WarehouseMarketListing `json:"warehouseListings"`
}

// HomeBanner 顶部轮播。
type HomeBanner struct {
	Image string `json:"image"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

// HomeLevelProgress 等级进度。
type HomeLevelProgress struct {
	CurrentLevelName string `json:"currentLevelName"`
	NextLevelName    string `json:"nextLevelName"`
	NeedActiveCount  int    `json:"needActiveCount" dc:"距离下一级还差多少有效用户数"`
	NeedTurnover     string `json:"needTurnover" dc:"距离下一级还差多少团队营业额（元）"`
	IsTopLevel       bool   `json:"isTopLevel" dc:"已是最高等级"`
}

// HomeWalletBriefs 三钱包余额简版。
type HomeWalletBriefs struct {
	Coupon  string `json:"coupon"`
	Reward  string `json:"reward"`
	Promote string `json:"promote"`
}
