
package user

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"

	"gbaseadmin/app/member/internal/dao"
	"gbaseadmin/app/member/internal/middleware"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/app/member/internal/model/do"
	"gbaseadmin/app/member/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}

type sUser struct{}

func normalizeUserIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
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

// Create 创建会员用户
func (s *sUser) Create(ctx context.Context, in *model.UserCreateInput) error {
	id := snowflake.Generate()
	middleware.ApplyTenantScopeToWrite(ctx, &in.TenantID, &in.MerchantID)
	if err := middleware.EnsureTenantMerchantAccessible(ctx, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.MemberUser.Ctx(ctx).Data(do.MemberUser{
		Id:        id,
		ParentId: in.ParentID,
		Username: in.Username,
		Password: string(hashedPassword),
		Nickname: in.Nickname,
		Phone: in.Phone,
		Avatar: in.Avatar,
		RealName: in.RealName,
		LevelId: in.LevelID,
		LevelExpireAt: in.LevelExpireAt,
		TeamCount: in.TeamCount,
		DirectCount: in.DirectCount,
		ActiveCount: in.ActiveCount,
		TeamTurnover: in.TeamTurnover,
		IsActive: in.IsActive,
		IsQualified: in.IsQualified,
		InviteCode: in.InviteCode,
		RegisterIp: in.RegisterIP,
		LastLoginAt: in.LastLoginAt,
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

// Update 更新会员用户
func (s *sUser) Update(ctx context.Context, in *model.UserUpdateInput) error {
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
	data := do.MemberUser{
		ParentId: in.ParentID,
		Username: in.Username,
		Nickname: in.Nickname,
		Phone: in.Phone,
		Avatar: in.Avatar,
		RealName: in.RealName,
		LevelId: in.LevelID,
		LevelExpireAt: in.LevelExpireAt,
		TeamCount: in.TeamCount,
		DirectCount: in.DirectCount,
		ActiveCount: in.ActiveCount,
		TeamTurnover: in.TeamTurnover,
		IsActive: in.IsActive,
		IsQualified: in.IsQualified,
		InviteCode: in.InviteCode,
		RegisterIp: in.RegisterIP,
		LastLoginAt: in.LastLoginAt,
		Remark: in.Remark,
		Sort: in.Sort,
		Status: in.Status,
	}
	if in.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hashed)
	}
	if err := middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberUser.Ctx(ctx), in.ID, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId, "会员用户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberUser.Ctx(ctx), in.ID, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberUser.Ctx(ctx).Where(dao.MemberUser.Columns().Id, in.ID).Where(dao.MemberUser.Columns().DeletedAt, nil).Data(data).Update()
	return err
}

// Delete 软删除会员用户
func (s *sUser) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	deleteIDs, err := s.collectDeleteIDs(ctx, []snowflake.JsonInt64{id})
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), deleteIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId, "会员用户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), deleteIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId); err != nil {
		return err
	}
	_, err = dao.MemberUser.Ctx(ctx).WhereIn(dao.MemberUser.Columns().Id, deleteIDs).Delete()
	return err
}

// BatchDelete 批量软删除会员用户
func (s *sUser) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	deleteIDs, err := s.collectDeleteIDs(ctx, ids)
	if err != nil {
		return err
	}
	if len(deleteIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), deleteIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId, "会员用户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), deleteIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId); err != nil {
		return err
	}
	_, err = dao.MemberUser.Ctx(ctx).WhereIn(dao.MemberUser.Columns().Id, deleteIDs).Delete()
	return err
}

// collectDeleteIDs 汇总批量删除所需的节点 ID，并补齐所有子节点
func (s *sUser) collectDeleteIDs(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := normalizeUserIDs(ids)
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
func (s *sUser) collectChildIDs(ctx context.Context, parentID snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	return s.doCollectChildIDs(ctx, parentID, 0)
}

func (s *sUser) doCollectChildIDs(ctx context.Context, parentID snowflake.JsonInt64, depth int) ([]snowflake.JsonInt64, error) {
	if depth > 20 {
		return nil, gerror.New("子树深度超过 20 层上限，请检查数据完整性")
	}
	var childIDs []snowflake.JsonInt64
	m := dao.MemberUser.Ctx(ctx).
		Where(dao.MemberUser.Columns().ParentId, parentID).
		Where(dao.MemberUser.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId)
	result, err := m.Fields(dao.MemberUser.Columns().Id).
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

// Detail 获取会员用户详情
func (s *sUser) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.UserDetailOutput, err error) {
	if err = middleware.EnsureTenantScopedRowAccessible(ctx, dao.MemberUser.Ctx(ctx), id, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId, "会员用户"); err != nil {
		return nil, err
	}
	if err = middleware.EnsureDataScopedRowAccessible(ctx, dao.MemberUser.Ctx(ctx), id, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId); err != nil {
		return nil, err
	}
	out = &model.UserDetailOutput{}
	err = dao.MemberUser.Ctx(ctx).Where(dao.MemberUser.Columns().Id, id).Where(dao.MemberUser.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("会员用户不存在或已删除")
	}
	// 查询上级会员关联显示
	if out.ParentID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_user").Where("id", out.ParentID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("username")
		if err == nil {
			out.UserUsername = val.String()
		}
	}
	// 查询当前等级关联显示
	if out.LevelID != 0 {
		refQuery := g.DB().Ctx(ctx).Model("member_level").Where("id", out.LevelID)
		refQuery = refQuery.Where("deleted_at", nil)
		refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
		val, err := refQuery.Value("name")
		if err == nil {
			out.LevelName = val.String()
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
func (s *sUser) applyListFilter(ctx context.Context, in *model.UserListInput) *gdb.Model {
	m := dao.MemberUser.Ctx(ctx).Where(dao.MemberUser.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.MemberUser.Columns().Username, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().Nickname, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().RealName, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().Phone, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Username != "" {
		m = m.Where(dao.MemberUser.Columns().Username, in.Username)
	}
	if in.InviteCode != "" {
		m = m.Where(dao.MemberUser.Columns().InviteCode, in.InviteCode)
	}
	if in.Nickname != "" {
		m = m.WhereLike(dao.MemberUser.Columns().Nickname, "%"+in.Nickname+"%")
	}
	if in.RealName != "" {
		m = m.WhereLike(dao.MemberUser.Columns().RealName, "%"+in.RealName+"%")
	}
	if in.Phone != "" {
		m = m.Where(dao.MemberUser.Columns().Phone, in.Phone)
	}
	if in.ParentID != nil {
		m = m.Where(dao.MemberUser.Columns().ParentId, *in.ParentID)
	}
	if in.LevelID != nil {
		m = m.Where(dao.MemberUser.Columns().LevelId, *in.LevelID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberUser.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberUser.Columns().MerchantId, *in.MerchantID)
	}
	if in.IsActive != nil {
		m = m.Where(dao.MemberUser.Columns().IsActive, *in.IsActive)
	}
	if in.IsQualified != nil {
		m = m.Where(dao.MemberUser.Columns().IsQualified, *in.IsQualified)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberUser.Columns().Status, *in.Status)
	}
	if in.LevelExpireAtStart != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().LevelExpireAt, in.LevelExpireAtStart)
	}
	if in.LevelExpireAtEnd != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().LevelExpireAt, in.LevelExpireAtEnd)
	}
	if in.LastLoginAtStart != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().LastLoginAt, in.LastLoginAtStart)
	}
	if in.LastLoginAtEnd != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().LastLoginAt, in.LastLoginAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId)
	return m
}

// fillRefFields 批量填充关联显示字段（避免 N+1 查询）
func (s *sUser) fillRefFields(ctx context.Context, list []*model.UserListOutput) {
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
			refQuery := g.DB().Ctx(ctx).Model("member_user").
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
					if val, ok := refMap[int64(item.ParentID)]; ok {
						item.UserUsername = val
					}
				}
			}
		}
	}
	{
		idSet := make(map[int64]struct{})
		for _, item := range list {
			if item.LevelID != 0 {
				idSet[int64(item.LevelID)] = struct{}{}
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
					if val, ok := refMap[int64(item.LevelID)]; ok {
						item.LevelName = val
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

// List 获取会员用户列表
func (s *sUser) List(ctx context.Context, in *model.UserListInput) (list []*model.UserListOutput, total int, err error) {
	if in == nil {
		in = &model.UserListInput{}
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
func (s *sUser) isAllowedOrderField(field string) bool {
	allowed := map[string]bool{
		dao.MemberUser.Columns().Id:        true,
		dao.MemberUser.Columns().CreatedAt: true,
		dao.MemberUser.Columns().Sort:      true,
		dao.MemberUser.Columns().Status:    true,
		dao.MemberUser.Columns().Username: true,
		dao.MemberUser.Columns().Nickname: true,
		dao.MemberUser.Columns().Phone: true,
		dao.MemberUser.Columns().RealName: true,
		dao.MemberUser.Columns().Remark: true,
	}
	return allowed[field]
}

func (s *sUser) applyListOrder(m *gdb.Model, orderBy, orderDir string) *gdb.Model {
	if orderBy != "" && s.isAllowedOrderField(orderBy) {
		if orderDir == "desc" {
			return m.OrderDesc(orderBy)
		}
		return m.OrderAsc(orderBy)
	}
	return m.OrderAsc(dao.MemberUser.Columns().Sort).OrderDesc(dao.MemberUser.Columns().Id)
}

// Export 导出会员用户（不分页）
func (s *sUser) Export(ctx context.Context, in *model.UserListInput) (list []*model.UserListOutput, err error) {
	if in == nil {
		in = &model.UserListInput{}
	}
	m := s.applyListFilter(ctx, in)
	err = s.applyListOrder(m, in.OrderBy, in.OrderDir).Limit(10000).Scan(&list)
	if err != nil {
		return
	}
	s.fillRefFields(ctx, list)
	return
}

// Tree 获取会员用户树形结构
func (s *sUser) Tree(ctx context.Context, in *model.UserTreeInput) (tree []*model.UserTreeOutput, err error) {
	var list []*model.UserTreeOutput
	if in == nil {
		in = &model.UserTreeInput{}
	}
	m := dao.MemberUser.Ctx(ctx).Where(dao.MemberUser.Columns().DeletedAt, nil)
	m = middleware.ApplyTenantScopeToModel(ctx, m, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId)
	if in.Keyword != "" {
		keywordBuilder := m.Builder()
		keywordBuilder = keywordBuilder.WhereLike(dao.MemberUser.Columns().Username, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().Nickname, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().RealName, "%"+in.Keyword+"%")
		keywordBuilder = keywordBuilder.WhereOrLike(dao.MemberUser.Columns().Phone, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Username != "" {
		m = m.Where(dao.MemberUser.Columns().Username, in.Username)
	}
	if in.InviteCode != "" {
		m = m.Where(dao.MemberUser.Columns().InviteCode, in.InviteCode)
	}
	if in.Nickname != "" {
		m = m.WhereLike(dao.MemberUser.Columns().Nickname, "%"+in.Nickname+"%")
	}
	if in.RealName != "" {
		m = m.WhereLike(dao.MemberUser.Columns().RealName, "%"+in.RealName+"%")
	}
	if in.Phone != "" {
		m = m.Where(dao.MemberUser.Columns().Phone, in.Phone)
	}
	if in.ParentID != nil {
		m = m.Where(dao.MemberUser.Columns().ParentId, *in.ParentID)
	}
	if in.LevelID != nil {
		m = m.Where(dao.MemberUser.Columns().LevelId, *in.LevelID)
	}
	if in.TenantID != nil {
		m = m.Where(dao.MemberUser.Columns().TenantId, *in.TenantID)
	}
	if in.MerchantID != nil {
		m = m.Where(dao.MemberUser.Columns().MerchantId, *in.MerchantID)
	}
	if in.IsActive != nil {
		m = m.Where(dao.MemberUser.Columns().IsActive, *in.IsActive)
	}
	if in.IsQualified != nil {
		m = m.Where(dao.MemberUser.Columns().IsQualified, *in.IsQualified)
	}
	if in.Status != nil {
		m = m.Where(dao.MemberUser.Columns().Status, *in.Status)
	}
	if in.LevelExpireAtStart != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().LevelExpireAt, in.LevelExpireAtStart)
	}
	if in.LevelExpireAtEnd != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().LevelExpireAt, in.LevelExpireAtEnd)
	}
	if in.LastLoginAtStart != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().LastLoginAt, in.LastLoginAtStart)
	}
	if in.LastLoginAtEnd != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().LastLoginAt, in.LastLoginAtEnd)
	}
	if in.StartTime != "" {
		m = m.WhereGTE(dao.MemberUser.Columns().CreatedAt, in.StartTime)
	}
	if in.EndTime != "" {
		m = m.WhereLTE(dao.MemberUser.Columns().CreatedAt, in.EndTime)
	}
	// 数据权限过滤
	m = middleware.ApplyDataScope(ctx, m, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId)
	err = m.OrderAsc(dao.MemberUser.Columns().Sort).Limit(5000).Scan(&list)
	if err != nil {
		return
	}

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.UserTreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.UserTreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.UserTreeOutput, 0)
	for _, item := range list {
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
	// 批量填充当前等级关联显示
	{
		idSet := make(map[int64]struct{})
		var collectIDs func(items []*model.UserTreeOutput)
		collectIDs = func(items []*model.UserTreeOutput) {
			for _, item := range items {
				if item.LevelID != 0 {
					idSet[int64(item.LevelID)] = struct{}{}
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
			refQuery := g.DB().Ctx(ctx).Model("member_level").
				Fields("id", "name")
			refQuery = refQuery.Where("deleted_at", nil)
			refQuery = middleware.ApplyTenantScopeToModel(ctx, refQuery, "tenant_id", "merchant_id")
			rows, queryErr := refQuery.WhereIn("id", ids).All()
			if queryErr == nil {
				refMap := make(map[int64]string, len(rows))
				for _, row := range rows {
					refMap[row["id"].Int64()] = row["name"].String()
				}
				var fillRef func(items []*model.UserTreeOutput)
				fillRef = func(items []*model.UserTreeOutput) {
					for _, item := range items {
						if val, ok := refMap[int64(item.LevelID)]; ok {
							item.LevelName = val
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
	// 批量填充租户关联显示
	{
		idSet := make(map[int64]struct{})
		var collectIDs func(items []*model.UserTreeOutput)
		collectIDs = func(items []*model.UserTreeOutput) {
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
				var fillRef func(items []*model.UserTreeOutput)
				fillRef = func(items []*model.UserTreeOutput) {
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
		var collectIDs func(items []*model.UserTreeOutput)
		collectIDs = func(items []*model.UserTreeOutput) {
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
				var fillRef func(items []*model.UserTreeOutput)
				fillRef = func(items []*model.UserTreeOutput) {
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

// BatchUpdate 批量编辑会员用户
func (s *sUser) BatchUpdate(ctx context.Context, in *model.UserBatchUpdateInput) error {
	data := do.MemberUser{}
	hasChange := false
	if in.IsActive != nil {
		data.IsActive = *in.IsActive
		hasChange = true
	}
	if in.IsQualified != nil {
		data.IsQualified = *in.IsQualified
		hasChange = true
	}
	if in.Status != nil {
		data.Status = *in.Status
		hasChange = true
	}
	if !hasChange {
		return nil
	}
	normalizedIDs := normalizeUserIDs(in.IDs)
	if len(normalizedIDs) == 0 {
		return nil
	}
	if err := middleware.EnsureTenantScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), normalizedIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().TenantId, dao.MemberUser.Columns().MerchantId, "会员用户"); err != nil {
		return err
	}
	if err := middleware.EnsureDataScopedRowsAccessible(ctx, dao.MemberUser.Ctx(ctx), normalizedIDs, dao.MemberUser.Columns().Id, dao.MemberUser.Columns().CreatedBy, dao.MemberUser.Columns().DeptId); err != nil {
		return err
	}
	_, err := dao.MemberUser.Ctx(ctx).WhereIn(dao.MemberUser.Columns().Id, normalizedIDs).Data(data).Update()
	return err
}
