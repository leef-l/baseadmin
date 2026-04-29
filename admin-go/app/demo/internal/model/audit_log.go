package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// AuditLog DTO 模型

// AuditLogCreateInput 创建体验审计日志输入
type AuditLogCreateInput struct {
	LogNo string `json:"logNo"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	Action int `json:"action"`
	TargetType int `json:"targetType"`
	TargetCode string `json:"targetCode"`
	RequestJSON string `json:"requestJSON"`
	Result int `json:"result"`
	ClientIP string `json:"clientIP"`
	OccurredAt *gtime.Time `json:"occurredAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// AuditLogUpdateInput 更新体验审计日志输入
type AuditLogUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	LogNo string `json:"logNo"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	Action int `json:"action"`
	TargetType int `json:"targetType"`
	TargetCode string `json:"targetCode"`
	RequestJSON string `json:"requestJSON"`
	Result int `json:"result"`
	ClientIP string `json:"clientIP"`
	OccurredAt *gtime.Time `json:"occurredAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// AuditLogDetailOutput 体验审计日志详情输出
type AuditLogDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	LogNo string `json:"logNo"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	UsersUsername string `json:"usersUsername"`
	Action int `json:"action"`
	TargetType int `json:"targetType"`
	TargetCode string `json:"targetCode"`
	RequestJSON string `json:"requestJSON"`
	Result int `json:"result"`
	ClientIP string `json:"clientIP"`
	OccurredAt *gtime.Time `json:"occurredAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// AuditLogListOutput 体验审计日志列表输出
type AuditLogListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	LogNo string `json:"logNo"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	UsersUsername string `json:"usersUsername"`
	Action int `json:"action"`
	TargetType int `json:"targetType"`
	TargetCode string `json:"targetCode"`
	RequestJSON string `json:"requestJSON"`
	Result int `json:"result"`
	ClientIP string `json:"clientIP"`
	OccurredAt *gtime.Time `json:"occurredAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// AuditLogListInput 体验审计日志列表查询输入
type AuditLogListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	LogNo string `json:"logNo"`
	TargetCode string `json:"targetCode"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	ClientIP string `json:"clientIP"`
	Action *int `json:"action"`
	TargetType *int `json:"targetType"`
	Result *int `json:"result"`
	OccurredAtStart string `json:"occurredAtStart"`
	OccurredAtEnd string `json:"occurredAtEnd"`
}
