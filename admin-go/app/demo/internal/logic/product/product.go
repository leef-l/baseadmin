
package product

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"strconv"
	"strings"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/demo/internal/dao"
	"gbaseadmin/app/demo/internal/middleware"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/model/do"
	"gbaseadmin/app/demo/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterProduct(New())
}

func New() *sProduct {
	return &sProduct{}
}

type sProduct struct{}

func normalizeProductIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	normalized := make([]snowflake.JsonInt64, 0, len(ids))
	for _, id := range ids {
		value := int64(id)
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, id)
		if len(normalized) >= 500 {
			break
		}
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

// Create 创建体验商品
func (s *sProduct) Create(ctx context.Context, in *model.ProductCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoProduct.Ctx(ctx).Data(do.DemoProduct{
		Id:        id,
		CategoryId: in.CategoryID,
		SkuNo: in.SkuNo,
		Name: in.Name,
		Cover: in.Cover,
		ManualFile: in.ManualFile,
		DetailContent: in.DetailContent,
		SpecJson: in.SpecJSON,
		WebsiteUrl: in.WebsiteURL,
		Type: in.Type,
		IsRecommend: in.IsRecommend,
		SalePrice: in.SalePrice,
		StockNum: in.StockNum,
		WeightNum: in.WeightNum,
		Sort: in.Sort,
		Icon: in.Icon,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验商品
func (s *sProduct) Update(ctx context.Context, in *model.ProductUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.DemoProduct{
		CategoryId: in.CategoryID,
		SkuNo: in.SkuNo,
		Name: in.Name,
		Cover: in.Cover,
		ManualFile: in.ManualFile,
		DetailContent: in.DetailContent,
		SpecJson: in.SpecJSON,
		WebsiteUrl: in.WebsiteURL,
		Type: in.Type,
		IsRecommend: in.IsRecommend,
		SalePrice: in.SalePrice,
		StockNum: in.StockNum,
		WeightNum: in.WeightNum,
		Sort: in.Sort,
		Icon: in.Icon,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.DemoProduct.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.DemoProduct.Table()).Ctx(ctx).
			Where(dao.DemoProduct.Columns().Id, in.ID).
			Where(dao.DemoProduct.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("体验商品不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.DemoProduct.Table()).Ctx(ctx), in.ID, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId, "体验商品"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.DemoProduct.Table()).Ctx(ctx), in.ID, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.DemoProduct.Table()).Ctx(ctx).
			Where(dao.DemoProduct.Columns().Id, in.ID).
			Where(dao.DemoProduct.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除体验商品
func (s *sProduct) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoProduct.Ctx(ctx), id, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId, "体验商品"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoProduct.Ctx(ctx), id, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoProduct.Ctx(ctx).Where(dao.DemoProduct.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验商品
func (s *sProduct) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeProductIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoProduct.Ctx(ctx), normalizedIDs, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId, "体验商品"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoProduct.Ctx(ctx), normalizedIDs, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoProduct.Ctx(ctx).WhereIn(dao.DemoProduct.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验商品详情
func (s *sProduct) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ProductDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoProduct.Ctx(ctx), id, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId, "体验商品"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoProduct.Ctx(ctx), id, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.ProductDetailOutput{}
	err = dao.DemoProduct.Ctx(ctx).Where(dao.DemoProduct.Columns().Id, id).Where(dao.DemoProduct.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("体验商品不存在或已删除")
	}
	// 查询商品分类关联显示
	if out.CategoryID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_category").Where("id", out.CategoryID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.CategoryName = val.String()
		}
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sProduct) applyListFilter(ctx context.Context, in *model.ProductListInput) *gdb.Model {
	m := dao.DemoProduct.Ctx(ctx).Where(dao.DemoProduct.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId)
	if in.SkuNo != "" {
		m = m.Where(dao.DemoProduct.Columns().SkuNo, in.SkuNo)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.DemoProduct.Columns().Name, "%"+in.Name+"%")
	}
	if in.CategoryID != nil {
		m = m.Where(dao.DemoProduct.Columns().CategoryId, *in.CategoryID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoProduct.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoProduct.Columns().MerchantId, *in.MerchantID)
	}
	if in.Type != nil {
		m = m.Where(dao.DemoProduct.Columns().Type, *in.Type)
	}
	if in.IsRecommend != nil {
		m = m.Where(dao.DemoProduct.Columns().IsRecommend, *in.IsRecommend)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoProduct.Columns().Status, *in.Status)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoProduct.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoProduct.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sProduct) fillRefFields(ctx context.Context, list []*model.ProductListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.CategoryID != 0 {
				idSet[int64(item.CategoryID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_category").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.CategoryID)]; ok {
						item.CategoryName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.TenantID != 0 {
				idSet[int64(item.TenantID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_tenant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.TenantID)]; ok {
						item.TenantName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.MerchantID != 0 {
				idSet[int64(item.MerchantID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_merchant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.MerchantID)]; ok {
						item.MerchantName = val
					}
				}
			}
		}
	}
}

// List 获取体验商品列表
func (s *sProduct) List(ctx context.Context, in *model.ProductListInput) (list []*model.ProductListOutput, total int, err error) {
	if in == nil {
		in = &model.ProductListInput{}
	}
	// PageSize 上限保护
	if in.PageSize <= 0 {
		in.PageSize = 10
	} else if in.PageSize > 500 {
		in.PageSize = 500
	}
	if in.PageNum <= 0 {
		in.PageNum = 1
	}
	m := s.applyListFilter(ctx, in)
	total, err = m.Count()
	if err != nil {
		return
	}
	// 动态排序（白名单校验防止 SQL 注入）
	m = s.applyListOrder(m, in.OrderBy, in.OrderDir)
	err = m.Page(in.PageNum, in.PageSize).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// isAllowedOrderField 校验排序字段是否在允许列表中
func (s *sProduct) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoProduct.Columns().Id:        true,
		dao.DemoProduct.Columns().CreatedAt: true,
		dao.DemoProduct.Columns().Sort:      true,
		dao.DemoProduct.Columns().Status:    true,
		dao.DemoProduct.Columns().SkuNo: true,
		dao.DemoProduct.Columns().Name: true,
		dao.DemoProduct.Columns().SalePrice: true,
	}
	return allowed[field]
}

func (s *sProduct) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.DemoProduct.Columns().Sort).OrderDesc(dao.DemoProduct.Columns().Id)
}

// Export 导出体验商品（不分页）
func (s *sProduct) Export(ctx context.Context, in *model.ProductListInput) (list []*model.ProductListOutput, err error) {
	if in == nil {
		in = &model.ProductListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验商品
func (s *sProduct) BatchUpdate(ctx context.Context, in *model.ProductBatchUpdateInput) error {
	data := do.DemoProduct{}
	hasChange := false
	if in.Type != nil {
		data.Type = *in.Type
		hasChange = true
	}
	if in.IsRecommend != nil {
		data.IsRecommend = *in.IsRecommend
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeProductIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoProduct.Ctx(ctx), normalizedIDs, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().TenantId, dao.DemoProduct.Columns().MerchantId, "体验商品"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoProduct.Ctx(ctx), normalizedIDs, dao.DemoProduct.Columns().Id, dao.DemoProduct.Columns().CreatedBy, dao.DemoProduct.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoProduct.Ctx(ctx).WhereIn(dao.DemoProduct.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验商品
func (s *sProduct) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
	const maxImportFileSize = 10 << 20 // 10MB
	const maxImportRows = 5000
	if file.Size > maxImportFileSize {
		return 0, 0, fmt.Errorf("文件大小超过限制（最大10MB）")
	}
	f, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	// 跳过表头
	if _, err = reader.Read(); err != nil {
		return 0, 0, fmt.Errorf("读取CSV表头失败: %w", err)
	}

	rowCount := 0
	for {
		record, readErr := reader.Read()
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return success, fail, fmt.Errorf("读取CSV数据失败: %w", readErr)
		}
		if len(record) == 0 {
			continue
		}
		rowCount++
		if rowCount > maxImportRows {
			return success, fail, fmt.Errorf("导入数据超过 %d 行上限，已处理 %d 条成功、%d 条失败", maxImportRows, success, fail)
		}
		// 逐行插入
		id := snowflake.Generate()
		data := do.DemoProduct{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.CategoryId = v
			}
		}
		idx++
		if idx < len(record) {
			data.SkuNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Name = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Cover = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.ManualFile = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.DetailContent = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.SpecJson = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.WebsiteUrl = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Type = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.IsRecommend = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.SalePrice = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.StockNum = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.WeightNum = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Sort = v
			}
		}
		idx++
		if idx < len(record) {
			data.Icon = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Status = v
			}
		}
		idx++
		tenantID := snowflake.JsonInt64(0)
		merchantID := snowflake.JsonInt64(0)
		middleware.ApplyTenantScopeToWrite(ctx, &tenantID, &merchantID)
		if err := middleware.EnsureTenantMerchantAccessible(ctx, tenantID, merchantID); err != nil {
			fail++
			continue
		}
		data.TenantId = tenantID
		data.MerchantId = merchantID
		if fkVal, ok := data.CategoryId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("demo_category").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.DemoProduct.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

