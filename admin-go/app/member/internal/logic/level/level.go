
package level

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

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterLevel(New())
}

func New() *sLevel {
	return &sLevel{}
}

type sLevel struct{}

func normalizeLevelIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建会员等级配置
func (s *sLevel) Create(ctx context.Context, in *model.LevelCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberLevel.Ctx(ctx).Data(do.MemberLevel{
		Id:        id,
		Name: in.Name,
		LevelNo: in.LevelNo,
		Icon: in.Icon,
		DurationDays: in.DurationDays,
		NeedActiveCount: in.NeedActiveCount,
		NeedTeamTurnover: in.NeedTeamTurnover,
		IsTop: in.IsTop,
		AutoDeploy: in.AutoDeploy,
		Remark: in.Remark,
		Sort: in.Sort,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新会员等级配置
func (s *sLevel) Update(ctx context.Context, in *model.LevelUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	data := do.MemberLevel{
		Name: in.Name,
		LevelNo: in.LevelNo,
		Icon: in.Icon,
		DurationDays: in.DurationDays,
		NeedActiveCount: in.NeedActiveCount,
		NeedTeamTurnover: in.NeedTeamTurnover,
		IsTop: in.IsTop,
		AutoDeploy: in.AutoDeploy,
		Remark: in.Remark,
		Sort: in.Sort,
		Status: in.Status,
	}
	// 含金额字段，使用事务 + 行锁，权限检查在行锁内防止 TOCTOU
	err := dao.MemberLevel.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		lockedRow, err := tx.Model(dao.MemberLevel.Table()).Ctx(ctx).
			Where(dao.MemberLevel.Columns().Id, in.ID).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			LockUpdate().
			One()
		if err != nil {
			return err
		}
		if lockedRow.IsEmpty() {
			return gerror.New("会员等级配置不存在或已删除")
		}
		if err := middleware.EnsureTenantScopedRowAccessible(ctx, tx.Model(dao.MemberLevel.Table()).Ctx(ctx), in.ID, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId, "会员等级配置"); err != nil {
			return err
		}
		if err := middleware.EnsureDataScopedRowAccessible(ctx, tx.Model(dao.MemberLevel.Table()).Ctx(ctx), in.ID, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId); err != nil {
			return err
		}
		_, err = tx.Model(dao.MemberLevel.Table()).Ctx(ctx).
			Where(dao.MemberLevel.Columns().Id, in.ID).
			Where(dao.MemberLevel.Columns().DeletedAt, nil).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除会员等级配置
func (s *sLevel) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberLevel.Ctx(ctx), id, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId, "会员等级配置"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberLevel.Ctx(ctx), id, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevel.Ctx(ctx).Where(dao.MemberLevel.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除会员等级配置
func (s *sLevel) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeLevelIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberLevel.Ctx(ctx), normalizedIDs, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId, "会员等级配置"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberLevel.Ctx(ctx), normalizedIDs, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevel.Ctx(ctx).WhereIn(dao.MemberLevel.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取会员等级配置详情
func (s *sLevel) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.LevelDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberLevel.Ctx(ctx), id, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId, "会员等级配置"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberLevel.Ctx(ctx), id, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.LevelDetailOutput{}
	err = dao.MemberLevel.Ctx(ctx).Where(dao.MemberLevel.Columns().Id, id).Where(dao.MemberLevel.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("会员等级配置不存在或已删除")
	}
	// 查询租户关联显示
	if out.TenantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_tenant").Where("id", out.TenantID)
		val, err := refQuery.Value("name")
		if err == nil {
			out.TenantName = val.String()
		}
	}
	// 查询商户关联显示
	if out.MerchantID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("system_merchant").Where("id", out.MerchantID)
		val, err := refQuery.Value("name")
		if err == nil {
			out.MerchantName = val.String()
		}
	}
	return
}

// applyListFilter 应用列表通用过滤条件
func (s *sLevel) applyListFilter(ctx context.Context, in *model.LevelListInput) *gdb.Model {
	m := dao.MemberLevel.Ctx(ctx).Where(dao.MemberLevel.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId)
	if in.Name != "" {
		m = m.WhereLike(dao.MemberLevel.Columns().Name, "%"+in.Name+"%")
	}
	if in.LevelNo != "" {
		m = m.Where(dao.MemberLevel.Columns().LevelNo, in.LevelNo)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberLevel.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberLevel.Columns().MerchantId, *in.MerchantID)
	}
	if in.IsTop != nil {
		m = m.Where(dao.MemberLevel.Columns().IsTop, *in.IsTop)
	}
	if in.AutoDeploy != nil {
		m = m.Where(dao.MemberLevel.Columns().AutoDeploy, *in.AutoDeploy)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberLevel.Columns().Status, *in.Status)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberLevel.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberLevel.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sLevel) fillRefFields(ctx context.Context, list []*model.LevelListOutput) {
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

// List 获取会员等级配置列表
func (s *sLevel) List(ctx context.Context, in *model.LevelListInput) (list []*model.LevelListOutput, total int, err error) {
	if in == nil {
		in = &model.LevelListInput{}
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
func (s *sLevel) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberLevel.Columns().Id:        true,
		dao.MemberLevel.Columns().CreatedAt: true,
		dao.MemberLevel.Columns().Sort:      true,
		dao.MemberLevel.Columns().Status:    true,
		dao.MemberLevel.Columns().Name: true,
		dao.MemberLevel.Columns().NeedTeamTurnover: true,
		dao.MemberLevel.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sLevel) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.MemberLevel.Columns().Sort).OrderDesc(dao.MemberLevel.Columns().Id)
}

// Export 导出会员等级配置（不分页）
func (s *sLevel) Export(ctx context.Context, in *model.LevelListInput) (list []*model.LevelListOutput, err error) {
	if in == nil {
		in = &model.LevelListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑会员等级配置
func (s *sLevel) BatchUpdate(ctx context.Context, in *model.LevelBatchUpdateInput) error {
	data := do.MemberLevel{}
	hasChange := false
	if in.IsTop != nil {
		data.IsTop = *in.IsTop
		hasChange = true
	}
	if in.AutoDeploy != nil {
		data.AutoDeploy = *in.AutoDeploy
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeLevelIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberLevel.Ctx(ctx), normalizedIDs, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().TenantId, dao.MemberLevel.Columns().MerchantId, "会员等级配置"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberLevel.Ctx(ctx), normalizedIDs, dao.MemberLevel.Columns().Id, dao.MemberLevel.Columns().CreatedBy, dao.MemberLevel.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberLevel.Ctx(ctx).WhereIn(dao.MemberLevel.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入会员等级配置
func (s *sLevel) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.MemberLevel{
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
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.LevelNo = v
			}
		}
		idx++
		if idx < len(record) {
			data.Icon = strings.TrimSpace(record[idx])
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.DurationDays = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.NeedActiveCount = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); parseErr == nil {
				data.NeedTeamTurnover = int64(math.Round(v * 100))
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.IsTop = v
			}
		}
		idx++
		if idx < len(record) {
			if v, parseErr := strconv.Atoi(strings.TrimSpace(record[idx])); parseErr == nil {
				data.AutoDeploy = v
			}
		}
		idx++
		if idx < len(record) {
			data.Remark = strings.TrimSpace(record[idx])
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
		if _, insertErr := dao.MemberLevel.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

