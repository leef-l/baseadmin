package work_order

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

func csvSafeWorkOrder(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var WorkOrder = cWorkOrder{}

type cWorkOrder struct{}

// Create 创建体验工单
func (c *cWorkOrder) Create(ctx context.Context, req *v1.WorkOrderCreateReq) (res *v1.WorkOrderCreateRes, err error) {
	err = service.WorkOrder().Create(ctx, &model.WorkOrderCreateInput{
		TicketNo: req.TicketNo,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		OrderID: req.OrderID,
		Title: req.Title,
		Priority: req.Priority,
		SourceType: req.SourceType,
		Description: req.Description,
		AttachmentFile: req.AttachmentFile,
		DueAt: req.DueAt,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验工单
func (c *cWorkOrder) Update(ctx context.Context, req *v1.WorkOrderUpdateReq) (res *v1.WorkOrderUpdateRes, err error) {
	err = service.WorkOrder().Update(ctx, &model.WorkOrderUpdateInput{
		ID: req.ID,
		TicketNo: req.TicketNo,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		OrderID: req.OrderID,
		Title: req.Title,
		Priority: req.Priority,
		SourceType: req.SourceType,
		Description: req.Description,
		AttachmentFile: req.AttachmentFile,
		DueAt: req.DueAt,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验工单
func (c *cWorkOrder) Delete(ctx context.Context, req *v1.WorkOrderDeleteReq) (res *v1.WorkOrderDeleteRes, err error) {
	err = service.WorkOrder().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验工单
func (c *cWorkOrder) BatchDelete(ctx context.Context, req *v1.WorkOrderBatchDeleteReq) (res *v1.WorkOrderBatchDeleteRes, err error) {
	err = service.WorkOrder().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验工单
func (c *cWorkOrder) BatchUpdate(ctx context.Context, req *v1.WorkOrderBatchUpdateReq) (res *v1.WorkOrderBatchUpdateRes, err error) {
	err = service.WorkOrder().BatchUpdate(ctx, &model.WorkOrderBatchUpdateInput{
		IDs: req.IDs,
		Priority: req.Priority,
		SourceType: req.SourceType,
		Status: req.Status,
	})
	return
}

// Detail 获取体验工单详情
func (c *cWorkOrder) Detail(ctx context.Context, req *v1.WorkOrderDetailReq) (res *v1.WorkOrderDetailRes, err error) {
	res = &v1.WorkOrderDetailRes{}
	res.WorkOrderDetailOutput, err = service.WorkOrder().Detail(ctx, req.ID)
	return
}

// List 获取体验工单列表
func (c *cWorkOrder) List(ctx context.Context, req *v1.WorkOrderListReq) (res *v1.WorkOrderListRes, err error) {
	res = &v1.WorkOrderListRes{}
	res.List, res.Total, err = service.WorkOrder().List(ctx, &model.WorkOrderListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		TicketNo: req.TicketNo,
		Title: req.Title,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		OrderID: req.OrderID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Priority: req.Priority,
		SourceType: req.SourceType,
		Status: req.Status,
		DueAtStart: req.DueAtStart,
		DueAtEnd: req.DueAtEnd,
	})
	return
}
// Export 导出体验工单
func (c *cWorkOrder) Export(ctx context.Context, req *v1.WorkOrderExportReq) (res *v1.WorkOrderExportRes, err error) {
	list, err := service.WorkOrder().Export(ctx, &model.WorkOrderListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		TicketNo: req.TicketNo,
		Title: req.Title,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		OrderID: req.OrderID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		Priority: req.Priority,
		SourceType: req.SourceType,
		Status: req.Status,
		DueAtStart: req.DueAtStart,
		DueAtEnd: req.DueAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="work_order.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"工单号", "客户", "商品", "订单", "工单标题", "优先级", "来源", "问题描述", "附件", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWorkOrder(item.TicketNo),
			csvSafeWorkOrder(item.CustomerName),
			csvSafeWorkOrder(item.ProductSkuNo),
			csvSafeWorkOrder(item.OrderOrderNo),
			csvSafeWorkOrder(item.Title),
			fmt.Sprintf("%v", item.Priority),
			fmt.Sprintf("%v", item.SourceType),
			csvSafeWorkOrder(item.Description),
			csvSafeWorkOrder(item.AttachmentFile),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验工单
func (c *cWorkOrder) Import(ctx context.Context, req *v1.WorkOrderImportReq) (res *v1.WorkOrderImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.WorkOrder().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WorkOrderImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验工单导入模板
func (c *cWorkOrder) ImportTemplate(ctx context.Context, req *v1.WorkOrderImportTemplateReq) (res *v1.WorkOrderImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="work_order_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"工单号", "客户", "商品", "订单", "工单标题", "优先级", "来源", "问题描述", "附件", "状态"})
	w.Flush()
	return
}
