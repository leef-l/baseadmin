package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// ShopGoods DTO 模型

// ShopGoodsCreateInput 创建商城商品输入
type ShopGoodsCreateInput struct {
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Images string `json:"images"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"originalPrice"`
	Stock int `json:"stock"`
	Sales int `json:"sales"`
	Content string `json:"content"`
	Sort int `json:"sort"`
	IsRecommend int `json:"isRecommend"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopGoodsUpdateInput 更新商城商品输入
type ShopGoodsUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Images string `json:"images"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"originalPrice"`
	Stock int `json:"stock"`
	Sales int `json:"sales"`
	Content string `json:"content"`
	Sort int `json:"sort"`
	IsRecommend int `json:"isRecommend"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopGoodsDetailOutput 商城商品详情输出
type ShopGoodsDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	ShopCategoryName string `json:"shopCategoryName"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Images string `json:"images"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"originalPrice"`
	Stock int `json:"stock"`
	Sales int `json:"sales"`
	Content string `json:"content"`
	Sort int `json:"sort"`
	IsRecommend int `json:"isRecommend"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopGoodsListOutput 商城商品列表输出
type ShopGoodsListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	ShopCategoryName string `json:"shopCategoryName"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Images string `json:"images"`
	Price int64 `json:"price"`
	OriginalPrice int64 `json:"originalPrice"`
	Stock int `json:"stock"`
	Sales int `json:"sales"`
	Content string `json:"content"`
	Sort int `json:"sort"`
	IsRecommend int `json:"isRecommend"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopGoodsListInput 商城商品列表查询输入
type ShopGoodsListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Title string `json:"title"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsRecommend *int `json:"isRecommend"`
	Status *int `json:"status"`
}

// ShopGoodsBatchUpdateInput 批量编辑商城商品输入
type ShopGoodsBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	IsRecommend *int `json:"isRecommend"`
	Status *int `json:"status"`
}

