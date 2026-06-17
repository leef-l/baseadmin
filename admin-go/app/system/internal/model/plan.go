package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Plan DTO 模型

// PlanCreateInput 创建套餐表输入
type PlanCreateInput struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
	UserLimit int `json:"userLimit"`
	StorageMb int `json:"storageMb"`
	Features string `json:"features"`
	Price int `json:"price"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// PlanUpdateInput 更新套餐表输入
type PlanUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
	UserLimit int `json:"userLimit"`
	StorageMb int `json:"storageMb"`
	Features string `json:"features"`
	Price int `json:"price"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// PlanDetailOutput 套餐表详情输出
type PlanDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
	UserLimit int `json:"userLimit"`
	StorageMb int `json:"storageMb"`
	Features string `json:"features"`
	Price int `json:"price"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// PlanListOutput 套餐表列表输出
type PlanListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Description string `json:"description"`
	UserLimit int `json:"userLimit"`
	StorageMb int `json:"storageMb"`
	Features string `json:"features"`
	Price int `json:"price"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// PlanListInput 套餐表列表查询输入
type PlanListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	Code string `json:"code"`
	Name string `json:"name"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
	Description string `json:"description"`
}

// PlanBatchUpdateInput 批量编辑套餐表输入
type PlanBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

