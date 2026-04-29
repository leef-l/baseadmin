package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// ShopOrder API

// ShopOrderCreateReq 创建商城订单请求
type ShopOrderCreateReq struct {
	g.Meta `path:"/shop_order/create" method:"post" tags:"商城订单" summary:"创建商城订单"`
	OrderNo string `json:"orderNo" v:"required|max-length:64" dc:"订单号"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"购买会员"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"商品"`
	GoodsTitle string `json:"goodsTitle" v:"max-length:200" dc:"商品名称（快照）"`
	GoodsCover string `json:"goodsCover" v:"max-length:500" dc:"商品封面（快照）"`
	Quantity int `json:"quantity"  dc:"购买数量"`
	TotalPrice int64 `json:"totalPrice"  dc:"订单总价（分）"`
	PayWallet int `json:"payWallet"  dc:"支付钱包"`
	OrderStatus int `json:"orderStatus"  dc:"订单状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"订单备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopOrderCreateRes 创建商城订单响应
type ShopOrderCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopOrderUpdateReq 更新商城订单请求
type ShopOrderUpdateReq struct {
	g.Meta `path:"/shop_order/update" method:"put" tags:"商城订单" summary:"更新商城订单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城订单ID"`
	OrderNo string `json:"orderNo" v:"max-length:64" dc:"订单号"`
	UserID snowflake.JsonInt64 `json:"userID"  dc:"购买会员"`
	GoodsID snowflake.JsonInt64 `json:"goodsID"  dc:"商品"`
	GoodsTitle string `json:"goodsTitle" v:"max-length:200" dc:"商品名称（快照）"`
	GoodsCover string `json:"goodsCover" v:"max-length:500" dc:"商品封面（快照）"`
	Quantity int `json:"quantity"  dc:"购买数量"`
	TotalPrice int64 `json:"totalPrice"  dc:"订单总价（分）"`
	PayWallet int `json:"payWallet"  dc:"支付钱包"`
	OrderStatus int `json:"orderStatus"  dc:"订单状态"`
	Remark string `json:"remark" v:"max-length:500" dc:"订单备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// ShopOrderUpdateRes 更新商城订单响应
type ShopOrderUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopOrderDeleteReq 删除商城订单请求
type ShopOrderDeleteReq struct {
	g.Meta `path:"/shop_order/delete" method:"delete" tags:"商城订单" summary:"删除商城订单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城订单ID"`
}

// ShopOrderDeleteRes 删除商城订单响应
type ShopOrderDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopOrderBatchDeleteReq 批量删除商城订单请求
type ShopOrderBatchDeleteReq struct {
	g.Meta `path:"/shop_order/batch-delete" method:"delete" tags:"商城订单" summary:"批量删除商城订单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城订单ID列表"`
}

// ShopOrderBatchDeleteRes 批量删除商城订单响应
type ShopOrderBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ShopOrderBatchUpdateReq 批量编辑商城订单请求
type ShopOrderBatchUpdateReq struct {
	g.Meta `path:"/shop_order/batch-update" method:"put" tags:"商城订单" summary:"批量编辑商城订单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"商城订单ID列表"`
	PayWallet *int `json:"payWallet" dc:"支付钱包"`
	OrderStatus *int `json:"orderStatus" dc:"订单状态"`
	Status *int `json:"status" dc:"状态"`
}

// ShopOrderBatchUpdateRes 批量编辑商城订单响应
type ShopOrderBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ShopOrderDetailReq 获取商城订单详情请求
type ShopOrderDetailReq struct {
	g.Meta `path:"/shop_order/detail" method:"get" tags:"商城订单" summary:"获取商城订单详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"商城订单ID"`
}

// ShopOrderDetailRes 获取商城订单详情响应
type ShopOrderDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ShopOrderDetailOutput
}

// ShopOrderListReq 获取商城订单列表请求
type ShopOrderListReq struct {
	g.Meta    `path:"/shop_order/list" method:"get" tags:"商城订单" summary:"获取商城订单列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	OrderNo string `json:"orderNo" dc:"订单号"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"购买会员"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"商品"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	PayWallet *int `json:"payWallet" dc:"支付钱包"`
	OrderStatus *int `json:"orderStatus" dc:"订单状态"`
	Status *int `json:"status" dc:"状态"`
	GoodsTitle string `json:"goodsTitle" dc:"商品名称（快照）"`
}

// ShopOrderListRes 获取商城订单列表响应
type ShopOrderListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ShopOrderListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// ShopOrderExportReq 导出商城订单请求
type ShopOrderExportReq struct {
	g.Meta    `path:"/shop_order/export" method:"get" tags:"商城订单" summary:"导出商城订单"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	OrderNo string `json:"orderNo" dc:"订单号"`
	UserID *snowflake.JsonInt64 `json:"userID" dc:"购买会员"`
	GoodsID *snowflake.JsonInt64 `json:"goodsID" dc:"商品"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	PayWallet *int `json:"payWallet" dc:"支付钱包"`
	OrderStatus *int `json:"orderStatus" dc:"订单状态"`
	Status *int `json:"status" dc:"状态"`
	GoodsTitle string `json:"goodsTitle" dc:"商品名称（快照）"`
}

// ShopOrderExportRes 导出商城订单响应
type ShopOrderExportRes struct {
	g.Meta `mime:"text/csv"`
}

// ShopOrderImportReq 导入商城订单请求
type ShopOrderImportReq struct {
	g.Meta `path:"/shop_order/import" method:"post" mime:"multipart/form-data" tags:"商城订单" summary:"导入商城订单"`
}

// ShopOrderImportRes 导入商城订单响应
type ShopOrderImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// ShopOrderImportTemplateReq 下载商城订单导入模板
type ShopOrderImportTemplateReq struct {
	g.Meta `path:"/shop_order/import-template" method:"get" tags:"商城订单" summary:"下载商城订单导入模板"`
}

// ShopOrderImportTemplateRes 下载商城订单导入模板响应
type ShopOrderImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

