package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// WarehouseListing API

// WarehouseListingCreateReq 创建仓库挂卖记录请求
type WarehouseListingCreateReq struct {
	g.Meta `path:"/warehouse_listing/create" method:"post" tags:"仓库挂卖记录" summary:"创建仓库挂卖记录"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"仓库商品"`
	SellerID snowflake.JsonInt64 `json:"sellerID"  dc:"卖家"`
	ListingPrice int64 `json:"listingPrice"  dc:"挂卖价格（分，自动加价后）"`
	ListingStatus int `json:"listingStatus"  dc:"挂卖状态"`
	ListedAt *gtime.Time `json:"listedAt"  dc:"挂卖时间"`
	SoldAt *gtime.Time `json:"soldAt"  dc:"售出时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseListingCreateRes 创建仓库挂卖记录响应
type WarehouseListingCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseListingUpdateReq 更新仓库挂卖记录请求
type WarehouseListingUpdateReq struct {
	g.Meta `path:"/warehouse_listing/update" method:"put" tags:"仓库挂卖记录" summary:"更新仓库挂卖记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库挂卖记录ID"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"仓库商品"`
	SellerID snowflake.JsonInt64 `json:"sellerID"  dc:"卖家"`
	ListingPrice int64 `json:"listingPrice"  dc:"挂卖价格（分，自动加价后）"`
	ListingStatus int `json:"listingStatus"  dc:"挂卖状态"`
	ListedAt *gtime.Time `json:"listedAt"  dc:"挂卖时间"`
	SoldAt *gtime.Time `json:"soldAt"  dc:"售出时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseListingUpdateRes 更新仓库挂卖记录响应
type WarehouseListingUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseListingDeleteReq 删除仓库挂卖记录请求
type WarehouseListingDeleteReq struct {
	g.Meta `path:"/warehouse_listing/delete" method:"delete" tags:"仓库挂卖记录" summary:"删除仓库挂卖记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库挂卖记录ID"`
}

// WarehouseListingDeleteRes 删除仓库挂卖记录响应
type WarehouseListingDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseListingBatchDeleteReq 批量删除仓库挂卖记录请求
type WarehouseListingBatchDeleteReq struct {
	g.Meta `path:"/warehouse_listing/batch-delete" method:"delete" tags:"仓库挂卖记录" summary:"批量删除仓库挂卖记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库挂卖记录ID列表"`
}

// WarehouseListingBatchDeleteRes 批量删除仓库挂卖记录响应
type WarehouseListingBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseListingBatchUpdateReq 批量编辑仓库挂卖记录请求
type WarehouseListingBatchUpdateReq struct {
	g.Meta `path:"/warehouse_listing/batch-update" method:"put" tags:"仓库挂卖记录" summary:"批量编辑仓库挂卖记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库挂卖记录ID列表"`
	ListingStatus *int `json:"listingStatus" dc:"挂卖状态"`
	Status *int `json:"status" dc:"状态"`
}

// WarehouseListingBatchUpdateRes 批量编辑仓库挂卖记录响应
type WarehouseListingBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseListingDetailReq 获取仓库挂卖记录详情请求
type WarehouseListingDetailReq struct {
	g.Meta `path:"/warehouse_listing/detail" method:"get" tags:"仓库挂卖记录" summary:"获取仓库挂卖记录详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库挂卖记录ID"`
}

// WarehouseListingDetailRes 获取仓库挂卖记录详情响应
type WarehouseListingDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WarehouseListingDetailOutput
}

// WarehouseListingListReq 获取仓库挂卖记录列表请求
type WarehouseListingListReq struct {
	g.Meta    `path:"/warehouse_listing/list" method:"get" tags:"仓库挂卖记录" summary:"获取仓库挂卖记录列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"仓库商品"`
	SellerID *snowflake.JsonInt64 `json:"sellerID" dc:"卖家"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ListingStatus *int `json:"listingStatus" dc:"挂卖状态"`
	Status *int `json:"status" dc:"状态"`
	ListedAtStart string `json:"listedAtStart" dc:"挂卖时间开始时间"`
	ListedAtEnd string `json:"listedAtEnd" dc:"挂卖时间结束时间"`
	SoldAtStart string `json:"soldAtStart" dc:"售出时间开始时间"`
	SoldAtEnd string `json:"soldAtEnd" dc:"售出时间结束时间"`
}

// WarehouseListingListRes 获取仓库挂卖记录列表响应
type WarehouseListingListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WarehouseListingListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WarehouseListingExportReq 导出仓库挂卖记录请求
type WarehouseListingExportReq struct {
	g.Meta    `path:"/warehouse_listing/export" method:"get" tags:"仓库挂卖记录" summary:"导出仓库挂卖记录"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"仓库商品"`
	SellerID *snowflake.JsonInt64 `json:"sellerID" dc:"卖家"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	ListingStatus *int `json:"listingStatus" dc:"挂卖状态"`
	Status *int `json:"status" dc:"状态"`
	ListedAtStart string `json:"listedAtStart" dc:"挂卖时间开始时间"`
	ListedAtEnd string `json:"listedAtEnd" dc:"挂卖时间结束时间"`
	SoldAtStart string `json:"soldAtStart" dc:"售出时间开始时间"`
	SoldAtEnd string `json:"soldAtEnd" dc:"售出时间结束时间"`
}

// WarehouseListingExportRes 导出仓库挂卖记录响应
type WarehouseListingExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WarehouseListingImportReq 导入仓库挂卖记录请求
type WarehouseListingImportReq struct {
	g.Meta `path:"/warehouse_listing/import" method:"post" mime:"multipart/form-data" tags:"仓库挂卖记录" summary:"导入仓库挂卖记录"`
}

// WarehouseListingImportRes 导入仓库挂卖记录响应
type WarehouseListingImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WarehouseListingImportTemplateReq 下载仓库挂卖记录导入模板
type WarehouseListingImportTemplateReq struct {
	g.Meta `path:"/warehouse_listing/import-template" method:"get" tags:"仓库挂卖记录" summary:"下载仓库挂卖记录导入模板"`
}

// WarehouseListingImportTemplateRes 下载仓库挂卖记录导入模板响应
type WarehouseListingImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

