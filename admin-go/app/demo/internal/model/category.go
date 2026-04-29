package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Category DTO 模型

// CategoryCreateInput 创建体验分类输入
type CategoryCreateInput struct {
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CategoryUpdateInput 更新体验分类输入
type CategoryUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CategoryDetailOutput 体验分类详情输出
type CategoryDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	CategoryName string `json:"categoryName"`
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

// CategoryListOutput 体验分类列表输出
type CategoryListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	CategoryName string `json:"categoryName"`
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

// CategoryListInput 体验分类列表查询输入
type CategoryListInput struct {
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

// CategoryTreeInput 体验分类树形查询输入
type CategoryTreeInput struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Name string `json:"name"`
	ParentID *snowflake.JsonInt64 `json:"parentID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
}

// CategoryTreeOutput 体验分类树形输出
type CategoryTreeOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ParentID snowflake.JsonInt64 `json:"parentID"`
	CategoryName string `json:"categoryName"`
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
	Children  []*CategoryTreeOutput `json:"children"`
}

// CategoryBatchUpdateInput 批量编辑体验分类输入
type CategoryBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

