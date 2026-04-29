package order

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/demo/api/demo/v1"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/service"
)

var Order = cOrder{}

type cOrder struct{}

// Create 创建体验订单
func (c *cOrder) Create(ctx context.Context, req *v1.OrderCreateReq) (res *v1.OrderCreateRes, err error) {
	err = service.Order().Create(ctx, &model.OrderCreateInput{
		OrderNo: req.OrderNo,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		Quantity: req.Quantity,
		Amount: req.Amount,
		PayStatus: req.PayStatus,
		DeliverStatus: req.DeliverStatus,
		PaidAt: req.PaidAt,
		DeliverAt: req.DeliverAt,
		ReceiverPhone: req.ReceiverPhone,
		Address: req.Address,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新体验订单
func (c *cOrder) Update(ctx context.Context, req *v1.OrderUpdateReq) (res *v1.OrderUpdateRes, err error) {
	err = service.Order().Update(ctx, &model.OrderUpdateInput{
		ID: req.ID,
		OrderNo: req.OrderNo,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		Quantity: req.Quantity,
		Amount: req.Amount,
		PayStatus: req.PayStatus,
		DeliverStatus: req.DeliverStatus,
		PaidAt: req.PaidAt,
		DeliverAt: req.DeliverAt,
		ReceiverPhone: req.ReceiverPhone,
		Address: req.Address,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除体验订单
func (c *cOrder) Delete(ctx context.Context, req *v1.OrderDeleteReq) (res *v1.OrderDeleteRes, err error) {
	err = service.Order().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除体验订单
func (c *cOrder) BatchDelete(ctx context.Context, req *v1.OrderBatchDeleteReq) (res *v1.OrderBatchDeleteRes, err error) {
	err = service.Order().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑体验订单
func (c *cOrder) BatchUpdate(ctx context.Context, req *v1.OrderBatchUpdateReq) (res *v1.OrderBatchUpdateRes, err error) {
	err = service.Order().BatchUpdate(ctx, &model.OrderBatchUpdateInput{
		IDs: req.IDs,
		PayStatus: req.PayStatus,
		DeliverStatus: req.DeliverStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取体验订单详情
func (c *cOrder) Detail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error) {
	res = &v1.OrderDetailRes{}
	res.OrderDetailOutput, err = service.Order().Detail(ctx, req.ID)
	return
}

// List 获取体验订单列表
func (c *cOrder) List(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error) {
	res = &v1.OrderListRes{}
	res.List, res.Total, err = service.Order().List(ctx, &model.OrderListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		OrderNo: req.OrderNo,
		ReceiverPhone: req.ReceiverPhone,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		PayStatus: req.PayStatus,
		DeliverStatus: req.DeliverStatus,
		Status: req.Status,
		PaidAtStart: req.PaidAtStart,
		PaidAtEnd: req.PaidAtEnd,
		DeliverAtStart: req.DeliverAtStart,
		DeliverAtEnd: req.DeliverAtEnd,
	})
	return
}
// Export 导出体验订单
func (c *cOrder) Export(ctx context.Context, req *v1.OrderExportReq) (res *v1.OrderExportRes, err error) {
	list, err := service.Order().Export(ctx, &model.OrderListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		OrderNo: req.OrderNo,
		ReceiverPhone: req.ReceiverPhone,
		CustomerID: req.CustomerID,
		ProductID: req.ProductID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		PayStatus: req.PayStatus,
		DeliverStatus: req.DeliverStatus,
		Status: req.Status,
		PaidAtStart: req.PaidAtStart,
		PaidAtEnd: req.PaidAtEnd,
		DeliverAtStart: req.DeliverAtStart,
		DeliverAtEnd: req.DeliverAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="order.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头
	_ = w.Write([]string{"订单号", "客户", "商品", "购买数量", "订单金额", "支付状态", "发货状态", "支付时间", "发货时间", "收货电话", "收货地址", "备注", "状态", "租户", "商户", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			item.OrderNo,
			item.CustomerName,
			item.ProductSkuNo,
			fmt.Sprintf("%v", item.Quantity),
			fmt.Sprintf("%v", item.Amount),
			fmt.Sprintf("%v", item.PayStatus),
			fmt.Sprintf("%v", item.DeliverStatus),
			func() string { if item.PaidAt != nil { return item.PaidAt.String() }; return "" }(),
			func() string { if item.DeliverAt != nil { return item.DeliverAt.String() }; return "" }(),
			item.ReceiverPhone,
			item.Address,
			item.Remark,
			fmt.Sprintf("%v", item.Status),
			item.TenantName,
			item.MerchantName,
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入体验订单
func (c *cOrder) Import(ctx context.Context, req *v1.OrderImportReq) (res *v1.OrderImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.Order().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.OrderImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载体验订单导入模板
func (c *cOrder) ImportTemplate(ctx context.Context, req *v1.OrderImportTemplateReq) (res *v1.OrderImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="order_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"订单号", "客户", "商品", "购买数量", "订单金额", "支付状态", "发货状态", "收货电话", "收货地址", "备注", "状态", "租户", "商户"})
	w.Flush()
	return
}
