package role

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/snowflake"
	"gbaseadmin/utility/treeutil"
)

func init() {
	service.RegisterRole(New())
}

func New() *sRole {
	return &sRole{}
}

type sRole struct{}

// Create 创建角色表
func (s *sRole) Create(ctx context.Context, in *model.RoleCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeRoleCreateInput(in)
	if err := validateRoleFields(in.Title, in.DataScope, in.Status, in.IsAdmin, in.Sort); err != nil {
		return err
	}
	if err := s.ensureAdminRoleMutationAllowed(ctx, in.IsAdmin == 1); err != nil {
		return err
	}
	if err := s.ensureRoleTitleUnique(ctx, 0, in.Title); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, 0); err != nil {
		return err
	}
	id := snowflake.Generate()
	_, err := dao.Role.Ctx(ctx).Data(g.Map{
		dao.Role.Columns().Id:        id,
		dao.Role.Columns().ParentId:  in.ParentID,
		dao.Role.Columns().Title:     in.Title,
		dao.Role.Columns().DataScope: in.DataScope,
		dao.Role.Columns().Sort:      in.Sort,
		dao.Role.Columns().Status:    in.Status,
		dao.Role.Columns().IsAdmin:   in.IsAdmin,
		dao.Role.Columns().CreatedAt: gtime.Now(),
		dao.Role.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新角色表
func (s *sRole) Update(ctx context.Context, in *model.RoleUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeRoleUpdateInput(in)
	if err := validateRoleFields(in.Title, in.DataScope, in.Status, in.IsAdmin, in.Sort); err != nil {
		return err
	}
	if err := s.ensureRoleExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureRoleTitleUnique(ctx, in.ID, in.Title); err != nil {
		return err
	}
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
	}
	isAdminRole, err := s.isAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleMutationAllowed(ctx, isAdminRole || in.IsAdmin == 1); err != nil {
		return err
	}
	isBuiltinAdminRole, err := s.isBuiltinAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if isBuiltinAdminRole {
		if in.Status == 0 {
			return gerror.New("内置管理员角色不能禁用")
		}
		if in.IsAdmin != 1 {
			return gerror.New("内置管理员角色不能取消超级管理员标记")
		}
	}
	data := g.Map{
		dao.Role.Columns().ParentId:  in.ParentID,
		dao.Role.Columns().Title:     in.Title,
		dao.Role.Columns().DataScope: in.DataScope,
		dao.Role.Columns().Sort:      in.Sort,
		dao.Role.Columns().Status:    in.Status,
		dao.Role.Columns().IsAdmin:   in.IsAdmin,
		dao.Role.Columns().UpdatedAt: gtime.Now(),
	}
	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			Where(dao.Role.Columns().Id, in.ID).
			Where(dao.Role.Columns().DeletedAt, nil).
			Data(data).
			Update(); err != nil {
			return err
		}
		if in.DataScope == 5 {
			return nil
		}
		_, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).
			Where(dao.RoleDept.Columns().RoleId, in.ID).
			Delete()
		return err
	})
	if err == nil {
		s.clearRoleUserCaches(ctx, in.ID)
	}
	return err
}

// Delete 软删除角色表
func (s *sRole) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureRoleExists(ctx, id); err != nil {
		return err
	}
	isAdminRole, err := s.isAdminRole(ctx, id)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleMutationAllowed(ctx, isAdminRole); err != nil {
		return err
	}
	isBuiltinAdminRole, err := s.isBuiltinAdminRole(ctx, id)
	if err != nil {
		return err
	}
	if isBuiltinAdminRole {
		return gerror.New("内置管理员角色不能删除")
	}
	if err := s.ensureRoleDeletable(ctx, id); err != nil {
		return err
	}
	userIDs, err := s.getRoleUserIDs(ctx, id)
	if err != nil {
		return err
	}
	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			Where(dao.Role.Columns().Id, id).
			Where(dao.Role.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.Role.Columns().DeletedAt: gtime.Now(),
			}).
			Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).
			Where(dao.RoleMenu.Columns().RoleId, id).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).
			Where(dao.RoleDept.Columns().RoleId, id).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).
			Where(dao.UserRole.Columns().RoleId, id).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, userIDs...)
	}
	return err
}

// BatchDelete 批量删除角色表
func (s *sRole) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的角色")
	}
	if err := s.ensureRoleIDsExist(ctx, ids); err != nil {
		return err
	}
	order, err := s.collectBatchDeleteOrder(ctx, ids)
	if err != nil {
		return err
	}
	if len(order) == 0 {
		return nil
	}
	deleteIDs := batchutil.ToInt64s(order)
	for _, id := range order {
		isAdminRole, err := s.isAdminRole(ctx, id)
		if err != nil {
			return err
		}
		if err := s.ensureAdminRoleMutationAllowed(ctx, isAdminRole); err != nil {
			return err
		}
		isBuiltinAdminRole, err := s.isBuiltinAdminRole(ctx, id)
		if err != nil {
			return err
		}
		if isBuiltinAdminRole {
			return gerror.New("内置管理员角色不能删除")
		}
		if err := s.ensureRoleBatchDeletable(ctx, id, deleteIDs); err != nil {
			return err
		}
	}
	userIDs, err := s.getRoleUserIDsByRoleIDs(ctx, deleteIDs)
	if err != nil {
		return err
	}
	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			WhereIn(dao.Role.Columns().Id, deleteIDs).
			Where(dao.Role.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.Role.Columns().DeletedAt: gtime.Now(),
			}).
			Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).
			WhereIn(dao.RoleMenu.Columns().RoleId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).
			WhereIn(dao.RoleDept.Columns().RoleId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).
			WhereIn(dao.UserRole.Columns().RoleId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, userIDs...)
	}
	return err
}

// Detail 获取角色表详情
func (s *sRole) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.RoleDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("角色不存在或已删除")
	}
	out = &model.RoleDetailOutput{}
	err = dao.Role.Ctx(ctx).Where(dao.Role.Columns().Id, id).Where(dao.Role.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("角色不存在或已删除")
	}
	out.RoleTitle = shared.LookupTitle(ctx, "system_role", int64(out.ParentID))
	return
}

// List 获取角色表列表
func (s *sRole) List(ctx context.Context, in *model.RoleListInput) (list []*model.RoleListOutput, total int, err error) {
	if in == nil {
		in = &model.RoleListInput{}
	}
	normalizeRoleListInput(in)
	m := dao.Role.Ctx(ctx).Where(dao.Role.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		m = m.WhereLike(dao.Role.Columns().Title, "%"+in.Keyword+"%")
	}
	if in.DataScope > 0 {
		m = m.Where(dao.Role.Columns().DataScope, in.DataScope)
	}
	if in.Status != nil {
		m = m.Where(dao.Role.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.Role.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillParentTitles(ctx, list)
	return
}

// Tree 获取角色表树形结构
func (s *sRole) Tree(ctx context.Context, in *model.RoleTreeInput) (tree []*model.RoleTreeOutput, err error) {
	var list []*model.RoleTreeOutput
	m := dao.Role.Ctx(ctx).Where(dao.Role.Columns().DeletedAt, nil)
	if in != nil {
		normalizeRoleTreeInput(in)
		if in.Keyword != "" {
			m = m.WhereLike(dao.Role.Columns().Title, "%"+in.Keyword+"%")
		}
		if in.DataScope > 0 {
			m = m.Where(dao.Role.Columns().DataScope, in.DataScope)
		}
		if in.Status != nil {
			m = m.Where(dao.Role.Columns().Status, *in.Status)
		}
	}
	err = m.OrderAsc(dao.Role.Columns().Sort).Scan(&list)
	if err != nil {
		return
	}

	tree = treeutil.BuildForest(list, treeutil.TreeNodeAccessor[*model.RoleTreeOutput]{
		ID:       func(item *model.RoleTreeOutput) int64 { return int64(item.ID) },
		ParentID: func(item *model.RoleTreeOutput) int64 { return int64(item.ParentID) },
		Init: func(item *model.RoleTreeOutput) {
			item.Children = make([]*model.RoleTreeOutput, 0)
		},
		Append: func(parent *model.RoleTreeOutput, child *model.RoleTreeOutput) {
			parent.Children = append(parent.Children, child)
		},
	})
	return
}

func normalizeRoleCreateInput(in *model.RoleCreateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
}

func normalizeRoleUpdateInput(in *model.RoleUpdateInput) {
	if in == nil {
		return
	}
	in.Title = strings.TrimSpace(in.Title)
}

func normalizeRoleListInput(in *model.RoleListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func normalizeRoleTreeInput(in *model.RoleTreeInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}

func validateRoleFields(title string, dataScope, status, isAdmin, sort int) error {
	if title == "" {
		return gerror.New("角色名称不能为空")
	}
	if err := fieldvalid.Enum("数据范围", dataScope, 1, 2, 3, 4, 5); err != nil {
		return err
	}
	if err := fieldvalid.NonNegative("排序", sort); err != nil {
		return err
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	if err := fieldvalid.Binary("超级管理员", isAdmin); err != nil {
		return err
	}
	return nil
}

// GrantMenu 角色授权菜单（先删后插）
func (s *sRole) GrantMenu(ctx context.Context, in *model.RoleGrantMenuInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := s.ensureRoleExists(ctx, in.ID); err != nil {
		return err
	}
	isAdminRole, err := s.isAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleMutationAllowed(ctx, isAdminRole); err != nil {
		return err
	}
	menuIDs, err := s.normalizeMenuIDs(ctx, in.MenuIDs)
	if err != nil {
		return err
	}
	err = dao.RoleMenu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Where(dao.RoleMenu.Columns().RoleId, in.ID).Delete(); err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		data := make([]g.Map, 0, len(menuIDs))
		for _, menuID := range menuIDs {
			data = append(data, g.Map{
				dao.RoleMenu.Columns().RoleId: in.ID,
				dao.RoleMenu.Columns().MenuId: menuID,
			})
		}
		_, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Data(data).Insert()
		return err
	})
	if err == nil {
		s.clearRoleUserCaches(ctx, in.ID)
	}
	return err
}

// GetMenuIDs 获取角色已授权的菜单ID列表
func (s *sRole) GetMenuIDs(ctx context.Context, roleID snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	if err := s.ensureRoleExists(ctx, roleID); err != nil {
		return nil, err
	}
	var list []struct {
		MenuId int64 `json:"menuId"`
	}
	err := dao.RoleMenu.Ctx(ctx).Where(dao.RoleMenu.Columns().RoleId, roleID).Scan(&list)
	if err != nil {
		return nil, err
	}
	ids := make([]snowflake.JsonInt64, 0, len(list))
	for _, item := range list {
		ids = append(ids, snowflake.JsonInt64(item.MenuId))
	}
	return ids, nil
}

// GrantDept 角色授权数据权限
func (s *sRole) GrantDept(ctx context.Context, in *model.RoleGrantDeptInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	if err := fieldvalid.Enum("数据范围", in.DataScope, 1, 2, 3, 4, 5); err != nil {
		return err
	}
	if err := s.ensureRoleExists(ctx, in.ID); err != nil {
		return err
	}
	isAdminRole, err := s.isAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleMutationAllowed(ctx, isAdminRole); err != nil {
		return err
	}
	isBuiltinAdminRole, err := s.isBuiltinAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if isBuiltinAdminRole && in.DataScope != 1 {
		return gerror.New("内置管理员角色的数据范围必须保持为全部")
	}
	deptIDs, err := s.normalizeDeptIDs(ctx, in.DeptIDs)
	if err != nil {
		return err
	}
	if in.DataScope != 5 {
		deptIDs = nil
	}
	err = dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			Where(dao.Role.Columns().Id, in.ID).
			Where(dao.Role.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.Role.Columns().DataScope: in.DataScope,
				dao.Role.Columns().UpdatedAt: gtime.Now(),
			}).
			Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).Where(dao.RoleDept.Columns().RoleId, in.ID).Delete(); err != nil {
			return err
		}
		if in.DataScope != 5 || len(deptIDs) == 0 {
			return nil
		}
		data := make([]g.Map, 0, len(deptIDs))
		for _, deptID := range deptIDs {
			data = append(data, g.Map{
				dao.RoleDept.Columns().RoleId: in.ID,
				dao.RoleDept.Columns().DeptId: deptID,
			})
		}
		_, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).Data(data).Insert()
		return err
	})
	if err == nil {
		s.clearRoleUserCaches(ctx, in.ID)
	}
	return err
}

// GetDeptIDs 获取角色已授权的部门ID列表
func (s *sRole) GetDeptIDs(ctx context.Context, roleID snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	if err := s.ensureRoleExists(ctx, roleID); err != nil {
		return nil, err
	}
	var list []struct {
		DeptId int64 `json:"deptId"`
	}
	err := dao.RoleDept.Ctx(ctx).Where(dao.RoleDept.Columns().RoleId, roleID).Scan(&list)
	if err != nil {
		return nil, err
	}
	ids := make([]snowflake.JsonInt64, 0, len(list))
	for _, item := range list {
		ids = append(ids, snowflake.JsonInt64(item.DeptId))
	}
	return ids, nil
}

func (s *sRole) fillParentTitles(ctx context.Context, list []*model.RoleListOutput) {
	parentIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.ParentID != 0 {
			parentIDs = append(parentIDs, int64(item.ParentID))
		}
	}
	parentMap := shared.LoadTitleMap(ctx, "system_role", parentIDs)
	for _, item := range list {
		item.RoleTitle = parentMap[int64(item.ParentID)]
	}
}

func (s *sRole) clearRoleUserCaches(ctx context.Context, roleID snowflake.JsonInt64) {
	userIDs, err := s.getRoleUserIDs(ctx, roleID)
	if err != nil {
		return
	}
	authlogic.ClearUserCaches(ctx, userIDs...)
}

func (s *sRole) ensureRoleExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("角色不存在或已删除")
	}
	count, err := dao.Role.Ctx(ctx).
		Where(dao.Role.Columns().Id, id).
		Where(dao.Role.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("角色不存在或已删除")
	}
	return nil
}

func (s *sRole) ensureAdminRoleMutationAllowed(ctx context.Context, touchesAdminRole bool) error {
	if !touchesAdminRole {
		return nil
	}
	actorHasAdmin, err := shared.HasCurrentActorAdminRole(ctx)
	if err != nil {
		return err
	}
	return validateAdminRoleMutationAllowed(touchesAdminRole, actorHasAdmin)
}

func (s *sRole) ensureRoleIDsExist(ctx context.Context, ids []snowflake.JsonInt64) error {
	dbIDs := batchutil.ToInt64s(ids)
	count, err := dao.Role.Ctx(ctx).
		WhereIn(dao.Role.Columns().Id, dbIDs).
		Where(dao.Role.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count != len(dbIDs) {
		return gerror.New("包含不存在或已删除的角色")
	}
	return nil
}

func (s *sRole) normalizeMenuIDs(ctx context.Context, menuIDs []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := batchutil.CompactIDs(menuIDs)
	if len(normalized) == 0 {
		return nil, nil
	}
	count, err := dao.Menu.Ctx(ctx).
		WhereIn(dao.Menu.Columns().Id, batchutil.ToInt64s(normalized)).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return nil, err
	}
	if count != len(normalized) {
		return nil, gerror.New("菜单不存在或已删除")
	}
	return normalized, nil
}

func (s *sRole) normalizeDeptIDs(ctx context.Context, deptIDs []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := batchutil.CompactIDs(deptIDs)
	if len(normalized) == 0 {
		return nil, nil
	}
	count, err := dao.Dept.Ctx(ctx).
		WhereIn(dao.Dept.Columns().Id, batchutil.ToInt64s(normalized)).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return nil, err
	}
	if count != len(normalized) {
		return nil, gerror.New("部门不存在或已删除")
	}
	return normalized, nil
}

func (s *sRole) getRoleUserIDs(ctx context.Context, roleID snowflake.JsonInt64) ([]int64, error) {
	var userRoles []struct {
		UserId int64 `json:"userId"`
	}
	if err := dao.UserRole.Ctx(ctx).
		Fields(dao.UserRole.Columns().UserId).
		Where(dao.UserRole.Columns().RoleId, roleID).
		Scan(&userRoles); err != nil {
		return nil, err
	}
	userIDs := make([]int64, 0, len(userRoles))
	for _, item := range userRoles {
		userIDs = append(userIDs, item.UserId)
	}
	return userIDs, nil
}

func (s *sRole) getRoleUserIDsByRoleIDs(ctx context.Context, roleIDs []int64) ([]int64, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var userRoles []struct {
		UserId int64 `json:"userId"`
	}
	if err := dao.UserRole.Ctx(ctx).
		Fields(dao.UserRole.Columns().UserId).
		WhereIn(dao.UserRole.Columns().RoleId, roleIDs).
		Scan(&userRoles); err != nil {
		return nil, err
	}
	seen := make(map[int64]struct{}, len(userRoles))
	userIDs := make([]int64, 0, len(userRoles))
	for _, item := range userRoles {
		if item.UserId <= 0 {
			continue
		}
		if _, ok := seen[item.UserId]; ok {
			continue
		}
		seen[item.UserId] = struct{}{}
		userIDs = append(userIDs, item.UserId)
	}
	if len(userIDs) == 0 {
		return nil, nil
	}
	return userIDs, nil
}

func (s *sRole) ensureRoleDeletable(ctx context.Context, id snowflake.JsonInt64) error {
	childCount, err := dao.Role.Ctx(ctx).
		Where(dao.Role.Columns().ParentId, id).
		Where(dao.Role.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前角色下存在子角色，不能直接删除")
	}
	return nil
}

func (s *sRole) ensureRoleTitleUnique(ctx context.Context, currentID snowflake.JsonInt64, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}
	m := dao.Role.Ctx(ctx).
		Where(dao.Role.Columns().Title, title).
		Where(dao.Role.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.Role.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("角色名称已存在")
	}
	return nil
}

func (s *sRole) ensureRoleBatchDeletable(ctx context.Context, id snowflake.JsonInt64, deleteIDs []int64) error {
	childModel := dao.Role.Ctx(ctx).
		Where(dao.Role.Columns().ParentId, id).
		Where(dao.Role.Columns().DeletedAt, nil)
	if len(deleteIDs) > 0 {
		childModel = childModel.WhereNotIn(dao.Role.Columns().Id, deleteIDs)
	}
	childCount, err := childModel.Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.New("当前角色下存在子角色，不能直接删除")
	}
	return nil
}

func (s *sRole) collectBatchDeleteOrder(ctx context.Context, ids []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	var rows []struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.Role.Ctx(ctx).
		Fields(dao.Role.Columns().Id, dao.Role.Columns().ParentId).
		Where(dao.Role.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil, err
	}
	treeRows := make([]batchutil.TreeRow, 0, len(rows))
	for _, row := range rows {
		treeRows = append(treeRows, batchutil.TreeRow{
			ID:       row.Id,
			ParentID: row.ParentId,
		})
	}
	return batchutil.ExpandTreeDeleteOrder(ids, treeRows), nil
}

func (s *sRole) isBuiltinAdminRole(ctx context.Context, roleID snowflake.JsonInt64) (bool, error) {
	if roleID <= 0 {
		return false, nil
	}
	count, err := g.DB().Ctx(ctx).
		Model("system_user_role ur").
		LeftJoin("system_users u", "u.id = ur.user_id").
		LeftJoin("system_role r", "r.id = ur.role_id").
		Where("ur.role_id", roleID).
		Where("u.username", "admin").
		Where("u.deleted_at", nil).
		Where("r.deleted_at", nil).
		Where("r.is_admin", 1).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sRole) isAdminRole(ctx context.Context, roleID snowflake.JsonInt64) (bool, error) {
	if roleID <= 0 {
		return false, nil
	}
	count, err := dao.Role.Ctx(ctx).
		Where(dao.Role.Columns().Id, roleID).
		Where(dao.Role.Columns().DeletedAt, nil).
		Where(dao.Role.Columns().IsAdmin, 1).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sRole) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	return treeutil.ValidateParent(parentID, currentID, func(id int64) (int64, int64, error) {
		var parent struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}
		if err := dao.Role.Ctx(ctx).
			Fields(dao.Role.Columns().Id, dao.Role.Columns().ParentId).
			Where(dao.Role.Columns().Id, id).
			Where(dao.Role.Columns().DeletedAt, nil).
			Scan(&parent); err != nil {
			return 0, 0, err
		}
		return parent.Id, parent.ParentId, nil
	}, treeutil.Messages{
		Self:         "上级角色不能选择自己",
		Missing:      "上级角色不存在或已删除",
		ChildLoop:    "不能将角色移动到自己的子级下",
		Cycle:        "角色层级存在循环引用",
		InvalidChain: "上级角色链路中存在无效节点",
	})
}

func validateAdminRoleMutationAllowed(touchesAdminRole bool, actorHasAdmin bool) error {
	if !touchesAdminRole || actorHasAdmin {
		return nil
	}
	return gerror.New("只有超级管理员可以操作超级管理员角色")
}
