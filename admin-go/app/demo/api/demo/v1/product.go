package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Product API

// ProductCreateReq 创建体验商品请求
type ProductCreateReq struct {
	g.Meta `path:"/product/create" method:"post" tags:"体验商品" summary:"创建体验商品"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"  dc:"商品分类"`
	SkuNo string `json:"skuNo" v:"required|max-length:50" dc:"SKU编号"`
	Name string `json:"name" v:"required|max-length:120" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"封面"`
	ManualFile string `json:"manualFile" v:"max-length:500" dc:"说明书文件"`
	DetailContent string `json:"detailContent" v:"max-length:65535" dc:"详情内容"`
	SpecJSON string `json:"specJSON" v:"max-length:65535" dc:"规格JSON"`
	WebsiteURL string `json:"websiteURL" v:"url|max-length:500" dc:"官网URL"`
	Type int `json:"type"  dc:"类型"`
	IsRecommend int `json:"isRecommend"  dc:"是否推荐"`
	SalePrice int `json:"salePrice"  dc:"销售价（分）"`
	StockNum int `json:"stockNum"  dc:"库存数量"`
	WeightNum int `json:"weightNum"  dc:"重量（克）"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Icon string `json:"icon" v:"max-length:100" dc:"图标"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ProductCreateRes 创建体验商品响应
type ProductCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ProductUpdateReq 更新体验商品请求
type ProductUpdateReq struct {
	g.Meta `path:"/product/update" method:"put" tags:"体验商品" summary:"更新体验商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验商品ID"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"  dc:"商品分类"`
	SkuNo string `json:"skuNo" v:"max-length:50" dc:"SKU编号"`
	Name string `json:"name" v:"max-length:120" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"封面"`
	ManualFile string `json:"manualFile" v:"max-length:500" dc:"说明书文件"`
	DetailContent string `json:"detailContent" v:"max-length:65535" dc:"详情内容"`
	SpecJSON string `json:"specJSON" v:"max-length:65535" dc:"规格JSON"`
	WebsiteURL string `json:"websiteURL" v:"url|max-length:500" dc:"官网URL"`
	Type int `json:"type"  dc:"类型"`
	IsRecommend int `json:"isRecommend"  dc:"是否推荐"`
	SalePrice int `json:"salePrice"  dc:"销售价（分）"`
	StockNum int `json:"stockNum"  dc:"库存数量"`
	WeightNum int `json:"weightNum"  dc:"重量（克）"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Icon string `json:"icon" v:"max-length:100" dc:"图标"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ProductUpdateRes 更新体验商品响应
type ProductUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ProductDeleteReq 删除体验商品请求
type ProductDeleteReq struct {
	g.Meta `path:"/product/delete" method:"delete" tags:"体验商品" summary:"删除体验商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验商品ID"`
}

// ProductDeleteRes 删除体验商品响应
type ProductDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ProductBatchDeleteReq 批量删除体验商品请求
type ProductBatchDeleteReq struct {
	g.Meta `path:"/product/batch-delete" method:"delete" tags:"体验商品" summary:"批量删除体验商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验商品ID列表"`
}

// ProductBatchDeleteRes 批量删除体验商品响应
type ProductBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ProductBatchUpdateReq 批量编辑体验商品请求
type ProductBatchUpdateReq struct {
	g.Meta `path:"/product/batch-update" method:"put" tags:"体验商品" summary:"批量编辑体验商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验商品ID列表"`
	Type *int `json:"type" dc:"类型"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ProductBatchUpdateRes 批量编辑体验商品响应
type ProductBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ProductDetailReq 获取体验商品详情请求
type ProductDetailReq struct {
	g.Meta `path:"/product/detail" method:"get" tags:"体验商品" summary:"获取体验商品详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验商品ID"`
}

// ProductDetailRes 获取体验商品详情响应
type ProductDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ProductDetailOutput
}

// ProductListReq 获取体验商品列表请求
type ProductListReq struct {
	g.Meta    `path:"/product/list" method:"get" tags:"体验商品" summary:"获取体验商品列表"`
	PageNum   int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	SkuNo string `json:"skuNo" dc:"SKU编号"`
	Name string `json:"name" dc:"商品名称"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID" dc:"商品分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Type *int `json:"type" dc:"类型"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ProductListRes 获取体验商品列表响应
type ProductListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ProductListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ProductExportReq 导出体验商品请求
type ProductExportReq struct {
	g.Meta    `path:"/product/export" method:"get" tags:"体验商品" summary:"导出体验商品"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	SkuNo string `json:"skuNo" dc:"SKU编号"`
	Name string `json:"name" dc:"商品名称"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID" dc:"商品分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Type *int `json:"type" dc:"类型"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ProductExportRes 导出体验商品响应
type ProductExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ProductImportReq 导入体验商品请求
type ProductImportReq struct {
	g.Meta `path:"/product/import" method:"post" mime:"multipart/form-data" tags:"体验商品" summary:"导入体验商品"`
}

// ProductImportRes 导入体验商品响应
type ProductImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// ProductImportTemplateReq 下载体验商品导入模板
type ProductImportTemplateReq struct {
	g.Meta `path:"/product/import-template" method:"get" tags:"体验商品" summary:"下载体验商品导入模板"`
}

// ProductImportTemplateRes 下载体验商品导入模板响应
type ProductImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

