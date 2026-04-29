package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Category API

// CategoryCreateReq 创建体验分类请求
type CategoryCreateReq struct {
	g.Meta `path:"/category/create" method:"post" tags:"体验分类" summary:"创建体验分类"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"父分类"`
	Name string `json:"name" v:"required|max-length:80" dc:"分类名称"`
	Icon string `json:"icon" v:"max-length:100" dc:"图标"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CategoryCreateRes 创建体验分类响应
type CategoryCreateRes struct {
	g.Meta `mime:"application/json"`
}

// CategoryUpdateReq 更新体验分类请求
type CategoryUpdateReq struct {
	g.Meta `path:"/category/update" method:"put" tags:"体验分类" summary:"更新体验分类"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验分类ID"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"父分类"`
	Name string `json:"name" v:"max-length:80" dc:"分类名称"`
	Icon string `json:"icon" v:"max-length:100" dc:"图标"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// CategoryUpdateRes 更新体验分类响应
type CategoryUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CategoryDeleteReq 删除体验分类请求
type CategoryDeleteReq struct {
	g.Meta `path:"/category/delete" method:"delete" tags:"体验分类" summary:"删除体验分类"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验分类ID"`
}

// CategoryDeleteRes 删除体验分类响应
type CategoryDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CategoryBatchDeleteReq 批量删除体验分类请求
type CategoryBatchDeleteReq struct {
	g.Meta `path:"/category/batch-delete" method:"delete" tags:"体验分类" summary:"批量删除体验分类"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验分类ID列表"`
}

// CategoryBatchDeleteRes 批量删除体验分类响应
type CategoryBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// CategoryBatchUpdateReq 批量编辑体验分类请求
type CategoryBatchUpdateReq struct {
	g.Meta `path:"/category/batch-update" method:"put" tags:"体验分类" summary:"批量编辑体验分类"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验分类ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// CategoryBatchUpdateRes 批量编辑体验分类响应
type CategoryBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// CategoryDetailReq 获取体验分类详情请求
type CategoryDetailReq struct {
	g.Meta `path:"/category/detail" method:"get" tags:"体验分类" summary:"获取体验分类详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验分类ID"`
}

// CategoryDetailRes 获取体验分类详情响应
type CategoryDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.CategoryDetailOutput
}

// CategoryListReq 获取体验分类列表请求
type CategoryListReq struct {
	g.Meta    `path:"/category/list" method:"get" tags:"体验分类" summary:"获取体验分类列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"父分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// CategoryListRes 获取体验分类列表响应
type CategoryListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CategoryListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// CategoryExportReq 导出体验分类请求
type CategoryExportReq struct {
	g.Meta    `path:"/category/export" method:"get" tags:"体验分类" summary:"导出体验分类"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"父分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// CategoryExportRes 导出体验分类响应
type CategoryExportRes struct {
	g.Meta `mime:"text/csv"`
}

// CategoryTreeReq 获取体验分类树形结构请求
type CategoryTreeReq struct {
	g.Meta    `path:"/category/tree" method:"get" tags:"体验分类" summary:"获取体验分类树形结构"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"父分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// CategoryTreeRes 获取体验分类树形结构响应
type CategoryTreeRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.CategoryTreeOutput `json:"list" dc:"树形数据"`
}
