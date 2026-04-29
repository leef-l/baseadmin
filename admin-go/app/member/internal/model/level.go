package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Level DTO 模型

// LevelCreateInput 创建会员等级配置输入
type LevelCreateInput struct {
	Name string `json:"name"`
	LevelNo int `json:"levelNo"`
	Icon string `json:"icon"`
	DurationDays int `json:"durationDays"`
	NeedActiveCount int `json:"needActiveCount"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"`
	IsTop int `json:"isTop"`
	AutoDeploy int `json:"autoDeploy"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// LevelUpdateInput 更新会员等级配置输入
type LevelUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	LevelNo int `json:"levelNo"`
	Icon string `json:"icon"`
	DurationDays int `json:"durationDays"`
	NeedActiveCount int `json:"needActiveCount"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"`
	IsTop int `json:"isTop"`
	AutoDeploy int `json:"autoDeploy"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// LevelDetailOutput 会员等级配置详情输出
type LevelDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	LevelNo int `json:"levelNo"`
	Icon string `json:"icon"`
	DurationDays int `json:"durationDays"`
	NeedActiveCount int `json:"needActiveCount"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"`
	IsTop int `json:"isTop"`
	AutoDeploy int `json:"autoDeploy"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// LevelListOutput 会员等级配置列表输出
type LevelListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	Name string `json:"name"`
	LevelNo int `json:"levelNo"`
	Icon string `json:"icon"`
	DurationDays int `json:"durationDays"`
	NeedActiveCount int `json:"needActiveCount"`
	NeedTeamTurnover int64 `json:"needTeamTurnover"`
	IsTop int `json:"isTop"`
	AutoDeploy int `json:"autoDeploy"`
	Remark string `json:"remark"`
	Sort int `json:"sort"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// LevelListInput 会员等级配置列表查询输入
type LevelListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Name string `json:"name"`
	LevelNo string `json:"levelNo"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	IsTop *int `json:"isTop"`
	AutoDeploy *int `json:"autoDeploy"`
	Status *int `json:"status"`
}

// LevelBatchUpdateInput 批量编辑会员等级配置输入
type LevelBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	IsTop *int `json:"isTop"`
	AutoDeploy *int `json:"autoDeploy"`
	Status *int `json:"status"`
}

