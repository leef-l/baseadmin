package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// WarehouseGoods API

// WarehouseGoodsCreateReq 创建仓库商品请求
type WarehouseGoodsCreateReq struct {
	g.Meta `path:"/warehouse_goods/create" method:"post" tags:"仓库商品" summary:"创建仓库商品"`
	GoodsNo string `json:"goodsNo" v:"required|max-length:64" dc:"商品编号"`
	Title string `json:"title" v:"required|max-length:200" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"商品封面"`
	InitPrice int64 `json:"initPrice"  dc:"初始价格（分）"`
	CurrentPrice int64 `json:"currentPrice"  dc:"当前价格（分）"`
	PriceRiseRate int `json:"priceRiseRate"  dc:"每次加价比例（百分比，如10=10%）"`
	PlatformFeeRate int `json:"platformFeeRate"  dc:"平台扣除比例（百分比，如5=5%）"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"  dc:"当前持有人"`
	TradeCount int `json:"tradeCount"  dc:"流转次数"`
	GoodsStatus int `json:"goodsStatus"  dc:"商品状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseGoodsCreateRes 创建仓库商品响应
type WarehouseGoodsCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseGoodsUpdateReq 更新仓库商品请求
type WarehouseGoodsUpdateReq struct {
	g.Meta `path:"/warehouse_goods/update" method:"put" tags:"仓库商品" summary:"更新仓库商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库商品ID"`
	GoodsNo string `json:"goodsNo" v:"max-length:64" dc:"商品编号"`
	Title string `json:"title" v:"max-length:200" dc:"商品名称"`
	Cover string `json:"cover" v:"max-length:500" dc:"商品封面"`
	InitPrice int64 `json:"initPrice"  dc:"初始价格（分）"`
	CurrentPrice int64 `json:"currentPrice"  dc:"当前价格（分）"`
	PriceRiseRate int `json:"priceRiseRate"  dc:"每次加价比例（百分比，如10=10%）"`
	PlatformFeeRate int `json:"platformFeeRate"  dc:"平台扣除比例（百分比，如5=5%）"`
	OwnerID snowflake.JsonInt64 `json:"ownerID"  dc:"当前持有人"`
	TradeCount int `json:"tradeCount"  dc:"流转次数"`
	GoodsStatus int `json:"goodsStatus"  dc:"商品状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Sort int `json:"sort"  dc:"排序（升序）"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseGoodsUpdateRes 更新仓库商品响应
type WarehouseGoodsUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseGoodsDeleteReq 删除仓库商品请求
type WarehouseGoodsDeleteReq struct {
	g.Meta `path:"/warehouse_goods/delete" method:"delete" tags:"仓库商品" summary:"删除仓库商品"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库商品ID"`
}

// WarehouseGoodsDeleteRes 删除仓库商品响应
type WarehouseGoodsDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseGoodsBatchDeleteReq 批量删除仓库商品请求
type WarehouseGoodsBatchDeleteReq struct {
	g.Meta `path:"/warehouse_goods/batch-delete" method:"delete" tags:"仓库商品" summary:"批量删除仓库商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库商品ID列表"`
}

// WarehouseGoodsBatchDeleteRes 批量删除仓库商品响应
type WarehouseGoodsBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseGoodsBatchUpdateReq 批量编辑仓库商品请求
type WarehouseGoodsBatchUpdateReq struct {
	g.Meta `path:"/warehouse_goods/batch-update" method:"put" tags:"仓库商品" summary:"批量编辑仓库商品"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库商品ID列表"`
	GoodsStatus *int `json:"goodsStatus" dc:"商品状态"`
	Status *int `json:"status" dc:"状态"`
}

// WarehouseGoodsBatchUpdateRes 批量编辑仓库商品响应
type WarehouseGoodsBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseGoodsDetailReq 获取仓库商品详情请求
type WarehouseGoodsDetailReq struct {
	g.Meta `path:"/warehouse_goods/detail" method:"get" tags:"仓库商品" summary:"获取仓库商品详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库商品ID"`
}

// WarehouseGoodsDetailRes 获取仓库商品详情响应
type WarehouseGoodsDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WarehouseGoodsDetailOutput
}

// WarehouseGoodsListReq 获取仓库商品列表请求
type WarehouseGoodsListReq struct {
	g.Meta    `path:"/warehouse_goods/list" method:"get" tags:"仓库商品" summary:"获取仓库商品列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	GoodsNo string `json:"goodsNo" dc:"商品编号"`
	Title string `json:"title" dc:"商品名称"`
	OwnerID *snowflake.JsonInt64 `json:"ownerID" dc:"当前持有人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	GoodsStatus *int `json:"goodsStatus" dc:"商品状态"`
	Status *int `json:"status" dc:"状态"`
}

// WarehouseGoodsListRes 获取仓库商品列表响应
type WarehouseGoodsListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WarehouseGoodsListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WarehouseGoodsExportReq 导出仓库商品请求
type WarehouseGoodsExportReq struct {
	g.Meta    `path:"/warehouse_goods/export" method:"get" tags:"仓库商品" summary:"导出仓库商品"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	GoodsNo string `json:"goodsNo" dc:"商品编号"`
	Title string `json:"title" dc:"商品名称"`
	OwnerID *snowflake.JsonInt64 `json:"ownerID" dc:"当前持有人"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	GoodsStatus *int `json:"goodsStatus" dc:"商品状态"`
	Status *int `json:"status" dc:"状态"`
}

// WarehouseGoodsExportRes 导出仓库商品响应
type WarehouseGoodsExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WarehouseGoodsImportReq 导入仓库商品请求
type WarehouseGoodsImportReq struct {
	g.Meta `path:"/warehouse_goods/import" method:"post" mime:"multipart/form-data" tags:"仓库商品" summary:"导入仓库商品"`
}

// WarehouseGoodsImportRes 导入仓库商品响应
type WarehouseGoodsImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WarehouseGoodsImportTemplateReq 下载仓库商品导入模板
type WarehouseGoodsImportTemplateReq struct {
	g.Meta `path:"/warehouse_goods/import-template" method:"get" tags:"仓库商品" summary:"下载仓库商品导入模板"`
}

// WarehouseGoodsImportTemplateRes 下载仓库商品导入模板响应
type WarehouseGoodsImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

