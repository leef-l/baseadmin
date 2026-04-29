
package campaign

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
	service.RegisterCampaign(New())
}

func New() *sCampaign {
	return &sCampaign{}
}

type sCampaign struct{}

func normalizeCampaignIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建体验活动
func (s *sCampaign) Create(ctx context.Context, in *model.CampaignCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.DemoCampaign.Ctx(ctx).Data(do.DemoCampaign{
		Id:        id,
		CampaignNo: in.CampaignNo,
		Title: in.Title,
		Banner: in.Banner,
		Type: in.Type,
		Channel: in.Channel,
		BudgetAmount: in.BudgetAmount,
		LandingUrl: in.LandingURL,
		RuleJson: in.RuleJSON,
		IntroContent: in.IntroContent,
		StartAt: in.StartAt,
		EndAt: in.EndAt,
		IsPublic: in.IsPublic,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新体验活动
func (s *sCampaign) Update(ctx context.Context, in *model.CampaignUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCampaign.Ctx(ctx), in.ID, dao.DemoCampaign.Columns().Id, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId, "体验活动"); err != nil {
		return err
	}
	data := do.DemoCampaign{
		CampaignNo: in.CampaignNo,
		Title: in.Title,
		Banner: in.Banner,
		Type: in.Type,
		Channel: in.Channel,
		BudgetAmount: in.BudgetAmount,
		LandingUrl: in.LandingURL,
		RuleJson: in.RuleJSON,
		IntroContent: in.IntroContent,
		StartAt: in.StartAt,
		EndAt: in.EndAt,
		IsPublic: in.IsPublic,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
	}
	// 含金额字段，使用事务 + 行锁保证并发安全
	err := dao.DemoCampaign.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// FOR UPDATE 行锁
		_, err := tx.Model(dao.DemoCampaign.Table()).Ctx(ctx).
			Where(dao.DemoCampaign.Columns().Id, in.ID).
			LockUpdate().
			Value(dao.DemoCampaign.Columns().Id)
		if err != nil {
			return err
		}
		_, err = tx.Model(dao.DemoCampaign.Table()).Ctx(ctx).
			Where(dao.DemoCampaign.Columns().Id, in.ID).
			Data(data).Update()
		return err
	})
	return err
}

// Delete 软删除体验活动
func (s *sCampaign) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCampaign.Ctx(ctx), id, dao.DemoCampaign.Columns().Id, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId, "体验活动"); err != nil {
		return err
	}
	_, err := dao.DemoCampaign.Ctx(ctx).Where(dao.DemoCampaign.Columns().Id, id).Delete()
	return err
}

// BatchDelete 批量软删除体验活动
func (s *sCampaign) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	normalizedIDs := normalizeCampaignIDs(ids)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoCampaign.Ctx(ctx), normalizedIDs, dao.DemoCampaign.Columns().Id, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId, "体验活动"); err != nil {
		return err
	}
	_, err := dao.DemoCampaign.Ctx(ctx).WhereIn(dao.DemoCampaign.Columns().Id, normalizedIDs).Delete()
	return err
}

// Detail 获取体验活动详情
func (s *sCampaign) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.CampaignDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.DemoCampaign.Ctx(ctx), id, dao.DemoCampaign.Columns().Id, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId, "体验活动"); err != nil {
		return nil, err
	}
	out = &model.CampaignDetailOutput{}
	err = dao.DemoCampaign.Ctx(ctx).Where(dao.DemoCampaign.Columns().Id, id).Where(dao.DemoCampaign.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out.ID == 0 {
		return nil, nil
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
func (s *sCampaign) applyListFilter(ctx context.Context, in *model.CampaignListInput) *gdb.Model {
	m := dao.DemoCampaign.Ctx(ctx).Where(dao.DemoCampaign.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId)
	if in.CampaignNo != "" {
		m = m.Where(dao.DemoCampaign.Columns().CampaignNo, in.CampaignNo)
	}
	if in.Title != "" {
		m = m.WhereLike(dao.DemoCampaign.Columns().Title, "%"+in.Title+"%")
	}
	if in.TenantID != nil {
		m = m.Where(dao.DemoCampaign.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.DemoCampaign.Columns().MerchantId, *in.MerchantID)
	}
	if in.Type != nil {
		m = m.Where(dao.DemoCampaign.Columns().Type, *in.Type)
	}
	if in.Channel != nil {
		m = m.Where(dao.DemoCampaign.Columns().Channel, *in.Channel)
	}
	if in.IsPublic != nil {
		m = m.Where(dao.DemoCampaign.Columns().IsPublic, *in.IsPublic)
	}
	if in.Status != nil {
		m = m.Where(dao.DemoCampaign.Columns().Status, *in.Status)
	}
	if in.StartAtStart != "" {
		m = m.WhereGTE(dao.DemoCampaign.Columns().StartAt, in.StartAtStart)
	}
	if in.StartAtEnd != "" {
		m = m.WhereLTE(dao.DemoCampaign.Columns().StartAt, in.StartAtEnd)
	}
	if in.EndAtStart != "" {
		m = m.WhereGTE(dao.DemoCampaign.Columns().EndAt, in.EndAtStart)
	}
	if in.EndAtEnd != "" {
		m = m.WhereLTE(dao.DemoCampaign.Columns().EndAt, in.EndAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.DemoCampaign.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.DemoCampaign.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.DemoCampaign.Columns().CreatedBy, dao.DemoCampaign.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sCampaign) fillRefFields(ctx context.Context, list []*model.CampaignListOutput) {
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

// List 获取体验活动列表
func (s *sCampaign) List(ctx context.Context, in *model.CampaignListInput) (list []*model.CampaignListOutput, total int, err error) {
	if in == nil {
		in = &model.CampaignListInput{}
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
func (s *sCampaign) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.DemoCampaign.Columns().Id:        true,
		dao.DemoCampaign.Columns().CreatedAt: true,
		dao.DemoCampaign.Columns().Status:    true,
		dao.DemoCampaign.Columns().CampaignNo: true,
		dao.DemoCampaign.Columns().Title: true,
		dao.DemoCampaign.Columns().BudgetAmount: true,
	}
	return allowed[field]
}

func (s *sCampaign) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderDesc(dao.DemoCampaign.Columns().Id)
}

// Export 导出体验活动（不分页）
func (s *sCampaign) Export(ctx context.Context, in *model.CampaignListInput) (list []*model.CampaignListOutput, err error) {
	if in == nil {
		in = &model.CampaignListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// BatchUpdate 批量编辑体验活动
func (s *sCampaign) BatchUpdate(ctx context.Context, in *model.CampaignBatchUpdateInput) error {
	data := do.DemoCampaign{}
	if in.Type != nil {
		data.Type = *in.Type
	}
	if in.Channel != nil {
		data.Channel = *in.Channel
	}
	if in.IsPublic != nil {
		data.IsPublic = *in.IsPublic
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	normalizedIDs := normalizeCampaignIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.DemoCampaign.Ctx(ctx), normalizedIDs, dao.DemoCampaign.Columns().Id, dao.DemoCampaign.Columns().TenantId, dao.DemoCampaign.Columns().MerchantId, "体验活动"); err != nil {
		return err
	}
	_, err := dao.DemoCampaign.Ctx(ctx).WhereIn(dao.DemoCampaign.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}

// Import 导入体验活动
func (s *sCampaign) Import(ctx context.Context, file *ghttp.UploadFile) (success int, fail int, err error) {
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
		data := do.DemoCampaign{
			Id: id,
			CreatedBy: middleware.GetUserID(ctx),
			DeptId: middleware.GetDeptID(ctx),
		}
		idx := 0
		if idx < len(record) {
			data.CampaignNo = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Title = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Banner = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Type = record[idx]
		}
		idx++
		if idx < len(record) {
			data.Channel = record[idx]
		}
		idx++
		if idx < len(record) {
			data.BudgetAmount = record[idx]
		}
		idx++
		if idx < len(record) {
			data.LandingUrl = record[idx]
		}
		idx++
		if idx < len(record) {
			data.RuleJson = record[idx]
		}
		idx++
		if idx < len(record) {
			data.IntroContent = record[idx]
		}
		idx++
		if idx < len(record) {
			data.IsPublic = record[idx]
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
		if _, insertErr := dao.DemoCampaign.Ctx(ctx).Data(data).Insert(); insertErr != nil {
			fail++
		} else {
			success++
		}
	}
	return
}

