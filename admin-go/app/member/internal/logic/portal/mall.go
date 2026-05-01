package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/model/entity"
)

// ----- 分类 -----

// MallCategoryNode 分类树节点。
type MallCategoryNode struct {
	ID       string
	ParentID string
	Name     string
	Icon     string
	Sort     int
	Children []*MallCategoryNode
}

// ListShopCategories 返回商城分类树（status=1 的全量）。
func (s *sPortalAuth) ListShopCategories(ctx context.Context) ([]*MallCategoryNode, error) {
	var rows []entity.MemberShopCategory
	if err := dao.MemberShopCategory.Ctx(ctx).
		Where(dao.MemberShopCategory.Columns().Status, 1).
		Where(dao.MemberShopCategory.Columns().DeletedAt, nil).
		OrderAsc(dao.MemberShopCategory.Columns().Sort).
		OrderAsc(dao.MemberShopCategory.Columns().Id).
		Scan(&rows); err != nil {
		return nil, err
	}
	nodeMap := make(map[uint64]*MallCategoryNode, len(rows))
	for _, row := range rows {
		nodeMap[row.Id] = &MallCategoryNode{
			ID:       fmt.Sprintf("%d", row.Id),
			ParentID: fmt.Sprintf("%d", row.ParentId),
			Name:     row.Name,
			Icon:     row.Icon,
			Sort:     row.Sort,
			Children: []*MallCategoryNode{},
		}
	}
	roots := make([]*MallCategoryNode, 0)
	for _, row := range rows {
		node := nodeMap[row.Id]
		if row.ParentId == 0 {
			roots = append(roots, node)
			continue
		}
		if parent, ok := nodeMap[row.ParentId]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			roots = append(roots, node)
		}
	}
	return roots, nil
}

// ----- 商品列表 -----

// MallGoodsListInput 商城商品分页入参。
type MallGoodsListInput struct {
	CategoryID  int64
	Keyword     string
	IsRecommend int
	PageNum     int
	PageSize    int
}

// MallGoodsListOutput 商品列表。
type MallGoodsListOutput struct {
	Total int
	List  []*MallGoodsListData
}

// MallGoodsListData 列表数据（金额已转字符串元）。
type MallGoodsListData struct {
	ID            string
	Title         string
	Cover         string
	Price         string
	OriginalPrice string
	Sales         int
	Stock         int
	IsRecommend   int
}

// ListShopGoods 商品分页。
func (s *sPortalAuth) ListShopGoods(ctx context.Context, in *MallGoodsListInput) (*MallGoodsListOutput, error) {
	if in == nil {
		in = &MallGoodsListInput{}
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
	m := dao.MemberShopGoods.Ctx(ctx).
		Where(dao.MemberShopGoods.Columns().Status, 1).
		Where(dao.MemberShopGoods.Columns().DeletedAt, nil)
	if in.CategoryID > 0 {
		m = m.Where(dao.MemberShopGoods.Columns().CategoryId, in.CategoryID)
	}
	if v := strings.TrimSpace(in.Keyword); v != "" {
		m = m.WhereLike(dao.MemberShopGoods.Columns().Title, "%"+v+"%")
	}
	if in.IsRecommend == 1 {
		m = m.Where(dao.MemberShopGoods.Columns().IsRecommend, 1)
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberShopGoods
	if err := m.OrderAsc(dao.MemberShopGoods.Columns().Sort).
		OrderDesc(dao.MemberShopGoods.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&rows); err != nil {
		return nil, err
	}
	out := &MallGoodsListOutput{Total: total, List: make([]*MallGoodsListData, 0, len(rows))}
	for _, row := range rows {
		out.List = append(out.List, &MallGoodsListData{
			ID:            fmt.Sprintf("%d", row.Id),
			Title:         row.Title,
			Cover:         row.Cover,
			Price:         formatCent(int64(row.Price)),
			OriginalPrice: formatCent(int64(row.OriginalPrice)),
			Sales:         int(row.Sales),
			Stock:         int(row.Stock),
			IsRecommend:   row.IsRecommend,
		})
	}
	return out, nil
}

// ----- 商品详情 -----

// MallGoodsDetailData 详情数据。
type MallGoodsDetailData struct {
	ID            string
	CategoryID    string
	CategoryName  string
	Title         string
	Cover         string
	Images        []string
	Price         string
	OriginalPrice string
	Stock         int
	Sales         int
	Content       string
	IsRecommend   int
	Status        int
}

// GetShopGoodsDetail 商品详情。
func (s *sPortalAuth) GetShopGoodsDetail(ctx context.Context, goodsID int64) (*MallGoodsDetailData, error) {
	if goodsID <= 0 {
		return nil, gerror.New("商品 ID 不能为空")
	}
	var row entity.MemberShopGoods
	if err := dao.MemberShopGoods.Ctx(ctx).
		Where(dao.MemberShopGoods.Columns().Id, goodsID).
		Where(dao.MemberShopGoods.Columns().DeletedAt, nil).
		Scan(&row); err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, gerror.New("商品不存在")
	}

	categoryName := ""
	if row.CategoryId > 0 {
		v, _ := dao.MemberShopCategory.Ctx(ctx).
			Where(dao.MemberShopCategory.Columns().Id, row.CategoryId).
			Where(dao.MemberShopCategory.Columns().DeletedAt, nil).
			Value(dao.MemberShopCategory.Columns().Name)
		if v != nil {
			categoryName = v.String()
		}
	}

	images := parseImagesJSON(row.Images)

	return &MallGoodsDetailData{
		ID:            fmt.Sprintf("%d", row.Id),
		CategoryID:    fmt.Sprintf("%d", row.CategoryId),
		CategoryName:  categoryName,
		Title:         row.Title,
		Cover:         row.Cover,
		Images:        images,
		Price:         formatCent(int64(row.Price)),
		OriginalPrice: formatCent(int64(row.OriginalPrice)),
		Stock:         int(row.Stock),
		Sales:         int(row.Sales),
		Content:       row.Content,
		IsRecommend:   row.IsRecommend,
		Status:        row.Status,
	}, nil
}

// parseImagesJSON 解析 member_shop_goods.images（JSON 数组）。
// 兼容老数据：如果不是 JSON 数组，按逗号分割。
func parseImagesJSON(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") {
		var arr []string
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			return arr
		}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// ----- 我的商城订单 -----

// MyShopOrdersInput 我的商城订单分页。
type MyShopOrdersInput struct {
	UserID   int64
	PageNum  int
	PageSize int
}

// MyShopOrdersOutput 商城订单列表。
type MyShopOrdersOutput struct {
	Total int
	List  []*ShopOrderItem
}

// ShopOrderItem 订单。
type ShopOrderItem struct {
	OrderID    string
	OrderNo    string
	GoodsID    string
	GoodsTitle string
	GoodsCover string
	Quantity   int
	TotalPrice string
	Status     int
	StatusText string
	Remark     string
	CreatedAt  string
}

// ListMyShopOrders 我的商城订单分页。
func (s *sPortalAuth) ListMyShopOrders(ctx context.Context, in *MyShopOrdersInput) (*MyShopOrdersOutput, error) {
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
	m := dao.MemberShopOrder.Ctx(ctx).
		Where(dao.MemberShopOrder.Columns().UserId, in.UserID).
		Where(dao.MemberShopOrder.Columns().DeletedAt, nil)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var rows []entity.MemberShopOrder
	if err := m.OrderDesc(dao.MemberShopOrder.Columns().Id).
		Page(pageNum, pageSize).
		Scan(&rows); err != nil {
		return nil, err
	}
	out := &MyShopOrdersOutput{Total: total, List: make([]*ShopOrderItem, 0, len(rows))}
	for _, row := range rows {
		out.List = append(out.List, &ShopOrderItem{
			OrderID:    fmt.Sprintf("%d", row.Id),
			OrderNo:    row.OrderNo,
			GoodsID:    fmt.Sprintf("%d", row.GoodsId),
			GoodsTitle: row.GoodsTitle,
			GoodsCover: row.GoodsCover,
			Quantity:   int(row.Quantity),
			TotalPrice: formatCent(int64(row.TotalPrice)),
			Status:     row.OrderStatus,
			StatusText: shopOrderStatusText(row.OrderStatus),
			Remark:     row.Remark,
			CreatedAt:  timeStr(row.CreatedAt),
		})
	}
	return out, nil
}

func shopOrderStatusText(status int) string {
	switch status {
	case 1:
		return "已完成"
	case 2:
		return "已取消"
	}
	return "未知"
}
