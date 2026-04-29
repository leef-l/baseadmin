package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// WorkOrder API

// WorkOrderCreateReq 创建体验工单请求
type WorkOrderCreateReq struct {
	g.Meta `path:"/work_order/create" method:"post" tags:"体验工单" summary:"创建体验工单"`
	TicketNo string `json:"ticketNo" v:"required|max-length:50" dc:"工单号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	ProductID snowflake.JsonInt64 `json:"productID"  dc:"商品"`
	OrderID snowflake.JsonInt64 `json:"orderID"  dc:"订单"`
	Title string `json:"title" v:"required|max-length:120" dc:"工单标题"`
	Priority int `json:"priority"  dc:"优先级"`
	SourceType int `json:"sourceType"  dc:"来源"`
	Description string `json:"description" v:"max-length:65535" dc:"问题描述"`
	AttachmentFile string `json:"attachmentFile" v:"max-length:500" dc:"附件"`
	DueAt *gtime.Time `json:"dueAt"  dc:"截止时间"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WorkOrderCreateRes 创建体验工单响应
type WorkOrderCreateRes struct {
	g.Meta `mime:"application/json"`
}

// WorkOrderUpdateReq 更新体验工单请求
type WorkOrderUpdateReq struct {
	g.Meta `path:"/work_order/update" method:"put" tags:"体验工单" summary:"更新体验工单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验工单ID"`
	TicketNo string `json:"ticketNo" v:"max-length:50" dc:"工单号"`
	CustomerID snowflake.JsonInt64 `json:"customerID"  dc:"客户"`
	ProductID snowflake.JsonInt64 `json:"productID"  dc:"商品"`
	OrderID snowflake.JsonInt64 `json:"orderID"  dc:"订单"`
	Title string `json:"title" v:"max-length:120" dc:"工单标题"`
	Priority int `json:"priority"  dc:"优先级"`
	SourceType int `json:"sourceType"  dc:"来源"`
	Description string `json:"description" v:"max-length:65535" dc:"问题描述"`
	AttachmentFile string `json:"attachmentFile" v:"max-length:500" dc:"附件"`
	DueAt *gtime.Time `json:"dueAt"  dc:"截止时间"`
	Status int `json:"status"  dc:"状态"`
	TenantID snowflake.JsonInt64 `json:"tenantID"  dc:"租户"`
	MerchantID snowflake.JsonInt64 `json:"merchantID"  dc:"商户"`
}

// WorkOrderUpdateRes 更新体验工单响应
type WorkOrderUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WorkOrderDeleteReq 删除体验工单请求
type WorkOrderDeleteReq struct {
	g.Meta `path:"/work_order/delete" method:"delete" tags:"体验工单" summary:"删除体验工单"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验工单ID"`
}

// WorkOrderDeleteRes 删除体验工单响应
type WorkOrderDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WorkOrderBatchDeleteReq 批量删除体验工单请求
type WorkOrderBatchDeleteReq struct {
	g.Meta `path:"/work_order/batch-delete" method:"delete" tags:"体验工单" summary:"批量删除体验工单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验工单ID列表"`
}

// WorkOrderBatchDeleteRes 批量删除体验工单响应
type WorkOrderBatchDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// WorkOrderBatchUpdateReq 批量编辑体验工单请求
type WorkOrderBatchUpdateReq struct {
	g.Meta `path:"/work_order/batch-update" method:"put" tags:"体验工单" summary:"批量编辑体验工单"`
	IDs    []snowflake.JsonInt64 `json:"ids" v:"required|max-length:500#ID列表不能为空|最多支持500条" dc:"体验工单ID列表"`
	Priority *int `json:"priority" dc:"优先级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	Status *int `json:"status" dc:"状态"`
}

// WorkOrderBatchUpdateRes 批量编辑体验工单响应
type WorkOrderBatchUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// WorkOrderDetailReq 获取体验工单详情请求
type WorkOrderDetailReq struct {
	g.Meta `path:"/work_order/detail" method:"get" tags:"体验工单" summary:"获取体验工单详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"体验工单ID"`
}

// WorkOrderDetailRes 获取体验工单详情响应
type WorkOrderDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.WorkOrderDetailOutput
}

// WorkOrderListReq 获取体验工单列表请求
type WorkOrderListReq struct {
	g.Meta    `path:"/work_order/list" method:"get" tags:"体验工单" summary:"获取体验工单列表"`
	PageNum   int    `json:"pageNum" d:"1" v:"min:1" dc:"页码"`
	PageSize  int    `json:"pageSize" d:"10" v:"between:1,500" dc:"每页数量"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	TicketNo string `json:"ticketNo" dc:"工单号"`
	Title string `json:"title" dc:"工单标题"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	ProductID *snowflake.JsonInt64 `json:"productID" dc:"商品"`
	OrderID *snowflake.JsonInt64 `json:"orderID" dc:"订单"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Priority *int `json:"priority" dc:"优先级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	Status *int `json:"status" dc:"状态"`
	DueAtStart string `json:"dueAtStart" dc:"截止时间开始时间"`
	DueAtEnd string `json:"dueAtEnd" dc:"截止时间结束时间"`
}

// WorkOrderListRes 获取体验工单列表响应
type WorkOrderListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.WorkOrderListOutput `json:"list" dc:"列表数据"`
	Total  int                               `json:"total" dc:"总数"`
}
// WorkOrderExportReq 导出体验工单请求
type WorkOrderExportReq struct {
	g.Meta    `path:"/work_order/export" method:"get" tags:"体验工单" summary:"导出体验工单"`
	OrderBy   string `json:"orderBy" dc:"排序字段"`
	OrderDir  string `json:"orderDir" d:"desc" v:"in:asc,desc" dc:"排序方向:asc/desc"`
	StartTime string `json:"startTime" dc:"开始时间"`
	EndTime   string `json:"endTime" dc:"结束时间"`
	Keyword string `json:"keyword" dc:"关键词"`
	TicketNo string `json:"ticketNo" dc:"工单号"`
	Title string `json:"title" dc:"工单标题"`
	CustomerID *snowflake.JsonInt64 `json:"customerID" dc:"客户"`
	ProductID *snowflake.JsonInt64 `json:"productID" dc:"商品"`
	OrderID *snowflake.JsonInt64 `json:"orderID" dc:"订单"`
	TenantID *snowflake.JsonInt64 `json:"tenantID" dc:"租户"`
	MerchantID *snowflake.JsonInt64 `json:"merchantID" dc:"商户"`
	Priority *int `json:"priority" dc:"优先级"`
	SourceType *int `json:"sourceType" dc:"来源"`
	Status *int `json:"status" dc:"状态"`
	DueAtStart string `json:"dueAtStart" dc:"截止时间开始时间"`
	DueAtEnd string `json:"dueAtEnd" dc:"截止时间结束时间"`
}

// WorkOrderExportRes 导出体验工单响应
type WorkOrderExportRes struct {
	g.Meta `mime:"text/csv"`
}

// WorkOrderImportReq 导入体验工单请求
type WorkOrderImportReq struct {
	g.Meta `path:"/work_order/import" method:"post" mime:"multipart/form-data" tags:"体验工单" summary:"导入体验工单"`
}

// WorkOrderImportRes 导入体验工单响应
type WorkOrderImportRes struct {
	g.Meta  `mime:"application/json"`
	Success int `json:"success" dc:"成功条数"`
	Fail    int `json:"fail" dc:"失败条数"`
}

// WorkOrderImportTemplateReq 下载体验工单导入模板
type WorkOrderImportTemplateReq struct {
	g.Meta `path:"/work_order/import-template" method:"get" tags:"体验工单" summary:"下载体验工单导入模板"`
}

// WorkOrderImportTemplateRes 下载体验工单导入模板响应
type WorkOrderImportTemplateRes struct {
	g.Meta `mime:"text/csv"`
}

