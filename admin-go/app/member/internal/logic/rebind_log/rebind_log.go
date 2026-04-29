
package rebind_log

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
	service.RegisterRebindLog(New())
}

func New() *sRebindLog {
	return &sRebindLog{}
}

type sRebindLog struct{}

func normalizeRebindLogIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建换绑上级日志
func (s *sRebindLog) Create(ctx context.Context, in *model.RebindLogCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberRebindLog.Ctx(ctx).Data(do.MemberRebindLog{
		Id:        id,
		UserId: in.UserID,
		OldParentId: in.OldParentID,
		NewParentId: in.NewParentID,
		Reason: in.Reason,
		OperatorId: in.OperatorID,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新换绑上级日志
func (s *sRebindLog) Update(ctx context.Context, in *model.RebindLogUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberRebindLog{
		UserId: in.UserID,
		OldParentId: in.OldParentID,
		NewParentId: in.NewParentID,
		Reason: in.Reason,
		OperatorId: in.OperatorID,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), in.ID, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().TenantId, dao.MemberRebindLog.Columns().MerchantId, "换绑上级日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), in.ID, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().CreatedBy, dao.MemberRebindLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberRebindLog.Ctx(ctx).Where(dao.MemberRebindLog.Columns().Id, in.ID).Where(dao.MemberRebindLog.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除换绑上级日志
func (s *sRebindLog) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), id, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().TenantId, dao.MemberRebindLog.Columns().MerchantId, "换绑上级日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), id, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().CreatedBy, dao.MemberRebindLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberRebindLog.Ctx(ctx).Where(dao.MemberRebindLog.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除换绑上级日志
func (s *sRebindLog) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeRebindLogIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), normalizedIDs, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().TenantId, dao.MemberRebindLog.Columns().MerchantId, "换绑上级日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), normalizedIDs, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().CreatedBy, dao.MemberRebindLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberRebindLog.Ctx(ctx).WhereIn(dao.MemberRebindLog.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取换绑上级日志详情
func (s *sRebindLog) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.RebindLogDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), id, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().TenantId, dao.MemberRebindLog.Columns().MerchantId, "换绑上级日志"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberRebindLog.Ctx(ctx), id, dao.MemberRebindLog.Columns().Id, dao.MemberRebindLog.Columns().CreatedBy, dao.MemberRebindLog.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.RebindLogDetailOutput{}
	err = dao.MemberRebindLog.Ctx(ctx).Where(dao.MemberRebindLog.Columns().Id, id).Where(dao.MemberRebindLog.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("换绑上级日志不存在或已删除")
	}
	// 查询会员关联显示
	if out.UserID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.UserID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.UserNickname = val.String()
		}
	}
	// 查询原上级关联显示
	if out.OldParentID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.OldParentID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.OldParentNickname = val.String()
		}
	}
	// 查询新上级关联显示
	if out.NewParentID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.NewParentID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("nickname")
		if err == nil {
			out.NewParentNickname = val.String()
		}
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
func (s *sRebindLog) applyListFilter(ctx context.Context, in *model.RebindLogListInput) *gdb.Model {
	m := dao.MemberRebindLog.Ctx(ctx).Where(dao.MemberRebindLog.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberRebindLog.Columns().TenantId, dao.MemberRebindLog.Columns().MerchantId)
	if in.UserID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().UserId, *in.UserID)
	}
	if in.OldParentID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().OldParentId, *in.OldParentID)
	}
	if in.NewParentID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().NewParentId, *in.NewParentID)
	}
	if in.OperatorID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().OperatorId, *in.OperatorID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberRebindLog.Columns().MerchantId, *in.MerchantID)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberRebindLog.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberRebindLog.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberRebindLog.Columns().CreatedBy, dao.MemberRebindLog.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sRebindLog) fillRefFields(ctx context.Context, list []*model.RebindLogListOutput) {
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
			if item.OldParentID != 0 {
				idSet[int64(item.OldParentID)] = struct{}{}
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
					if val, ok := refMap[int64(item.OldParentID)]; ok {
						item.OldParentNickname = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.NewParentID != 0 {
				idSet[int64(item.NewParentID)] = struct{}{}
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
					if val, ok := refMap[int64(item.NewParentID)]; ok {
						item.NewParentNickname = val
					}
				}
			}
		}
	}
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

// List 获取换绑上级日志列表
func (s *sRebindLog) List(ctx context.Context, in *model.RebindLogListInput) (list []*model.RebindLogListOutput, total int, err error) {
	if in == nil {
		in = &model.RebindLogListInput{}
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
func (s *sRebindLog) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberRebindLog.Columns().Id:        true,
		dao.MemberRebindLog.Columns().CreatedAt: true,
	}
	return allowed[field]
}

func (s *sRebindLog) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.MemberRebindLog.Columns().Id)
}

// Export 导出换绑上级日志（不分页）
func (s *sRebindLog) Export(ctx context.Context, in *model.RebindLogListInput) (list []*model.RebindLogListOutput, err error) {
	if in == nil {
		in = &model.RebindLogListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// Import 导入换绑上级日志
func (s *sRebindLog) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberRebindLog{
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
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.OldParentId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.NewParentId = v
			}
		}
		idx++
		if idx < len(record) {
			data.Reason = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.OperatorId = v
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
		if fkVal, ok := data.OldParentId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.NewParentId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.OperatorId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("system_users").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.MemberRebindLog.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

