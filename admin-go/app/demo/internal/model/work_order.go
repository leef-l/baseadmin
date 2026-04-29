package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// WorkOrder DTO 模型

// WorkOrderCreateInput 创建体验工单输入
type WorkOrderCreateInput struct {
	TicketNo string `json:"ticketNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	Title string `json:"title"`
	Priority int `json:"priority"`
	SourceType int `json:"sourceType"`
	Description string `json:"description"`
	AttachmentFile string `json:"attachmentFile"`
	DueAt *gtime.Time `json:"dueAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WorkOrderUpdateInput 更新体验工单输入
type WorkOrderUpdateInput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TicketNo string `json:"ticketNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	Title string `json:"title"`
	Priority int `json:"priority"`
	SourceType int `json:"sourceType"`
	Description string `json:"description"`
	AttachmentFile string `json:"attachmentFile"`
	DueAt *gtime.Time `json:"dueAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
}

// WorkOrderDetailOutput 体验工单详情输出
type WorkOrderDetailOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TicketNo string `json:"ticketNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	ProductSkuNo string `json:"productSkuNo"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	OrderOrderNo string `json:"orderOrderNo"`
	Title string `json:"title"`
	Priority int `json:"priority"`
	SourceType int `json:"sourceType"`
	Description string `json:"description"`
	AttachmentFile string `json:"attachmentFile"`
	DueAt *gtime.Time `json:"dueAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WorkOrderListOutput 体验工单列表输出
type WorkOrderListOutput struct {
	ID snowflake.JsonInt64 `json:"id"`
	TicketNo string `json:"ticketNo"`
	CustomerID snowflake.JsonInt64 `json:"customerID"`
	CustomerName string `json:"customerName"`
	ProductID snowflake.JsonInt64 `json:"productID"`
	ProductSkuNo string `json:"productSkuNo"`
	OrderID snowflake.JsonInt64 `json:"orderID"`
	OrderOrderNo string `json:"orderOrderNo"`
	Title string `json:"title"`
	Priority int `json:"priority"`
	SourceType int `json:"sourceType"`
	Description string `json:"description"`
	AttachmentFile string `json:"attachmentFile"`
	DueAt *gtime.Time `json:"dueAt"`
	Status int `json:"status"`
	TenantID snowflake.JsonInt64 `json:"tenantID"`
	TenantName string `json:"tenantName"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"`
	MerchantName string `json:"merchantName"`
	CreatedAt *gtime.Time `json:"createdAt"`
	UpdatedAt *gtime.Time `json:"updatedAt"`
}

// WorkOrderListInput 体验工单列表查询输入
type WorkOrderListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	OrderBy   string `json:"orderBy"`
	OrderDir  string `json:"orderDir"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Keyword string `json:"keyword"`
	TicketNo string `json:"ticketNo"`
	Title string `json:"title"`
	CustomerID *snowflake.JsonInt64 `json:"customerID"`
	ProductID *snowflake.JsonInt64 `json:"productID"`
	OrderID *snowflake.JsonInt64 `json:"orderID"`
	TenantID *snowflake.JsonInt64 `json:"tenantID"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID"`
	Priority *int `json:"priority"`
	SourceType *int `json:"sourceType"`
	Status *int `json:"status"`
	DueAtStart string `json:"dueAtStart"`
	DueAtEnd string `json:"dueAtEnd"`
}

// WorkOrderBatchUpdateInput 批量编辑体验工单输入
type WorkOrderBatchUpdateInput struct {
	IDs    []snowflake.JsonInt64 `json:"ids"`
	Priority *int `json:"priority"`
	SourceType *int `json:"sourceType"`
	Status *int `json:"status"`
}

