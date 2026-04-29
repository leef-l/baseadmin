package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Wallet DTO 模型

// WalletCreateInput 创建会员钱包输入
type WalletCreateInput struct {
	UserID snowflake.JsonInt64 `json:"userID"`
	WalletType int `json:"walletType"`
	Balance int64 `json:"balance"`
	TotalIncome int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	FrozenAmount int64 `json:"frozenAmount"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WalletUpdateInput 更新会员钱包输入
type WalletUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	WalletType int `json:"walletType"`
	Balance int64 `json:"balance"`
	TotalIncome int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	FrozenAmount int64 `json:"frozenAmount"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WalletDetailOutput 会员钱包详情输出
type WalletDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserUsername string `json:"userUsername"`
	WalletType int `json:"walletType"`
	Balance int64 `json:"balance"`
	TotalIncome int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	FrozenAmount int64 `json:"frozenAmount"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WalletListOutput 会员钱包列表输出
type WalletListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserUsername string `json:"userUsername"`
	WalletType int `json:"walletType"`
	Balance int64 `json:"balance"`
	TotalIncome int64 `json:"totalIncome"`
	TotalExpense int64 `json:"totalExpense"`
	FrozenAmount int64 `json:"frozenAmount"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WalletListInput 会员钱包列表查询输入
type WalletListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	WalletType *int `json:"walletType"`
	Status *int `json:"status"`
}

// WalletBatchUpdateInput 批量编辑会员钱包输入
type WalletBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	WalletType *int `json:"walletType"`
	Status *int `json:"status"`
}

