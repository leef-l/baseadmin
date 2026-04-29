
package team_export

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

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterTeamExport(New())
}

func New() *sTeamExport {
	return &sTeamExport{}
}

type sTeamExport struct{}

func normalizeTeamExportIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建团队数据导出
func (s *sTeamExport) Create(ctx context.Context, in *model.TeamExportCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberTeamExport.Ctx(ctx).Data(do.MemberTeamExport{
		Id:        id,
		UserId: in.UserID,
		TeamMemberCount: in.TeamMemberCount,
		ExportType: in.ExportType,
		FileUrl: in.FileURL,
		FileSize: in.FileSize,
		DeployStatus: in.DeployStatus,
		DeployDomain: in.DeployDomain,
		DeployedAt: in.DeployedAt,
		Remark: in.Remark,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新团队数据导出
func (s *sTeamExport) Update(ctx context.Context, in *model.TeamExportUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberTeamExport{
		UserId: in.UserID,
		TeamMemberCount: in.TeamMemberCount,
		ExportType: in.ExportType,
		FileUrl: in.FileURL,
		FileSize: in.FileSize,
		DeployStatus: in.DeployStatus,
		DeployDomain: in.DeployDomain,
		DeployedAt: in.DeployedAt,
		Remark: in.Remark,
		Status: in.Status,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), in.ID, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId, "团队数据导出"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), in.ID, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberTeamExport.Ctx(ctx).Where(dao.MemberTeamExport.Columns().Id, in.ID).Where(dao.MemberTeamExport.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除团队数据导出
func (s *sTeamExport) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), id, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId, "团队数据导出"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), id, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberTeamExport.Ctx(ctx).Where(dao.MemberTeamExport.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除团队数据导出
func (s *sTeamExport) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeTeamExportIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), normalizedIDs, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId, "团队数据导出"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), normalizedIDs, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberTeamExport.Ctx(ctx).WhereIn(dao.MemberTeamExport.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取团队数据导出详情
func (s *sTeamExport) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.TeamExportDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), id, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId, "团队数据导出"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), id, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.TeamExportDetailOutput{}
	err = dao.MemberTeamExport.Ctx(ctx).Where(dao.MemberTeamExport.Columns().Id, id).Where(dao.MemberTeamExport.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("团队数据导出不存在或已删除")
	}
	// 查询目标会员关联显示
	if out.UserID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.UserID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.UserNickname = val.String()
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
func (s *sTeamExport) applyListFilter(ctx context.Context, in *model.TeamExportListInput) *gdb.Model {
	m := dao.MemberTeamExport.Ctx(ctx).Where(dao.MemberTeamExport.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId)
	if in.UserID != nil {
		m = m.Where(dao.MemberTeamExport.Columns().UserId, *in.UserID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberTeamExport.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberTeamExport.Columns().MerchantId, *in.MerchantID)
	}
	if in.ExportType != nil {
		m = m.Where(dao.MemberTeamExport.Columns().ExportType, *in.ExportType)
	}
	if in.DeployStatus != nil {
		m = m.Where(dao.MemberTeamExport.Columns().DeployStatus, *in.DeployStatus)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberTeamExport.Columns().Status, *in.Status)
	}
	if in.DeployDomain != "" {
		m = m.WhereLike(dao.MemberTeamExport.Columns().DeployDomain, "%"+in.DeployDomain+"%")
	}
	if in.DeployedAtStart != "" {
		m = m.WhereGTE(dao.MemberTeamExport.Columns().DeployedAt, in.DeployedAtStart)
	}
	if in.DeployedAtEnd != "" {
		m = m.WhereLTE(dao.MemberTeamExport.Columns().DeployedAt, in.DeployedAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberTeamExport.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberTeamExport.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sTeamExport) fillRefFields(ctx context.Context, list []*model.TeamExportListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.UserID != 0 {
				idSet[int64(item.UserID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_user").
				Fields("id", "nickname")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, err := refQuery.WhereIn("id", ids).All()
			if err == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["nickname"].String()
				}
				for _, item := range list {
					if val, ok := refMap[int64(item.UserID)]; ok {
						item.UserNickname = val
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

// List 获取团队数据导出列表
func (s *sTeamExport) List(ctx context.Context, in *model.TeamExportListInput) (list []*model.TeamExportListOutput, total int, err error) {
	if in == nil {
		in = &model.TeamExportListInput{}
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
func (s *sTeamExport) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberTeamExport.Columns().Id:        true,
		dao.MemberTeamExport.Columns().CreatedAt: true,
		dao.MemberTeamExport.Columns().Status:    true,
		dao.MemberTeamExport.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sTeamExport) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.MemberTeamExport.Columns().Id)
}

// Export 导出团队数据导出（不分页）
func (s *sTeamExport) Export(ctx context.Context, in *model.TeamExportListInput) (list []*model.TeamExportListOutput, err error) {
	if in == nil {
		in = &model.TeamExportListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑团队数据导出
func (s *sTeamExport) BatchUpdate(ctx context.Context, in *model.TeamExportBatchUpdateInput) error {
	data := do.MemberTeamExport{}
	hasChange := false
	if in.ExportType != nil {
		data.ExportType = *in.ExportType
		hasChange = true
	}
	if in.DeployStatus != nil {
		data.DeployStatus = *in.DeployStatus
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeTeamExportIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), normalizedIDs, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().TenantId, dao.MemberTeamExport.Columns().MerchantId, "团队数据导出"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberTeamExport.Ctx(ctx), normalizedIDs, dao.MemberTeamExport.Columns().Id, dao.MemberTeamExport.Columns().CreatedBy, dao.MemberTeamExport.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberTeamExport.Ctx(ctx).WhereIn(dao.MemberTeamExport.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入团队数据导出
func (s *sTeamExport) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberTeamExport{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.UserId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.TeamMemberCount = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.ExportType = v
			}
		}
		idx++
		if idx < len(record) {
			data.FileUrl = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.FileSize = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.DeployStatus = v
			}
		}
		idx++
		if idx < len(record) {
			data.DeployDomain = strings.TrimSpace(record[idx])
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
		if fkVal, ok := data.UserId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.MemberTeamExport.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

