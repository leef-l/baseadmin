
package level_log

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
	service.RegisterLevelLog(New())
}

func New() *sLevelLog {
	return &sLevelLog{}
}

type sLevelLog struct{}

func normalizeLevelLogIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建等级变更日志
func (s *sLevelLog) Create(ctx context.Context, in *model.LevelLogCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberLevelLog.Ctx(ctx).Data(do.MemberLevelLog{
		Id:        id,
		UserId: in.UserID,
		OldLevelId: in.OldLevelID,
		NewLevelId: in.NewLevelID,
		ChangeType: in.ChangeType,
		ExpireAt: in.ExpireAt,
		Remark: in.Remark,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新等级变更日志
func (s *sLevelLog) Update(ctx context.Context, in *model.LevelLogUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberLevelLog{
		UserId: in.UserID,
		OldLevelId: in.OldLevelID,
		NewLevelId: in.NewLevelID,
		ChangeType: in.ChangeType,
		ExpireAt: in.ExpireAt,
		Remark: in.Remark,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), in.ID, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().TenantId, dao.MemberLevelLog.Columns().MerchantId, "等级变更日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), in.ID, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().CreatedBy, dao.MemberLevelLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevelLog.Ctx(ctx).Where(dao.MemberLevelLog.Columns().Id, in.ID).Where(dao.MemberLevelLog.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除等级变更日志
func (s *sLevelLog) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), id, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().TenantId, dao.MemberLevelLog.Columns().MerchantId, "等级变更日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), id, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().CreatedBy, dao.MemberLevelLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevelLog.Ctx(ctx).Where(dao.MemberLevelLog.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除等级变更日志
func (s *sLevelLog) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeLevelLogIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), normalizedIDs, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().TenantId, dao.MemberLevelLog.Columns().MerchantId, "等级变更日志"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), normalizedIDs, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().CreatedBy, dao.MemberLevelLog.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevelLog.Ctx(ctx).WhereIn(dao.MemberLevelLog.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取等级变更日志详情
func (s *sLevelLog) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.LevelLogDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), id, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().TenantId, dao.MemberLevelLog.Columns().MerchantId, "等级变更日志"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberLevelLog.Ctx(ctx), id, dao.MemberLevelLog.Columns().Id, dao.MemberLevelLog.Columns().CreatedBy, dao.MemberLevelLog.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.LevelLogDetailOutput{}
	err = dao.MemberLevelLog.Ctx(ctx).Where(dao.MemberLevelLog.Columns().Id, id).Where(dao.MemberLevelLog.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("等级变更日志不存在或已删除")
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
	// 查询变更前等级关联显示
	if out.OldLevelID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_level").Where("id", out.OldLevelID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.LevelName = val.String()
		}
	}
	// 查询变更后等级关联显示
	if out.NewLevelID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_level").Where("id", out.NewLevelID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.NewLevelName = val.String()
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
func (s *sLevelLog) applyListFilter(ctx context.Context, in *model.LevelLogListInput) *gdb.Model {
	m := dao.MemberLevelLog.Ctx(ctx).Where(dao.MemberLevelLog.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberLevelLog.Columns().TenantId, dao.MemberLevelLog.Columns().MerchantId)
	if in.UserID != nil {
		m = m.Where(dao.MemberLevelLog.Columns().UserId, *in.UserID)
	}
	if in.OldLevelID != nil {
		m = m.Where(dao.MemberLevelLog.Columns().OldLevelId, *in.OldLevelID)
	}
	if in.NewLevelID != nil {
		m = m.Where(dao.MemberLevelLog.Columns().NewLevelId, *in.NewLevelID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberLevelLog.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberLevelLog.Columns().MerchantId, *in.MerchantID)
	}
	if in.ChangeType != nil {
		m = m.Where(dao.MemberLevelLog.Columns().ChangeType, *in.ChangeType)
	}
	if in.ExpireAtStart != "" {
		m = m.WhereGTE(dao.MemberLevelLog.Columns().ExpireAt, in.ExpireAtStart)
	}
	if in.ExpireAtEnd != "" {
		m = m.WhereLTE(dao.MemberLevelLog.Columns().ExpireAt, in.ExpireAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberLevelLog.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberLevelLog.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberLevelLog.Columns().CreatedBy, dao.MemberLevelLog.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sLevelLog) fillRefFields(ctx context.Context, list []*model.LevelLogListOutput) {
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
			if item.OldLevelID != 0 {
				idSet[int64(item.OldLevelID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_level").
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
					if val, ok := refMap[int64(item.OldLevelID)]; ok {
						item.LevelName = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.NewLevelID != 0 {
				idSet[int64(item.NewLevelID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_level").
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
					if val, ok := refMap[int64(item.NewLevelID)]; ok {
						item.NewLevelName = val
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

// List 获取等级变更日志列表
func (s *sLevelLog) List(ctx context.Context, in *model.LevelLogListInput) (list []*model.LevelLogListOutput, total int, err error) {
	if in == nil {
		in = &model.LevelLogListInput{}
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
func (s *sLevelLog) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberLevelLog.Columns().Id:        true,
		dao.MemberLevelLog.Columns().CreatedAt: true,
		dao.MemberLevelLog.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sLevelLog) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.MemberLevelLog.Columns().Id)
}

// Export 导出等级变更日志（不分页）
func (s *sLevelLog) Export(ctx context.Context, in *model.LevelLogListInput) (list []*model.LevelLogListOutput, err error) {
	if in == nil {
		in = &model.LevelLogListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// Import 导入等级变更日志
func (s *sLevelLog) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberLevelLog{
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
				data.OldLevelId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseInt(strings.TrimSpace(record[idx]), 10, 64); parseErr == nil {
				data.NewLevelId = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.ChangeType = v
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
		if fkVal, ok := data.UserId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.OldLevelId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_level").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if fkVal, ok := data.NewLevelId.(int64); ok && fkVal > 0 {
			refQuery := g.DB().Ctx(ctx).Model("member_level").Where("id", fkVal)
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			if cnt, cntErr := refQuery.Count(); cntErr != nil || cnt == 0 {
				fail++
				continue
			}
		}
		if _, insertErr := dao.MemberLevelLog.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

