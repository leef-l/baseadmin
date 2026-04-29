package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Order DTO 模型

// OrderCreateInput 创建体验订单输入
type OrderCreateInput struct {
	OrderNo string `json:"orderNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	Quantity int `json:"quantity"`
	Amount int `json:"amount"`
	PayStatus int `json:"payStatus"`
	DeliverStatus int `json:"deliverStatus"`
	PaidAt *gtime.Time `json:"paidAt"`
	DeliverAt *gtime.Time `json:"deliverAt"`
	ReceiverPhone string `json:"receiverPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// OrderUpdateInput 更新体验订单输入
type OrderUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	Quantity int `json:"quantity"`
	Amount int `json:"amount"`
	PayStatus int `json:"payStatus"`
	DeliverStatus int `json:"deliverStatus"`
	PaidAt *gtime.Time `json:"paidAt"`
	DeliverAt *gtime.Time `json:"deliverAt"`
	ReceiverPhone string `json:"receiverPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// OrderDetailOutput 体验订单详情输出
type OrderDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	ProductSkuNo string `json:"productSkuNo"`
	Quantity int `json:"quantity"`
	Amount int `json:"amount"`
	PayStatus int `json:"payStatus"`
	DeliverStatus int `json:"deliverStatus"`
	PaidAt *gtime.Time `json:"paidAt"`
	DeliverAt *gtime.Time `json:"deliverAt"`
	ReceiverPhone string `json:"receiverPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// OrderListOutput 体验订单列表输出
type OrderListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	OrderNo string `json:"orderNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	ProductSkuNo string `json:"productSkuNo"`
	Quantity int `json:"quantity"`
	Amount int `json:"amount"`
	PayStatus int `json:"payStatus"`
	DeliverStatus int `json:"deliverStatus"`
	PaidAt *gtime.Time `json:"paidAt"`
	DeliverAt *gtime.Time `json:"deliverAt"`
	ReceiverPhone string `json:"receiverPhone"`
	Address string `json:"address"`
	Remark string `json:"remark"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// OrderListInput 体验订单列表查询输入
type OrderListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	OrderNo string `json:"orderNo"`
	ReceiverPhone string `json:"receiverPhone"`
	CustomerID *snowflake.JsonInt64 `json:"customerID"`
	ProductID *snowflake.JsonInt64 `json:"productID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	PayStatus *int `json:"payStatus"`
	DeliverStatus *int `json:"deliverStatus"`
	Status *int `json:"status"`
	PaidAtStart string `json:"paidAtStart"`
	PaidAtEnd string `json:"paidAtEnd"`
	DeliverAtStart string `json:"deliverAtStart"`
	DeliverAtEnd string `json:"deliverAtEnd"`
}

// OrderBatchUpdateInput 批量编辑体验订单输入
type OrderBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	PayStatus *int `json:"payStatus"`
	DeliverStatus *int `json:"deliverStatus"`
	Status *int `json:"status"`
}

