package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Product DTO 模型

// ProductCreateInput 创建体验商品输入
type ProductCreateInput struct {
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	SkuNo string `json:"skuNo"`
	Name string `json:"name"`
	Cover string `json:"cover"`
	ManualFile string `json:"manualFile"`
	DetailContent string `json:"detailContent"`
	SpecJSON string `json:"specJSON"`
	WebsiteURL string `json:"websiteURL"`
	Type int `json:"type"`
	IsRecommend int `json:"isRecommend"`
	SalePrice int `json:"salePrice"`
	StockNum int `json:"stockNum"`
	WeightNum int `json:"weightNum"`
	Sort int `json:"sort"`
	Icon string `json:"icon"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ProductUpdateInput 更新体验商品输入
type ProductUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	SkuNo string `json:"skuNo"`
	Name string `json:"name"`
	Cover string `json:"cover"`
	ManualFile string `json:"manualFile"`
	DetailContent string `json:"detailContent"`
	SpecJSON string `json:"specJSON"`
	WebsiteURL string `json:"websiteURL"`
	Type int `json:"type"`
	IsRecommend int `json:"isRecommend"`
	SalePrice int `json:"salePrice"`
	StockNum int `json:"stockNum"`
	WeightNum int `json:"weightNum"`
	Sort int `json:"sort"`
	Icon string `json:"icon"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ProductDetailOutput 体验商品详情输出
type ProductDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	CategoryName string `json:"categoryName"`
	SkuNo string `json:"skuNo"`
	Name string `json:"name"`
	Cover string `json:"cover"`
	ManualFile string `json:"manualFile"`
	DetailContent string `json:"detailContent"`
	SpecJSON string `json:"specJSON"`
	WebsiteURL string `json:"websiteURL"`
	Type int `json:"type"`
	IsRecommend int `json:"isRecommend"`
	SalePrice int `json:"salePrice"`
	StockNum int `json:"stockNum"`
	WeightNum int `json:"weightNum"`
	Sort int `json:"sort"`
	Icon string `json:"icon"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ProductListOutput 体验商品列表输出
type ProductListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"`
	CategoryName string `json:"categoryName"`
	SkuNo string `json:"skuNo"`
	Name string `json:"name"`
	Cover string `json:"cover"`
	ManualFile string `json:"manualFile"`
	DetailContent string `json:"detailContent"`
	SpecJSON string `json:"specJSON"`
	WebsiteURL string `json:"websiteURL"`
	Type int `json:"type"`
	IsRecommend int `json:"isRecommend"`
	SalePrice int `json:"salePrice"`
	StockNum int `json:"stockNum"`
	WeightNum int `json:"weightNum"`
	Sort int `json:"sort"`
	Icon string `json:"icon"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ProductListInput 体验商品列表查询输入
type ProductListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	SkuNo string `json:"skuNo"`
	Name string `json:"name"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Type *int `json:"type"`
	IsRecommend *int `json:"isRecommend"`
	Status *int `json:"status"`
}

// ProductBatchUpdateInput 批量编辑体验商品输入
type ProductBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Type *int `json:"type"`
	IsRecommend *int `json:"isRecommend"`
	Status *int `json:"status"`
}

