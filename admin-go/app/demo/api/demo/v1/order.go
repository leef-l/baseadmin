package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Order API

// OrderCreateReq 创建体验订单请求
type OrderCreateReq struct {
	g.Meta `path:"/order/create" method:"post" tags:"体验订单" summary:"创建体验订单"`
	OrderNo string `json:"orderNo" v:"required|max-length:50" dc:"订单号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	ProductID snowflake.JsonInt64 `json:"productID"  dc:"商品"`
	Quantity int `json:"quantity"  dc:"购买数量"`
	Amount int `json:"amount"  dc:"订单金额（分）"`
	PayStatus int `json:"payStatus"  dc:"支付状态"`
	DeliverStatus int `json:"deliverStatus"  dc:"发货状态"`
	PaidAt *gtime.Time `json:"paidAt"  dc:"支付时间"`
	DeliverAt *gtime.Time `json:"deliverAt"  dc:"发货时间"`
	ReceiverPhone string `json:"receiverPhone" v:"phone-loose|max-length:30" dc:"收货电话"`
	Address string `json:"address" v:"max-length:255" dc:"收货地址"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// OrderCreateRes 创建体验订单响应
type OrderCreateRes struct {
	g.Meta `mime:"application/json"`
}

// OrderUpdateReq 更新体验订单请求
type OrderUpdateReq struct {
	g.Meta `path:"/order/update" method:"put" tags:"体验订单" summary:"更新体验订单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验订单ID"`
	OrderNo string `json:"orderNo" v:"max-length:50" dc:"订单号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	ProductID snowflake.JsonInt64 `json:"productID"  dc:"商品"`
	Quantity int `json:"quantity"  dc:"购买数量"`
	Amount int `json:"amount"  dc:"订单金额（分）"`
	PayStatus int `json:"payStatus"  dc:"支付状态"`
	DeliverStatus int `json:"deliverStatus"  dc:"发货状态"`
	PaidAt *gtime.Time `json:"paidAt"  dc:"支付时间"`
	DeliverAt *gtime.Time `json:"deliverAt"  dc:"发货时间"`
	ReceiverPhone string `json:"receiverPhone" v:"phone-loose|max-length:30" dc:"收货电话"`
	Address string `json:"address" v:"max-length:255" dc:"收货地址"`
	Remark string `json:"remark" v:"max-length:65535" dc:"备注"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// OrderUpdateRes 更新体验订单响应
type OrderUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// OrderDeleteReq 删除体验订单请求
type OrderDeleteReq struct {
	g.Meta `path:"/order/delete" method:"delete" tags:"体验订单" summary:"删除体验订单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验订单ID"`
}

// OrderDeleteRes 删除体验订单响应
type OrderDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// OrderBatchDeleteReq 批量删除体验订单请求
type OrderBatchDeleteReq struct {
	g.Meta `path:"/order/batch-delete" method:"delete" tags:"体验订单" summary:"批量删除体验订单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验订单ID列表"`
}

// OrderBatchDeleteRes 批量删除体验订单响应
type OrderBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// OrderBatchUpdateReq 批量编辑体验订单请求
type OrderBatchUpdateReq struct {
	g.Meta `path:"/order/batch-update" method:"put" tags:"体验订单" summary:"批量编辑体验订单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required#ID列表不能为空" dc:"体验订单ID列表"`
	PayStatus *int `json:"payStatus" dc:"支付状态"`
	DeliverStatus *int `json:"deliverStatus" dc:"发货状态"`
	Status *int `json:"status" dc:"状态"`
}

// OrderBatchUpdateRes 批量编辑体验订单响应
type OrderBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// OrderDetailReq 获取体验订单详情请求
type OrderDetailReq struct {
	g.Meta `path:"/order/detail" method:"get" tags:"体验订单" summary:"获取体验订单详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验订单ID"`
}

// OrderDetailRes 获取体验订单详情响应
type OrderDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.OrderDetailOutput
}

// OrderListReq 获取体验订单列表请求
type OrderListReq struct {
	g.Meta    `path:"/order/list" method:"get" tags:"体验订单" summary:"获取体验订单列表"`
	PageNum   int    `json:"pageNum" d:"1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	OrderNo string `json:"orderNo" dc:"订单号"`
	ReceiverPhone string `json:"receiverPhone" dc:"收货电话"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	ProductID *snowflake.JsonInt64 `json:"productID" dc:"商品"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	PayStatus *int `json:"payStatus" dc:"支付状态"`
	DeliverStatus *int `json:"deliverStatus" dc:"发货状态"`
	Status *int `json:"status" dc:"状态"`
	PaidAtStart string `json:"paidAtStart" dc:"支付时间开始时间"`
	PaidAtEnd string `json:"paidAtEnd" dc:"支付时间结束时间"`
	DeliverAtStart string `json:"deliverAtStart" dc:"发货时间开始时间"`
	DeliverAtEnd string `json:"deliverAtEnd" dc:"发货时间结束时间"`
}

// OrderListRes 获取体验订单列表响应
type OrderListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.OrderListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// OrderExportReq 导出体验订单请求
type OrderExportReq struct {
	g.Meta    `path:"/order/export" method:"get" tags:"体验订单" summary:"导出体验订单"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	OrderNo string `json:"orderNo" dc:"订单号"`
	ReceiverPhone string `json:"receiverPhone" dc:"收货电话"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	ProductID *snowflake.JsonInt64 `json:"productID" dc:"商品"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	PayStatus *int `json:"payStatus" dc:"支付状态"`
	DeliverStatus *int `json:"deliverStatus" dc:"发货状态"`
	Status *int `json:"status" dc:"状态"`
	PaidAtStart string `json:"paidAtStart" dc:"支付时间开始时间"`
	PaidAtEnd string `json:"paidAtEnd" dc:"支付时间结束时间"`
	DeliverAtStart string `json:"deliverAtStart" dc:"发货时间开始时间"`
	DeliverAtEnd string `json:"deliverAtEnd" dc:"发货时间结束时间"`
}

// OrderExportRes 导出体验订单响应
type OrderExportRes struct {
	g.Meta `mime:"text/csv"`
}

// OrderImportReq 导入体验订单请求
type OrderImportReq struct {
	g.Meta `path:"/order/import" method:"post" mime:"multipart/form-data" tags:"体验订单" summary:"导入体验订单"`
}

// OrderImportRes 导入体验订单响应
type OrderImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// OrderImportTemplateReq 下载体验订单导入模板
type OrderImportTemplateReq struct {
	g.Meta `path:"/order/import-template" method:"get" tags:"体验订单" summary:"下载体验订单导入模板"`
}

// OrderImportTemplateRes 下载体验订单导入模板响应
type OrderImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

