package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Cron API

// CronCreateReq 创建定时任务表请求
type CronCreateReq struct {
	g.Meta `path:"/cron/create" method:"post" tags:"定时任务表" summary:"创建定时任务表"`
	Name string `json:"name" v:"required|max-length:100" dc:"任务名称"`
	Expression string `json:"expression" v:"required|max-length:50" dc:"Cron表达式"`
	Handler string `json:"handler" v:"required|max-length:100" dc:"处理器名称"`
	Params string `json:"params" v:"max-length:65535" dc:"参数（JSON）"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	LastRunAt *gtime.Time `json:"lastRunAt"  dc:"上次执行时间"`
	LastResult string `json:"lastResult" v:"max-length:255" dc:"上次执行结果"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CronCreateRes 创建定时任务表响应
type CronCreateRes struct {
	g.Meta `mime:"application/json"`
}

// CronUpdateReq 更新定时任务表请求
type CronUpdateReq struct {
	g.Meta `path:"/cron/update" method:"put" tags:"定时任务表" summary:"更新定时任务表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"定时任务表ID"`
	Name string `json:"name" v:"max-length:100" dc:"任务名称"`
	Expression string `json:"expression" v:"max-length:50" dc:"Cron表达式"`
	Handler string `json:"handler" v:"max-length:100" dc:"处理器名称"`
	Params string `json:"params" v:"max-length:65535" dc:"参数（JSON）"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	LastRunAt *gtime.Time `json:"lastRunAt"  dc:"上次执行时间"`
	LastResult string `json:"lastResult" v:"max-length:255" dc:"上次执行结果"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CronUpdateRes 更新定时任务表响应
type CronUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CronDeleteReq 删除定时任务表请求
type CronDeleteReq struct {
	g.Meta `path:"/cron/delete" method:"delete" tags:"定时任务表" summary:"删除定时任务表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"定时任务表ID"`
}

// CronDeleteRes 删除定时任务表响应
type CronDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CronBatchDeleteReq 批量删除定时任务表请求
type CronBatchDeleteReq struct {
	g.Meta `path:"/cron/batch-delete" method:"delete" tags:"定时任务表" summary:"批量删除定时任务表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"定时任务表ID列表"`
}

// CronBatchDeleteRes 批量删除定时任务表响应
type CronBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CronBatchUpdateReq 批量编辑定时任务表请求
type CronBatchUpdateReq struct {
	g.Meta `path:"/cron/batch-update" method:"put" tags:"定时任务表" summary:"批量编辑定时任务表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"定时任务表ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// CronBatchUpdateRes 批量编辑定时任务表响应
type CronBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CronDetailReq 获取定时任务表详情请求
type CronDetailReq struct {
	g.Meta `path:"/cron/detail" method:"get" tags:"定时任务表" summary:"获取定时任务表详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"定时任务表ID"`
}

// CronDetailRes 获取定时任务表详情响应
type CronDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.CronDetailOutput
}

// CronListReq 获取定时任务表列表请求
type CronListReq struct {
	g.Meta    `path:"/cron/list" method:"get" tags:"定时任务表" summary:"获取定时任务表列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Name string `json:"name" dc:"任务名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
	LastRunAtStart string `json:"lastRunAtStart" dc:"上次执行时间开始时间"`
	LastRunAtEnd string `json:"lastRunAtEnd" dc:"上次执行时间结束时间"`
}

// CronListRes 获取定时任务表列表响应
type CronListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CronListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// CronExportReq 导出定时任务表请求
type CronExportReq struct {
	g.Meta    `path:"/cron/export" method:"get" tags:"定时任务表" summary:"导出定时任务表"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Name string `json:"name" dc:"任务名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Remark string `json:"remark" dc:"备注"`
	LastRunAtStart string `json:"lastRunAtStart" dc:"上次执行时间开始时间"`
	LastRunAtEnd string `json:"lastRunAtEnd" dc:"上次执行时间结束时间"`
}

// CronExportRes 导出定时任务表响应
type CronExportRes struct {
	g.Meta `mime:"text/csv"`
}

// CronImportReq 导入定时任务表请求
type CronImportReq struct {
	g.Meta `path:"/cron/import" method:"post" mime:"multipart/form-data" tags:"定时任务表" summary:"导入定时任务表"`
}

// CronImportRes 导入定时任务表响应
type CronImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// CronImportTemplateReq 下载定时任务表导入模板
type CronImportTemplateReq struct {
	g.Meta `path:"/cron/import-template" method:"get" tags:"定时任务表" summary:"下载定时任务表导入模板"`
}

// CronImportTemplateRes 下载定时任务表导入模板响应
type CronImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

