package shop_order

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/service"
)

func csvSafeShopOrder(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var ShopOrder = cShopOrder{}

type cShopOrder struct{}

// Create 创建商城订单
func (c *cShopOrder) Create(ctx context.Context, req *v1.ShopOrderCreateReq) (res *v1.ShopOrderCreateRes, err error) {
	err = service.ShopOrder().Create(ctx, &model.ShopOrderCreateInput{
		OrderNo: req.OrderNo,
		UserID: req.UserID,
		GoodsID: req.GoodsID,
		GoodsTitle: req.GoodsTitle,
		GoodsCover: req.GoodsCover,
		Quantity: req.Quantity,
		TotalPrice: req.TotalPrice,
		PayWallet: req.PayWallet,
		OrderStatus: req.OrderStatus,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新商城订单
func (c *cShopOrder) Update(ctx context.Context, req *v1.ShopOrderUpdateReq) (res *v1.ShopOrderUpdateRes, err error) {
	err = service.ShopOrder().Update(ctx, &model.ShopOrderUpdateInput{
		ID: req.ID,
		OrderNo: req.OrderNo,
		UserID: req.UserID,
		GoodsID: req.GoodsID,
		GoodsTitle: req.GoodsTitle,
		GoodsCover: req.GoodsCover,
		Quantity: req.Quantity,
		TotalPrice: req.TotalPrice,
		PayWallet: req.PayWallet,
		OrderStatus: req.OrderStatus,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除商城订单
func (c *cShopOrder) Delete(ctx context.Context, req *v1.ShopOrderDeleteReq) (res *v1.ShopOrderDeleteRes, err error) {
	err = service.ShopOrder().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除商城订单
func (c *cShopOrder) BatchDelete(ctx context.Context, req *v1.ShopOrderBatchDeleteReq) (res *v1.ShopOrderBatchDeleteRes, err error) {
	err = service.ShopOrder().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑商城订单
func (c *cShopOrder) BatchUpdate(ctx context.Context, req *v1.ShopOrderBatchUpdateReq) (res *v1.ShopOrderBatchUpdateRes, err error) {
	err = service.ShopOrder().BatchUpdate(ctx, &model.ShopOrderBatchUpdateInput{
		IDs: req.IDs,
		PayWallet: req.PayWallet,
		OrderStatus: req.OrderStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取商城订单详情
func (c *cShopOrder) Detail(ctx context.Context, req *v1.ShopOrderDetailReq) (res *v1.ShopOrderDetailRes, err error) {
	res = &v1.ShopOrderDetailRes{}
	res.ShopOrderDetailOutput, err = service.ShopOrder().Detail(ctx, req.ID)
	return
}

// List 获取商城订单列表
func (c *cShopOrder) List(ctx context.Context, req *v1.ShopOrderListReq) (res *v1.ShopOrderListRes, err error) {
	res = &v1.ShopOrderListRes{}
	res.List, res.Total, err = service.ShopOrder().List(ctx, &model.ShopOrderListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		OrderNo: req.OrderNo,
		UserID: req.UserID,
		GoodsID: req.GoodsID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		PayWallet: req.PayWallet,
		OrderStatus: req.OrderStatus,
		Status: req.Status,
		GoodsTitle: req.GoodsTitle,
	})
	return
}
// Export 导出商城订单
func (c *cShopOrder) Export(ctx context.Context, req *v1.ShopOrderExportReq) (res *v1.ShopOrderExportRes, err error) {
	list, err := service.ShopOrder().Export(ctx, &model.ShopOrderListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		OrderNo: req.OrderNo,
		UserID: req.UserID,
		GoodsID: req.GoodsID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		PayWallet: req.PayWallet,
		OrderStatus: req.OrderStatus,
		Status: req.Status,
		GoodsTitle: req.GoodsTitle,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="shop_order.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"订单号", "购买会员", "商品", "商品名称", "商品封面", "购买数量", "订单总价", "支付钱包", "订单状态", "订单备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeShopOrder(item.OrderNo),
			csvSafeShopOrder(item.UserNickname),
			csvSafeShopOrder(item.ShopGoodsTitle),
			csvSafeShopOrder(item.GoodsTitle),
			csvSafeShopOrder(item.GoodsCover),
			fmt.Sprintf("%v", item.Quantity),
			fmt.Sprintf("%.2f", float64(item.TotalPrice)/100),
			fmt.Sprintf("%v", item.PayWallet),
			fmt.Sprintf("%v", item.OrderStatus),
			csvSafeShopOrder(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入商城订单
func (c *cShopOrder) Import(ctx context.Context, req *v1.ShopOrderImportReq) (res *v1.ShopOrderImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.ShopOrder().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.ShopOrderImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载商城订单导入模板
func (c *cShopOrder) ImportTemplate(ctx context.Context, req *v1.ShopOrderImportTemplateReq) (res *v1.ShopOrderImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="shop_order_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"订单号", "购买会员", "商品", "商品名称", "商品封面", "购买数量", "订单总价", "支付钱包", "订单状态", "订单备注", "状态"})
	w.Flush()
	return
}
