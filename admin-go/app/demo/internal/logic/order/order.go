
package order

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
	service.RegisterOrder(New())
}

func New() *sOrder {
	return &sOrder{}
}

type sOrder struct{}

func normalizeOrderIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建体验订单
func (s *sOrder) Create(ctx context.Context, in *model.OrderCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoOrder.Ctx(ctx).Data(do.DemoOrder{
		Id:        id,
		OrderNo: in.OrderNo,
		CustomerId: in.CustomerID,
		ProductId: in.ProductID,
		Quantity: in.Quantity,
		Amount: in.Amount,
		PayStatus: in.PayStatus,
		DeliverStatus: in.DeliverStatus,
		PaidAt: in.PaidAt,
		DeliverAt: in.DeliverAt,
		ReceiverPhone: in.ReceiverPhone,
		Address: in.Address,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验订单
func (s *sOrder) Update(ctx context.Context, in *model.OrderUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.DemoOrder{
		OrderNo: in.OrderNo,
		CustomerId: in.CustomerID,
		ProductId: in.ProductID,
		Quantity: in.Quantity,
		Amount: in.Amount,
		PayStatus: in.PayStatus,
		DeliverStatus: in.DeliverStatus,
		PaidAt: in.PaidAt,
		DeliverAt: in.DeliverAt,
		ReceiverPhone: in.ReceiverPhone,
		Address: in.Address,
		Remark: in.Remark,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.DemoOrder.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.DemoOrder.Table()).Ctx(ctx).
			Where(dao.DemoOrder.Columns().Id, in.ID).
			Where(dao.DemoOrder.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("体验订单不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.DemoOrder.Table()).Ctx(ctx), in.ID, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId, "体验订单"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.DemoOrder.Table()).Ctx(ctx), in.ID, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.DemoOrder.Table()).Ctx(ctx).
			Where(dao.DemoOrder.Columns().Id, in.ID).
			Where(dao.DemoOrder.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除体验订单
func (s *sOrder) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoOrder.Ctx(ctx), id, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId, "体验订单"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoOrder.Ctx(ctx), id, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoOrder.Ctx(ctx).Where(dao.DemoOrder.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验订单
func (s *sOrder) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeOrderIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoOrder.Ctx(ctx), normalizedIDs, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId, "体验订单"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoOrder.Ctx(ctx), normalizedIDs, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoOrder.Ctx(ctx).WhereIn(dao.DemoOrder.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验订单详情
func (s *sOrder) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.OrderDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoOrder.Ctx(ctx), id, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId, "体验订单"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoOrder.Ctx(ctx), id, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.OrderDetailOutput{}
	err = dao.DemoOrder.Ctx(ctx).Where(dao.DemoOrder.Columns().Id, id).Where(dao.DemoOrder.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("体验订单不存在或已删除")
	}
	// 查询客户关联显示
	if out.CustomerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", out.CustomerID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.CustomerName = val.String()
		}
	}
	// 查询商品关联显示
	if out.ProductID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_product").Where("id", out.ProductID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("sku_no")
		if err == nil {
			out.ProductSkuNo = val.String()
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
func (s *sOrder) applyListFilter(ctx context.Context, in *model.OrderListInput) *gdb.Model {
	m := dao.DemoOrder.Ctx(ctx).Where(dao.DemoOrder.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.DemoOrder.Columns().ReceiverPhone, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoOrder.Columns().Address, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoOrder.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.OrderNo != "" {
		m = m.Where(dao.DemoOrder.Columns().OrderNo, in.OrderNo)
	}
	if in.ReceiverPhone != "" {
		m = m.WhereLike(dao.DemoOrder.Columns().ReceiverPhone, "%"+in.ReceiverPhone+"%")
	}
	if in.CustomerID != nil {
		m = m.Where(dao.DemoOrder.Columns().CustomerId, *in.CustomerID)
	}
	if in.ProductID != nil {
		m = m.Where(dao.DemoOrder.Columns().ProductId, *in.ProductID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoOrder.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoOrder.Columns().MerchantId, *in.MerchantID)
	}
	if in.PayStatus != nil {
		m = m.Where(dao.DemoOrder.Columns().PayStatus, *in.PayStatus)
	}
	if in.DeliverStatus != nil {
		m = m.Where(dao.DemoOrder.Columns().DeliverStatus, *in.DeliverStatus)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoOrder.Columns().Status, *in.Status)
	}
	if in.PaidAtStart != "" {
		m = m.WhereGTE(dao.DemoOrder.Columns().PaidAt, in.PaidAtStart)
	}
	if in.PaidAtEnd != "" {
		m = m.WhereLTE(dao.DemoOrder.Columns().PaidAt, in.PaidAtEnd)
	}
	if in.DeliverAtStart != "" {
		m = m.WhereGTE(dao.DemoOrder.Columns().DeliverAt, in.DeliverAtStart)
	}
	if in.DeliverAtEnd != "" {
		m = m.WhereLTE(dao.DemoOrder.Columns().DeliverAt, in.DeliverAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoOrder.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoOrder.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sOrder) fillRefFields(ctx context.Context, list []*model.OrderListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.CustomerID != 0 {
				idSet[int64(item.CustomerID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_customer").
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
					if val, ok := refMap[int64(item.CustomerID)]; ok {
						item.CustomerName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.ProductID != 0 {
				idSet[int64(item.ProductID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_product").
				Fields("id", "sku_no")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["sku_no"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.ProductID)]; ok {
						item.ProductSkuNo = val
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

// List 获取体验订单列表
func (s *sOrder) List(ctx context.Context, in *model.OrderListInput) (list []*model.OrderListOutput, total int, err error) {
	if in == nil {
		in = &model.OrderListInput{}
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
func (s *sOrder) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoOrder.Columns().Id:        true,
		dao.DemoOrder.Columns().CreatedAt: true,
		dao.DemoOrder.Columns().Status:    true,
		dao.DemoOrder.Columns().OrderNo: true,
		dao.DemoOrder.Columns().Amount: true,
		dao.DemoOrder.Columns().ReceiverPhone: true,
		dao.DemoOrder.Columns().Address: true,
		dao.DemoOrder.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sOrder) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoOrder.Columns().Id)
}

// Export 导出体验订单（不分页）
func (s *sOrder) Export(ctx context.Context, in *model.OrderListInput) (list []*model.OrderListOutput, err error) {
	if in == nil {
		in = &model.OrderListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验订单
func (s *sOrder) BatchUpdate(ctx context.Context, in *model.OrderBatchUpdateInput) error {
	data := do.DemoOrder{}
	hasChange := false
	if in.PayStatus != nil {
		data.PayStatus = *in.PayStatus
		hasChange = true
	}
	if in.DeliverStatus != nil {
		data.DeliverStatus = *in.DeliverStatus
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeOrderIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoOrder.Ctx(ctx), normalizedIDs, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().TenantId, dao.DemoOrder.Columns().MerchantId, "体验订单"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoOrder.Ctx(ctx), normalizedIDs, dao.DemoOrder.Columns().Id, dao.DemoOrder.Columns().CreatedBy, dao.DemoOrder.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoOrder.Ctx(ctx).WhereIn(dao.DemoOrder.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验订单
func (s *sOrder) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.DemoOrder{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.OrderNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.CustomerId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.ProductId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Quantity = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.Amount = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.PayStatus = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.DeliverStatus = v
			}
		}
		idx++
		if idx < len(record) {
			data.ReceiverPhone = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Address = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Remark = strings.TrimSpace(record[idx])
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
		if fkVal, ok := data.CustomerId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.ProductId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("demo_product").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.DemoOrder.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

