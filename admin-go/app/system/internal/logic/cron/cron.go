
package cron

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
	service.RegisterCron(New())
}

func New() *sCron {
	return &sCron{}
}

type sCron struct{}

func normalizeCronIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建定时任务表
func (s *sCron) Create(ctx context.Context, in *model.CronCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.SystemCron.Ctx(ctx).Data(do.SystemCron{
		Id:        id,
		Name: in.Name,
		Expression: in.Expression,
		Handler: in.Handler,
		Params: in.Params,
		Remark: in.Remark,
		Status: in.Status,
		LastRunAt: in.LastRunAt,
		LastResult: in.LastResult,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新定时任务表
func (s *sCron) Update(ctx context.Context, in *model.CronUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.SystemCron{
		Name: in.Name,
		Expression: in.Expression,
		Handler: in.Handler,
		Params: in.Params,
		Remark: in.Remark,
		Status: in.Status,
		LastRunAt: in.LastRunAt,
		LastResult: in.LastResult,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), in.ID, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId, "定时任务表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), in.ID, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemCron.Ctx(ctx).Where(dao.SystemCron.Columns().Id, in.ID).Where(dao.SystemCron.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除定时任务表
func (s *sCron) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), id, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId, "定时任务表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), id, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemCron.Ctx(ctx).Where(dao.SystemCron.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除定时任务表
func (s *sCron) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeCronIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemCron.Ctx(ctx), normalizedIDs, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId, "定时任务表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemCron.Ctx(ctx), normalizedIDs, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemCron.Ctx(ctx).WhereIn(dao.SystemCron.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取定时任务表详情
func (s *sCron) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.CronDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), id, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId, "定时任务表"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.SystemCron.Ctx(ctx), id, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.CronDetailOutput{}
	err = dao.SystemCron.Ctx(ctx).Where(dao.SystemCron.Columns().Id, id).Where(dao.SystemCron.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("定时任务表不存在或已删除")
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
func (s *sCron) applyListFilter(ctx context.Context, in *model.CronListInput) *gdb.Model {
	m := dao.SystemCron.Ctx(ctx).Where(dao.SystemCron.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.SystemCron.Columns().Name, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.SystemCron.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Name != "" {
		m = m.WhereLike(dao.SystemCron.Columns().Name, "%"+in.Name+"%")
	}
	if in.TenantID != nil {
		m = m.Where(dao.SystemCron.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.SystemCron.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.SystemCron.Columns().Status, *in.Status)
	}
	if in.Remark != "" {
		m = m.WhereLike(dao.SystemCron.Columns().Remark, "%"+in.Remark+"%")
	}
	if in.LastRunAtStart != "" {
		m = m.WhereGTE(dao.SystemCron.Columns().LastRunAt, in.LastRunAtStart)
	}
	if in.LastRunAtEnd != "" {
		m = m.WhereLTE(dao.SystemCron.Columns().LastRunAt, in.LastRunAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.SystemCron.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.SystemCron.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sCron) fillRefFields(ctx context.Context, list []*model.CronListOutput) {
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

// List 获取定时任务表列表
func (s *sCron) List(ctx context.Context, in *model.CronListInput) (list []*model.CronListOutput, total int, err error) {
	if in == nil {
		in = &model.CronListInput{}
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
func (s *sCron) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.SystemCron.Columns().Id:        true,
		dao.SystemCron.Columns().CreatedAt: true,
		dao.SystemCron.Columns().Status:    true,
		dao.SystemCron.Columns().Name: true,
		dao.SystemCron.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sCron) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.SystemCron.Columns().Id)
}

// Export 导出定时任务表（不分页）
func (s *sCron) Export(ctx context.Context, in *model.CronListInput) (list []*model.CronListOutput, err error) {
	if in == nil {
		in = &model.CronListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑定时任务表
func (s *sCron) BatchUpdate(ctx context.Context, in *model.CronBatchUpdateInput) error {
	data := do.SystemCron{}
	hasChange := false
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeCronIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.SystemCron.Ctx(ctx), normalizedIDs, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().TenantId, dao.SystemCron.Columns().MerchantId, "定时任务表"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.SystemCron.Ctx(ctx), normalizedIDs, dao.SystemCron.Columns().Id, dao.SystemCron.Columns().CreatedBy, dao.SystemCron.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.SystemCron.Ctx(ctx).WhereIn(dao.SystemCron.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入定时任务表
func (s *sCron) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.SystemCron{
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
			data.Expression = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Handler = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.Params = strings.TrimSpace(record[idx])
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
		if idx < len(record) {
			data.LastResult = strings.TrimSpace(record[idx])
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
		if _, insertErr := dao.SystemCron.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

