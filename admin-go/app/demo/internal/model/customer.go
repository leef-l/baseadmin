package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Customer DTO 模型

// CustomerCreateInput 创建体验客户输入
type CustomerCreateInput struct {
	Avatar string `json:"avatar"`
	Name string `json:"name"`
	CustomerNo string `json:"customerNo"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Gender int `json:"gender"`
	Level int `json:"level"`
	SourceType int `json:"sourceType"`
	IsVip int `json:"isVip"`
	RegisteredAt *gtime.Time `json:"registeredAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CustomerUpdateInput 更新体验客户输入
type CustomerUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
	CustomerNo string `json:"customerNo"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Gender int `json:"gender"`
	Level int `json:"level"`
	SourceType int `json:"sourceType"`
	IsVip int `json:"isVip"`
	RegisteredAt *gtime.Time `json:"registeredAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// CustomerDetailOutput 体验客户详情输出
type CustomerDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
	CustomerNo string `json:"customerNo"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Gender int `json:"gender"`
	Level int `json:"level"`
	SourceType int `json:"sourceType"`
	IsVip int `json:"isVip"`
	RegisteredAt *gtime.Time `json:"registeredAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CustomerListOutput 体验客户列表输出
type CustomerListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Avatar string `json:"avatar"`
	Name string `json:"name"`
	CustomerNo string `json:"customerNo"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Gender int `json:"gender"`
	Level int `json:"level"`
	SourceType int `json:"sourceType"`
	IsVip int `json:"isVip"`
	RegisteredAt *gtime.Time `json:"registeredAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// CustomerListInput 体验客户列表查询输入
type CustomerListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	CustomerNo string `json:"customerNo"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Gender *int `json:"gender"`
	Level *int `json:"level"`
	SourceType *int `json:"sourceType"`
	IsVip *int `json:"isVip"`
	Status *int `json:"status"`
	RegisteredAtStart string `json:"registeredAtStart"`
	RegisteredAtEnd string `json:"registeredAtEnd"`
}

// CustomerBatchUpdateInput 批量编辑体验客户输入
type CustomerBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Gender *int `json:"gender"`
	Level *int `json:"level"`
	SourceType *int `json:"sourceType"`
	IsVip *int `json:"isVip"`
	Status *int `json:"status"`
}

