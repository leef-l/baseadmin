
package plan

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

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/middleware"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterPlan(New())
}

func New() *sPlan {
	return &sPlan{}
}

type sPlan struct{}

func normalizePlanIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建套餐表
func (s *sPlan) Create(ctx context.Context, in *model.PlanCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.SystemPlan.Ctx(ctx).Data(do.SystemPlan{
		Id:        id,
		Name: in.Name,
		Code: in.Code,
		Description: in.Description,
		UserLimit: in.UserLimit,
		StorageMb: in.StorageMb,
		Features: in.Features,
		Price: in.Price,
		Sort: in.Sort,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新套餐表
func (s *sPlan) Update(ctx context.Context, in *model.PlanUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.SystemPlan{
		Name: in.Name,
		Code: in.Code,
		Description: in.Description,
		UserLimit: in.UserLimit,
		StorageMb: in.StorageMb,
		Features: in.Features,
		Price: in.Price,
		Sort: in.Sort,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.SystemPlan.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.SystemPlan.Table()).Ctx(ctx).
			Where(dao.SystemPlan.Columns().Id, in.ID).
			Where(dao.SystemPlan.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("套餐表不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.SystemPlan.Table()).Ctx(ctx), in.ID, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId, "套餐表"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.SystemPlan.Table()).Ctx(ctx), in.ID, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.SystemPlan.Table()).Ctx(ctx).
			Where(dao.SystemPlan.Columns().Id, in.ID).
			Where(dao.SystemPlan.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除套餐表
func (s *sPlan) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemPlan.Ctx(ctx), id, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId, "套餐表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemPlan.Ctx(ctx), id, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemPlan.Ctx(ctx).Where(dao.SystemPlan.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除套餐表
func (s *sPlan) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizePlanIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemPlan.Ctx(ctx), normalizedIDs, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId, "套餐表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemPlan.Ctx(ctx), normalizedIDs, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemPlan.Ctx(ctx).WhereIn(dao.SystemPlan.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取套餐表详情
func (s *sPlan) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.PlanDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemPlan.Ctx(ctx), id, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId, "套餐表"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemPlan.Ctx(ctx), id, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.PlanDetailOutput{}
	err = dao.SystemPlan.Ctx(ctx).Where(dao.SystemPlan.Columns().Id, id).Where(dao.SystemPlan.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("套餐表不存在或已删除")
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
func (s *sPlan) applyListFilter(ctx context.Context, in *model.PlanListInput) *gdb.Model {
	m := dao.SystemPlan.Ctx(ctx).Where(dao.SystemPlan.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.SystemPlan.Columns().Name, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.SystemPlan.Columns().Description, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Code != "" {
		m = m.Where(dao.SystemPlan.Columns().Code, in.Code)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.SystemPlan.Columns().Name, "%"+in.Name+"%")
	}
	if in.TenantID != nil {
		m = m.Where(dao.SystemPlan.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.SystemPlan.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.SystemPlan.Columns().Status, *in.Status)
	}
	if in.Description != "" {
		m = m.WhereLike(dao.SystemPlan.Columns().Description, "%"+in.Description+"%")
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.SystemPlan.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.SystemPlan.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sPlan) fillRefFields(ctx context.Context, list []*model.PlanListOutput) {
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

// List 获取套餐表列表
func (s *sPlan) List(ctx context.Context, in *model.PlanListInput) (list []*model.PlanListOutput, total int, err error) {
	if in == nil {
		in = &model.PlanListInput{}
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
func (s *sPlan) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.SystemPlan.Columns().Id:        true,
		dao.SystemPlan.Columns().CreatedAt: true,
		dao.SystemPlan.Columns().Sort:      true,
		dao.SystemPlan.Columns().Status:    true,
		dao.SystemPlan.Columns().Name: true,
		dao.SystemPlan.Columns().Description: true,
		dao.SystemPlan.Columns().Price: true,
	}
	return allowed[field]
}

func (s *sPlan) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.SystemPlan.Columns().Sort).OrderDesc(dao.SystemPlan.Columns().Id)
}

// Export 导出套餐表（不分页）
func (s *sPlan) Export(ctx context.Context, in *model.PlanListInput) (list []*model.PlanListOutput, err error) {
	if in == nil {
		in = &model.PlanListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑套餐表
func (s *sPlan) BatchUpdate(ctx context.Context, in *model.PlanBatchUpdateInput) error {
	data := do.SystemPlan{}
	hasChange := false
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizePlanIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemPlan.Ctx(ctx), normalizedIDs, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().TenantId, dao.SystemPlan.Columns().MerchantId, "套餐表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemPlan.Ctx(ctx), normalizedIDs, dao.SystemPlan.Columns().Id, dao.SystemPlan.Columns().CreatedBy, dao.SystemPlan.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemPlan.Ctx(ctx).WhereIn(dao.SystemPlan.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入套餐表
func (s *sPlan) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.SystemPlan{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.Name = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Code = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Description = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.UserLimit = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.StorageMb = v
			}
		}
		idx++
		if idx < len(record) {
			data.Features = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.Price = int64(math.Round(v * 100))
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
		if _, insertErr := dao.SystemPlan.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

