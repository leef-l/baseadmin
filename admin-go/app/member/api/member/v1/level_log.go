package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// LevelLog API

// LevelLogCreateReq 创建等级变更日志请求
type LevelLogCreateReq struct {
	g.Meta `path:"/level_log/create" method:"post" tags:"等级变更日志" summary:"创建等级变更日志"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"  dc:"变更前等级"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"  dc:"变更后等级"`
	ChangeType int `json:"changeType"  dc:"变更类型"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"新等级到期时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"变更说明"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// LevelLogCreateRes 创建等级变更日志响应
type LevelLogCreateRes struct {
	g.Meta `mime:"application/json"`
}

// LevelLogUpdateReq 更新等级变更日志请求
type LevelLogUpdateReq struct {
	g.Meta `path:"/level_log/update" method:"put" tags:"等级变更日志" summary:"更新等级变更日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"等级变更日志ID"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"会员"`
	OldLevelID snowflake.JsonInt64 `json:"oldLevelID"  dc:"变更前等级"`
	NewLevelID snowflake.JsonInt64 `json:"newLevelID"  dc:"变更后等级"`
	ChangeType int `json:"changeType"  dc:"变更类型"`
	ExpireAt *gtime.Time `json:"expireAt"  dc:"新等级到期时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"变更说明"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// LevelLogUpdateRes 更新等级变更日志响应
type LevelLogUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// LevelLogDeleteReq 删除等级变更日志请求
type LevelLogDeleteReq struct {
	g.Meta `path:"/level_log/delete" method:"delete" tags:"等级变更日志" summary:"删除等级变更日志"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"等级变更日志ID"`
}

// LevelLogDeleteRes 删除等级变更日志响应
type LevelLogDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// LevelLogBatchDeleteReq 批量删除等级变更日志请求
type LevelLogBatchDeleteReq struct {
	g.Meta `path:"/level_log/batch-delete" method:"delete" tags:"等级变更日志" summary:"批量删除等级变更日志"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"等级变更日志ID列表"`
}

// LevelLogBatchDeleteRes 批量删除等级变更日志响应
type LevelLogBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// LevelLogDetailReq 获取等级变更日志详情请求
type LevelLogDetailReq struct {
	g.Meta `path:"/level_log/detail" method:"get" tags:"等级变更日志" summary:"获取等级变更日志详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"等级变更日志ID"`
}

// LevelLogDetailRes 获取等级变更日志详情响应
type LevelLogDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.LevelLogDetailOutput
}

// LevelLogListReq 获取等级变更日志列表请求
type LevelLogListReq struct {
	g.Meta    `path:"/level_log/list" method:"get" tags:"等级变更日志" summary:"获取等级变更日志列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	OldLevelID *snowflake.JsonInt64 `json:"oldLevelID" dc:"变更前等级"`
	NewLevelID *snowflake.JsonInt64 `json:"newLevelID" dc:"变更后等级"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ChangeType *int `json:"changeType" dc:"变更类型"`
	ExpireAtStart string `json:"expireAtStart" dc:"新等级到期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"新等级到期时间结束时间"`
}

// LevelLogListRes 获取等级变更日志列表响应
type LevelLogListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.LevelLogListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// LevelLogExportReq 导出等级变更日志请求
type LevelLogExportReq struct {
	g.Meta    `path:"/level_log/export" method:"get" tags:"等级变更日志" summary:"导出等级变更日志"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"会员"`
	OldLevelID *snowflake.JsonInt64 `json:"oldLevelID" dc:"变更前等级"`
	NewLevelID *snowflake.JsonInt64 `json:"newLevelID" dc:"变更后等级"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ChangeType *int `json:"changeType" dc:"变更类型"`
	ExpireAtStart string `json:"expireAtStart" dc:"新等级到期时间开始时间"`
	ExpireAtEnd string `json:"expireAtEnd" dc:"新等级到期时间结束时间"`
}

// LevelLogExportRes 导出等级变更日志响应
type LevelLogExportRes struct {
	g.Meta `mime:"text/csv"`
}

// LevelLogImportReq 导入等级变更日志请求
type LevelLogImportReq struct {
	g.Meta `path:"/level_log/import" method:"post" mime:"multipart/form-data" tags:"等级变更日志" summary:"导入等级变更日志"`
}

// LevelLogImportRes 导入等级变更日志响应
type LevelLogImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// LevelLogImportTemplateReq 下载等级变更日志导入模板
type LevelLogImportTemplateReq struct {
	g.Meta `path:"/level_log/import-template" method:"get" tags:"等级变更日志" summary:"下载等级变更日志导入模板"`
}

// LevelLogImportTemplateRes 下载等级变更日志导入模板响应
type LevelLogImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

