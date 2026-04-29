package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// WarehouseTrade API

// WarehouseTradeCreateReq 创建仓库交易记录请求
type WarehouseTradeCreateReq struct {
	g.Meta `path:"/warehouse_trade/create" method:"post" tags:"仓库交易记录" summary:"创建仓库交易记录"`
	TradeNo string `json:"tradeNo" v:"required|max-length:64" dc:"交易编号"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"仓库商品"`
	ListingID snowflake.JsonInt64 `json:"listingID"  dc:"挂卖记录"`
	SellerID snowflake.JsonInt64 `json:"sellerID"  dc:"卖家"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"  dc:"买家"`
	TradePrice int64 `json:"tradePrice"  dc:"成交价格（分）"`
	PlatformFee int64 `json:"platformFee"  dc:"平台扣除费用（分）"`
	SellerIncome int64 `json:"sellerIncome"  dc:"卖家实收（分）"`
	TradeStatus int `json:"tradeStatus"  dc:"交易状态"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"  dc:"确认时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseTradeCreateRes 创建仓库交易记录响应
type WarehouseTradeCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseTradeUpdateReq 更新仓库交易记录请求
type WarehouseTradeUpdateReq struct {
	g.Meta `path:"/warehouse_trade/update" method:"put" tags:"仓库交易记录" summary:"更新仓库交易记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库交易记录ID"`
	TradeNo string `json:"tradeNo" v:"max-length:64" dc:"交易编号"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"仓库商品"`
	ListingID snowflake.JsonInt64 `json:"listingID"  dc:"挂卖记录"`
	SellerID snowflake.JsonInt64 `json:"sellerID"  dc:"卖家"`
	BuyerID snowflake.JsonInt64 `json:"buyerID"  dc:"买家"`
	TradePrice int64 `json:"tradePrice"  dc:"成交价格（分）"`
	PlatformFee int64 `json:"platformFee"  dc:"平台扣除费用（分）"`
	SellerIncome int64 `json:"sellerIncome"  dc:"卖家实收（分）"`
	TradeStatus int `json:"tradeStatus"  dc:"交易状态"`
	ConfirmedAt *gtime.Time `json:"confirmedAt"  dc:"确认时间"`
	Remark string `json:"remark" v:"max-length:500" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WarehouseTradeUpdateRes 更新仓库交易记录响应
type WarehouseTradeUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseTradeDeleteReq 删除仓库交易记录请求
type WarehouseTradeDeleteReq struct {
	g.Meta `path:"/warehouse_trade/delete" method:"delete" tags:"仓库交易记录" summary:"删除仓库交易记录"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库交易记录ID"`
}

// WarehouseTradeDeleteRes 删除仓库交易记录响应
type WarehouseTradeDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseTradeBatchDeleteReq 批量删除仓库交易记录请求
type WarehouseTradeBatchDeleteReq struct {
	g.Meta `path:"/warehouse_trade/batch-delete" method:"delete" tags:"仓库交易记录" summary:"批量删除仓库交易记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库交易记录ID列表"`
}

// WarehouseTradeBatchDeleteRes 批量删除仓库交易记录响应
type WarehouseTradeBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseTradeBatchUpdateReq 批量编辑仓库交易记录请求
type WarehouseTradeBatchUpdateReq struct {
	g.Meta `path:"/warehouse_trade/batch-update" method:"put" tags:"仓库交易记录" summary:"批量编辑仓库交易记录"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"仓库交易记录ID列表"`
	TradeStatus *int `json:"tradeStatus" dc:"交易状态"`
	Status *int `json:"status" dc:"状态"`
}

// WarehouseTradeBatchUpdateRes 批量编辑仓库交易记录响应
type WarehouseTradeBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WarehouseTradeDetailReq 获取仓库交易记录详情请求
type WarehouseTradeDetailReq struct {
	g.Meta `path:"/warehouse_trade/detail" method:"get" tags:"仓库交易记录" summary:"获取仓库交易记录详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"仓库交易记录ID"`
}

// WarehouseTradeDetailRes 获取仓库交易记录详情响应
type WarehouseTradeDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WarehouseTradeDetailOutput
}

// WarehouseTradeListReq 获取仓库交易记录列表请求
type WarehouseTradeListReq struct {
	g.Meta    `path:"/warehouse_trade/list" method:"get" tags:"仓库交易记录" summary:"获取仓库交易记录列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TradeNo string `json:"tradeNo" dc:"交易编号"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"仓库商品"`
	ListingID *snowflake.JsonInt64 `json:"listingID" dc:"挂卖记录"`
	SellerID *snowflake.JsonInt64 `json:"sellerID" dc:"卖家"`
	BuyerID *snowflake.JsonInt64 `json:"buyerID" dc:"买家"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	TradeStatus *int `json:"tradeStatus" dc:"交易状态"`
	Status *int `json:"status" dc:"状态"`
	ConfirmedAtStart string `json:"confirmedAtStart" dc:"确认时间开始时间"`
	ConfirmedAtEnd string `json:"confirmedAtEnd" dc:"确认时间结束时间"`
}

// WarehouseTradeListRes 获取仓库交易记录列表响应
type WarehouseTradeListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WarehouseTradeListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WarehouseTradeExportReq 导出仓库交易记录请求
type WarehouseTradeExportReq struct {
	g.Meta    `path:"/warehouse_trade/export" method:"get" tags:"仓库交易记录" summary:"导出仓库交易记录"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	TradeNo string `json:"tradeNo" dc:"交易编号"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"仓库商品"`
	ListingID *snowflake.JsonInt64 `json:"listingID" dc:"挂卖记录"`
	SellerID *snowflake.JsonInt64 `json:"sellerID" dc:"卖家"`
	BuyerID *snowflake.JsonInt64 `json:"buyerID" dc:"买家"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	TradeStatus *int `json:"tradeStatus" dc:"交易状态"`
	Status *int `json:"status" dc:"状态"`
	ConfirmedAtStart string `json:"confirmedAtStart" dc:"确认时间开始时间"`
	ConfirmedAtEnd string `json:"confirmedAtEnd" dc:"确认时间结束时间"`
}

// WarehouseTradeExportRes 导出仓库交易记录响应
type WarehouseTradeExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WarehouseTradeImportReq 导入仓库交易记录请求
type WarehouseTradeImportReq struct {
	g.Meta `path:"/warehouse_trade/import" method:"post" mime:"multipart/form-data" tags:"仓库交易记录" summary:"导入仓库交易记录"`
}

// WarehouseTradeImportRes 导入仓库交易记录响应
type WarehouseTradeImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WarehouseTradeImportTemplateReq 下载仓库交易记录导入模板
type WarehouseTradeImportTemplateReq struct {
	g.Meta `path:"/warehouse_trade/import-template" method:"get" tags:"仓库交易记录" summary:"下载仓库交易记录导入模板"`
}

// WarehouseTradeImportTemplateRes 下载仓库交易记录导入模板响应
type WarehouseTradeImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

