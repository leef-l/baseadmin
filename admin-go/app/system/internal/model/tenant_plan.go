package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// TenantPlan DTO 模型

// TenantPlanCreateInput 创建租户套餐订阅表输入
type TenantPlanCreateInput struct {
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	PlanID snowflake.JsonInt64 `json:"planID"`
	StartAt *gtime.Time `json:"startAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Status int `json:"status"`
	Remark string `json:"remark"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// TenantPlanUpdateInput 更新租户套餐订阅表输入
type TenantPlanUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	PlanID snowflake.JsonInt64 `json:"planID"`
	StartAt *gtime.Time `json:"startAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Status int `json:"status"`
	Remark string `json:"remark"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// TenantPlanDetailOutput 租户套餐订阅表详情输出
type TenantPlanDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	PlanID snowflake.JsonInt64 `json:"planID"`
	PlanName string `json:"planName"`
	StartAt *gtime.Time `json:"startAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Status int `json:"status"`
	Remark string `json:"remark"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// TenantPlanListOutput 租户套餐订阅表列表输出
type TenantPlanListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	PlanID snowflake.JsonInt64 `json:"planID"`
	PlanName string `json:"planName"`
	StartAt *gtime.Time `json:"startAt"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Status int `json:"status"`
	Remark string `json:"remark"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// TenantPlanListInput 租户套餐订阅表列表查询输入
type TenantPlanListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	PlanID *snowflake.JsonInt64 `json:"planID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
	Remark string `json:"remark"`
	StartAtStart string `json:"startAtStart"`
	StartAtEnd string `json:"startAtEnd"`
	ExpireAtStart string `json:"expireAtStart"`
	ExpireAtEnd string `json:"expireAtEnd"`
}

// TenantPlanBatchUpdateInput 批量编辑租户套餐订阅表输入
type TenantPlanBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

