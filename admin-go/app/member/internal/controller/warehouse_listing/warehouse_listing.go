package warehouse_listing

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

func csvSafeWarehouseListing(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var WarehouseListing = cWarehouseListing{}

type cWarehouseListing struct{}

// Create 创建仓库挂卖记录
func (c *cWarehouseListing) Create(ctx context.Context, req *v1.WarehouseListingCreateReq) (res *v1.WarehouseListingCreateRes, err error) {
	err = service.WarehouseListing().Create(ctx, &model.WarehouseListingCreateInput{
		GoodsID: req.GoodsID,
		SellerID: req.SellerID,
		ListingPrice: req.ListingPrice,
		ListingStatus: req.ListingStatus,
		ListedAt: req.ListedAt,
		SoldAt: req.SoldAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新仓库挂卖记录
func (c *cWarehouseListing) Update(ctx context.Context, req *v1.WarehouseListingUpdateReq) (res *v1.WarehouseListingUpdateRes, err error) {
	err = service.WarehouseListing().Update(ctx, &model.WarehouseListingUpdateInput{
		ID: req.ID,
		GoodsID: req.GoodsID,
		SellerID: req.SellerID,
		ListingPrice: req.ListingPrice,
		ListingStatus: req.ListingStatus,
		ListedAt: req.ListedAt,
		SoldAt: req.SoldAt,
		Remark: req.Remark,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除仓库挂卖记录
func (c *cWarehouseListing) Delete(ctx context.Context, req *v1.WarehouseListingDeleteReq) (res *v1.WarehouseListingDeleteRes, err error) {
	err = service.WarehouseListing().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除仓库挂卖记录
func (c *cWarehouseListing) BatchDelete(ctx context.Context, req *v1.WarehouseListingBatchDeleteReq) (res *v1.WarehouseListingBatchDeleteRes, err error) {
	err = service.WarehouseListing().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑仓库挂卖记录
func (c *cWarehouseListing) BatchUpdate(ctx context.Context, req *v1.WarehouseListingBatchUpdateReq) (res *v1.WarehouseListingBatchUpdateRes, err error) {
	err = service.WarehouseListing().BatchUpdate(ctx, &model.WarehouseListingBatchUpdateInput{
		IDs: req.IDs,
		ListingStatus: req.ListingStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取仓库挂卖记录详情
func (c *cWarehouseListing) Detail(ctx context.Context, req *v1.WarehouseListingDetailReq) (res *v1.WarehouseListingDetailRes, err error) {
	res = &v1.WarehouseListingDetailRes{}
	res.WarehouseListingDetailOutput, err = service.WarehouseListing().Detail(ctx, req.ID)
	return
}

// List 获取仓库挂卖记录列表
func (c *cWarehouseListing) List(ctx context.Context, req *v1.WarehouseListingListReq) (res *v1.WarehouseListingListRes, err error) {
	res = &v1.WarehouseListingListRes{}
	res.List, res.Total, err = service.WarehouseListing().List(ctx, &model.WarehouseListingListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		GoodsID: req.GoodsID,
		SellerID: req.SellerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ListingStatus: req.ListingStatus,
		Status: req.Status,
		ListedAtStart: req.ListedAtStart,
		ListedAtEnd: req.ListedAtEnd,
		SoldAtStart: req.SoldAtStart,
		SoldAtEnd: req.SoldAtEnd,
	})
	return
}
// Export 导出仓库挂卖记录
func (c *cWarehouseListing) Export(ctx context.Context, req *v1.WarehouseListingExportReq) (res *v1.WarehouseListingExportRes, err error) {
	list, err := service.WarehouseListing().Export(ctx, &model.WarehouseListingListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		GoodsID: req.GoodsID,
		SellerID: req.SellerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		ListingStatus: req.ListingStatus,
		Status: req.Status,
		ListedAtStart: req.ListedAtStart,
		ListedAtEnd: req.ListedAtEnd,
		SoldAtStart: req.SoldAtStart,
		SoldAtEnd: req.SoldAtEnd,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_listing.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"仓库商品", "卖家", "挂卖价格", "挂卖状态", "备注", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWarehouseListing(item.WarehouseGoodsTitle),
			csvSafeWarehouseListing(item.UserNickname),
			fmt.Sprintf("%.2f", float64(item.ListingPrice)/100),
			fmt.Sprintf("%v", item.ListingStatus),
			csvSafeWarehouseListing(item.Remark),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入仓库挂卖记录
func (c *cWarehouseListing) Import(ctx context.Context, req *v1.WarehouseListingImportReq) (res *v1.WarehouseListingImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.WarehouseListing().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WarehouseListingImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载仓库挂卖记录导入模板
func (c *cWarehouseListing) ImportTemplate(ctx context.Context, req *v1.WarehouseListingImportTemplateReq) (res *v1.WarehouseListingImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_listing_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"仓库商品", "卖家", "挂卖价格", "挂卖状态", "备注", "状态"})
	w.Flush()
	return
}
