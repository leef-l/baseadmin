package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// ShopCategory DTO 模型

// ShopCategoryCreateInput 创建商城商品分类输入
type ShopCategoryCreateInput struct {
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopCategoryUpdateInput 更新商城商品分类输入
type ShopCategoryUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ShopCategoryDetailOutput 商城商品分类详情输出
type ShopCategoryDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	ShopCategoryName string `json:"shopCategoryName"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopCategoryListOutput 商城商品分类列表输出
type ShopCategoryListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	ShopCategoryName string `json:"shopCategoryName"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ShopCategoryListInput 商城商品分类列表查询输入
type ShopCategoryListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Name string `json:"name"`
	ParentID *snowflake.JsonInt64 `json:"parentID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
}

// ShopCategoryTreeInput 商城商品分类树形查询输入
type ShopCategoryTreeInput struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Name string `json:"name"`
	ParentID *snowflake.JsonInt64 `json:"parentID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
}

// ShopCategoryTreeOutput 商城商品分类树形输出
type ShopCategoryTreeOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	ShopCategoryName string `json:"shopCategoryName"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time                `json:"createdAt"`
	UpdatedAt *gtime.Time                `json:"updatedAt"`
	Children  []*ShopCategoryTreeOutput `json:"children"`
}

// ShopCategoryBatchUpdateInput 批量编辑商城商品分类输入
type ShopCategoryBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

