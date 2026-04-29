package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// ShopGoods API

// ShopGoodsCreateReq 创建商城商品请求
type ShopGoodsCreateReq struct {
	g.Meta `path:"/shop_goods/create" method:"post" tags:"商城商品" summary:"创建商城商品"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"  dc:"商品分类"`
	Title string `json:"title" v:"required|max-length:200" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"封面图"`
	Images string `json:"images" v:"max-length:65535" dc:"商品图片（JSON数组）"`
	Price int64 `json:"price"  dc:"售价（分，优惠券余额支付）"`
	OriginalPrice int64 `json:"originalPrice"  dc:"原价（分）"`
	Stock int `json:"stock"  dc:"库存"`
	Sales int `json:"sales"  dc:"销量"`
	Content string `json:"content" v:"max-length:65535" dc:"商品详情"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	IsRecommend int `json:"isRecommend"  dc:"是否推荐"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopGoodsCreateRes 创建商城商品响应
type ShopGoodsCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopGoodsUpdateReq 更新商城商品请求
type ShopGoodsUpdateReq struct {
	g.Meta `path:"/shop_goods/update" method:"put" tags:"商城商品" summary:"更新商城商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品ID"`
	CategoryID snowflake.JsonInt64 `json:"categoryID"  dc:"商品分类"`
	Title string `json:"title" v:"max-length:200" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"封面图"`
	Images string `json:"images" v:"max-length:65535" dc:"商品图片（JSON数组）"`
	Price int64 `json:"price"  dc:"售价（分，优惠券余额支付）"`
	OriginalPrice int64 `json:"originalPrice"  dc:"原价（分）"`
	Stock int `json:"stock"  dc:"库存"`
	Sales int `json:"sales"  dc:"销量"`
	Content string `json:"content" v:"max-length:65535" dc:"商品详情"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	IsRecommend int `json:"isRecommend"  dc:"是否推荐"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopGoodsUpdateRes 更新商城商品响应
type ShopGoodsUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopGoodsDeleteReq 删除商城商品请求
type ShopGoodsDeleteReq struct {
	g.Meta `path:"/shop_goods/delete" method:"delete" tags:"商城商品" summary:"删除商城商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品ID"`
}

// ShopGoodsDeleteRes 删除商城商品响应
type ShopGoodsDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopGoodsBatchDeleteReq 批量删除商城商品请求
type ShopGoodsBatchDeleteReq struct {
	g.Meta `path:"/shop_goods/batch-delete" method:"delete" tags:"商城商品" summary:"批量删除商城商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城商品ID列表"`
}

// ShopGoodsBatchDeleteRes 批量删除商城商品响应
type ShopGoodsBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopGoodsBatchUpdateReq 批量编辑商城商品请求
type ShopGoodsBatchUpdateReq struct {
	g.Meta `path:"/shop_goods/batch-update" method:"put" tags:"商城商品" summary:"批量编辑商城商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城商品ID列表"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ShopGoodsBatchUpdateRes 批量编辑商城商品响应
type ShopGoodsBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopGoodsDetailReq 获取商城商品详情请求
type ShopGoodsDetailReq struct {
	g.Meta `path:"/shop_goods/detail" method:"get" tags:"商城商品" summary:"获取商城商品详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城商品ID"`
}

// ShopGoodsDetailRes 获取商城商品详情响应
type ShopGoodsDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ShopGoodsDetailOutput
}

// ShopGoodsListReq 获取商城商品列表请求
type ShopGoodsListReq struct {
	g.Meta    `path:"/shop_goods/list" method:"get" tags:"商城商品" summary:"获取商城商品列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Title string `json:"title" dc:"商品名称"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID" dc:"商品分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ShopGoodsListRes 获取商城商品列表响应
type ShopGoodsListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ShopGoodsListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ShopGoodsExportReq 导出商城商品请求
type ShopGoodsExportReq struct {
	g.Meta    `path:"/shop_goods/export" method:"get" tags:"商城商品" summary:"导出商城商品"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Title string `json:"title" dc:"商品名称"`
	CategoryID *snowflake.JsonInt64 `json:"categoryID" dc:"商品分类"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	IsRecommend *int `json:"isRecommend" dc:"是否推荐"`
	Status *int `json:"status" dc:"状态"`
}

// ShopGoodsExportRes 导出商城商品响应
type ShopGoodsExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ShopGoodsImportReq 导入商城商品请求
type ShopGoodsImportReq struct {
	g.Meta `path:"/shop_goods/import" method:"post" mime:"multipart/form-data" tags:"商城商品" summary:"导入商城商品"`
}

// ShopGoodsImportRes 导入商城商品响应
type ShopGoodsImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// ShopGoodsImportTemplateReq 下载商城商品导入模板
type ShopGoodsImportTemplateReq struct {
	g.Meta `path:"/shop_goods/import-template" method:"get" tags:"商城商品" summary:"下载商城商品导入模板"`
}

// ShopGoodsImportTemplateRes 下载商城商品导入模板响应
type ShopGoodsImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

