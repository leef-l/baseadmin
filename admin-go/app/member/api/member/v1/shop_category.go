package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// ShopCategory API

// ShopCategoryCreateReq 创建商城商品分类请求
type ShopCategoryCreateReq struct {
	g.Meta `path:"/shop_category/create" method:"post" tags:"商城商品分类" summary:"创建商城商品分类"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"上级分类"`
	Name string `json:"name" v:"required|max-length:50" dc:"分类名称"`
	Icon string `json:"icon" v:"max-length:500" dc:"分类图标"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopCategoryCreateRes 创建商城商品分类响应
type ShopCategoryCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopCategoryUpdateReq 更新商城商品分类请求
type ShopCategoryUpdateReq struct {
	g.Meta `path:"/shop_category/update" method:"put" tags:"商城商品分类" summary:"更新商城商品分类"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品分类ID"`
	ParentID snowflake.JsonInt64 `json:"parentID"  dc:"上级分类"`
	Name string `json:"name" v:"max-length:50" dc:"分类名称"`
	Icon string `json:"icon" v:"max-length:500" dc:"分类图标"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopCategoryUpdateRes 更新商城商品分类响应
type ShopCategoryUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopCategoryDeleteReq 删除商城商品分类请求
type ShopCategoryDeleteReq struct {
	g.Meta `path:"/shop_category/delete" method:"delete" tags:"商城商品分类" summary:"删除商城商品分类"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品分类ID"`
}

// ShopCategoryDeleteRes 删除商城商品分类响应
type ShopCategoryDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopCategoryBatchDeleteReq 批量删除商城商品分类请求
type ShopCategoryBatchDeleteReq struct {
	g.Meta `path:"/shop_category/batch-delete" method:"delete" tags:"商城商品分类" summary:"批量删除商城商品分类"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城商品分类ID列表"`
}

// ShopCategoryBatchDeleteRes 批量删除商城商品分类响应
type ShopCategoryBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopCategoryBatchUpdateReq 批量编辑商城商品分类请求
type ShopCategoryBatchUpdateReq struct {
	g.Meta `path:"/shop_category/batch-update" method:"put" tags:"商城商品分类" summary:"批量编辑商城商品分类"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城商品分类ID列表"`
	Status *int `json:"status" dc:"状态"`
}

// ShopCategoryBatchUpdateRes 批量编辑商城商品分类响应
type ShopCategoryBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopCategoryDetailReq 获取商城商品分类详情请求
type ShopCategoryDetailReq struct {
	g.Meta `path:"/shop_category/detail" method:"get" tags:"商城商品分类" summary:"获取商城商品分类详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品分类ID"`
}

// ShopCategoryDetailRes 获取商城商品分类详情响应
type ShopCategoryDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ShopCategoryDetailOutput
}

// ShopCategoryListReq 获取商城商品分类列表请求
type ShopCategoryListReq struct {
	g.Meta    `path:"/shop_category/list" method:"get" tags:"商城商品分类" summary:"获取商城商品分类列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// ShopCategoryListRes 获取商城商品分类列表响应
type ShopCategoryListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ShopCategoryListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ShopCategoryExportReq 导出商城商品分类请求
type ShopCategoryExportReq struct {
	g.Meta    `path:"/shop_category/export" method:"get" tags:"商城商品分类" summary:"导出商城商品分类"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// ShopCategoryExportRes 导出商城商品分类响应
type ShopCategoryExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ShopCategoryTreeReq 获取商城商品分类树形结构请求
type ShopCategoryTreeReq struct {
	g.Meta    `path:"/shop_category/tree" method:"get" tags:"商城商品分类" summary:"获取商城商品分类树形结构"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Name string `json:"name" dc:"分类名称"`
	ParentID *snowflake.JsonInt64 `json:"parentID" dc:"上级分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Status *int `json:"status" dc:"状态"`
}

// ShopCategoryTreeRes 获取商城商品分类树形结构响应
type ShopCategoryTreeRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ShopCategoryTreeOutput `json:"list" dc:"树形数据"`
}
