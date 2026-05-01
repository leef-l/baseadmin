package warehouse_goods

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

func csvSafeWarehouseGoods(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var WarehouseGoods = cWarehouseGoods{}

type cWarehouseGoods struct{}

// Create 创建仓库商品
//
// 业务规则：管理员后台手动给指定会员（owner_id）创建持有中的仓库商品。
//   - 未填 GoodsStatus 时默认 1=持有中
//   - 未填 CurrentPrice 时回退到 InitPrice，避免出现 0 价导致挂卖价计算异常
func (c *cWarehouseGoods) Create(ctx context.Context, req *v1.WarehouseGoodsCreateReq) (res *v1.WarehouseGoodsCreateRes, err error) {
	goodsStatus := req.GoodsStatus
	if goodsStatus == 0 {
		goodsStatus = 1
	}
	currentPrice := req.CurrentPrice
	if currentPrice == 0 {
		currentPrice = req.InitPrice
	}
	err = service.WarehouseGoods().Create(ctx, &model.WarehouseGoodsCreateInput{
		GoodsNo:         req.GoodsNo,
		Title:           req.Title,
		Cover:           req.Cover,
		InitPrice:       req.InitPrice,
		CurrentPrice:    currentPrice,
		PriceRiseRate:   req.PriceRiseRate,
		PlatformFeeRate: req.PlatformFeeRate,
		OwnerID:         req.OwnerID,
		TradeCount:      req.TradeCount,
		GoodsStatus:     goodsStatus,
		Remark:          req.Remark,
		Sort:            req.Sort,
		Status:          req.Status,
		TenantID:        req.TenantID,
		MerchantID:      req.MerchantID,
	})
	return
}

// Update 更新仓库商品
func (c *cWarehouseGoods) Update(ctx context.Context, req *v1.WarehouseGoodsUpdateReq) (res *v1.WarehouseGoodsUpdateRes, err error) {
	err = service.WarehouseGoods().Update(ctx, &model.WarehouseGoodsUpdateInput{
		ID: req.ID,
		GoodsNo: req.GoodsNo,
		Title: req.Title,
		Cover: req.Cover,
		InitPrice: req.InitPrice,
		CurrentPrice: req.CurrentPrice,
		PriceRiseRate: req.PriceRiseRate,
		PlatformFeeRate: req.PlatformFeeRate,
		OwnerID: req.OwnerID,
		TradeCount: req.TradeCount,
		GoodsStatus: req.GoodsStatus,
		Remark: req.Remark,
		Sort: req.Sort,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除仓库商品
func (c *cWarehouseGoods) Delete(ctx context.Context, req *v1.WarehouseGoodsDeleteReq) (res *v1.WarehouseGoodsDeleteRes, err error) {
	err = service.WarehouseGoods().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除仓库商品
func (c *cWarehouseGoods) BatchDelete(ctx context.Context, req *v1.WarehouseGoodsBatchDeleteReq) (res *v1.WarehouseGoodsBatchDeleteRes, err error) {
	err = service.WarehouseGoods().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑仓库商品
func (c *cWarehouseGoods) BatchUpdate(ctx context.Context, req *v1.WarehouseGoodsBatchUpdateReq) (res *v1.WarehouseGoodsBatchUpdateRes, err error) {
	err = service.WarehouseGoods().BatchUpdate(ctx, &model.WarehouseGoodsBatchUpdateInput{
		IDs: req.IDs,
		GoodsStatus: req.GoodsStatus,
		Status: req.Status,
	})
	return
}

// Detail 获取仓库商品详情
func (c *cWarehouseGoods) Detail(ctx context.Context, req *v1.WarehouseGoodsDetailReq) (res *v1.WarehouseGoodsDetailRes, err error) {
	res = &v1.WarehouseGoodsDetailRes{}
	res.WarehouseGoodsDetailOutput, err = service.WarehouseGoods().Detail(ctx, req.ID)
	return
}

// List 获取仓库商品列表
func (c *cWarehouseGoods) List(ctx context.Context, req *v1.WarehouseGoodsListReq) (res *v1.WarehouseGoodsListRes, err error) {
	res = &v1.WarehouseGoodsListRes{}
	res.List, res.Total, err = service.WarehouseGoods().List(ctx, &model.WarehouseGoodsListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		GoodsNo: req.GoodsNo,
		Title: req.Title,
		OwnerID: req.OwnerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		GoodsStatus: req.GoodsStatus,
		Status: req.Status,
	})
	return
}
// Export 导出仓库商品
func (c *cWarehouseGoods) Export(ctx context.Context, req *v1.WarehouseGoodsExportReq) (res *v1.WarehouseGoodsExportRes, err error) {
	list, err := service.WarehouseGoods().Export(ctx, &model.WarehouseGoodsListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Keyword: req.Keyword,
		GoodsNo: req.GoodsNo,
		Title: req.Title,
		OwnerID: req.OwnerID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		GoodsStatus: req.GoodsStatus,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_goods.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"商品编号", "商品名称", "商品封面", "初始价格", "当前价格", "每次加价比例", "平台扣除比例", "当前持有人", "流转次数", "商品状态", "备注", "排序", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeWarehouseGoods(item.GoodsNo),
			csvSafeWarehouseGoods(item.Title),
			csvSafeWarehouseGoods(item.Cover),
			fmt.Sprintf("%.2f", float64(item.InitPrice)/100),
			fmt.Sprintf("%.2f", float64(item.CurrentPrice)/100),
			fmt.Sprintf("%v", item.PriceRiseRate),
			fmt.Sprintf("%v", item.PlatformFeeRate),
			csvSafeWarehouseGoods(item.UserNickname),
			fmt.Sprintf("%v", item.TradeCount),
			fmt.Sprintf("%v", item.GoodsStatus),
			csvSafeWarehouseGoods(item.Remark),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入仓库商品
func (c *cWarehouseGoods) Import(ctx context.Context, req *v1.WarehouseGoodsImportReq) (res *v1.WarehouseGoodsImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.WarehouseGoods().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.WarehouseGoodsImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载仓库商品导入模板
func (c *cWarehouseGoods) ImportTemplate(ctx context.Context, req *v1.WarehouseGoodsImportTemplateReq) (res *v1.WarehouseGoodsImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="warehouse_goods_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"商品编号", "商品名称", "商品封面", "初始价格", "当前价格", "每次加价比例", "平台扣除比例", "当前持有人", "流转次数", "商品状态", "备注", "排序", "状态"})
	w.Flush()
	return
}
