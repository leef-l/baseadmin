package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// RebindLog API

// RebindLogCreateReq 创建换绑上级日志请求
type RebindLogCreateReq struct {
	g.Meta `path:"/rebind_log/create" method:"post" tags:"换绑上级日志" summary:"创建换绑上级日志"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"  dc:"原上级"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"  dc:"新上级"`
	Reason string `json:"reason" v:"max-length:500" dc:"换绑原因"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"  dc:"操作人"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// RebindLogCreateRes 创建换绑上级日志响应
type RebindLogCreateRes struct {
	g.Meta `mime:"application/json"`
}

// RebindLogUpdateReq 更新换绑上级日志请求
type RebindLogUpdateReq struct {
	g.Meta `path:"/rebind_log/update" method:"put" tags:"换绑上级日志" summary:"更新换绑上级日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"换绑上级日志ID"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	OldParentID snowflake.JsonInt64 `json:"oldParentID"  dc:"原上级"`
	NewParentID snowflake.JsonInt64 `json:"newParentID"  dc:"新上级"`
	Reason string `json:"reason" v:"max-length:500" dc:"换绑原因"`
	OperatorID snowflake.JsonInt64 `json:"operatorID"  dc:"操作人"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// RebindLogUpdateRes 更新换绑上级日志响应
type RebindLogUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// RebindLogDeleteReq 删除换绑上级日志请求
type RebindLogDeleteReq struct {
	g.Meta `path:"/rebind_log/delete" method:"delete" tags:"换绑上级日志" summary:"删除换绑上级日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"换绑上级日志ID"`
}

// RebindLogDeleteRes 删除换绑上级日志响应
type RebindLogDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// RebindLogBatchDeleteReq 批量删除换绑上级日志请求
type RebindLogBatchDeleteReq struct {
	g.Meta `path:"/rebind_log/batch-delete" method:"delete" tags:"换绑上级日志" summary:"批量删除换绑上级日志"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"换绑上级日志ID列表"`
}

// RebindLogBatchDeleteRes 批量删除换绑上级日志响应
type RebindLogBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// RebindLogDetailReq 获取换绑上级日志详情请求
type RebindLogDetailReq struct {
	g.Meta `path:"/rebind_log/detail" method:"get" tags:"换绑上级日志" summary:"获取换绑上级日志详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"换绑上级日志ID"`
}

// RebindLogDetailRes 获取换绑上级日志详情响应
type RebindLogDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.RebindLogDetailOutput
}

// RebindLogListReq 获取换绑上级日志列表请求
type RebindLogListReq struct {
	g.Meta    `path:"/rebind_log/list" method:"get" tags:"换绑上级日志" summary:"获取换绑上级日志列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	OldParentID *snowflake.JsonInt64 `json:"oldParentID" dc:"原上级"`
	NewParentID *snowflake.JsonInt64 `json:"newParentID" dc:"新上级"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID" dc:"操作人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
}

// RebindLogListRes 获取换绑上级日志列表响应
type RebindLogListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.RebindLogListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// RebindLogExportReq 导出换绑上级日志请求
type RebindLogExportReq struct {
	g.Meta    `path:"/rebind_log/export" method:"get" tags:"换绑上级日志" summary:"导出换绑上级日志"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	OldParentID *snowflake.JsonInt64 `json:"oldParentID" dc:"原上级"`
	NewParentID *snowflake.JsonInt64 `json:"newParentID" dc:"新上级"`
	OperatorID *snowflake.JsonInt64 `json:"operatorID" dc:"操作人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
}

// RebindLogExportRes 导出换绑上级日志响应
type RebindLogExportRes struct {
	g.Meta `mime:"text/csv"`
}

// RebindLogImportReq 导入换绑上级日志请求
type RebindLogImportReq struct {
	g.Meta `path:"/rebind_log/import" method:"post" mime:"multipart/form-data" tags:"换绑上级日志" summary:"导入换绑上级日志"`
}

// RebindLogImportRes 导入换绑上级日志响应
type RebindLogImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// RebindLogImportTemplateReq 下载换绑上级日志导入模板
type RebindLogImportTemplateReq struct {
	g.Meta `path:"/rebind_log/import-template" method:"get" tags:"换绑上级日志" summary:"下载换绑上级日志导入模板"`
}

// RebindLogImportTemplateRes 下载换绑上级日志导入模板响应
type RebindLogImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

