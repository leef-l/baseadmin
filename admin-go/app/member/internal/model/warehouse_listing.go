package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// WarehouseListing DTO 模型

// WarehouseListingCreateInput 创建仓库挂卖记录输入
type WarehouseListingCreateInput struct {
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	ListingPrice int64 `json:"listingPrice"`
	ListingStatus int `json:"listingStatus"`
	ListedAt *gtime.Time `json:"listedAt"`
	SoldAt *gtime.Time `json:"soldAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseListingUpdateInput 更新仓库挂卖记录输入
type WarehouseListingUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	ListingPrice int64 `json:"listingPrice"`
	ListingStatus int `json:"listingStatus"`
	ListedAt *gtime.Time `json:"listedAt"`
	SoldAt *gtime.Time `json:"soldAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WarehouseListingDetailOutput 仓库挂卖记录详情输出
type WarehouseListingDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	WarehouseGoodsTitle string `json:"warehouseGoodsTitle"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	UserNickname string `json:"userNickname"`
	ListingPrice int64 `json:"listingPrice"`
	ListingStatus int `json:"listingStatus"`
	ListedAt *gtime.Time `json:"listedAt"`
	SoldAt *gtime.Time `json:"soldAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseListingListOutput 仓库挂卖记录列表输出
type WarehouseListingListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"`
	WarehouseGoodsTitle string `json:"warehouseGoodsTitle"`
	SellerID snowflake.JsonInt64 `json:"sellerID"`
	UserNickname string `json:"userNickname"`
	ListingPrice int64 `json:"listingPrice"`
	ListingStatus int `json:"listingStatus"`
	ListedAt *gtime.Time `json:"listedAt"`
	SoldAt *gtime.Time `json:"soldAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WarehouseListingListInput 仓库挂卖记录列表查询输入
type WarehouseListingListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID"`
	SellerID *snowflake.JsonInt64 `json:"sellerID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	ListingStatus *int `json:"listingStatus"`
	Status *int `json:"status"`
	ListedAtStart string `json:"listedAtStart"`
	ListedAtEnd string `json:"listedAtEnd"`
	SoldAtStart string `json:"soldAtStart"`
	SoldAtEnd string `json:"soldAtEnd"`
}

// WarehouseListingBatchUpdateInput 批量编辑仓库挂卖记录输入
type WarehouseListingBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	ListingStatus *int `json:"listingStatus"`
	Status *int `json:"status"`
}

