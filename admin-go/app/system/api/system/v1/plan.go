package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Plan API

// PlanCreateReq 创建套餐表请求
type PlanCreateReq struct {
	g.Meta `path:"/plan/create" method:"post" tags:"套餐表" summary:"创建套餐表"`
	Name string `json:"name" v:"required|max-length:80" dc:"套餐名称"`
	Code string `json:"code" v:"required|max-length:50" dc:"套餐编码"`
	Description string `json:"description" v:"max-length:500" dc:"套餐描述"`
	UserLimit int `json:"userLimit"  dc:"用户数量上限（-1=无限）"`
	StorageMb int `json:"storageMb"  dc:"存储空间上限MB（-1=无限）"`
	Features string `json:"features" v:"max-length:65535" dc:"功能开关列表（JSON数组）"`
	Price int `json:"price"  dc:"价格（分）"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// PlanCreateRes 创建套餐表响应
type PlanCreateRes struct {
	g.Meta `mime:"application/json"`
}

// PlanUpdateReq 更新套餐表请求
type PlanUpdateReq struct {
	g.Meta `path:"/plan/update" method:"put" tags:"套餐表" summary:"更新套餐表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"套餐表ID"`
	Name string `json:"name" v:"max-length:80" dc:"套餐名称"`
	Code string `json:"code" v:"max-length:50" dc:"套餐编码"`
	Description string `json:"description" v:"max-length:500" dc:"套餐描述"`
	UserLimit int `json:"userLimit"  dc:"用户数量上限（-1=无限）"`
	StorageMb int `json:"storageMb"  dc:"存储空间上限MB（-1=无限）"`
	Features string `json:"features" v:"max-length:65535" dc:"功能开关列表（JSON数组）"`
	Price int `json:"price"  dc:"价格（分）"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// PlanUpdateRes 更新套餐表响应
type PlanUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// PlanDeleteReq 删除套餐表请求
type PlanDeleteReq struct {
	g.Meta `path:"/plan/delete" method:"delete" tags:"套餐表" summary:"删除套餐表"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"套餐表ID"`
}

// PlanDeleteRes 删除套餐表响应
type PlanDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// PlanBatchDeleteReq 批量删除套餐表请求
type PlanBatchDeleteReq struct {
	g.Meta `path:"/plan/batch-delete" method:"delete" tags:"套餐表" summary:"批量删除套餐表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"套餐表ID列表"`
}

// PlanBatchDeleteRes 批量删除套餐表响应
type PlanBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// PlanBatchUpdateReq 批量编辑套餐表请求
type PlanBatchUpdateReq struct {
	g.Meta `path:"/plan/batch-update" method:"put" tags:"套餐表" summary:"批量编辑套餐表"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"套餐表ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// PlanBatchUpdateRes 批量编辑套餐表响应
type PlanBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// PlanDetailReq 获取套餐表详情请求
type PlanDetailReq struct {
	g.Meta `path:"/plan/detail" method:"get" tags:"套餐表" summary:"获取套餐表详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"套餐表ID"`
}

// PlanDetailRes 获取套餐表详情响应
type PlanDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.PlanDetailOutput
}

// PlanListReq 获取套餐表列表请求
type PlanListReq struct {
	g.Meta    `path:"/plan/list" method:"get" tags:"套餐表" summary:"获取套餐表列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Code string `json:"code" dc:"套餐编码"`
	Name string `json:"name" dc:"套餐名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Description string `json:"description" dc:"套餐描述"`
}

// PlanListRes 获取套餐表列表响应
type PlanListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.PlanListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// PlanExportReq 导出套餐表请求
type PlanExportReq struct {
	g.Meta    `path:"/plan/export" method:"get" tags:"套餐表" summary:"导出套餐表"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	Code string `json:"code" dc:"套餐编码"`
	Name string `json:"name" dc:"套餐名称"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
	Description string `json:"description" dc:"套餐描述"`
}

// PlanExportRes 导出套餐表响应
type PlanExportRes struct {
	g.Meta `mime:"text/csv"`
}

// PlanImportReq 导入套餐表请求
type PlanImportReq struct {
	g.Meta `path:"/plan/import" method:"post" mime:"multipart/form-data" tags:"套餐表" summary:"导入套餐表"`
}

// PlanImportRes 导入套餐表响应
type PlanImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// PlanImportTemplateReq 下载套餐表导入模板
type PlanImportTemplateReq struct {
	g.Meta `path:"/plan/import-template" method:"get" tags:"套餐表" summary:"下载套餐表导入模板"`
}

// PlanImportTemplateRes 下载套餐表导入模板响应
type PlanImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

