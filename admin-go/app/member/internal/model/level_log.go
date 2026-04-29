package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// LevelLog DTO 模型

// LevelLogCreateInput 创建等级变更日志输入
type LevelLogCreateInput struct {
	UserID snowflake.JsonInt64 `json:"userID"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"`
	ChangeType int `json:"changeType"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// LevelLogUpdateInput 更新等级变更日志输入
type LevelLogUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"`
	ChangeType int `json:"changeType"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// LevelLogDetailOutput 等级变更日志详情输出
type LevelLogDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"`
	LevelName string `json:"levelName"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"`
	NewLevelName string `json:"newLevelName"`
	ChangeType int `json:"changeType"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// LevelLogListOutput 等级变更日志列表输出
type LevelLogListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"`
	LevelName string `json:"levelName"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"`
	NewLevelName string `json:"newLevelName"`
	ChangeType int `json:"changeType"`
	ExpireAt *gtime.Time `json:"expireAt"`
	Remark string `json:"remark"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// LevelLogListInput 等级变更日志列表查询输入
type LevelLogListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	OldLevelID *snowflake.JsonInt64 `json:"oldLevelID"`
	NewLevelID *snowflake.JsonInt64 `json:"newLevelID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	ChangeType *int `json:"changeType"`
	ExpireAtStart string `json:"expireAtStart"`
	ExpireAtEnd string `json:"expireAtEnd"`
}
