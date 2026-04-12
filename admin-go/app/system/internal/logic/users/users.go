package users

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterUsers(New())
}

func New() *sUsers {
	return &sUsers{}
}

type sUsers struct{}

// Create 创建用户表
func (s *sUsers) Create(ctx context.Context, in *model.UsersCreateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeUsersWriteInput(in)
	if err := s.ensureUsersUniqueFields(ctx, 0, in.Username, in.Email); err != nil {
		return err
	}
	if err := s.ensureDeptExists(ctx, in.DeptID); err != nil {
		return err
	}
	if err := s.ensureDeptAccessible(ctx, in.DeptID); err != nil {
		return err
	}
	roleIDs, err := s.normalizeRoleIDs(ctx, in.RoleIDs)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleGrantAllowed(ctx, roleIDs); err != nil {
		return err
	}
	if err := s.ensureRoleIDsAssignable(ctx, roleIDs); err != nil {
		return err
	}
	id := snowflake.Generate()
	if err := password.ValidatePolicy(in.Password); err != nil {
		return gerror.New(err.Error())
	}
	hashedPassword, err := password.Hash(in.Password)
	if err != nil {
		return err
	}
	return dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = tx.Model(dao.Users.Table()).Ctx(ctx).Data(do.Users{
			Id:       id,
			Username: in.Username,
			Password: hashedPassword,
			Nickname: in.Nickname,
			Email:    in.Email,
			Avatar:   in.Avatar,
			Status:   in.Status,
			DeptId:   in.DeptID,
		}).Insert()
		if err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		data := make([]g.Map, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			data = append(data, g.Map{
				dao.UserRole.Columns().UserId: id,
				dao.UserRole.Columns().RoleId: roleID,
			})
		}
		_, err = tx.Model(dao.UserRole.Table()).Ctx(ctx).Data(data).Insert()
		return err
	})
}

// Update 更新用户表
func (s *sUsers) Update(ctx context.Context, in *model.UsersUpdateInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeUsersUpdateInput(in)
	if in.Username == "" {
		return gerror.New("登录用户名不能为空")
	}
	if err := s.ensureUserExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureUserAccessible(ctx, in.ID); err != nil {
		return err
	}
	// 内置管理员不可禁用
	isBuiltinAdmin, err := s.isBuiltinAdmin(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureBuiltinAdminManageAllowed(ctx, isBuiltinAdmin); err != nil {
		return err
	}
	hasAdminRole, err := s.userHasAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleUserManageAllowed(ctx, hasAdminRole); err != nil {
		return err
	}
	if err := s.ensureUserRolesManageable(ctx, in.ID); err != nil {
		return err
	}
	if in.Status == 0 && isBuiltinAdmin {
		return gerror.New("内置管理员账号不能禁用")
	}
	if isBuiltinAdmin && in.Username != "admin" {
		return gerror.New("内置管理员账号登录名不能修改")
	}
	if err := s.ensureUsersUniqueFields(ctx, in.ID, in.Username, in.Email); err != nil {
		return err
	}
	if err := s.ensureDeptExists(ctx, in.DeptID); err != nil {
		return err
	}
	if err := s.ensureDeptAccessible(ctx, in.DeptID); err != nil {
		return err
	}
	var (
		roleIDs        []snowflake.JsonInt64
		shouldSyncRole = in.RoleIDs != nil
	)
	if shouldSyncRole {
		roleIDs, err = s.normalizeRoleIDs(ctx, in.RoleIDs)
		if err != nil {
			return err
		}
		if err := s.ensureAdminRoleGrantAllowed(ctx, roleIDs); err != nil {
			return err
		}
		if err := s.ensureRoleIDsAssignable(ctx, roleIDs); err != nil {
			return err
		}
		if isBuiltinAdmin {
			if err := s.ensureBuiltinAdminRoleAssignment(ctx, roleIDs); err != nil {
				return err
			}
		}
	}
	data := do.Users{
		Username: in.Username,
		Nickname: in.Nickname,
		Email:    in.Email,
		Avatar:   in.Avatar,
		Status:   in.Status,
		DeptId:   in.DeptID,
	}
	if in.Password != "" {
		if err := password.ValidatePolicy(in.Password); err != nil {
			return gerror.New(err.Error())
		}
		hashedPassword, err := password.Hash(in.Password)
		if err != nil {
			return err
		}
		data.Password = hashedPassword
	}
	err = dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).
			Where(dao.Users.Columns().Id, in.ID).
			Where(dao.Users.Columns().DeletedAt, nil).
			Data(data).
			Update(); err != nil {
			return err
		}
		if !shouldSyncRole {
			return nil
		}
		if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).Where(dao.UserRole.Columns().UserId, in.ID).Delete(); err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		roleData := make([]g.Map, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			roleData = append(roleData, g.Map{
				dao.UserRole.Columns().UserId: in.ID,
				dao.UserRole.Columns().RoleId: roleID,
			})
		}
		_, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).Data(roleData).Insert()
		return err
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, int64(in.ID))
	}
	return err
}

// isBuiltinAdmin 检查用户是否为内置管理员
func (s *sUsers) isBuiltinAdmin(ctx context.Context, id snowflake.JsonInt64) (bool, error) {
	val, err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, id).Where(dao.Users.Columns().DeletedAt, nil).Value(dao.Users.Columns().Username)
	if err != nil {
		return false, err
	}
	return val.String() == "admin", nil
}

func (s *sUsers) ensureUsersUniqueFields(ctx context.Context, currentID snowflake.JsonInt64, username, email string) error {
	username = strings.TrimSpace(username)
	if username != "" {
		m := dao.Users.Ctx(ctx).
			Where(dao.Users.Columns().Username, username).
			Where(dao.Users.Columns().DeletedAt, nil)
		if currentID > 0 {
			m = m.WhereNot(dao.Users.Columns().Id, currentID)
		}
		count, err := m.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.New("登录用户名已存在")
		}
	}
	email = strings.TrimSpace(email)
	if email != "" {
		m := dao.Users.Ctx(ctx).
			Where(dao.Users.Columns().Email, email).
			Where(dao.Users.Columns().DeletedAt, nil)
		if currentID > 0 {
			m = m.WhereNot(dao.Users.Columns().Id, currentID)
		}
		count, err := m.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.New("邮箱已存在")
		}
	}
	return nil
}

// Delete 软删除用户表
func (s *sUsers) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	if err := s.ensureUserExists(ctx, id); err != nil {
		return err
	}
	if err := s.ensureUserAccessible(ctx, id); err != nil {
		return err
	}
	hasAdminRole, err := s.userHasAdminRole(ctx, id)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleUserManageAllowed(ctx, hasAdminRole); err != nil {
		return err
	}
	if err := s.ensureUserRolesManageable(ctx, id); err != nil {
		return err
	}
	// 内置管理员不可删除
	isAdmin, err := s.isBuiltinAdmin(ctx, id)
	if err != nil {
		return err
	}
	if isAdmin {
		return gerror.New("内置管理员账号不能删除")
	}
	err = dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).
			Where(dao.Users.Columns().Id, id).
			Where(dao.Users.Columns().DeletedAt, nil).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).
			Where(dao.UserRole.Columns().UserId, id).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserDept.Table()).Ctx(ctx).
			Where(dao.UserDept.Columns().UserId, id).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, int64(id))
	}
	return err
}

// BatchDelete 批量删除用户表
func (s *sUsers) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的用户")
	}
	deleteIDs := batchutil.ToInt64s(ids)
	var users []struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
		DeptId   int64  `json:"deptId"`
	}
	if err := dao.Users.Ctx(ctx).
		Fields(dao.Users.Columns().Id, dao.Users.Columns().Username, dao.Users.Columns().DeptId).
		WhereIn(dao.Users.Columns().Id, deleteIDs).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&users); err != nil {
		return err
	}
	if len(users) != len(deleteIDs) {
		return gerror.New("包含不存在或已删除的用户")
	}
	for _, item := range users {
		if item.Username == "admin" {
			return gerror.New("内置管理员账号不能删除")
		}
		allowed, err := shared.CanAccessUser(ctx, item.Id, item.DeptId)
		if err != nil {
			return err
		}
		if !allowed {
			return gerror.New("包含无权操作的用户")
		}
		hasAdminRole, err := s.userHasAdminRole(ctx, snowflake.JsonInt64(item.Id))
		if err != nil {
			return err
		}
		if err := s.ensureAdminRoleUserManageAllowed(ctx, hasAdminRole); err != nil {
			return err
		}
		if err := s.ensureUserRolesManageable(ctx, snowflake.JsonInt64(item.Id)); err != nil {
			return err
		}
	}
	err := dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).
			WhereIn(dao.Users.Columns().Id, deleteIDs).
			Where(dao.Users.Columns().DeletedAt, nil).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).
			WhereIn(dao.UserRole.Columns().UserId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		if _, err := tx.Model(dao.UserDept.Table()).Ctx(ctx).
			WhereIn(dao.UserDept.Columns().UserId, deleteIDs).
			Delete(); err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		authlogic.ClearUserCaches(ctx, deleteIDs...)
	}
	return err
}

// Detail 获取用户表详情
func (s *sUsers) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.UsersDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("用户不存在或已删除")
	}
	out = &model.UsersDetailOutput{}
	err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, id).Where(dao.Users.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("用户不存在或已删除")
	}
	allowed, err := shared.CanAccessUser(ctx, int64(out.ID), int64(out.DeptID))
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, gerror.New("用户不存在或已删除")
	}
	out.DeptTitle = shared.LookupTitle(ctx, "system_dept", int64(out.DeptID))
	// 查询用户角色ID列表
	var roles []struct {
		RoleId int64 `json:"roleId"`
	}
	if err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, id).Scan(&roles); err != nil {
		return nil, err
	}
	out.RoleIDs = make([]snowflake.JsonInt64, 0, len(roles))
	for _, r := range roles {
		out.RoleIDs = append(out.RoleIDs, snowflake.JsonInt64(r.RoleId))
	}
	return
}

// List 获取用户表列表
func (s *sUsers) List(ctx context.Context, in *model.UsersListInput) (list []*model.UsersListOutput, total int, err error) {
	if in == nil {
		in = &model.UsersListInput{}
	}
	normalizeUsersListInput(in)
	m := dao.Users.Ctx(ctx).Where(dao.Users.Columns().DeletedAt, nil)
	m, err = shared.ApplyDeptScopeToUserModel(ctx, m, dao.Users.Columns().Id, dao.Users.Columns().DeptId)
	if err != nil {
		return nil, 0, err
	}
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.Users.Columns().Username, "%"+in.Keyword+"%").
			WhereOrLike(dao.Users.Columns().Nickname, "%"+in.Keyword+"%").
			WhereOrLike(dao.Users.Columns().Email, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
	if in.Status != nil {
		m = m.Where(dao.Users.Columns().Status, *in.Status)
	}
	if in.Username != "" {
		m = m.WhereLike(dao.Users.Columns().Username, "%"+in.Username+"%")
	}
	if in.Nickname != "" {
		m = m.WhereLike(dao.Users.Columns().Nickname, "%"+in.Nickname+"%")
	}
	if in.Email != "" {
		m = m.WhereLike(dao.Users.Columns().Email, "%"+in.Email+"%")
	}
	if in.DeptId > 0 {
		m = m.Where(dao.Users.Columns().DeptId, in.DeptId)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.Users.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	s.fillDeptTitles(ctx, list)
	if err = s.fillRoleTitles(ctx, list); err != nil {
		return nil, 0, err
	}
	return
}

// ResetPassword 重置用户密码
func (s *sUsers) ResetPassword(ctx context.Context, in *model.UsersResetPasswordInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeUsersResetPasswordInput(in)
	if in.Password == "" {
		return gerror.New("新密码不能为空")
	}
	if err := password.ValidatePolicy(in.Password); err != nil {
		return gerror.New(err.Error())
	}
	if err := s.ensureUserExists(ctx, in.ID); err != nil {
		return err
	}
	if err := s.ensureUserAccessible(ctx, in.ID); err != nil {
		return err
	}
	isBuiltinAdmin, err := s.isBuiltinAdmin(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureBuiltinAdminManageAllowed(ctx, isBuiltinAdmin); err != nil {
		return err
	}
	hasAdminRole, err := s.userHasAdminRole(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureAdminRoleUserManageAllowed(ctx, hasAdminRole); err != nil {
		return err
	}
	if err := s.ensureUserRolesManageable(ctx, in.ID); err != nil {
		return err
	}
	hashedPassword, err := password.Hash(in.Password)
	if err != nil {
		return err
	}
	_, err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, in.ID).
		Where(dao.Users.Columns().DeletedAt, nil).
		Data(do.Users{
			Password: hashedPassword,
		}).
		Update()
	if err == nil {
		authlogic.ClearUserCaches(ctx, int64(in.ID))
	}
	return err
}

func (s *sUsers) fillDeptTitles(ctx context.Context, list []*model.UsersListOutput) {
	deptIDs := make([]int64, 0, len(list))
	for _, item := range list {
		if item.DeptID != 0 {
			deptIDs = append(deptIDs, int64(item.DeptID))
		}
	}
	deptMap := shared.LoadTitleMap(ctx, "system_dept", deptIDs)
	for _, item := range list {
		item.DeptTitle = deptMap[int64(item.DeptID)]
	}
}

func (s *sUsers) fillRoleTitles(ctx context.Context, list []*model.UsersListOutput) error {
	userIDs := make([]int64, 0, len(list))
	for _, item := range list {
		userIDs = append(userIDs, int64(item.ID))
		item.RoleTitles = make([]string, 0)
	}
	if len(userIDs) == 0 {
		return nil
	}
	var userRoles []struct {
		UserId int64 `json:"userId"`
		RoleId int64 `json:"roleId"`
	}
	if err := dao.UserRole.Ctx(ctx).
		WhereIn(dao.UserRole.Columns().UserId, userIDs).
		Scan(&userRoles); err != nil {
		return err
	}
	if len(userRoles) == 0 {
		return nil
	}
	roleSet := make(map[int64]struct{})
	userRoleMap := make(map[int64][]int64)
	for _, item := range userRoles {
		roleSet[item.RoleId] = struct{}{}
		userRoleMap[item.UserId] = append(userRoleMap[item.UserId], item.RoleId)
	}
	roleIDs := make([]int64, 0, len(roleSet))
	for id := range roleSet {
		roleIDs = append(roleIDs, id)
	}
	roleMap := shared.LoadTitleMap(ctx, "system_role", roleIDs)
	for _, item := range list {
		for _, roleID := range userRoleMap[int64(item.ID)] {
			appendUniqueRoleTitle(item, roleMap[roleID])
		}
	}
	return nil
}

func appendUniqueRoleTitle(item *model.UsersListOutput, title string) {
	if item == nil || title == "" {
		return
	}
	for _, existing := range item.RoleTitles {
		if existing == title {
			return
		}
	}
	item.RoleTitles = append(item.RoleTitles, title)
}

func (s *sUsers) ensureUserExists(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("用户不存在或已删除")
	}
	count, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, id).
		Where(dao.Users.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("用户不存在或已删除")
	}
	return nil
}

func (s *sUsers) ensureUserAccessible(ctx context.Context, id snowflake.JsonInt64) error {
	if id <= 0 {
		return gerror.New("用户不存在或已删除")
	}
	var row struct {
		Id     int64 `json:"id"`
		DeptId int64 `json:"deptId"`
	}
	if err := dao.Users.Ctx(ctx).
		Fields(dao.Users.Columns().Id, dao.Users.Columns().DeptId).
		Where(dao.Users.Columns().Id, id).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&row); err != nil {
		return err
	}
	if row.Id == 0 {
		return gerror.New("用户不存在或已删除")
	}
	allowed, err := shared.CanAccessUser(ctx, row.Id, row.DeptId)
	if err != nil {
		return err
	}
	if !allowed {
		return gerror.New("用户不存在或已删除")
	}
	return nil
}

func (s *sUsers) ensureDeptExists(ctx context.Context, deptID snowflake.JsonInt64) error {
	if deptID == 0 {
		return nil
	}
	count, err := dao.Dept.Ctx(ctx).
		Where(dao.Dept.Columns().Id, deptID).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("所选部门不存在或已删除")
	}
	return nil
}

func (s *sUsers) ensureDeptAccessible(ctx context.Context, deptID snowflake.JsonInt64) error {
	if deptID == 0 {
		scope, err := shared.ResolveDataAccessScope(ctx)
		if err != nil {
			return err
		}
		if scope.All {
			return nil
		}
		return gerror.New("部门不存在或无权操作")
	}
	allowed, err := shared.CanAccessDept(ctx, int64(deptID))
	if err != nil {
		return err
	}
	if !allowed {
		return gerror.New("部门不存在或无权操作")
	}
	return nil
}

func (s *sUsers) ensureAdminRoleGrantAllowed(ctx context.Context, roleIDs []snowflake.JsonInt64) error {
	assignsAdminRole, err := s.containsAdminRole(ctx, roleIDs)
	if err != nil {
		return err
	}
	if !assignsAdminRole {
		return nil
	}
	actorHasAdmin, err := shared.HasCurrentActorAdminRole(ctx)
	if err != nil {
		return err
	}
	return validateAdminRoleGrantAllowed(assignsAdminRole, actorHasAdmin)
}

func (s *sUsers) ensureRoleIDsAssignable(ctx context.Context, roleIDs []snowflake.JsonInt64) error {
	if len(roleIDs) == 0 {
		return nil
	}
	assignableRoleIDs, err := shared.LoadCurrentActorAssignableRoleIDs(ctx)
	if err != nil {
		return err
	}
	if shared.RoleIDsWithinScope(batchutil.ToInt64s(roleIDs), assignableRoleIDs) {
		return nil
	}
	return gerror.New("包含无权分配的角色")
}

func (s *sUsers) ensureBuiltinAdminManageAllowed(ctx context.Context, isBuiltinAdmin bool) error {
	if !isBuiltinAdmin {
		return nil
	}
	actorHasAdmin, err := shared.HasCurrentActorAdminRole(ctx)
	if err != nil {
		return err
	}
	return validateBuiltinAdminManageAllowed(isBuiltinAdmin, actorHasAdmin)
}

func (s *sUsers) ensureAdminRoleUserManageAllowed(ctx context.Context, hasAdminRole bool) error {
	if !hasAdminRole {
		return nil
	}
	actorHasAdmin, err := shared.HasCurrentActorAdminRole(ctx)
	if err != nil {
		return err
	}
	return validateAdminRoleUserManageAllowed(hasAdminRole, actorHasAdmin)
}

func (s *sUsers) containsAdminRole(ctx context.Context, roleIDs []snowflake.JsonInt64) (bool, error) {
	if len(roleIDs) == 0 {
		return false, nil
	}
	count, err := dao.Role.Ctx(ctx).
		WhereIn(dao.Role.Columns().Id, batchutil.ToInt64s(roleIDs)).
		Where(dao.Role.Columns().DeletedAt, nil).
		Where(dao.Role.Columns().IsAdmin, 1).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sUsers) userHasAdminRole(ctx context.Context, userID snowflake.JsonInt64) (bool, error) {
	if userID <= 0 {
		return false, nil
	}
	count, err := dao.UserRole.Ctx(ctx).
		LeftJoin(dao.Role.Table(), dao.Role.Table()+"."+dao.Role.Columns().Id+"="+dao.UserRole.Table()+"."+dao.UserRole.Columns().RoleId).
		Where(dao.UserRole.Table()+"."+dao.UserRole.Columns().UserId, userID).
		Where(dao.Role.Table()+"."+dao.Role.Columns().DeletedAt, nil).
		Where(dao.Role.Table()+"."+dao.Role.Columns().IsAdmin, 1).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *sUsers) ensureUserRolesManageable(ctx context.Context, userID snowflake.JsonInt64) error {
	roleIDs, err := s.loadUserActiveRoleIDs(ctx, userID)
	if err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	assignableRoleIDs, err := shared.LoadCurrentActorAssignableRoleIDs(ctx)
	if err != nil {
		return err
	}
	if shared.RoleIDsWithinScope(roleIDs, assignableRoleIDs) {
		return nil
	}
	return gerror.New("用户存在不可管理的角色")
}

func (s *sUsers) loadUserActiveRoleIDs(ctx context.Context, userID snowflake.JsonInt64) ([]int64, error) {
	if userID <= 0 {
		return nil, nil
	}
	var rows []struct {
		RoleId int64 `json:"roleId"`
	}
	if err := dao.UserRole.Ctx(ctx).
		LeftJoin(dao.Role.Table(), dao.Role.Table()+"."+dao.Role.Columns().Id+"="+dao.UserRole.Table()+"."+dao.UserRole.Columns().RoleId).
		Fields(dao.UserRole.Table()+"."+dao.UserRole.Columns().RoleId+" AS roleId").
		Where(dao.UserRole.Table()+"."+dao.UserRole.Columns().UserId, userID).
		Where(dao.Role.Table()+"."+dao.Role.Columns().DeletedAt, nil).
		Where(dao.Role.Table()+"."+dao.Role.Columns().Status, 1).
		Scan(&rows); err != nil {
		return nil, err
	}
	roleIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		roleIDs = append(roleIDs, row.RoleId)
	}
	return roleIDs, nil
}

func (s *sUsers) normalizeRoleIDs(ctx context.Context, roleIDs []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := compactRoleIDs(roleIDs)
	if len(normalized) == 0 {
		return nil, nil
	}
	dbRoleIDs := make([]int64, 0, len(normalized))
	for _, roleID := range normalized {
		dbRoleIDs = append(dbRoleIDs, int64(roleID))
	}
	if !shared.ContainsAllIDs(ctx, dao.Role.Table(), dbRoleIDs) {
		return nil, gerror.New("包含不存在或已删除的角色")
	}
	return normalized, nil
}

func (s *sUsers) ensureBuiltinAdminRoleAssignment(ctx context.Context, roleIDs []snowflake.JsonInt64) error {
	if len(roleIDs) == 0 {
		return gerror.New("内置管理员账号必须保留至少一个超级管理员角色")
	}
	dbRoleIDs := make([]int64, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		dbRoleIDs = append(dbRoleIDs, int64(roleID))
	}
	count, err := dao.Role.Ctx(ctx).
		WhereIn(dao.Role.Columns().Id, dbRoleIDs).
		Where(dao.Role.Columns().DeletedAt, nil).
		Where(dao.Role.Columns().IsAdmin, 1).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("内置管理员账号必须保留至少一个超级管理员角色")
	}
	return nil
}

func compactRoleIDs(roleIDs []snowflake.JsonInt64) []snowflake.JsonInt64 {
	if len(roleIDs) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(roleIDs))
	normalized := make([]snowflake.JsonInt64, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		id := int64(roleID)
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, roleID)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func normalizeUsersWriteInput(in *model.UsersCreateInput) {
	if in == nil {
		return
	}
	in.Username = strings.TrimSpace(in.Username)
	in.Nickname = strings.TrimSpace(in.Nickname)
	in.Email = strings.TrimSpace(in.Email)
	in.Avatar = strings.TrimSpace(in.Avatar)
}

func normalizeUsersUpdateInput(in *model.UsersUpdateInput) {
	if in == nil {
		return
	}
	in.Username = strings.TrimSpace(in.Username)
	in.Password = strings.TrimSpace(in.Password)
	in.Nickname = strings.TrimSpace(in.Nickname)
	in.Email = strings.TrimSpace(in.Email)
	in.Avatar = strings.TrimSpace(in.Avatar)
}

func normalizeUsersListInput(in *model.UsersListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Username = strings.TrimSpace(in.Username)
	in.Nickname = strings.TrimSpace(in.Nickname)
	in.Email = strings.TrimSpace(in.Email)
}

func normalizeUsersResetPasswordInput(in *model.UsersResetPasswordInput) {
	if in == nil {
		return
	}
	in.Password = strings.TrimSpace(in.Password)
}

func validateAdminRoleGrantAllowed(assignsAdminRole bool, actorHasAdmin bool) error {
	if !assignsAdminRole || actorHasAdmin {
		return nil
	}
	return gerror.New("只有超级管理员可以分配超级管理员角色")
}

func validateBuiltinAdminManageAllowed(isBuiltinAdmin bool, actorHasAdmin bool) error {
	if !isBuiltinAdmin || actorHasAdmin {
		return nil
	}
	return gerror.New("内置管理员账号只能由超级管理员操作")
}

func validateAdminRoleUserManageAllowed(hasAdminRole bool, actorHasAdmin bool) error {
	if !hasAdminRole || actorHasAdmin {
		return nil
	}
	return gerror.New("超级管理员账号只能由超级管理员操作")
}
