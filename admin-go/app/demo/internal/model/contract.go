package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Contract DTO 模型

// ContractCreateInput 创建体验合同输入
type ContractCreateInput struct {
	ContractNo string `json:"contractNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	Title string `json:"title"`
	ContractFile string `json:"contractFile"`
	SignImage string `json:"signImage"`
	ContractAmount int `json:"contractAmount"`
	SignPassword string `json:"signPassword"`
	SignedAt *gtime.Time `json:"signedAt"`
	ExpiresAt *gtime.Time `json:"expiresAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ContractUpdateInput 更新体验合同输入
type ContractUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ContractNo string `json:"contractNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	Title string `json:"title"`
	ContractFile string `json:"contractFile"`
	SignImage string `json:"signImage"`
	ContractAmount int `json:"contractAmount"`
	SignPassword string `json:"signPassword"`
	SignedAt *gtime.Time `json:"signedAt"`
	ExpiresAt *gtime.Time `json:"expiresAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// ContractDetailOutput 体验合同详情输出
type ContractDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ContractNo string `json:"contractNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	OrderOrderNo string `json:"orderOrderNo"`
	Title string `json:"title"`
	ContractFile string `json:"contractFile"`
	SignImage string `json:"signImage"`
	ContractAmount int `json:"contractAmount"`
	SignedAt *gtime.Time `json:"signedAt"`
	ExpiresAt *gtime.Time `json:"expiresAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ContractListOutput 体验合同列表输出
type ContractListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	ContractNo string `json:"contractNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	OrderOrderNo string `json:"orderOrderNo"`
	Title string `json:"title"`
	ContractFile string `json:"contractFile"`
	SignImage string `json:"signImage"`
	ContractAmount int `json:"contractAmount"`
	SignedAt *gtime.Time `json:"signedAt"`
	ExpiresAt *gtime.Time `json:"expiresAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// ContractListInput 体验合同列表查询输入
type ContractListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	ContractNo string `json:"contractNo"`
	Title string `json:"title"`
	CustomerID *snowflake.JsonInt64 `json:"customerID"`
	OrderID *snowflake.JsonInt64 `json:"orderID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Status *int `json:"status"`
	SignedAtStart string `json:"signedAtStart"`
	SignedAtEnd string `json:"signedAtEnd"`
	ExpiresAtStart string `json:"expiresAtStart"`
	ExpiresAtEnd string `json:"expiresAtEnd"`
}

// ContractBatchUpdateInput 批量编辑体验合同输入
type ContractBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Status *int `json:"status"`
}

