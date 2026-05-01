package portal

import (
	"context"
	"strconv"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/logic/portal"
	"gbaseadmin/app/member/internal/middleware"
)

// Warehouse 控制器（C 端仓库 / 寄售）。
var Warehouse = cWarehouse{}

type cWarehouse struct{}

// MyHoldings 我的仓库列表。
func (c *cWarehouse) MyHoldings(ctx context.Context, req *v1.MyWarehouseReq) (res *v1.MyWarehouseRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().ListMyWarehouse(ctx, &portal.MyWarehouseInput{
		UserID:   memberID,
		Status:   req.Status,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MyWarehouseRes{Total: out.Total, List: make([]*v1.MyWarehouseGoodsItem, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.MyWarehouseGoodsItem{
			ID:               item.ID,
			GoodsNo:          item.GoodsNo,
			Title:            item.Title,
			Cover:            item.Cover,
			InitPrice:        item.InitPrice,
			CurrentPrice:     item.CurrentPrice,
			NextListingPrice: item.NextListingPrice,
			PriceRiseRate:    item.PriceRiseRate,
			PlatformFeeRate:  item.PlatformFeeRate,
			TradeCount:       item.TradeCount,
			GoodsStatus:      item.GoodsStatus,
			GoodsStatusText:  item.GoodsStatusText,
			ActiveListingID:  item.ActiveListingID,
		})
	}
	return res, nil
}

// ListGoods 挂卖（一键挂卖，价格系统自动算）。
func (c *cWarehouse) ListGoods(ctx context.Context, req *v1.WarehouseListGoodsReq) (res *v1.WarehouseListGoodsRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	goodsID, _ := strconv.ParseInt(req.GoodsID, 10, 64)
	out, err := portal.AuthLogic().ListWarehouseGoods(ctx, &portal.ListWarehouseGoodsInput{
		UserID:  memberID,
		GoodsID: goodsID,
	})
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseListGoodsRes{
		ListingID:    out.ListingID,
		ListingPrice: portalFormatYuan(out.ListingPrice),
	}, nil
}

// Market 仓库市场。
func (c *cWarehouse) Market(ctx context.Context, req *v1.WarehouseMarketReq) (res *v1.WarehouseMarketRes, err error) {
	out, err := portal.AuthLogic().ListMarket(ctx, &portal.MarketListInput{
		Keyword:  req.Keyword,
		OrderBy:  req.OrderBy,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.WarehouseMarketRes{Total: out.Total, List: make([]*v1.WarehouseMarketListing, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.WarehouseMarketListing{
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

// PlaceTrade 买家下单。
func (c *cWarehouse) PlaceTrade(ctx context.Context, req *v1.WarehousePlaceTradeReq) (res *v1.WarehousePlaceTradeRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	listingID, _ := strconv.ParseInt(req.ListingID, 10, 64)
	out, err := portal.AuthLogic().PlaceWarehouseTrade(ctx, &portal.PlaceWarehouseTradeInput{
		BuyerID:   memberID,
		ListingID: listingID,
	})
	if err != nil {
		return nil, err
	}
	return &v1.WarehousePlaceTradeRes{
		TradeID: out.TradeID,
		TradeNo: out.TradeNo,
	}, nil
}

// ConfirmTrade 卖家确认。
func (c *cWarehouse) ConfirmTrade(ctx context.Context, req *v1.WarehouseConfirmTradeReq) (res *v1.WarehouseConfirmTradeRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	tradeID, _ := strconv.ParseInt(req.TradeID, 10, 64)
	out, err := portal.AuthLogic().ConfirmWarehouseTrade(ctx, &portal.ConfirmWarehouseTradeInput{
		SellerID: memberID,
		TradeID:  tradeID,
	})
	if err != nil {
		return nil, err
	}
	return &v1.WarehouseConfirmTradeRes{
		TradeID:      out.TradeID,
		TradePrice:   portalFormatYuan(out.TradePrice),
		PlatformFee:  portalFormatYuan(out.PlatformFee),
		SellerIncome: portalFormatYuan(out.SellerIncome),
	}, nil
}

// MyTrades 我的交易记录。
func (c *cWarehouse) MyTrades(ctx context.Context, req *v1.MyTradesReq) (res *v1.MyTradesRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().ListMyTrades(ctx, &portal.MyTradesInput{
		UserID:   memberID,
		Role:     req.Role,
		Status:   req.Status,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MyTradesRes{Total: out.Total, List: make([]*v1.TradeRecord, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.TradeRecord{
			TradeID:         item.TradeID,
			TradeNo:         item.TradeNo,
			GoodsID:         item.GoodsID,
			GoodsNo:         item.GoodsNo,
			GoodsTitle:      item.GoodsTitle,
			GoodsCover:      item.GoodsCover,
			TradePrice:      item.TradePrice,
			PlatformFee:     item.PlatformFee,
			SellerIncome:    item.SellerIncome,
			TradeStatus:     item.TradeStatus,
			TradeStatusText: item.TradeStatusText,
			Counterparty:    item.Counterparty,
			CreatedAt:       item.CreatedAt,
			ConfirmedAt:     item.ConfirmedAt,
		})
	}
	return res, nil
}
