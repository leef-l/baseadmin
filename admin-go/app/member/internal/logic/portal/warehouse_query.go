package portal

// 这里提供仓库 C 端查询能力（非事务，纯读）：
//   - MyWarehouseList 我的仓库列表（按 owner_id + 状态）
//   - MarketList 市场所有挂卖中
//   - MyTrades 我的交易记录（买入 / 卖出）
//
// 三个事务（挂卖 / 下单 / 确认）在 warehouse.go。

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
)

// ----- 我的仓库 -----

// MyWarehouseInput 我的仓库列表入参。
type MyWarehouseInput struct {
	UserID   int64
	Status   int // 0=全部 1=持有 2=挂卖中 3=交易中
	PageNum  int
	PageSize int
}

// MyWarehouseOutput 列表。
type MyWarehouseOutput struct {
	Total int
	List  []*MyWarehouseGoods
}

// MyWarehouseGoods 我持有的仓库商品。
type MyWarehouseGoods struct {
	ID               string
	GoodsNo          string
	Title            string
	Cover            string
	InitPrice        string
	CurrentPrice     string
	NextListingPrice string
	PriceRiseRate    int
	PlatformFeeRate  int
	TradeCount       int
	GoodsStatus      int
	GoodsStatusText  string
	ActiveListingID  string
}

// ListMyWarehouse 我的仓库列表。
func (s *sPortalAuth) ListMyWarehouse(ctx context.Context, in *MyWarehouseInput) (*MyWarehouseOutput, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	pageNum := in.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	m := dao.MemberWarehouseGoods.Ctx(ctx).
		Where(dao.MemberWarehouseGoods.Columns().OwnerId, in.UserID).
		Where(dao.MemberWarehouseGoods.Columns().DeletedAt, nil).
		Where(dao.MemberWarehouseGoods.Columns().Status, 1)
	if in.Status > 0 {
		m = m.Where(dao.MemberWarehouseGoods.Columns().GoodsStatus, in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberWarehouseGoods
	if err := m.OrderDesc(dao.MemberWarehouseGoods.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&rows); err != nil {
		return nil, err
	}

	// 收集挂卖中商品对应的 active listing ID
	listedIDs := make([]uint64, 0)
	for _, row := range rows {
		if row.GoodsStatus == GoodsStatusListed {
			listedIDs = append(listedIDs, row.Id)
		}
	}
	listingMap := make(map[uint64]uint64, len(listedIDs))
	if len(listedIDs) > 0 {
		var listings []entity.MemberWarehouseListing
		if err := dao.MemberWarehouseListing.Ctx(ctx).
			WhereIn(dao.MemberWarehouseListing.Columns().GoodsId, listedIDs).
			Where(dao.MemberWarehouseListing.Columns().ListingStatus, ListingStatusActive).
			Where(dao.MemberWarehouseListing.Columns().DeletedAt, nil).
			Scan(&listings); err == nil {
			for _, l := range listings {
				listingMap[l.GoodsId] = l.Id
			}
		}
	}

	out := &MyWarehouseOutput{Total: total, List: make([]*MyWarehouseGoods, 0, len(rows))}
	for _, row := range rows {
		nextPrice := computeListingPrice(int64(row.CurrentPrice), int64(row.PriceRiseRate))
		item := &MyWarehouseGoods{
			ID:               fmt.Sprintf("%d", row.Id),
			GoodsNo:          row.GoodsNo,
			Title:            row.Title,
			Cover:            row.Cover,
			InitPrice:        formatCent(int64(row.InitPrice)),
			CurrentPrice:     formatCent(int64(row.CurrentPrice)),
			NextListingPrice: formatCent(nextPrice),
			PriceRiseRate:    int(row.PriceRiseRate),
			PlatformFeeRate:  int(row.PlatformFeeRate),
			TradeCount:       int(row.TradeCount),
			GoodsStatus:      row.GoodsStatus,
			GoodsStatusText:  warehouseGoodsStatusText(row.GoodsStatus),
		}
		if id, ok := listingMap[row.Id]; ok {
			item.ActiveListingID = fmt.Sprintf("%d", id)
		}
		out.List = append(out.List, item)
	}
	return out, nil
}

// ----- 市场列表 -----

// MarketListInput 市场列表入参。
type MarketListInput struct {
	Keyword  string
	OrderBy  string // price_asc / price_desc / latest
	PageNum  int
	PageSize int
}

// MarketListOutput 市场列表。
type MarketListOutput struct {
	Total int
	List  []*MarketItem
}

// MarketItem 市场单项。
type MarketItem struct {
	ListingID    string
	GoodsID      string
	GoodsNo      string
	Title        string
	Cover        string
	ListingPrice string
	SellerName   string
	TradeCount   int
	ListedAt     string
}

// ListMarket 仓库市场分页（所有 listing_status=1 的挂卖）。
func (s *sPortalAuth) ListMarket(ctx context.Context, in *MarketListInput) (*MarketListOutput, error) {
	if in == nil {
		in = &MarketListInput{}
	}
	pageNum := in.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	listingT := dao.MemberWarehouseListing.Table()
	goodsT := dao.MemberWarehouseGoods.Table()
	m := dao.MemberWarehouseListing.Ctx(ctx).As("l").
		LeftJoin(goodsT+" g", "g.id = l.goods_id").
		Where("l.listing_status", ListingStatusActive).
		Where("l.deleted_at", nil).
		Where("g.deleted_at", nil).
		Where("g.status", 1)
	if v := strings.TrimSpace(in.Keyword); v != "" {
		m = m.WhereLike("g.title", "%"+v+"%")
	}
	total, err := m.Count("l.id")
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(strings.TrimSpace(in.OrderBy)) {
	case "price_asc":
		m = m.OrderAsc("l.listing_price")
	case "price_desc":
		m = m.OrderDesc("l.listing_price")
	default:
		m = m.OrderDesc("l.listed_at").OrderDesc("l.id")
	}

	type row struct {
		ListingId     uint64 `json:"listing_id"`
		GoodsId       uint64 `json:"goods_id"`
		GoodsNo       string `json:"goods_no"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		ListingPrice  uint64 `json:"listing_price"`
		SellerId      uint64 `json:"seller_id"`
		TradeCount    uint   `json:"trade_count"`
		ListedAt      string `json:"listed_at"`
	}
	var rows []row
	if err := m.Fields(
		"l.id AS listing_id",
		"l.goods_id",
		"g.goods_no",
		"g.title",
		"g.cover",
		"l.listing_price",
		"l.seller_id",
		"g.trade_count",
		"l.listed_at",
	).Page(pageNum, pageSize).Scan(&rows); err != nil {
		// LeftJoin + Scan 失败时尝试拆分查询（fallback）
		_ = listingT
		return nil, err
	}

	// 批量查卖家昵称
	sellerIDs := make([]uint64, 0, len(rows))
	for _, r := range rows {
		if r.SellerId > 0 {
			sellerIDs = append(sellerIDs, r.SellerId)
		}
	}
	sellerNameMap := make(map[uint64]string, len(sellerIDs))
	if len(sellerIDs) > 0 {
		var users []entity.MemberUser
		if err := dao.MemberUser.Ctx(ctx).
			Fields(dao.MemberUser.Columns().Id, dao.MemberUser.Columns().Nickname).
			WhereIn(dao.MemberUser.Columns().Id, sellerIDs).
			Where(dao.MemberUser.Columns().DeletedAt, nil).
			Scan(&users); err == nil {
			for _, u := range users {
				sellerNameMap[u.Id] = u.Nickname
			}
		}
	}

	out := &MarketListOutput{Total: total, List: make([]*MarketItem, 0, len(rows))}
	for _, r := range rows {
		out.List = append(out.List, &MarketItem{
			ListingID:    fmt.Sprintf("%d", r.ListingId),
			GoodsID:      fmt.Sprintf("%d", r.GoodsId),
			GoodsNo:      r.GoodsNo,
			Title:        r.Title,
			Cover:        r.Cover,
			ListingPrice: formatCent(int64(r.ListingPrice)),
			SellerName:   sellerNameMap[r.SellerId],
			TradeCount:   int(r.TradeCount),
			ListedAt:     r.ListedAt,
		})
	}
	return out, nil
}

// ----- 我的交易记录 -----

// MyTradesInput 我的交易（买入或卖出）。
type MyTradesInput struct {
	UserID   int64
	Role     string // buyer / seller
	Status   int    // 0=全部 1=待确认 2=已完成 3=已取消
	PageNum  int
	PageSize int
}

// MyTradesOutput 列表。
type MyTradesOutput struct {
	Total int
	List  []*TradeItem
}

// TradeItem 交易记录。
type TradeItem struct {
	TradeID         string
	TradeNo         string
	GoodsID         string
	GoodsNo         string
	GoodsTitle      string
	GoodsCover      string
	TradePrice      string
	PlatformFee     string
	SellerIncome    string
	TradeStatus     int
	TradeStatusText string
	Counterparty    string
	CreatedAt       string
	ConfirmedAt     string
}

// ListMyTrades 列出我的交易记录。
func (s *sPortalAuth) ListMyTrades(ctx context.Context, in *MyTradesInput) (*MyTradesOutput, error) {
	if in == nil || in.UserID <= 0 {
		return nil, gerror.New("会员未登录")
	}
	role := strings.TrimSpace(in.Role)
	if role != "seller" {
		role = "buyer"
	}
	pageNum := in.PageNum
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	m := dao.MemberWarehouseTrade.Ctx(ctx).
		Where(dao.MemberWarehouseTrade.Columns().DeletedAt, nil)
	if role == "buyer" {
		m = m.Where(dao.MemberWarehouseTrade.Columns().BuyerId, in.UserID)
	} else {
		m = m.Where(dao.MemberWarehouseTrade.Columns().SellerId, in.UserID)
	}
	if in.Status > 0 {
		m = m.Where(dao.MemberWarehouseTrade.Columns().TradeStatus, in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberWarehouseTrade
	if err := m.OrderDesc(dao.MemberWarehouseTrade.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&rows); err != nil {
		return nil, err
	}

	// 批量查商品 + 对手昵称
	goodsIDs := make([]uint64, 0, len(rows))
	userIDs := make([]uint64, 0, len(rows))
	for _, r := range rows {
		goodsIDs = append(goodsIDs, r.GoodsId)
		if role == "buyer" {
			userIDs = append(userIDs, r.SellerId)
		} else {
			userIDs = append(userIDs, r.BuyerId)
		}
	}

	goodsMap := make(map[uint64]entity.MemberWarehouseGoods)
	if len(goodsIDs) > 0 {
		var gs []entity.MemberWarehouseGoods
		if err := dao.MemberWarehouseGoods.Ctx(ctx).
			WhereIn(dao.MemberWarehouseGoods.Columns().Id, goodsIDs).
			Scan(&gs); err == nil {
			for _, g := range gs {
				goodsMap[g.Id] = g
			}
		}
	}
	userMap := make(map[uint64]string)
	if len(userIDs) > 0 {
		var users []entity.MemberUser
		if err := dao.MemberUser.Ctx(ctx).
			Fields(dao.MemberUser.Columns().Id, dao.MemberUser.Columns().Nickname).
			WhereIn(dao.MemberUser.Columns().Id, userIDs).
			Scan(&users); err == nil {
			for _, u := range users {
				userMap[u.Id] = u.Nickname
			}
		}
	}

	out := &MyTradesOutput{Total: total, List: make([]*TradeItem, 0, len(rows))}
	for _, r := range rows {
		gd := goodsMap[r.GoodsId]
		var counterparty string
		if role == "buyer" {
			counterparty = userMap[r.SellerId]
		} else {
			counterparty = userMap[r.BuyerId]
		}
		out.List = append(out.List, &TradeItem{
			TradeID:         fmt.Sprintf("%d", r.Id),
			TradeNo:         r.TradeNo,
			GoodsID:         fmt.Sprintf("%d", r.GoodsId),
			GoodsNo:         gd.GoodsNo,
			GoodsTitle:      gd.Title,
			GoodsCover:      gd.Cover,
			TradePrice:      formatCent(int64(r.TradePrice)),
			PlatformFee:     formatCent(int64(r.PlatformFee)),
			SellerIncome:    formatCent(int64(r.SellerIncome)),
			TradeStatus:     r.TradeStatus,
			TradeStatusText: tradeStatusText(r.TradeStatus),
			Counterparty:    counterparty,
			CreatedAt:       timeStr(r.CreatedAt),
			ConfirmedAt:     timeStr(r.ConfirmedAt),
		})
	}
	return out, nil
}

// warehouseGoodsStatusText 商品状态文案。
func warehouseGoodsStatusText(status int) string {
	switch status {
	case GoodsStatusHolding:
		return "持有中"
	case GoodsStatusListed:
		return "挂卖中"
	case GoodsStatusTrading:
		return "交易中"
	}
	return "未知"
}

// tradeStatusText 交易状态文案。
func tradeStatusText(status int) string {
	switch status {
	case TradeStatusPending:
		return "待卖家确认"
	case TradeStatusConfirmed:
		return "已确认完成"
	case TradeStatusCanceled:
		return "已取消"
	}
	return "未知"
}
