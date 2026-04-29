
package shop_category

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterShopCategory(New())
}

func New() *sShopCategory {
	return &sShopCategory{}
}

type sShopCategory struct{}

func normalizeShopCategoryIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建商城商品分类
func (s *sShopCategory) Create(ctx context.Context, in *model.ShopCategoryCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	_, err := dao.MemberShopCategory.Ctx(ctx).Data(do.MemberShopCategory{
		Id:        id,
		ParentId: in.ParentID,
		Name: in.Name,
		Icon: in.Icon,
		Sort: in.Sort,
		Status: in.Status,
		TenantId: in.TenantID,
		MerchantId: in.MerchantID,
		CreatedBy: middleware.GetUserID(ctx),
		DeptId: middleware.GetDeptID(ctx),
	}).Insert()
	return err
}

// Update 更新商城商品分类
func (s *sShopCategory) Update(ctx context.Context, in *model.ShopCategoryUpdateInput) error {
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	if in.ParentID == in.ID {
		return gerror.New("不能将自身设为父节点")
	}
	if int64(in.ParentID) != 0 {
		childIDs, collectErr := s.collectChildIDs(ctx, in.ID)
		if collectErr != nil {
			return collectErr
		}
		for _, cid := range childIDs {
			if cid == in.ParentID {
				return gerror.New("不能将子节点设为父节点，会形成循环引用")
			}
		}
	}
	data := do.MemberShopCategory{
		ParentId: in.ParentID,
		Name: in.Name,
		Icon: in.Icon,
		Sort: in.Sort,
		Status: in.Status,
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), in.ID, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId, "商城商品分类"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), in.ID, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberShopCategory.Ctx(ctx).Where(dao.MemberShopCategory.Columns().Id, in.ID).Where(dao.MemberShopCategory.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除商城商品分类
func (s *sShopCategory) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	deleteIDs, err := s.collectDeleteIDs(ctx, []snowflake.JsonInt64{id})
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), deleteIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId, "商城商品分类"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), deleteIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId); err != nil {
		return err
	}
	_, err = dao.MemberShopCategory.Ctx(ctx).WhereIn(dao.MemberShopCategory.Columns().Id, deleteIDs).Delete()
	return err
}

// BatchDelete 批量软删除商城商品分类
func (s *sShopCategory) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	deleteIDs, err := s.collectDeleteIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), deleteIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId, "商城商品分类"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), deleteIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId); err != nil {
		return err
	}
	_, err = dao.MemberShopCategory.Ctx(ctx).WhereIn(dao.MemberShopCategory.Columns().Id, deleteIDs).Delete()
	return err
}

// collectDeleteIDs 汇总批量删除所需的节点 ID，并补齐所有子节点
func (s *sShopCategory) collectDeleteIDs(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := normalizeShopCategoryIDs(ids)
	if len(normalized) == 0 {
		return nil, nil
	}
	const maxCollect = 10000
	collected := make([]snowflake.JsonInt64, 0, len(normalized))
	seen := make(map[int64]struct{}, len(normalized))
	for _, id := range normalized {
		if _, ok := seen[int64(id)]; !ok {
			seen[int64(id)] = struct{}{}
			collected = append(collected, id)
		}
		childIDs, err := s.collectChildIDs(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, childID := range childIDs {
			if _, ok := seen[int64(childID)]; ok {
				continue
			}
			seen[int64(childID)] = struct{}{}
			collected = append(collected, childID)
			if len(collected) > maxCollect {
				return nil, gerror.Newf("子树节点过多（超过 %d），请分批删除", maxCollect)
			}
		}
	}
	return collected, nil
}

// collectChildIDs 递归收集所有子节点 ID（最大深度 20 层防止无限递归）
func (s *sShopCategory) collectChildIDs(ctx context.Context, parentID snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	return s.doCollectChildIDs(ctx, parentID, 0)
}

func (s *sShopCategory) doCollectChildIDs(ctx context.Context, parentID snowflake.JsonInt64, depth int) ([]snowflake.JsonInt64, error) {
	if depth > 20 {
		return nil, gerror.New("子树深度超过 20 层上限，请检查数据完整性")
	}
	var childIDs []snowflake.JsonInt64
	m := dao.MemberShopCategory.Ctx(ctx).
		Where(dao.MemberShopCategory.Columns().ParentId, parentID).
		Where(dao.MemberShopCategory.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId)
	result, err := m.Fields(dao.MemberShopCategory.Columns().Id).
		Array()
	if err != nil || len(result) == 0 {
		return childIDs, err
	}
	for _, v := range result {
		cid := snowflake.JsonInt64(v.Int64())
		childIDs = append(childIDs, cid)
		grandChildren, err := s.doCollectChildIDs(ctx, cid, depth+1)
		if err != nil {
			return childIDs, err
		}
		childIDs = append(childIDs, grandChildren...)
	}
	return childIDs, nil
}

// Detail 获取商城商品分类详情
func (s *sShopCategory) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ShopCategoryDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), id, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId, "商城商品分类"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), id, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.ShopCategoryDetailOutput{}
	err = dao.MemberShopCategory.Ctx(ctx).Where(dao.MemberShopCategory.Columns().Id, id).Where(dao.MemberShopCategory.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("商城商品分类不存在或已删除")
	}
	// 查询上级分类关联显示
	if out.ParentID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_shop_category").Where("id", out.ParentID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.ShopCategoryName = val.String()
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
func (s *sShopCategory) applyListFilter(ctx context.Context, in *model.ShopCategoryListInput) *gdb.Model {
	m := dao.MemberShopCategory.Ctx(ctx).Where(dao.MemberShopCategory.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId)
	if in.Name != "" {
		m = m.WhereLike(dao.MemberShopCategory.Columns().Name, "%"+in.Name+"%")
	}
	if in.ParentID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().ParentId, *in.ParentID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberShopCategory.Columns().Status, *in.Status)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberShopCategory.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberShopCategory.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sShopCategory) fillRefFields(ctx context.Context, list []*model.ShopCategoryListOutput) {
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.ParentID != 0 {
				idSet[int64(item.ParentID)] = struct{}{}
			}
		}
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("member_shop_category").
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
					if val, ok := refMap[int64(item.ParentID)]; ok {
						item.ShopCategoryName = val
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

// List 获取商城商品分类列表
func (s *sShopCategory) List(ctx context.Context, in *model.ShopCategoryListInput) (list []*model.ShopCategoryListOutput, total int, err error) {
	if in == nil {
		in = &model.ShopCategoryListInput{}
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
func (s *sShopCategory) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberShopCategory.Columns().Id:        true,
		dao.MemberShopCategory.Columns().CreatedAt: true,
		dao.MemberShopCategory.Columns().Sort:      true,
		dao.MemberShopCategory.Columns().Status:    true,
		dao.MemberShopCategory.Columns().Name: true,
	}
	return allowed[field]
}

func (s *sShopCategory) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.MemberShopCategory.Columns().Sort).OrderDesc(dao.MemberShopCategory.Columns().Id)
}

// Export 导出商城商品分类（不分页）
func (s *sShopCategory) Export(ctx context.Context, in *model.ShopCategoryListInput) (list []*model.ShopCategoryListOutput, err error) {
	if in == nil {
		in = &model.ShopCategoryListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// Tree 获取商城商品分类树形结构
func (s *sShopCategory) Tree(ctx context.Context, in *model.ShopCategoryTreeInput) (tree []*model.ShopCategoryTreeOutput, err error) {
	var list []*model.ShopCategoryTreeOutput
	if in == nil {
		in = &model.ShopCategoryTreeInput{}
	}
	m := dao.MemberShopCategory.Ctx(ctx).Where(dao.MemberShopCategory.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId)
	if in.Name != "" {
		m = m.WhereLike(dao.MemberShopCategory.Columns().Name, "%"+in.Name+"%")
	}
	if in.ParentID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().ParentId, *in.ParentID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberShopCategory.Columns().MerchantId, *in.MerchantID)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberShopCategory.Columns().Status, *in.Status)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberShopCategory.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberShopCategory.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId)
	err = m.OrderAsc(dao.MemberShopCategory.Columns().Sort).Limit(5000).Scan(&list)
	if err != nil {
		return
	}

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.ShopCategoryTreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.ShopCategoryTreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.ShopCategoryTreeOutput, 0)
	for _, item := range list {
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
	// 批量填充租户关联显示
	{
		idSet := make(map[int64]struct{})
		var collectIDs func(items []*model.ShopCategoryTreeOutput)
		collectIDs = func(items []*model.ShopCategoryTreeOutput) {
			for _, item := range items {
				if item.TenantID != 0 {
					idSet[int64(item.TenantID)] = struct{}{}
				}
				if len(item.Children) > 0 {
					collectIDs(item.Children)
				}
			}
		}
		collectIDs(list)
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_tenant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, queryErr := refQuery.WhereIn("id", ids).All()
			if queryErr == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				var fillRef func(items []*model.ShopCategoryTreeOutput)
				fillRef = func(items []*model.ShopCategoryTreeOutput) {
					for _, item := range items {
						if val, ok := refMap[int64(item.TenantID)]; ok {
							item.TenantName = val
						}
						if len(item.Children) > 0 {
							fillRef(item.Children)
						}
					}
				}
				fillRef(list)
			}
		}
	}
	// 批量填充商户关联显示
	{
		idSet := make(map[int64]struct{})
		var collectIDs func(items []*model.ShopCategoryTreeOutput)
		collectIDs = func(items []*model.ShopCategoryTreeOutput) {
			for _, item := range items {
				if item.MerchantID != 0 {
					idSet[int64(item.MerchantID)] = struct{}{}
				}
				if len(item.Children) > 0 {
					collectIDs(item.Children)
				}
			}
		}
		collectIDs(list)
		if len(idSet) > 0 {
			ids := make([]int64, 0, len(idSet))
			for id := range idSet {
				ids = append(ids, id)
			}
			refQuery := g.DB().Ctx(ctx).Model("system_merchant").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, queryErr := refQuery.WhereIn("id", ids).All()
			if queryErr == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				var fillRef func(items []*model.ShopCategoryTreeOutput)
				fillRef = func(items []*model.ShopCategoryTreeOutput) {
					for _, item := range items {
						if val, ok := refMap[int64(item.MerchantID)]; ok {
							item.MerchantName = val
						}
						if len(item.Children) > 0 {
							fillRef(item.Children)
						}
					}
				}
				fillRef(list)
			}
		}
	}
	return
}

// BatchUpdate 批量编辑商城商品分类
func (s *sShopCategory) BatchUpdate(ctx context.Context, in *model.ShopCategoryBatchUpdateInput) error {
	data := do.MemberShopCategory{}
	hasChange := false
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeShopCategoryIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), normalizedIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().TenantId, dao.MemberShopCategory.Columns().MerchantId, "商城商品分类"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberShopCategory.Ctx(ctx), normalizedIDs, dao.MemberShopCategory.Columns().Id, dao.MemberShopCategory.Columns().CreatedBy, dao.MemberShopCategory.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberShopCategory.Ctx(ctx).WhereIn(dao.MemberShopCategory.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}
