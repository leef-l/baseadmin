package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// TeamExport API

// TeamExportCreateReq 创建团队数据导出请求
type TeamExportCreateReq struct {
	g.Meta `path:"/team_export/create" method:"post" tags:"团队数据导出" summary:"创建团队数据导出"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"目标会员"`
	TeamMemberCount int `json:"teamMemberCount"  dc:"团队成员数"`
	ExportType int `json:"exportType"  dc:"导出类型"`
	FileURL string `json:"fileURL" v:"url|max-length:500" dc:"导出文件地址"`
	FileSize int64 `json:"fileSize"  dc:"文件大小（字节）"`
	DeployStatus int `json:"deployStatus"  dc:"部署状态"`
	DeployDomain string `json:"deployDomain" v:"max-length:200" dc:"部署域名"`
	DeployedAt *gtime.Time `json:"deployedAt"  dc:"部署完成时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// TeamExportCreateRes 创建团队数据导出响应
type TeamExportCreateRes struct {
	g.Meta `mime:"application/json"`
}

// TeamExportUpdateReq 更新团队数据导出请求
type TeamExportUpdateReq struct {
	g.Meta `path:"/team_export/update" method:"put" tags:"团队数据导出" summary:"更新团队数据导出"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"团队数据导出ID"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"目标会员"`
	TeamMemberCount int `json:"teamMemberCount"  dc:"团队成员数"`
	ExportType int `json:"exportType"  dc:"导出类型"`
	FileURL string `json:"fileURL" v:"url|max-length:500" dc:"导出文件地址"`
	FileSize int64 `json:"fileSize"  dc:"文件大小（字节）"`
	DeployStatus int `json:"deployStatus"  dc:"部署状态"`
	DeployDomain string `json:"deployDomain" v:"max-length:200" dc:"部署域名"`
	DeployedAt *gtime.Time `json:"deployedAt"  dc:"部署完成时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// TeamExportUpdateRes 更新团队数据导出响应
type TeamExportUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// TeamExportDeleteReq 删除团队数据导出请求
type TeamExportDeleteReq struct {
	g.Meta `path:"/team_export/delete" method:"delete" tags:"团队数据导出" summary:"删除团队数据导出"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"团队数据导出ID"`
}

// TeamExportDeleteRes 删除团队数据导出响应
type TeamExportDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TeamExportBatchDeleteReq 批量删除团队数据导出请求
type TeamExportBatchDeleteReq struct {
	g.Meta `path:"/team_export/batch-delete" method:"delete" tags:"团队数据导出" summary:"批量删除团队数据导出"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"团队数据导出ID列表"`
}

// TeamExportBatchDeleteRes 批量删除团队数据导出响应
type TeamExportBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// TeamExportBatchUpdateReq 批量编辑团队数据导出请求
type TeamExportBatchUpdateReq struct {
	g.Meta `path:"/team_export/batch-update" method:"put" tags:"团队数据导出" summary:"批量编辑团队数据导出"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"团队数据导出ID列表"`
	ExportType *int `json:"exportType" dc:"导出类型"`
	DeployStatus *int `json:"deployStatus" dc:"部署状态"`
	Status *int `json:"status" dc:"状态"`
}

// TeamExportBatchUpdateRes 批量编辑团队数据导出响应
type TeamExportBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// TeamExportDetailReq 获取团队数据导出详情请求
type TeamExportDetailReq struct {
	g.Meta `path:"/team_export/detail" method:"get" tags:"团队数据导出" summary:"获取团队数据导出详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"团队数据导出ID"`
}

// TeamExportDetailRes 获取团队数据导出详情响应
type TeamExportDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.TeamExportDetailOutput
}

// TeamExportListReq 获取团队数据导出列表请求
type TeamExportListReq struct {
	g.Meta    `path:"/team_export/list" method:"get" tags:"团队数据导出" summary:"获取团队数据导出列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"目标会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ExportType *int `json:"exportType" dc:"导出类型"`
	DeployStatus *int `json:"deployStatus" dc:"部署状态"`
	Status *int `json:"status" dc:"状态"`
	DeployDomain string `json:"deployDomain" dc:"部署域名"`
	DeployedAtStart string `json:"deployedAtStart" dc:"部署完成时间开始时间"`
	DeployedAtEnd string `json:"deployedAtEnd" dc:"部署完成时间结束时间"`
}

// TeamExportListRes 获取团队数据导出列表响应
type TeamExportListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.TeamExportListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// TeamExportExportReq 导出团队数据导出请求
type TeamExportExportReq struct {
	g.Meta    `path:"/team_export/export" method:"get" tags:"团队数据导出" summary:"导出团队数据导出"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"目标会员"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ExportType *int `json:"exportType" dc:"导出类型"`
	DeployStatus *int `json:"deployStatus" dc:"部署状态"`
	Status *int `json:"status" dc:"状态"`
	DeployDomain string `json:"deployDomain" dc:"部署域名"`
	DeployedAtStart string `json:"deployedAtStart" dc:"部署完成时间开始时间"`
	DeployedAtEnd string `json:"deployedAtEnd" dc:"部署完成时间结束时间"`
}

// TeamExportExportRes 导出团队数据导出响应
type TeamExportExportRes struct {
	g.Meta `mime:"text/csv"`
}

// TeamExportImportReq 导入团队数据导出请求
type TeamExportImportReq struct {
	g.Meta `path:"/team_export/import" method:"post" mime:"multipart/form-data" tags:"团队数据导出" summary:"导入团队数据导出"`
}

// TeamExportImportRes 导入团队数据导出响应
type TeamExportImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// TeamExportImportTemplateReq 下载团队数据导出导入模板
type TeamExportImportTemplateReq struct {
	g.Meta `path:"/team_export/import-template" method:"get" tags:"团队数据导出" summary:"下载团队数据导出导入模板"`
}

// TeamExportImportTemplateRes 下载团队数据导出导入模板响应
type TeamExportImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

