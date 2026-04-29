package shop_goods

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

func csvSafeShopGoods(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == 0 || strings.ContainsAny(s[:1], "=+-@\t\r") {
		return "'" + s
	}
	return s
}

var ShopGoods = cShopGoods{}

type cShopGoods struct{}

// Create 创建商城商品
func (c *cShopGoods) Create(ctx context.Context, req *v1.ShopGoodsCreateReq) (res *v1.ShopGoodsCreateRes, err error) {
	err = service.ShopGoods().Create(ctx, &model.ShopGoodsCreateInput{
		CategoryID: req.CategoryID,
		Title: req.Title,
		Cover: req.Cover,
		Images: req.Images,
		Price: req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock: req.Stock,
		Sales: req.Sales,
		Content: req.Content,
		Sort: req.Sort,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Update 更新商城商品
func (c *cShopGoods) Update(ctx context.Context, req *v1.ShopGoodsUpdateReq) (res *v1.ShopGoodsUpdateRes, err error) {
	err = service.ShopGoods().Update(ctx, &model.ShopGoodsUpdateInput{
		ID: req.ID,
		CategoryID: req.CategoryID,
		Title: req.Title,
		Cover: req.Cover,
		Images: req.Images,
		Price: req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock: req.Stock,
		Sales: req.Sales,
		Content: req.Content,
		Sort: req.Sort,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
	})
	return
}

// Delete 删除商城商品
func (c *cShopGoods) Delete(ctx context.Context, req *v1.ShopGoodsDeleteReq) (res *v1.ShopGoodsDeleteRes, err error) {
	err = service.ShopGoods().Delete(ctx, req.ID)
	return
}

// BatchDelete 批量删除商城商品
func (c *cShopGoods) BatchDelete(ctx context.Context, req *v1.ShopGoodsBatchDeleteReq) (res *v1.ShopGoodsBatchDeleteRes, err error) {
	err = service.ShopGoods().BatchDelete(ctx, req.IDs)
	return
}

// BatchUpdate 批量编辑商城商品
func (c *cShopGoods) BatchUpdate(ctx context.Context, req *v1.ShopGoodsBatchUpdateReq) (res *v1.ShopGoodsBatchUpdateRes, err error) {
	err = service.ShopGoods().BatchUpdate(ctx, &model.ShopGoodsBatchUpdateInput{
		IDs: req.IDs,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	return
}

// Detail 获取商城商品详情
func (c *cShopGoods) Detail(ctx context.Context, req *v1.ShopGoodsDetailReq) (res *v1.ShopGoodsDetailRes, err error) {
	res = &v1.ShopGoodsDetailRes{}
	res.ShopGoodsDetailOutput, err = service.ShopGoods().Detail(ctx, req.ID)
	return
}

// List 获取商城商品列表
func (c *cShopGoods) List(ctx context.Context, req *v1.ShopGoodsListReq) (res *v1.ShopGoodsListRes, err error) {
	res = &v1.ShopGoodsListRes{}
	res.List, res.Total, err = service.ShopGoods().List(ctx, &model.ShopGoodsListInput{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title: req.Title,
		CategoryID: req.CategoryID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	return
}
// Export 导出商城商品
func (c *cShopGoods) Export(ctx context.Context, req *v1.ShopGoodsExportReq) (res *v1.ShopGoodsExportRes, err error) {
	list, err := service.ShopGoods().Export(ctx, &model.ShopGoodsListInput{
		OrderBy:   req.OrderBy,
		OrderDir:  req.OrderDir,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Title: req.Title,
		CategoryID: req.CategoryID,
		TenantID: req.TenantID,
		MerchantID: req.MerchantID,
		IsRecommend: req.IsRecommend,
		Status: req.Status,
	})
	if err != nil {
		return
	}
	// CSV 导出（使用 csv.Writer 防止注入和格式问题）
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="shop_goods.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	// 表头（与导入模板列对齐，末尾追加只读列）
	_ = w.Write([]string{"商品分类", "商品名称", "封面图", "商品图片", "售价", "原价", "库存", "销量", "商品详情", "排序", "是否推荐", "状态", "创建时间"})
	// 数据行
	for _, item := range list {
		_ = w.Write([]string{
			csvSafeShopGoods(item.ShopCategoryName),
			csvSafeShopGoods(item.Title),
			csvSafeShopGoods(item.Cover),
			csvSafeShopGoods(item.Images),
			fmt.Sprintf("%.2f", float64(item.Price)/100),
			fmt.Sprintf("%.2f", float64(item.OriginalPrice)/100),
			fmt.Sprintf("%v", item.Stock),
			fmt.Sprintf("%v", item.Sales),
			csvSafeShopGoods(item.Content),
			fmt.Sprintf("%v", item.Sort),
			fmt.Sprintf("%v", item.IsRecommend),
			fmt.Sprintf("%v", item.Status),
			func() string { if item.CreatedAt != nil { return item.CreatedAt.String() }; return "" }(),
		})
	}
	w.Flush()
	return
}

// Import 导入商城商品
func (c *cShopGoods) Import(ctx context.Context, req *v1.ShopGoodsImportReq) (res *v1.ShopGoodsImportRes, err error) {
	r := g.RequestFromCtx(ctx)
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, fmt.Errorf("请上传文件")
	}
	success, fail, err := service.ShopGoods().Import(ctx, file)
	if err != nil {
		return nil, err
	}
	res = &v1.ShopGoodsImportRes{Success: success, Fail: fail}
	return
}

// ImportTemplate 下载商城商品导入模板
func (c *cShopGoods) ImportTemplate(ctx context.Context, req *v1.ShopGoodsImportTemplateReq) (res *v1.ShopGoodsImportTemplateRes, err error) {
	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Response.Header().Set("Content-Disposition", `attachment; filename="shop_goods_template.csv"`)
	r.Response.Write("\xEF\xBB\xBF") // UTF-8 BOM
	w := csv.NewWriter(r.Response.Writer)
	_ = w.Write([]string{"商品分类", "商品名称", "封面图", "商品图片", "售价", "原价", "库存", "销量", "商品详情", "排序", "是否推荐", "状态"})
	w.Flush()
	return
}
