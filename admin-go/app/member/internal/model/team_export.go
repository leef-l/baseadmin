package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// TeamExport DTO 模型

// TeamExportCreateInput 创建团队数据导出输入
type TeamExportCreateInput struct {
	UserID snowflake.JsonInt64 `json:"userID"`
	TeamMemberCount int `json:"teamMemberCount"`
	ExportType int `json:"exportType"`
	FileURL string `json:"fileURL"`
	FileSize int64 `json:"fileSize"`
	DeployStatus int `json:"deployStatus"`
	DeployDomain string `json:"deployDomain"`
	DeployedAt *gtime.Time `json:"deployedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// TeamExportUpdateInput 更新团队数据导出输入
type TeamExportUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	TeamMemberCount int `json:"teamMemberCount"`
	ExportType int `json:"exportType"`
	FileURL string `json:"fileURL"`
	FileSize int64 `json:"fileSize"`
	DeployStatus int `json:"deployStatus"`
	DeployDomain string `json:"deployDomain"`
	DeployedAt *gtime.Time `json:"deployedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// TeamExportDetailOutput 团队数据导出详情输出
type TeamExportDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	TeamMemberCount int `json:"teamMemberCount"`
	ExportType int `json:"exportType"`
	FileURL string `json:"fileURL"`
	FileSize int64 `json:"fileSize"`
	DeployStatus int `json:"deployStatus"`
	DeployDomain string `json:"deployDomain"`
	DeployedAt *gtime.Time `json:"deployedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// TeamExportListOutput 团队数据导出列表输出
type TeamExportListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	UserID snowflake.JsonInt64 `json:"userID"`
	UserNickname string `json:"userNickname"`
	TeamMemberCount int `json:"teamMemberCount"`
	ExportType int `json:"exportType"`
	FileURL string `json:"fileURL"`
	FileSize int64 `json:"fileSize"`
	DeployStatus int `json:"deployStatus"`
	DeployDomain string `json:"deployDomain"`
	DeployedAt *gtime.Time `json:"deployedAt"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// TeamExportListInput 团队数据导出列表查询输入
type TeamExportListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	UserID *snowflake.JsonInt64 `json:"userID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	ExportType *int `json:"exportType"`
	DeployStatus *int `json:"deployStatus"`
	Status *int `json:"status"`
	DeployDomain string `json:"deployDomain"`
	DeployedAtStart string `json:"deployedAtStart"`
	DeployedAtEnd string `json:"deployedAtEnd"`
}

// TeamExportBatchUpdateInput 批量编辑团队数据导出输入
type TeamExportBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	ExportType *int `json:"exportType"`
	DeployStatus *int `json:"deployStatus"`
	Status *int `json:"status"`
}

