package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// WarehouseGoods DTO 模型

// WarehouseGoodsCreateInput 创建仓库商品输入
type WarehouseGoodsCreateInput struct {
	GoodsNo string `json:"goodsNo"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	InitPrice int64 `json:"initPrice"`
	CurrentPrice int64 `json:"currentPrice"`
	PriceRiseRate int `json:"priceRiseRate"`
	PlatformFeeRate int `json:"platformFeeRate"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"`
	TradeCount int `json:"tradeCount"`
	GoodsStatus int `json:"goodsStatus"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseGoodsUpdateInput 更新仓库商品输入
type WarehouseGoodsUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsNo string `json:"goodsNo"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	InitPrice int64 `json:"initPrice"`
	CurrentPrice int64 `json:"currentPrice"`
	PriceRiseRate int `json:"priceRiseRate"`
	PlatformFeeRate int `json:"platformFeeRate"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"`
	TradeCount int `json:"tradeCount"`
	GoodsStatus int `json:"goodsStatus"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseGoodsDetailOutput 仓库商品详情输出
type WarehouseGoodsDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsNo string `json:"goodsNo"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	InitPrice int64 `json:"initPrice"`
	CurrentPrice int64 `json:"currentPrice"`
	PriceRiseRate int `json:"priceRiseRate"`
	PlatformFeeRate int `json:"platformFeeRate"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"`
	UserNickname string `json:"userNickname"`
	TradeCount int `json:"tradeCount"`
	GoodsStatus int `json:"goodsStatus"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseGoodsListOutput 仓库商品列表输出
type WarehouseGoodsListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsNo string `json:"goodsNo"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	InitPrice int64 `json:"initPrice"`
	CurrentPrice int64 `json:"currentPrice"`
	PriceRiseRate int `json:"priceRiseRate"`
	PlatformFeeRate int `json:"platformFeeRate"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"`
	UserNickname string `json:"userNickname"`
	TradeCount int `json:"tradeCount"`
	GoodsStatus int `json:"goodsStatus"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseGoodsListInput 仓库商品列表查询输入
type WarehouseGoodsListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	GoodsNo string `json:"goodsNo"`
	Title string `json:"title"`
	OwnerID *snowflake.JsonInt64 `json:"ownerID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	GoodsStatus *int `json:"goodsStatus"`
	Status *int `json:"status"`
}

// WarehouseGoodsBatchUpdateInput 批量编辑仓库商品输入
type WarehouseGoodsBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	GoodsStatus *int `json:"goodsStatus"`
	Status *int `json:"status"`
}

