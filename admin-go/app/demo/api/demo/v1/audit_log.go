package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// AuditLog API

// AuditLogCreateReq 创建体验审计日志请求
type AuditLogCreateReq struct {
	g.Meta `path:"/audit_log/create" method:"post" tags:"体验审计日志" summary:"创建体验审计日志"`
	LogNo string `json:"logNo" v:"required|max-length:50" dc:"日志编号"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"  dc:"操作人"`
	Action int `json:"action"  dc:"动作"`
	TargetType int `json:"targetType"  dc:"对象类型"`
	TargetCode string `json:"targetCode" v:"max-length:80" dc:"对象编号"`
	RequestJSON string `json:"requestJSON" v:"max-length:65535" dc:"请求JSON"`
	Result int `json:"result"  dc:"结果"`
	ClientIP string `json:"clientIP" v:"max-length:50" dc:"客户端IP"`
	OccurredAt *gtime.Time `json:"occurredAt"  dc:"发生时间"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// AuditLogCreateRes 创建体验审计日志响应
type AuditLogCreateRes struct {
	g.Meta `mime:"application/json"`
}

// AuditLogUpdateReq 更新体验审计日志请求
type AuditLogUpdateReq struct {
	g.Meta `path:"/audit_log/update" method:"put" tags:"体验审计日志" summary:"更新体验审计日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验审计日志ID"`
	LogNo string `json:"logNo" v:"max-length:50" dc:"日志编号"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"  dc:"操作人"`
	Action int `json:"action"  dc:"动作"`
	TargetType int `json:"targetType"  dc:"对象类型"`
	TargetCode string `json:"targetCode" v:"max-length:80" dc:"对象编号"`
	RequestJSON string `json:"requestJSON" v:"max-length:65535" dc:"请求JSON"`
	Result int `json:"result"  dc:"结果"`
	ClientIP string `json:"clientIP" v:"max-length:50" dc:"客户端IP"`
	OccurredAt *gtime.Time `json:"occurredAt"  dc:"发生时间"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// AuditLogUpdateRes 更新体验审计日志响应
type AuditLogUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// AuditLogDeleteReq 删除体验审计日志请求
type AuditLogDeleteReq struct {
	g.Meta `path:"/audit_log/delete" method:"delete" tags:"体验审计日志" summary:"删除体验审计日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验审计日志ID"`
}

// AuditLogDeleteRes 删除体验审计日志响应
type AuditLogDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// AuditLogBatchDeleteReq 批量删除体验审计日志请求
type AuditLogBatchDeleteReq struct {
	g.Meta `path:"/audit_log/batch-delete" method:"delete" tags:"体验审计日志" summary:"批量删除体验审计日志"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验审计日志ID列表"`
}

// AuditLogBatchDeleteRes 批量删除体验审计日志响应
type AuditLogBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// AuditLogDetailReq 获取体验审计日志详情请求
type AuditLogDetailReq struct {
	g.Meta `path:"/audit_log/detail" method:"get" tags:"体验审计日志" summary:"获取体验审计日志详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验审计日志ID"`
}

// AuditLogDetailRes 获取体验审计日志详情响应
type AuditLogDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.AuditLogDetailOutput
}

// AuditLogListReq 获取体验审计日志列表请求
type AuditLogListReq struct {
	g.Meta    `path:"/audit_log/list" method:"get" tags:"体验审计日志" summary:"获取体验审计日志列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	LogNo string `json:"logNo" dc:"日志编号"`
	TargetCode string `json:"targetCode" dc:"对象编号"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID" dc:"操作人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ClientIP string `json:"clientIP" dc:"客户端IP"`
	Action *int `json:"action" dc:"动作"`
	TargetType *int `json:"targetType" dc:"对象类型"`
	Result *int `json:"result" dc:"结果"`
	OccurredAtStart string `json:"occurredAtStart" dc:"发生时间开始时间"`
	OccurredAtEnd string `json:"occurredAtEnd" dc:"发生时间结束时间"`
}

// AuditLogListRes 获取体验审计日志列表响应
type AuditLogListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.AuditLogListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// AuditLogExportReq 导出体验审计日志请求
type AuditLogExportReq struct {
	g.Meta    `path:"/audit_log/export" method:"get" tags:"体验审计日志" summary:"导出体验审计日志"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	LogNo string `json:"logNo" dc:"日志编号"`
	TargetCode string `json:"targetCode" dc:"对象编号"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID" dc:"操作人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ClientIP string `json:"clientIP" dc:"客户端IP"`
	Action *int `json:"action" dc:"动作"`
	TargetType *int `json:"targetType" dc:"对象类型"`
	Result *int `json:"result" dc:"结果"`
	OccurredAtStart string `json:"occurredAtStart" dc:"发生时间开始时间"`
	OccurredAtEnd string `json:"occurredAtEnd" dc:"发生时间结束时间"`
}

// AuditLogExportRes 导出体验审计日志响应
type AuditLogExportRes struct {
	g.Meta `mime:"text/csv"`
}

// AuditLogImportReq 导入体验审计日志请求
type AuditLogImportReq struct {
	g.Meta `path:"/audit_log/import" method:"post" mime:"multipart/form-data" tags:"体验审计日志" summary:"导入体验审计日志"`
}

// AuditLogImportRes 导入体验审计日志响应
type AuditLogImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// AuditLogImportTemplateReq 下载体验审计日志导入模板
type AuditLogImportTemplateReq struct {
	g.Meta `path:"/audit_log/import-template" method:"get" tags:"体验审计日志" summary:"下载体验审计日志导入模板"`
}

// AuditLogImportTemplateRes 下载体验审计日志导入模板响应
type AuditLogImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

