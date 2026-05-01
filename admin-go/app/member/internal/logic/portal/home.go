package portal

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
)

// HomeData 首页聚合返回。
type HomeData struct {
	Banners           []*HomeBannerItem
	LevelProgress     *HomeLevelProgressData
	WalletBriefs      *HomeWalletBriefsData
	RecommendedGoods  []*MallGoodsListData
	WarehouseListings []*MarketItem
}

// HomeBannerItem 首页轮播。
type HomeBannerItem struct {
	Image string
	Link  string
	Title string
}

// HomeLevelProgressData 等级进度。
type HomeLevelProgressData struct {
	CurrentLevelName string
	NextLevelName    string
	NeedActiveCount  int
	NeedTurnover     string // 元
	IsTopLevel       bool
}

// HomeWalletBriefsData 三钱包简版（只余额，单位元）。
type HomeWalletBriefsData struct {
	Coupon  string
	Reward  string
	Promote string
}

// GetHome 首页聚合数据。
//
// 不强求每个子模块都有数据，缺什么返回空数组/空对象，前端容错渲染。
func (s *sPortalAuth) GetHome(ctx context.Context, userID int64) (*HomeData, error) {
	if userID <= 0 {
		return nil, gerror.New("会员未登录")
	}

	out := &HomeData{
		Banners:           []*HomeBannerItem{}, // 后续可从 system 配置或独立表加载
		LevelProgress:     &HomeLevelProgressData{},
		WalletBriefs:      &HomeWalletBriefsData{},
		RecommendedGoods:  []*MallGoodsListData{},
		WarehouseListings: []*MarketItem{},
	}

	// 等级进度
	progress, err := s.computeLevelProgress(ctx, userID)
	if err == nil && progress != nil {
		out.LevelProgress = progress
	}

	// 三钱包简版
	wallets, err := s.GetMyWallets(ctx, userID)
	if err == nil && wallets != nil {
		out.WalletBriefs.Coupon = wallets.Coupon.Balance
		out.WalletBriefs.Reward = wallets.Reward.Balance
		out.WalletBriefs.Promote = wallets.Promote.Balance
	}

	// 推荐商品（前 8 个推荐位）
	rec, err := s.ListShopGoods(ctx, &MallGoodsListInput{
		IsRecommend: 1,
		PageNum:     1,
		PageSize:    8,
	})
	if err == nil && rec != nil {
		out.RecommendedGoods = rec.List
	}

	// 仓库市场动态（最近 4 条挂卖）
	market, err := s.ListMarket(ctx, &MarketListInput{
		OrderBy:  "latest",
		PageNum:  1,
		PageSize: 4,
	})
	if err == nil && market != nil {
		out.WarehouseListings = market.List
	}

	return out, nil
}

// computeLevelProgress 计算当前会员距离下一级还差多少。
func (s *sPortalAuth) computeLevelProgress(ctx context.Context, userID int64) (*HomeLevelProgressData, error) {
	var u entity.MemberUser
	if err := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().Id, userID).
		Where(dao.MemberUser.Columns().DeletedAt, nil).
		Scan(&u); err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return nil, nil
	}

	out := &HomeLevelProgressData{}

	// 当前等级名
	var currentLevelNo uint
	if u.LevelId > 0 {
		var lv entity.MemberLevel
		if err := dao.MemberLevel.Ctx(ctx).
			Where(dao.MemberLevel.Columns().Id, u.LevelId).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			Scan(&lv); err == nil {
			out.CurrentLevelName = lv.Name
			currentLevelNo = lv.LevelNo
			if lv.IsTop == 1 {
				out.IsTopLevel = true
			}
		}
	}
	if out.IsTopLevel {
		return out, nil
	}

	// 下一等级
	var nextLevel entity.MemberLevel
	if err := dao.MemberLevel.Ctx(ctx).
		Where(dao.MemberLevel.Columns().Status, 1).
		Where(dao.MemberLevel.Columns().DeletedAt, nil).
		Where(dao.MemberLevel.Columns().LevelNo+" > ?", currentLevelNo).
		OrderAsc(dao.MemberLevel.Columns().LevelNo).
		Limit(1).
		Scan(&nextLevel); err != nil {
		return out, nil
	}
	if nextLevel.Id == 0 {
		out.IsTopLevel = true
		return out, nil
	}

	out.NextLevelName = nextLevel.Name
	if int(nextLevel.NeedActiveCount) > int(u.ActiveCount) {
		out.NeedActiveCount = int(nextLevel.NeedActiveCount) - int(u.ActiveCount)
	}
	if int64(nextLevel.NeedTeamTurnover) > int64(u.TeamTurnover) {
		out.NeedTurnover = formatCent(int64(nextLevel.NeedTeamTurnover) - int64(u.TeamTurnover))
	} else {
		out.NeedTurnover = "0.00"
	}
	return out, nil
}
