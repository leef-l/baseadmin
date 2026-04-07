package role

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/snowflake"
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
	if err := s.ensureParentValid(ctx, in.ParentID, in.ID); err != nil {
		return err
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
	err := dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			Where(dao.Role.Columns().Id, in.ID).
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

// Detail 获取角色表详情
func (s *sRole) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.RoleDetailOutput, err error) {
	out = &model.RoleDetailOutput{}
	err = dao.Role.Ctx(ctx).Where(dao.Role.Columns().Id, id).Where(dao.Role.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	// 查询上级角色ID，0 表示顶级角色关联显示
	if out.ParentID != 0 {
		val, _ := g.DB().Ctx(ctx).Model("system_role").Where("id", out.ParentID).Where("deleted_at", nil).Value("title")
		out.RoleTitle = val.String()
	}
	return
}

// List 获取角色表列表
func (s *sRole) List(ctx context.Context, in *model.RoleListInput) (list []*model.RoleListOutput, total int, err error) {
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

	// 使用 map 迭代方式组装树
	nodeMap := make(map[int64]*model.RoleTreeOutput, len(list))
	for _, item := range list {
		item.Children = make([]*model.RoleTreeOutput, 0)
		nodeMap[int64(item.ID)] = item
	}

	tree = make([]*model.RoleTreeOutput, 0)
	for _, item := range list {
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
	return
}

// GrantMenu 角色授权菜单（先删后插）
func (s *sRole) GrantMenu(ctx context.Context, in *model.RoleGrantMenuInput) error {
	err := dao.RoleMenu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Where(dao.RoleMenu.Columns().RoleId, in.ID).Delete(); err != nil {
			return err
		}
		if len(in.MenuIDs) == 0 {
			return nil
		}
		data := make([]g.Map, 0, len(in.MenuIDs))
		for _, menuID := range in.MenuIDs {
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
	err := dao.Role.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).Where(dao.Role.Columns().Id, in.ID).Data(g.Map{
			dao.Role.Columns().DataScope: in.DataScope,
		}).Update(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.RoleDept.Table()).Ctx(ctx).Where(dao.RoleDept.Columns().RoleId, in.ID).Delete(); err != nil {
			return err
		}
		if in.DataScope != 5 || len(in.DeptIDs) == 0 {
			return nil
		}
		data := make([]g.Map, 0, len(in.DeptIDs))
		for _, deptID := range in.DeptIDs {
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
	parentSet := make(map[int64]struct{})
	for _, item := range list {
		if item.ParentID != 0 {
			parentSet[int64(item.ParentID)] = struct{}{}
		}
	}
	if len(parentSet) == 0 {
		return
	}
	parentIDs := make([]int64, 0, len(parentSet))
	for id := range parentSet {
		parentIDs = append(parentIDs, id)
	}
	rows, err := g.DB().Ctx(ctx).Model("system_role").
		Fields("id", "title").
		Where("deleted_at", nil).
		WhereIn("id", parentIDs).
		All()
	if err != nil {
		return
	}
	parentMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		parentMap[row["id"].Int64()] = row["title"].String()
	}
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

func (s *sRole) ensureParentValid(ctx context.Context, parentID, currentID snowflake.JsonInt64) error {
	if parentID == 0 {
		return nil
	}
	if currentID != 0 && parentID == currentID {
		return gerror.New("上级角色不能选择自己")
	}
	var parent struct {
		Id       int64 `json:"id"`
		ParentId int64 `json:"parentId"`
	}
	if err := dao.Role.Ctx(ctx).
		Fields(dao.Role.Columns().Id, dao.Role.Columns().ParentId).
		Where(dao.Role.Columns().Id, parentID).
		Where(dao.Role.Columns().DeletedAt, nil).
		Scan(&parent); err != nil {
		return err
	}
	if parent.Id == 0 {
		return gerror.New("上级角色不存在或已删除")
	}
	seen := map[int64]struct{}{int64(parentID): {}}
	for parent.ParentId != 0 {
		if currentID != 0 && parent.ParentId == int64(currentID) {
			return gerror.New("不能将角色移动到自己的子级下")
		}
		if _, ok := seen[parent.ParentId]; ok {
			return gerror.New("角色层级存在循环引用")
		}
		seen[parent.ParentId] = struct{}{}
		next := struct {
			Id       int64 `json:"id"`
			ParentId int64 `json:"parentId"`
		}{}
		if err := dao.Role.Ctx(ctx).
			Fields(dao.Role.Columns().Id, dao.Role.Columns().ParentId).
			Where(dao.Role.Columns().Id, parent.ParentId).
			Where(dao.Role.Columns().DeletedAt, nil).
			Scan(&next); err != nil {
			return err
		}
		if next.Id == 0 {
			return gerror.New("上级角色链路中存在无效节点")
		}
		parent = next
	}
	return nil
}
