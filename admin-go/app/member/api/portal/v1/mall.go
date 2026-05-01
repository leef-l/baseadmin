package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 分类 -----

// MallCategoryListReq 获取商城分类树（带商品数量）。
type MallCategoryListReq struct {
	g.Meta `path:"/mall/categories" method:"get" tags:"会员-商城" summary:"商城分类树"`
}

// MallCategoryListRes 分类树。
type MallCategoryListRes struct {
	g.Meta `mime:"application/json"`
	List   []*MallCategoryItem `json:"list"`
}

// MallCategoryItem 分类节点。
type MallCategoryItem struct {
	ID       string              `json:"id"`
	ParentID string              `json:"parentId"`
	Name     string              `json:"name"`
	Icon     string              `json:"icon"`
	Sort     int                 `json:"sort"`
	Children []*MallCategoryItem `json:"children"`
}

// ----- 商品列表 -----

// MallGoodsListReq 商城商品分页（按分类筛选 / 推荐 / 关键字）。
type MallGoodsListReq struct {
	g.Meta      `path:"/mall/goods" method:"get" tags:"会员-商城" summary:"商城商品列表"`
	CategoryID  string `json:"categoryId" dc:"分类 ID（可选）"`
	Keyword     string `json:"keyword" dc:"标题关键字"`
	IsRecommend int    `json:"isRecommend" v:"in:0,1" dc:"是否推荐 1=仅推荐"`
	PageNum     int    `json:"pageNum" d:"1"`
	PageSize    int    `json:"pageSize" d:"20"`
}

// MallGoodsListRes 商品列表。
type MallGoodsListRes struct {
	g.Meta `mime:"application/json"`
	Total  int                  `json:"total"`
	List   []*MallGoodsListItem `json:"list"`
}

// MallGoodsListItem 列表项。
type MallGoodsListItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Cover         string `json:"cover"`
	Price         string `json:"price" dc:"售价（元，优惠券钱包支付）"`
	OriginalPrice string `json:"originalPrice"`
	Sales         int    `json:"sales"`
	Stock         int    `json:"stock"`
	IsRecommend   int    `json:"isRecommend"`
}

// ----- 商品详情 -----

// MallGoodsDetailReq 商品详情。
type MallGoodsDetailReq struct {
	g.Meta `path:"/mall/goods/detail" method:"get" tags:"会员-商城" summary:"商品详情"`
	ID     string `json:"id" v:"required#商品 ID 不能为空"`
}

// MallGoodsDetailRes 商品详情。
type MallGoodsDetailRes struct {
	g.Meta        `mime:"application/json"`
	ID            string   `json:"id"`
	CategoryID    string   `json:"categoryId"`
	CategoryName  string   `json:"categoryName"`
	Title         string   `json:"title"`
	Cover         string   `json:"cover"`
	Images        []string `json:"images"`
	Price         string   `json:"price"`
	OriginalPrice string   `json:"originalPrice"`
	Stock         int      `json:"stock"`
	Sales         int      `json:"sales"`
	Content       string   `json:"content" dc:"商品详情富文本 HTML"`
	IsRecommend   int      `json:"isRecommend"`
	Status        int      `json:"status"`
}

// ----- 下单 -----

// MallPlaceOrderReq 下单（扣优惠券钱包）。
type MallPlaceOrderReq struct {
	g.Meta   `path:"/mall/order/place" method:"post" tags:"会员-商城" summary:"商城下单"`
	GoodsID  string `json:"goodsId" v:"required#商品 ID 不能为空"`
	Quantity int    `json:"quantity" v:"min:1|max:999" d:"1" dc:"购买数量"`
	Remark   string `json:"remark" v:"max-length:500" dc:"备注"`
}

// MallPlaceOrderRes 下单响应。
type MallPlaceOrderRes struct {
	g.Meta     `mime:"application/json"`
	OrderID    string `json:"orderId"`
	OrderNo    string `json:"orderNo"`
	TotalPrice string `json:"totalPrice" dc:"订单总价（元）"`
}

// ----- 我的商城订单 -----

// MallMyOrdersReq 我的商城订单分页。
type MallMyOrdersReq struct {
	g.Meta   `path:"/mall/orders" method:"get" tags:"会员-商城" summary:"我的商城订单"`
	PageNum  int `json:"pageNum" d:"1"`
	PageSize int `json:"pageSize" d:"20"`
}

// MallMyOrdersRes 订单列表。
type MallMyOrdersRes struct {
	g.Meta `mime:"application/json"`
	Total  int                `json:"total"`
	List   []*MallOrderRecord `json:"list"`
}

// MallOrderRecord 订单。
type MallOrderRecord struct {
	OrderID    string `json:"orderId"`
	OrderNo    string `json:"orderNo"`
	GoodsID    string `json:"goodsId"`
	GoodsTitle string `json:"goodsTitle"`
	GoodsCover string `json:"goodsCover"`
	Quantity   int    `json:"quantity"`
	TotalPrice string `json:"totalPrice"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Remark     string `json:"remark"`
	CreatedAt  string `json:"createdAt"`
}
