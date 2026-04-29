
package audit_log

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

	"gbaseadmin/app/demo/internal/dao"
	"gbaseadmin/app/demo/internal/middleware"
	"gbaseadmin/app/demo/internal/model"
	"gbaseadmin/app/demo/internal/model/do"
	"gbaseadmin/app/demo/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterAuditLog(New())
}

func New() *sAuditLog {
	return &sAuditLog{}
}

type sAuditLog struct{}

func normalizeAuditLogIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建体验审计日志
func (s *sAuditLog) Create(ctx context.Context, in *model.AuditLogCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoAuditLog.Ctx(ctx).Data(do.DemoAuditLog{
		Id:        id,
		LogNo: in.LogNo,
		OperatorId: in.OperatorID,
		Action: in.Action,
		TargetType: in.TargetType,
		TargetCode: in.TargetCode,
		RequestJson: in.RequestJSON,
		Result: in.Result,
		ClientIp: in.ClientIP,
		OccurredAt: in.OccurredAt,
		Remark: in.Remark,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验审计日志
func (s *sAuditLog) Update(ctx context.Context, in *model.AuditLogUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.DemoAuditLog{
		LogNo: in.LogNo,
		OperatorId: in.OperatorID,
		Action: in.Action,
		TargetType: in.TargetType,
		TargetCode: in.TargetCode,
		RequestJson: in.RequestJSON,
		Result: in.Result,
		ClientIp: in.ClientIP,
		OccurredAt: in.OccurredAt,
		Remark: in.Remark,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), in.ID, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().TenantId, dao.DemoAuditLog.Columns().MerchantId, "体验审计日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), in.ID, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().CreatedBy, dao.DemoAuditLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoAuditLog.Ctx(ctx).Where(dao.DemoAuditLog.Columns().Id, in.ID).Where(dao.DemoAuditLog.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除体验审计日志
func (s *sAuditLog) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), id, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().TenantId, dao.DemoAuditLog.Columns().MerchantId, "体验审计日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), id, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().CreatedBy, dao.DemoAuditLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoAuditLog.Ctx(ctx).Where(dao.DemoAuditLog.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验审计日志
func (s *sAuditLog) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeAuditLogIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), normalizedIDs, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().TenantId, dao.DemoAuditLog.Columns().MerchantId, "体验审计日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), normalizedIDs, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().CreatedBy, dao.DemoAuditLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.DemoAuditLog.Ctx(ctx).WhereIn(dao.DemoAuditLog.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验审计日志详情
func (s *sAuditLog) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.AuditLogDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), id, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().TenantId, dao.DemoAuditLog.Columns().MerchantId, "体验审计日志"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.DemoAuditLog.Ctx(ctx), id, dao.DemoAuditLog.Columns().Id, dao.DemoAuditLog.Columns().CreatedBy, dao.DemoAuditLog.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.AuditLogDetailOutput{}
	err = dao.DemoAuditLog.Ctx(ctx).Where(dao.DemoAuditLog.Columns().Id, id).Where(dao.DemoAuditLog.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("体验审计日志不存在或已删除")
	}
	// 查询操作人关联显示
	if out.OperatorID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_users").Where("id", out.OperatorID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("username")
		if err == nil {
			out.UsersUsername = val.String()
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
func (s *sAuditLog) applyListFilter(ctx context.Context, in *model.AuditLogListInput) *gdb.Model {
	m := dao.DemoAuditLog.Ctx(ctx).Where(dao.DemoAuditLog.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoAuditLog.Columns().TenantId, dao.DemoAuditLog.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.DemoAuditLog.Columns().Remark, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.LogNo != "" {
		m = m.Where(dao.DemoAuditLog.Columns().LogNo, in.LogNo)
	}
	if in.TargetCode != "" {
		m = m.Where(dao.DemoAuditLog.Columns().TargetCode, in.TargetCode)
	}
	if in.OperatorID != nil {
		m = m.Where(dao.DemoAuditLog.Columns().OperatorId, *in.OperatorID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoAuditLog.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoAuditLog.Columns().MerchantId, *in.MerchantID)
	}
	if in.ClientIP != "" {
		m = m.Where(dao.DemoAuditLog.Columns().ClientIp, in.ClientIP)
	}
	if in.Action != nil {
		m = m.Where(dao.DemoAuditLog.Columns().Action, *in.Action)
	}
	if in.TargetType != nil {
		m = m.Where(dao.DemoAuditLog.Columns().TargetType, *in.TargetType)
	}
	if in.Result != nil {
		m = m.Where(dao.DemoAuditLog.Columns().Result, *in.Result)
	}
	if in.OccurredAtStart != "" {
		m = m.WhereGTE(dao.DemoAuditLog.Columns().OccurredAt, in.OccurredAtStart)
	}
	if in.OccurredAtEnd != "" {
		m = m.WhereLTE(dao.DemoAuditLog.Columns().OccurredAt, in.OccurredAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoAuditLog.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoAuditLog.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoAuditLog.Columns().CreatedBy, dao.DemoAuditLog.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sAuditLog) fillRefFields(ctx context.Context, list []*model.AuditLogListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.OperatorID != 0 {
				idSet[int64(item.OperatorID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_users").
				Fields("id", "username")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["username"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.OperatorID)]; ok {
						item.UsersUsername = val
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

// List 获取体验审计日志列表
func (s *sAuditLog) List(ctx context.Context, in *model.AuditLogListInput) (list []*model.AuditLogListOutput, total int, err error) {
	if in == nil {
		in = &model.AuditLogListInput{}
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
func (s *sAuditLog) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoAuditLog.Columns().Id:        true,
		dao.DemoAuditLog.Columns().CreatedAt: true,
		dao.DemoAuditLog.Columns().LogNo: true,
		dao.DemoAuditLog.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sAuditLog) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoAuditLog.Columns().Id)
}

// Export 导出体验审计日志（不分页）
func (s *sAuditLog) Export(ctx context.Context, in *model.AuditLogListInput) (list []*model.AuditLogListOutput, err error) {
	if in == nil {
		in = &model.AuditLogListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// Import 导入体验审计日志
func (s *sAuditLog) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.DemoAuditLog{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.LogNo = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.OperatorId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Action = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.TargetType = v
			}
		}
		idx++
		if idx < len(record) {
			data.TargetCode = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			data.RequestJson = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.Result = v
			}
		}
		idx++
		if idx < len(record) {
			data.ClientIp = strings.TrimSpace(record[idx])
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
		if fkVal, ok := data.OperatorId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("system_users").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.DemoAuditLog.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

