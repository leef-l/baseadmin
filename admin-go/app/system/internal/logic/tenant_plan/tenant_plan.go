
package tenant_plan

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/middleware"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterTenantPlan(New())
}

func New() *sTenantPlan {
	return &sTenantPlan{}
}

type sTenantPlan struct{}

func normalizeTenantPlanIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建租户套餐订阅表
func (s *sTenantPlan) Create(ctx context.Context, in *model.TenantPlanCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.SystemTenantPlan.Ctx(ctx).Data(do.SystemTenantPlan{
		Id:        id,
		TenantId: in.TenantID,
		PlanId: in.PlanID,
		StartAt: in.StartAt,
		ExpireAt: in.ExpireAt,
		Status: in.Status,
		Remark: in.Remark,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新租户套餐订阅表
func (s *sTenantPlan) Update(ctx context.Context, in *model.TenantPlanUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.SystemTenantPlan{
		PlanId: in.PlanID,
		StartAt: in.StartAt,
		ExpireAt: in.ExpireAt,
		Status: in.Status,
		Remark: in.Remark,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), in.ID, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId, "租户套餐订阅表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), in.ID, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemTenantPlan.Ctx(ctx).Where(dao.SystemTenantPlan.Columns().Id, in.ID).Where(dao.SystemTenantPlan.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除租户套餐订阅表
func (s *sTenantPlan) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), id, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId, "租户套餐订阅表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), id, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemTenantPlan.Ctx(ctx).Where(dao.SystemTenantPlan.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除租户套餐订阅表
func (s *sTenantPlan) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeTenantPlanIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), normalizedIDs, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId, "租户套餐订阅表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), normalizedIDs, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemTenantPlan.Ctx(ctx).WhereIn(dao.SystemTenantPlan.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取租户套餐订阅表详情
func (s *sTenantPlan) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TenantPlanDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), id, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId, "租户套餐订阅表"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), id, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.TenantPlanDetailOutput{}
	err = dao.SystemTenantPlan.Ctx(ctx).Where(dao.SystemTenantPlan.Columns().Id, id).Where(dao.SystemTenantPlan.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("租户套餐订阅表不存在或已删除")
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
	// 查询套餐ID关联显示
	if out.PlanID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_plan").Where("id", out.PlanID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.PlanName = val.String()
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
func (s *sTenantPlan) applyListFilter(ctx context.Context, in *model.TenantPlanListInput) *gdb.Model {
	m := dao.SystemTenantPlan.Ctx(ctx).Where(dao.SystemTenantPlan.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId)
	if in.TenantID != nil {
		m = m.Where(dao.SystemTenantPlan.Columns().TenantId, *in.TenantID)
	}
	if in.PlanID != nil {
		m = m.Where(dao.SystemTenantPlan.Columns().PlanId, *in.PlanID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.SystemTenantPlan.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.SystemTenantPlan.Columns().Status, *in.Status)
	}
	if in.Remark != "" {
		m = m.WhereLike(dao.SystemTenantPlan.Columns().Remark, "%"+in.Remark+"%")
	}
	if in.StartAtStart != "" {
		m = m.WhereGTE(dao.SystemTenantPlan.Columns().StartAt, in.StartAtStart)
	}
	if in.StartAtEnd != "" {
		m = m.WhereLTE(dao.SystemTenantPlan.Columns().StartAt, in.StartAtEnd)
	}
	if in.ExpireAtStart != "" {
		m = m.WhereGTE(dao.SystemTenantPlan.Columns().ExpireAt, in.ExpireAtStart)
	}
	if in.ExpireAtEnd != "" {
		m = m.WhereLTE(dao.SystemTenantPlan.Columns().ExpireAt, in.ExpireAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.SystemTenantPlan.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.SystemTenantPlan.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sTenantPlan) fillRefFields(ctx context.Context, list []*model.TenantPlanListOutput) {
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
			if item.PlanID != 0 {
				idSet[int64(item.PlanID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_plan").
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
					if val, ok := refMap[int64(item.PlanID)]; ok {
						item.PlanName = val
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

// List 获取租户套餐订阅表列表
func (s *sTenantPlan) List(ctx context.Context, in *model.TenantPlanListInput) (list []*model.TenantPlanListOutput, total int, err error) {
	if in == nil {
		in = &model.TenantPlanListInput{}
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
func (s *sTenantPlan) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.SystemTenantPlan.Columns().Id:        true,
		dao.SystemTenantPlan.Columns().CreatedAt: true,
		dao.SystemTenantPlan.Columns().Status:    true,
		dao.SystemTenantPlan.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sTenantPlan) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.SystemTenantPlan.Columns().Id)
}

// Export 导出租户套餐订阅表（不分页）
func (s *sTenantPlan) Export(ctx context.Context, in *model.TenantPlanListInput) (list []*model.TenantPlanListOutput, err error) {
	if in == nil {
		in = &model.TenantPlanListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑租户套餐订阅表
func (s *sTenantPlan) BatchUpdate(ctx context.Context, in *model.TenantPlanBatchUpdateInput) error {
	data := do.SystemTenantPlan{}
	hasChange := false
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeTenantPlanIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), normalizedIDs, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().TenantId, dao.SystemTenantPlan.Columns().MerchantId, "租户套餐订阅表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemTenantPlan.Ctx(ctx), normalizedIDs, dao.SystemTenantPlan.Columns().Id, dao.SystemTenantPlan.Columns().CreatedBy, dao.SystemTenantPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemTenantPlan.Ctx(ctx).WhereIn(dao.SystemTenantPlan.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入租户套餐订阅表
func (s *sTenantPlan) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.SystemTenantPlan{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.PlanId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Status = v
			}
		}
		idx++
		if idx < len(record) {
			data.Remark = strings.TrimSpace(record[idx])
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
		if fkVal, ok := data.PlanId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("system_plan").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.SystemTenantPlan.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

