package warehouse_trade

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

func csvSafeWarehouseTrade(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var WarehouseTrade = cWarehouseTrade{}

type cWarehouseTrade struct{}

// Create 创建仓库交易记录
func (c *cWarehouseTrade) Create(ctx context.Context, req *v1.WarehouseTradeCreateReq) (res *v1.WarehouseTradeCreateRes, err error) {
	err = service.WarehouseTrade().Create(ctx, &model.WarehouseTradeCreateInput{
		TradeNo: req.TradeNo,
		GoodsID: req.GoodsID,
		ListingID: req.ListingID,
		SellerID: req.SellerID,
		BuyerID: req.BuyerID,
		TradePrice: req.TradePrice,
		PlatformFee: req.PlatformFee,
		SellerIncome: req.SellerIncome,
		TradeStatus: req.TradeStatus,
		ConfirmedAt: req.ConfirmedAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新仓库交易记录
func (c *cWarehouseTrade) Update(ctx context.Context, req *v1.WarehouseTradeUpdateReq) (res *v1.WarehouseTradeUpdateRes, err error) {
	err = service.WarehouseTrade().Update(ctx, &model.WarehouseTradeUpdateInput{
		ID: req.ID,
		TradeNo: req.TradeNo,
		GoodsID: req.GoodsID,
		ListingID: req.ListingID,
		SellerID: req.SellerID,
		BuyerID: req.BuyerID,
		TradePrice: req.TradePrice,
		PlatformFee: req.PlatformFee,
		SellerIncome: req.SellerIncome,
		TradeStatus: req.TradeStatus,
		ConfirmedAt: req.ConfirmedAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除仓库交易记录
func (c *cWarehouseTrade) Delete(ctx context.Context, req *v1.WarehouseTradeDeleteReq) (res *v1.WarehouseTradeDeleteRes, err error) {
	err = service.WarehouseTrade().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除仓库交易记录
func (c *cWarehouseTrade) BatchDelete(ctx context.Context, req *v1.WarehouseTradeBatchDeleteReq) (res *v1.WarehouseTradeBatchDeleteRes, err error) {
	err = service.WarehouseTrade().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑仓库交易记录
func (c *cWarehouseTrade) BatchUpdate(ctx context.Context, req *v1.WarehouseTradeBatchUpdateReq) (res *v1.WarehouseTradeBatchUpdateRes, err error) {
	err = service.WarehouseTrade().BatchUpdate(ctx, &model.WarehouseTradeBatchUpdateInput{
		IDs: req.IDs,
		TradeStatus: req.TradeStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取仓库交易记录详情
func (c *cWarehouseTrade) Detail(ctx context.Context, req *v1.WarehouseTradeDetailReq) (res *v1.WarehouseTradeDetailRes, err error) {
	res = &v1.WarehouseTradeDetailRes{}
	res.WarehouseTradeDetailOutput, err = service.WarehouseTrade().Detail(ctx, req.ID)
	return
}

// List 获取仓库交易记录列表
func (c *cWarehouseTrade) List(ctx context.Context, req *v1.WarehouseTradeListReq) (res *v1.WarehouseTradeListRes, err error) {
	res = &v1.WarehouseTradeListRes{}
	res.List, res.Total, err = service.WarehouseTrade().List(ctx, &model.WarehouseTradeListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TradeNo: req.TradeNo,
		GoodsID: req.GoodsID,
		ListingID: req.ListingID,
		SellerID: req.SellerID,
		BuyerID: req.BuyerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		TradeStatus: req.TradeStatus,
		Status: req.Status,
		ConfirmedAtStart: req.ConfirmedAtStart,
		ConfirmedAtEnd: req.ConfirmedAtEnd,
	})
	return
}
// Export 导出仓库交易记录
func (c *cWarehouseTrade) Export(ctx context.Context, req *v1.WarehouseTradeExportReq) (res *v1.WarehouseTradeExportRes, err error) {
	list, err := service.WarehouseTrade().Export(ctx, &model.WarehouseTradeListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		TradeNo: req.TradeNo,
		GoodsID: req.GoodsID,
		ListingID: req.ListingID,
		SellerID: req.SellerID,
		BuyerID: req.BuyerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		TradeStatus: req.TradeStatus,
		Status: req.Status,
		ConfirmedAtStart: req.ConfirmedAtStart,
		ConfirmedAtEnd: req.ConfirmedAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_trade.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"交易编号", "仓库商品", "挂卖记录", "卖家", "买家", "成交价格", "平台扣除费用", "卖家实收", "交易状态", "备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWarehouseTrade(item.TradeNo),
			csvSafeWarehouseTrade(item.WarehouseGoodsTitle),
			csvSafeWarehouseTrade(item.WarehouseListingID),
			csvSafeWarehouseTrade(item.UserNickname),
			csvSafeWarehouseTrade(item.BuyerNickname),
			fmt.Sprintf("%.2f", float64(item.TradePrice)/100),
			fmt.Sprintf("%.2f", float64(item.PlatformFee)/100),
			fmt.Sprintf("%.2f", float64(item.SellerIncome)/100),
			fmt.Sprintf("%v", item.TradeStatus),
			csvSafeWarehouseTrade(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入仓库交易记录
func (c *cWarehouseTrade) Import(ctx context.Context, req *v1.WarehouseTradeImportReq) (res *v1.WarehouseTradeImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.WarehouseTrade().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WarehouseTradeImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载仓库交易记录导入模板
func (c *cWarehouseTrade) ImportTemplate(ctx context.Context, req *v1.WarehouseTradeImportTemplateReq) (res *v1.WarehouseTradeImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_trade_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"交易编号", "仓库商品", "挂卖记录", "卖家", "买家", "成交价格", "平台扣除费用", "卖家实收", "交易状态", "备注", "状态"})
	w.Flush()
	return
}
