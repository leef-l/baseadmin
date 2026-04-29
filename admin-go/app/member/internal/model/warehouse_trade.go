package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// WarehouseTrade DTO 模型

// WarehouseTradeCreateInput 创建仓库交易记录输入
type WarehouseTradeCreateInput struct {
	TradeNo string `json:"tradeNo"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	ListingID snowflake.JsonInt64 `json:"listingID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"`
	TradePrice int64 `json:"tradePrice"`
	PlatformFee int64 `json:"platformFee"`
	SellerIncome int64 `json:"sellerIncome"`
	TradeStatus int `json:"tradeStatus"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseTradeUpdateInput 更新仓库交易记录输入
type WarehouseTradeUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TradeNo string `json:"tradeNo"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	ListingID snowflake.JsonInt64 `json:"listingID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"`
	TradePrice int64 `json:"tradePrice"`
	PlatformFee int64 `json:"platformFee"`
	SellerIncome int64 `json:"sellerIncome"`
	TradeStatus int `json:"tradeStatus"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseTradeDetailOutput 仓库交易记录详情输出
type WarehouseTradeDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TradeNo string `json:"tradeNo"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	WarehouseGoodsTitle string `json:"warehouseGoodsTitle"`
	ListingID snowflake.JsonInt64 `json:"listingID"`
	WarehouseListingID string `json:"warehouseListingID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	UserNickname string `json:"userNickname"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"`
	BuyerNickname string `json:"buyerNickname"`
	TradePrice int64 `json:"tradePrice"`
	PlatformFee int64 `json:"platformFee"`
	SellerIncome int64 `json:"sellerIncome"`
	TradeStatus int `json:"tradeStatus"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseTradeListOutput 仓库交易记录列表输出
type WarehouseTradeListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TradeNo string `json:"tradeNo"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	WarehouseGoodsTitle string `json:"warehouseGoodsTitle"`
	ListingID snowflake.JsonInt64 `json:"listingID"`
	WarehouseListingID string `json:"warehouseListingID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	UserNickname string `json:"userNickname"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"`
	BuyerNickname string `json:"buyerNickname"`
	TradePrice int64 `json:"tradePrice"`
	PlatformFee int64 `json:"platformFee"`
	SellerIncome int64 `json:"sellerIncome"`
	TradeStatus int `json:"tradeStatus"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseTradeListInput 仓库交易记录列表查询输入
type WarehouseTradeListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TradeNo string `json:"tradeNo"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID"`
	ListingID *snowflake.JsonInt64 `json:"listingID"`
	SellerID *snowflake.JsonInt64 `json:"sellerID"`
	BuyerID *snowflake.JsonInt64 `json:"buyerID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	TradeStatus *int `json:"tradeStatus"`
	Status *int `json:"status"`
	ConfirmedAtStart string `json:"confirmedAtStart"`
	ConfirmedAtEnd string `json:"confirmedAtEnd"`
}

// WarehouseTradeBatchUpdateInput 批量编辑仓库交易记录输入
type WarehouseTradeBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	TradeStatus *int `json:"tradeStatus"`
	Status *int `json:"status"`
}

