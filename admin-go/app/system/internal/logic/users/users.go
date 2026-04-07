package users

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
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
	normalizeUsersWriteInput(in)
	if err := s.ensureDeptExists(ctx, in.DeptID); err != nil {
		return err
	}
	roleIDs, err := s.normalizeRoleIDs(ctx, in.RoleIDs)
	if err != nil {
		return err
	}
	id := snowflake.Generate()
	hashedPassword, err := password.Hash(in.Password)
	if err != nil {
		return err
	}
	now := gtime.Now()
	return dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = tx.Model(dao.Users.Table()).Ctx(ctx).Data(g.Map{
			dao.Users.Columns().Id:        id,
			dao.Users.Columns().Username:  in.Username,
			dao.Users.Columns().Password:  hashedPassword,
			dao.Users.Columns().Nickname:  in.Nickname,
			dao.Users.Columns().Email:     in.Email,
			dao.Users.Columns().Avatar:    in.Avatar,
			dao.Users.Columns().Status:    in.Status,
			dao.Users.Columns().DeptId:    in.DeptID,
			dao.Users.Columns().CreatedAt: now,
			dao.Users.Columns().UpdatedAt: now,
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
	normalizeUsersUpdateInput(in)
	// 内置管理员不可禁用
	if in.Status == 0 {
		isAdmin, err := s.isBuiltinAdmin(ctx, in.ID)
		if err != nil {
			return err
		}
		if isAdmin {
			return gerror.New("内置管理员账号不能禁用")
		}
	}
	if err := s.ensureDeptExists(ctx, in.DeptID); err != nil {
		return err
	}
	roleIDs, err := s.normalizeRoleIDs(ctx, in.RoleIDs)
	if err != nil {
		return err
	}
	data := g.Map{
		dao.Users.Columns().Username:  in.Username,
		dao.Users.Columns().Nickname:  in.Nickname,
		dao.Users.Columns().Email:     in.Email,
		dao.Users.Columns().Avatar:    in.Avatar,
		dao.Users.Columns().Status:    in.Status,
		dao.Users.Columns().DeptId:    in.DeptID,
		dao.Users.Columns().UpdatedAt: gtime.Now(),
	}
	if in.Password != "" {
		hashedPassword, err := password.Hash(in.Password)
		if err != nil {
			return err
		}
		data[dao.Users.Columns().Password] = hashedPassword
	}
	err = dao.Users.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).Where(dao.Users.Columns().Id, in.ID).Data(data).Update(); err != nil {
			return err
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

// Delete 软删除用户表
func (s *sUsers) Delete(ctx context.Context, id snowflake.JsonInt64) error {
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
			Data(g.Map{
				dao.Users.Columns().DeletedAt: gtime.Now(),
			}).
			Update(); err != nil {
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

// Detail 获取用户表详情
func (s *sUsers) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.UsersDetailOutput, err error) {
	out = &model.UsersDetailOutput{}
	err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, id).Where(dao.Users.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	// 查询部门名称
	if out.DeptID != 0 {
		val, _ := g.DB().Ctx(ctx).Model("system_dept").Where("id", out.DeptID).Where("deleted_at", nil).Value("title")
		out.DeptTitle = val.String()
	}
	// 查询用户角色ID列表
	var roles []struct {
		RoleId int64 `json:"roleId"`
	}
	_ = dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, id).Scan(&roles)
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
	s.fillRoleTitles(ctx, list)
	return
}

// ResetPassword 重置用户密码
func (s *sUsers) ResetPassword(ctx context.Context, in *model.UsersResetPasswordInput) error {
	hashedPassword, err := password.Hash(in.Password)
	if err != nil {
		return err
	}
	_, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, in.ID).Data(g.Map{
		dao.Users.Columns().Password:  hashedPassword,
		dao.Users.Columns().UpdatedAt: gtime.Now(),
	}).Update()
	if err == nil {
		authlogic.ClearUserCaches(ctx, int64(in.ID))
	}
	return err
}

func (s *sUsers) fillDeptTitles(ctx context.Context, list []*model.UsersListOutput) {
	deptSet := make(map[int64]struct{})
	for _, item := range list {
		if item.DeptID != 0 {
			deptSet[int64(item.DeptID)] = struct{}{}
		}
	}
	if len(deptSet) == 0 {
		return
	}
	deptIDs := make([]int64, 0, len(deptSet))
	for id := range deptSet {
		deptIDs = append(deptIDs, id)
	}
	rows, err := g.DB().Ctx(ctx).Model("system_dept").
		Fields("id", "title").
		Where("deleted_at", nil).
		WhereIn("id", deptIDs).
		All()
	if err != nil {
		return
	}
	deptMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		deptMap[row["id"].Int64()] = row["title"].String()
	}
	for _, item := range list {
		item.DeptTitle = deptMap[int64(item.DeptID)]
	}
}

func (s *sUsers) fillRoleTitles(ctx context.Context, list []*model.UsersListOutput) {
	userIDs := make([]int64, 0, len(list))
	for _, item := range list {
		userIDs = append(userIDs, int64(item.ID))
		item.RoleTitles = make([]string, 0)
	}
	if len(userIDs) == 0 {
		return
	}
	var userRoles []struct {
		UserId int64 `json:"userId"`
		RoleId int64 `json:"roleId"`
	}
	if err := dao.UserRole.Ctx(ctx).
		WhereIn(dao.UserRole.Columns().UserId, userIDs).
		Scan(&userRoles); err != nil || len(userRoles) == 0 {
		return
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
	rows, err := g.DB().Ctx(ctx).Model("system_role").
		Fields("id", "title").
		Where("deleted_at", nil).
		WhereIn("id", roleIDs).
		All()
	if err != nil {
		return
	}
	roleMap := make(map[int64]string, len(rows))
	for _, row := range rows {
		roleMap[row["id"].Int64()] = row["title"].String()
	}
	for _, item := range list {
		for _, roleID := range userRoleMap[int64(item.ID)] {
			appendUniqueRoleTitle(item, roleMap[roleID])
		}
	}
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

func (s *sUsers) normalizeRoleIDs(ctx context.Context, roleIDs []snowflake.JsonInt64) ([]snowflake.JsonInt64, error) {
	normalized := compactRoleIDs(roleIDs)
	if len(normalized) == 0 {
		return nil, nil
	}
	dbRoleIDs := make([]int64, 0, len(normalized))
	for _, roleID := range normalized {
		dbRoleIDs = append(dbRoleIDs, int64(roleID))
	}
	var rows []struct {
		Id int64 `json:"id"`
	}
	if err := dao.Role.Ctx(ctx).
		Fields(dao.Role.Columns().Id).
		WhereIn(dao.Role.Columns().Id, dbRoleIDs).
		Where(dao.Role.Columns().DeletedAt, nil).
		Scan(&rows); err != nil {
		return nil, err
	}
	if len(rows) != len(normalized) {
		return nil, gerror.New("包含不存在或已删除的角色")
	}
	return normalized, nil
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
