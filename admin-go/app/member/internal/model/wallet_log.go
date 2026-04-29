package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// WalletLog DTO 模型

// WalletLogCreateInput 创建钱包流水记录输入
type WalletLogCreateInput struct {
	UserID snowflake.JsonInt64 `json:"userID"`
	WalletType int `json:"walletType"`
	ChangeType int `json:"changeType"`
	ChangeAmount int64 `json:"changeAmount"`
	BeforeBalance int64 `json:"beforeBalance"`
	AfterBalance int64 `json:"afterBalance"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WalletLogUpdateInput 更新钱包流水记录输入
type WalletLogUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	WalletType int `json:"walletType"`
	ChangeType int `json:"changeType"`
	ChangeAmount int64 `json:"changeAmount"`
	BeforeBalance int64 `json:"beforeBalance"`
	AfterBalance int64 `json:"afterBalance"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WalletLogDetailOutput 钱包流水记录详情输出
type WalletLogDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	WalletType int `json:"walletType"`
	ChangeType int `json:"changeType"`
	ChangeAmount int64 `json:"changeAmount"`
	BeforeBalance int64 `json:"beforeBalance"`
	AfterBalance int64 `json:"afterBalance"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WalletLogListOutput 钱包流水记录列表输出
type WalletLogListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	WalletType int `json:"walletType"`
	ChangeType int `json:"changeType"`
	ChangeAmount int64 `json:"changeAmount"`
	BeforeBalance int64 `json:"beforeBalance"`
	AfterBalance int64 `json:"afterBalance"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WalletLogListInput 钱包流水记录列表查询输入
type WalletLogListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	RelatedOrderNo string `json:"relatedOrderNo"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	WalletType *int `json:"walletType"`
	ChangeType *int `json:"changeType"`
}
