
package work_order

import (
	"context"
	"encoding/csv"
	"io"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
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
	service.RegisterWorkOrder(New())
}

func New() *sWorkOrder {
	return &sWorkOrder{}
}

type sWorkOrder struct{}

func normalizeWorkOrderIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

// Create 创建体验工单
func (s *sWorkOrder) Create(ctx context.Context, in *model.WorkOrderCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoWorkOrder.Ctx(ctx).Data(do.DemoWorkOrder{
		Id:        id,
		TicketNo: in.TicketNo,
		CustomerId: in.CustomerID,
		ProductId: in.ProductID,
		OrderId: in.OrderID,
		Title: in.Title,
		Priority: in.Priority,
		SourceType: in.SourceType,
		Description: in.Description,
		AttachmentFile: in.AttachmentFile,
		DueAt: in.DueAt,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验工单
func (s *sWorkOrder) Update(ctx context.Context, in *model.WorkOrderUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoWorkOrder.Ctx(ctx), in.ID, dao.DemoWorkOrder.Columns().Id, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId, "体验工单"); err != nil {
		return err
	}
	data := do.DemoWorkOrder{
		TicketNo: in.TicketNo,
		CustomerId: in.CustomerID,
		ProductId: in.ProductID,
		OrderId: in.OrderID,
		Title: in.Title,
		Priority: in.Priority,
		SourceType: in.SourceType,
		Description: in.Description,
		AttachmentFile: in.AttachmentFile,
		DueAt: in.DueAt,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
	}
	_, err := dao.DemoWorkOrder.Ctx(ctx).Where(dao.DemoWorkOrder.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 软删除体验工单
func (s *sWorkOrder) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoWorkOrder.Ctx(ctx), id, dao.DemoWorkOrder.Columns().Id, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId, "体验工单"); err != nil {
		return err
	}
	_, err := dao.DemoWorkOrder.Ctx(ctx).Where(dao.DemoWorkOrder.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验工单
func (s *sWorkOrder) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeWorkOrderIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoWorkOrder.Ctx(ctx), normalizedIDs, dao.DemoWorkOrder.Columns().Id, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId, "体验工单"); err != nil {
		return err
	}
	_, err := dao.DemoWorkOrder.Ctx(ctx).WhereIn(dao.DemoWorkOrder.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验工单详情
func (s *sWorkOrder) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.WorkOrderDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoWorkOrder.Ctx(ctx), id, dao.DemoWorkOrder.Columns().Id, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId, "体验工单"); err != nil {
		return nil, err
	}
	out = &model.WorkOrderDetailOutput{}
	err = dao.DemoWorkOrder.Ctx(ctx).Where(dao.DemoWorkOrder.Columns().Id, id).Where(dao.DemoWorkOrder.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out.ID == 0 {
		return nil, nil
	}
	// 查询客户关联显示
	if out.CustomerID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_customer").Where("id", out.CustomerID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.CustomerName = val.String()
		}
	}
	// 查询商品关联显示
	if out.ProductID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_product").Where("id", out.ProductID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("sku_no")
		if err == nil {
			out.ProductSkuNo = val.String()
		}
	}
	// 查询订单关联显示
	if out.OrderID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("demo_order").Where("id", out.OrderID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("order_no")
		if err == nil {
			out.OrderOrderNo = val.String()
		}
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		refQuery = refQuery.Where("deleted_at", nil)
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sWorkOrder) applyListFilter(ctx context.Context, in *model.WorkOrderListInput) *gdb.Model {
	m := dao.DemoWorkOrder.Ctx(ctx).Where(dao.DemoWorkOrder.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.DemoWorkOrder.Columns().Title, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.DemoWorkOrder.Columns().Description, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.TicketNo != "" {
		m = m.Where(dao.DemoWorkOrder.Columns().TicketNo, in.TicketNo)
	}
	if in.Title != "" {
		m = m.WhereLike(dao.DemoWorkOrder.Columns().Title, "%"+in.Title+"%")
	}
	if in.CustomerID != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().CustomerId, *in.CustomerID)
	}
	if in.ProductID != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().ProductId, *in.ProductID)
	}
	if in.OrderID != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().OrderId, *in.OrderID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().MerchantId, *in.MerchantID)
	}
	if in.Priority != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().Priority, *in.Priority)
	}
	if in.SourceType != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().SourceType, *in.SourceType)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoWorkOrder.Columns().Status, *in.Status)
	}
	if in.DueAtStart != "" {
		m = m.WhereGTE(dao.DemoWorkOrder.Columns().DueAt, in.DueAtStart)
	}
	if in.DueAtEnd != "" {
		m = m.WhereLTE(dao.DemoWorkOrder.Columns().DueAt, in.DueAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoWorkOrder.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoWorkOrder.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoWorkOrder.Columns().CreatedBy, dao.DemoWorkOrder.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sWorkOrder) fillRefFields(ctx context.Context, list []*model.WorkOrderListOutput) {
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
			if item.OrderID != 0 {
				idSet[int64(item.OrderID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("demo_order").
				Fields("id", "order_no")
			refQuery = refQuery.Where("deleted_at", nil)
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["order_no"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.OrderID)]; ok {
						item.OrderOrderNo = val
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

// List 获取体验工单列表
func (s *sWorkOrder) List(ctx context.Context, in *model.WorkOrderListInput) (list []*model.WorkOrderListOutput, total int, err error) {
	if in == nil {
		in = &model.WorkOrderListInput{}
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
func (s *sWorkOrder) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoWorkOrder.Columns().Id:        true,
		dao.DemoWorkOrder.Columns().CreatedAt: true,
		dao.DemoWorkOrder.Columns().Status:    true,
		dao.DemoWorkOrder.Columns().TicketNo: true,
		dao.DemoWorkOrder.Columns().Title: true,
		dao.DemoWorkOrder.Columns().Description: true,
	}
	return allowed[field]
}

func (s *sWorkOrder) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoWorkOrder.Columns().Id)
}

// Export 导出体验工单（不分页）
func (s *sWorkOrder) Export(ctx context.Context, in *model.WorkOrderListInput) (list []*model.WorkOrderListOutput, err error) {
	if in == nil {
		in = &model.WorkOrderListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验工单
func (s *sWorkOrder) BatchUpdate(ctx context.Context, in *model.WorkOrderBatchUpdateInput) error {
	data := do.DemoWorkOrder{}
	if in.Priority != nil {
		data.Priority = *in.Priority
	}
	if in.SourceType != nil {
		data.SourceType = *in.SourceType
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	normalizedIDs := normalizeWorkOrderIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoWorkOrder.Ctx(ctx), normalizedIDs, dao.DemoWorkOrder.Columns().Id, dao.DemoWorkOrder.Columns().TenantId, dao.DemoWorkOrder.Columns().MerchantId, "体验工单"); err != nil {
		return err
	}
	_, err := dao.DemoWorkOrder.Ctx(ctx).WhereIn(dao.DemoWorkOrder.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验工单
func (s *sWorkOrder) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
	f, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// 跳过表头
	if _, err = reader.Read(); err != nil {
		return 0, 0, fmt.Errorf("读取CSV表头失败: %w", err)
	}

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
		// 逐行插入
		id := snowflake.Generate()
		data := do.DemoWorkOrder{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.TicketNo = record[idx]
		}
		idx++
		if idx < len(record) {
			data.CustomerId = record[idx]
		}
		idx++
		if idx < len(record) {
			data.ProductId = record[idx]
		}
		idx++
		if idx < len(record) {
			data.OrderId = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Title = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Priority = record[idx]
		}
		idx++
		if idx < len(record) {
			data.SourceType = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Description = record[idx]
		}
		idx++
		if idx < len(record) {
			data.AttachmentFile = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Status = record[idx]
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
		if _, insertErr := dao.DemoWorkOrder.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

