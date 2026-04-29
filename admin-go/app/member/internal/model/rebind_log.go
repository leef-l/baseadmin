package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// RebindLog DTO 模型

// RebindLogCreateInput 创建换绑上级日志输入
type RebindLogCreateInput struct {
	UserID snowflake.JsonInt64 `json:"userID"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"`
	Reason string `json:"reason"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// RebindLogUpdateInput 更新换绑上级日志输入
type RebindLogUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"`
	Reason string `json:"reason"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// RebindLogDetailOutput 换绑上级日志详情输出
type RebindLogDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"`
	OldParentNickname string `json:"oldParentNickname"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"`
	NewParentNickname string `json:"newParentNickname"`
	Reason string `json:"reason"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	UsersUsername string `json:"usersUsername"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// RebindLogListOutput 换绑上级日志列表输出
type RebindLogListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"`
	OldParentNickname string `json:"oldParentNickname"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"`
	NewParentNickname string `json:"newParentNickname"`
	Reason string `json:"reason"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"`
	UsersUsername string `json:"usersUsername"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// RebindLogListInput 换绑上级日志列表查询输入
type RebindLogListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	OldParentID *snowflake.JsonInt64 `json:"oldParentID"`
	NewParentID *snowflake.JsonInt64 `json:"newParentID"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
}
