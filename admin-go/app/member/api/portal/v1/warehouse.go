package v1

import "github.com/gogf/gf/v2/frame/g"

// ----- 我的仓库（持有 / 挂卖中 / 交易中） -----

// MyWarehouseReq 我的仓库列表（按状态筛选）。
type MyWarehouseReq struct {
	g.Meta   `path:"/warehouse/my" method:"get" tags:"会员-仓库" summary:"我的仓库"`
	Status   int `json:"status" v:"in:0,1,2,3" dc:"商品状态 0=全部 1=持有中 2=挂卖中 3=交易中"`
	PageNum  int `json:"pageNum" d:"1"`
	PageSize int `json:"pageSize" d:"20"`
}

// MyWarehouseRes 我的仓库列表。
type MyWarehouseRes struct {
	g.Meta `mime:"application/json"`
	Total  int                   `json:"total"`
	List   []*MyWarehouseGoodsItem `json:"list"`
}

// MyWarehouseGoodsItem 我持有的仓库商品。
type MyWarehouseGoodsItem struct {
	ID                string `json:"id"`
	GoodsNo           string `json:"goodsNo"`
	Title             string `json:"title"`
	Cover             string `json:"cover"`
	InitPrice         string `json:"initPrice"`
	CurrentPrice      string `json:"currentPrice" dc:"我持有时的价格（元）"`
	NextListingPrice  string `json:"nextListingPrice" dc:"下次挂卖价（按 price_rise_rate 自动算，元）"`
	PriceRiseRate     int    `json:"priceRiseRate" dc:"加价比例（百分比）"`
	PlatformFeeRate   int    `json:"platformFeeRate" dc:"平台抽成比例（百分比）"`
	TradeCount        int    `json:"tradeCount"`
	GoodsStatus       int    `json:"goodsStatus"`
	GoodsStatusText   string `json:"goodsStatusText"`
	ActiveListingID   string `json:"activeListingId,omitempty" dc:"挂卖中时返回当前挂卖记录 ID，便于取消"`
}

// ----- 挂卖（一键挂卖，价格系统算） -----

// WarehouseListGoodsReq 挂卖请求（卖家不可改价）。
type WarehouseListGoodsReq struct {
	g.Meta  `path:"/warehouse/list" method:"post" tags:"会员-仓库" summary:"挂卖商品"`
	GoodsID string `json:"goodsId" v:"required#商品 ID 不能为空"`
}

// WarehouseListGoodsRes 挂卖响应。
type WarehouseListGoodsRes struct {
	g.Meta       `mime:"application/json"`
	ListingID    string `json:"listingId"`
	ListingPrice string `json:"listingPrice" dc:"自动加价后的挂卖价（元）"`
}

// ----- 仓库市场（所有挂卖中的商品） -----

// WarehouseMarketReq 市场列表。
type WarehouseMarketReq struct {
	g.Meta   `path:"/warehouse/market" method:"get" tags:"会员-仓库" summary:"仓库市场"`
	Keyword  string `json:"keyword" dc:"商品名称关键字"`
	OrderBy  string `json:"orderBy" v:"in:price_asc,price_desc,latest" d:"latest" dc:"排序：price_asc/price_desc/latest"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"20"`
}

// WarehouseMarketRes 市场列表。
type WarehouseMarketRes struct {
	g.Meta `mime:"application/json"`
	Total  int                       `json:"total"`
	List   []*WarehouseMarketListing `json:"list"`
}

// WarehouseMarketListing 单个挂卖商品。
type WarehouseMarketListing struct {
	ListingID    string `json:"listingId"`
	GoodsID      string `json:"goodsId"`
	GoodsNo      string `json:"goodsNo"`
	Title        string `json:"title"`
	Cover        string `json:"cover"`
	ListingPrice string `json:"listingPrice"`
	SellerName   string `json:"sellerName"`
	TradeCount   int    `json:"tradeCount"`
	ListedAt     string `json:"listedAt"`
}

// ----- 买家下单 -----

// WarehousePlaceTradeReq 买家下单。
type WarehousePlaceTradeReq struct {
	g.Meta    `path:"/warehouse/trade/place" method:"post" tags:"会员-仓库" summary:"购买仓库商品"`
	ListingID string `json:"listingId" v:"required#挂卖 ID 不能为空"`
}

// WarehousePlaceTradeRes 下单响应。
type WarehousePlaceTradeRes struct {
	g.Meta  `mime:"application/json"`
	TradeID string `json:"tradeId"`
	TradeNo string `json:"tradeNo"`
}

// ----- 卖家确认 -----

// WarehouseConfirmTradeReq 卖家确认。
type WarehouseConfirmTradeReq struct {
	g.Meta  `path:"/warehouse/trade/confirm" method:"post" tags:"会员-仓库" summary:"卖家确认交易"`
	TradeID string `json:"tradeId" v:"required#交易 ID 不能为空"`
}

// WarehouseConfirmTradeRes 确认响应。
type WarehouseConfirmTradeRes struct {
	g.Meta       `mime:"application/json"`
	TradeID      string `json:"tradeId"`
	TradePrice   string `json:"tradePrice" dc:"成交价（元）"`
	PlatformFee  string `json:"platformFee" dc:"平台抽成（元）"`
	SellerIncome string `json:"sellerIncome" dc:"卖家奖金入账（元）"`
}

// ----- 我的交易记录 -----

// MyTradesReq 我的交易记录（买入 / 卖出）。
type MyTradesReq struct {
	g.Meta   `path:"/warehouse/my-trades" method:"get" tags:"会员-仓库" summary:"我的交易记录"`
	Role     string `json:"role" v:"in:buyer,seller" d:"buyer" dc:"buyer=买入 seller=卖出"`
	Status   int    `json:"status" v:"in:0,1,2,3" dc:"交易状态 0=全部 1=待确认 2=已完成 3=已取消"`
	PageNum  int    `json:"pageNum" d:"1"`
	PageSize int    `json:"pageSize" d:"20"`
}

// MyTradesRes 交易列表。
type MyTradesRes struct {
	g.Meta `mime:"application/json"`
	Total  int            `json:"total"`
	List   []*TradeRecord `json:"list"`
}

// TradeRecord 单条交易。
type TradeRecord struct {
	TradeID         string `json:"tradeId"`
	TradeNo         string `json:"tradeNo"`
	GoodsID         string `json:"goodsId"`
	GoodsNo         string `json:"goodsNo"`
	GoodsTitle      string `json:"goodsTitle"`
	GoodsCover      string `json:"goodsCover"`
	TradePrice      string `json:"tradePrice"`
	PlatformFee     string `json:"platformFee"`
	SellerIncome    string `json:"sellerIncome"`
	TradeStatus     int    `json:"tradeStatus"`
	TradeStatusText string `json:"tradeStatusText"`
	Counterparty    string `json:"counterparty" dc:"对方昵称"`
	CreatedAt       string `json:"createdAt"`
	ConfirmedAt     string `json:"confirmedAt"`
}
