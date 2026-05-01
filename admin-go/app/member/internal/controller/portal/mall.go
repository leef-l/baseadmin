package portal

import (
	"context"
	"strconv"

	v1 "gbaseadmin/app/member/api/portal/v1"
	"gbaseadmin/app/member/internal/logic/portal"
	"gbaseadmin/app/member/internal/middleware"
)

// Mall 控制器（C 端商城）。
var Mall = cMall{}

type cMall struct{}

// Categories 分类树。
func (c *cMall) Categories(ctx context.Context, req *v1.MallCategoryListReq) (res *v1.MallCategoryListRes, err error) {
	roots, err := portal.AuthLogic().ListShopCategories(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.MallCategoryListRes{List: convertCategoryNodes(roots)}
	return
}

func convertCategoryNodes(nodes []*portal.MallCategoryNode) []*v1.MallCategoryItem {
	out := make([]*v1.MallCategoryItem, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, &v1.MallCategoryItem{
			ID:       n.ID,
			ParentID: n.ParentID,
			Name:     n.Name,
			Icon:     n.Icon,
			Sort:     n.Sort,
			Children: convertCategoryNodes(n.Children),
		})
	}
	return out
}

// Goods 商品分页。
func (c *cMall) Goods(ctx context.Context, req *v1.MallGoodsListReq) (res *v1.MallGoodsListRes, err error) {
	categoryID, _ := strconv.ParseInt(req.CategoryID, 10, 64)
	out, err := portal.AuthLogic().ListShopGoods(ctx, &portal.MallGoodsListInput{
		CategoryID:  categoryID,
		Keyword:     req.Keyword,
		IsRecommend: req.IsRecommend,
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MallGoodsListRes{Total: out.Total, List: convertGoodsItems(out.List)}
	return
}

func convertGoodsItems(items []*portal.MallGoodsListData) []*v1.MallGoodsListItem {
	out := make([]*v1.MallGoodsListItem, 0, len(items))
	for _, it := range items {
		out = append(out, &v1.MallGoodsListItem{
			ID:            it.ID,
			Title:         it.Title,
			Cover:         it.Cover,
			Price:         it.Price,
			OriginalPrice: it.OriginalPrice,
			Sales:         it.Sales,
			Stock:         it.Stock,
			IsRecommend:   it.IsRecommend,
		})
	}
	return out
}

// GoodsDetail 商品详情。
func (c *cMall) GoodsDetail(ctx context.Context, req *v1.MallGoodsDetailReq) (res *v1.MallGoodsDetailRes, err error) {
	id, _ := strconv.ParseInt(req.ID, 10, 64)
	out, err := portal.AuthLogic().GetShopGoodsDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.MallGoodsDetailRes{
		ID:            out.ID,
		CategoryID:    out.CategoryID,
		CategoryName:  out.CategoryName,
		Title:         out.Title,
		Cover:         out.Cover,
		Images:        out.Images,
		Price:         out.Price,
		OriginalPrice: out.OriginalPrice,
		Stock:         out.Stock,
		Sales:         out.Sales,
		Content:       out.Content,
		IsRecommend:   out.IsRecommend,
		Status:        out.Status,
	}, nil
}

// PlaceOrder 下单。
func (c *cMall) PlaceOrder(ctx context.Context, req *v1.MallPlaceOrderReq) (res *v1.MallPlaceOrderRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	goodsID, _ := strconv.ParseInt(req.GoodsID, 10, 64)
	out, err := portal.AuthLogic().PlaceShopOrder(ctx, &portal.PlaceShopOrderInput{
		UserID:   memberID,
		GoodsID:  goodsID,
		Quantity: req.Quantity,
		Remark:   req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.MallPlaceOrderRes{
		OrderID:    out.OrderID,
		OrderNo:    out.OrderNo,
		TotalPrice: portalFormatYuan(out.TotalPrice),
	}, nil
}

// MyOrders 我的商城订单。
func (c *cMall) MyOrders(ctx context.Context, req *v1.MallMyOrdersReq) (res *v1.MallMyOrdersRes, err error) {
	memberID := int64(middleware.CurrentMemberID(ctx))
	out, err := portal.AuthLogic().ListMyShopOrders(ctx, &portal.MyShopOrdersInput{
		UserID:   memberID,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.MallMyOrdersRes{Total: out.Total, List: make([]*v1.MallOrderRecord, 0, len(out.List))}
	for _, item := range out.List {
		res.List = append(res.List, &v1.MallOrderRecord{
			OrderID:    item.OrderID,
			OrderNo:    item.OrderNo,
			GoodsID:    item.GoodsID,
			GoodsTitle: item.GoodsTitle,
			GoodsCover: item.GoodsCover,
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
			Status:     item.Status,
			StatusText: item.StatusText,
			Remark:     item.Remark,
			CreatedAt:  item.CreatedAt,
		})
	}
	return res, nil
}

// portalFormatYuan 把分转字符串元（仅 controller 局部使用，避免依赖 logic 内部 helper）。
func portalFormatYuan(cent int64) string {
	negative := cent < 0
	if negative {
		cent = -cent
	}
	yuan := cent / 100
	fen := cent % 100
	sign := ""
	if negative {
		sign = "-"
	}
	return sign + strconv.FormatInt(yuan, 10) + "." + padFen(fen)
}

func padFen(fen int64) string {
	if fen < 10 {
		return "0" + strconv.FormatInt(fen, 10)
	}
	return strconv.FormatInt(fen, 10)
}
