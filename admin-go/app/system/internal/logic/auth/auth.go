package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/cache"
	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterAuth(New())
}

func New() *sAuth {
	return &sAuth{}
}

type sAuth struct{}

type permissionRow struct {
	Permission string `json:"permission"`
}

const (
	authLoginFailLimit  = 5
	authLoginFailWindow = 10 * time.Minute
	authInfoCacheTTL    = time.Minute
	authMenusCacheTTL   = time.Minute
)

func normalizeAuthLoginInput(in *model.AuthLoginInput) {
	if in == nil {
		return
	}
	in.Username = strings.TrimSpace(in.Username)
}

// Login 用户登录
func (s *sAuth) Login(ctx context.Context, in *model.AuthLoginInput) (out *model.AuthLoginOutput, err error) {
	normalizeAuthLoginInput(in)
	if s.isLoginRateLimited(ctx, in.Username) {
		return nil, gerror.New("登录失败次数过多，请10分钟后再试")
	}

	// 查询用户
	var user struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		DeptId   int64  `json:"deptId"`
		Status   int    `json:"status"`
	}

	err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Username, in.Username).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&user)

	if err != nil {
		return nil, gerror.New("用户名或密码错误")
	}
	if user.Id == 0 {
		s.recordLoginFailure(ctx, in.Username)
		return nil, gerror.New("用户名或密码错误")
	}

	// 校验状态
	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

	// 校验密码
	if !password.Verify(user.Password, in.Password) {
		s.recordLoginFailure(ctx, in.Username)
		return nil, gerror.New("用户名或密码错误")
	}
	s.clearLoginFailures(ctx, in.Username)
	if password.NeedsRehash(user.Password) {
		if upgraded, hashErr := password.Hash(in.Password); hashErr == nil {
			_, _ = dao.Users.Ctx(ctx).
				Where(dao.Users.Columns().Id, user.Id).
				Data(g.Map{
					dao.Users.Columns().Password:  upgraded,
					dao.Users.Columns().UpdatedAt: gtime.Now(),
				}).
				Update()
		}
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.Id, user.Username, user.DeptId)
	if err != nil {
		return nil, gerror.New("生成Token失败")
	}

	out = &model.AuthLoginOutput{
		Token:    token,
		UserID:   snowflake.JsonInt64(user.Id),
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
	return
}

// Info 获取当前用户信息
func (s *sAuth) Info(ctx context.Context, userID snowflake.JsonInt64) (out *model.AuthInfoOutput, err error) {
	if cached, ok := s.getCachedInfo(ctx, userID); ok {
		return cached, nil
	}

	out, err = s.loadInfo(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.setCachedInfo(ctx, userID, out)
	return out, nil
}

func (s *sAuth) loadInfo(ctx context.Context, userID snowflake.JsonInt64) (out *model.AuthInfoOutput, err error) {
	var user struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		DeptId   int64  `json:"deptId"`
		Status   int    `json:"status"`
	}
	err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userID).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&user)
	if err != nil {
		return nil, err
	}

	out = &model.AuthInfoOutput{
		UserID:   snowflake.JsonInt64(user.Id),
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		DeptID:   snowflake.JsonInt64(user.DeptId),
		Status:   user.Status,
		Roles:    make([]string, 0),
		Perms:    make([]string, 0),
	}

	// 查询用户角色
	var userRoles []struct {
		RoleId int64 `json:"roleId"`
	}
	_ = dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userID).Scan(&userRoles)

	if len(userRoles) > 0 {
		roleIDs := make([]int64, 0, len(userRoles))
		for _, ur := range userRoles {
			roleIDs = append(roleIDs, ur.RoleId)
		}
		roleIDs = compactInt64s(roleIDs)

		if len(roleIDs) > 0 {
			// 查询角色名称
			var roles []struct {
				Title string `json:"title"`
			}
			_ = g.DB().Ctx(ctx).Model("system_role").
				Where("id", roleIDs).
				Where("deleted_at", nil).
				Where("status", 1).
				Scan(&roles)
			for _, r := range roles {
				out.Roles = append(out.Roles, r.Title)
			}

			// 检查是否有超级管理员角色
			isAdmin := false
			adminCount, _ := g.DB().Ctx(ctx).Model("system_role").
				Where("id", roleIDs).
				Where("deleted_at", nil).
				Where("status", 1).
				Where("is_admin", 1).
				Count()
			isAdmin = adminCount > 0

			if isAdmin {
				// 超级管理员获取所有权限
				var perms []permissionRow
				_ = g.DB().Ctx(ctx).Model("system_menu").
					Where("deleted_at", nil).
					Where("status", 1).
					WhereNot("permission", "").
					Scan(&perms)
				out.Perms = compactPermissions(collectPermissions(perms))
			} else {
				// 查询角色关联的菜单权限标识
				var menuIDs []struct {
					MenuId int64 `json:"menuId"`
				}
				_ = dao.RoleMenu.Ctx(ctx).WhereIn(dao.RoleMenu.Columns().RoleId, roleIDs).Scan(&menuIDs)

				if len(menuIDs) > 0 {
					mIDs := make([]int64, 0, len(menuIDs))
					for _, m := range menuIDs {
						mIDs = append(mIDs, m.MenuId)
					}
					mIDs = compactInt64s(mIDs)
					if len(mIDs) > 0 {
						var perms []permissionRow
						_ = g.DB().Ctx(ctx).Model("system_menu").
							Where("id", mIDs).
							Where("deleted_at", nil).
							Where("status", 1).
							WhereNot("permission", "").
							Scan(&perms)
						out.Perms = compactPermissions(collectPermissions(perms))
					}
				}
			}
		}
	}

	return
}

// ChangePassword 修改密码
func (s *sAuth) ChangePassword(ctx context.Context, in *model.AuthChangePasswordInput) error {
	// 查询当前密码
	currentPassword, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, in.UserID).
		Value(dao.Users.Columns().Password)
	if err != nil {
		return err
	}

	// 校验旧密码
	if !password.Verify(currentPassword.String(), in.OldPassword) {
		return gerror.New("旧密码错误")
	}

	// 加密新密码
	hashedNew, err := password.Hash(in.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	_, err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, in.UserID).
		Data(dao.Users.Columns().Password, hashedNew).
		Update()
	if err == nil {
		s.clearAuthCache(ctx, int64(in.UserID))
	}
	return err
}

// Menus 获取当前用户的菜单树（动态路由）
func (s *sAuth) Menus(ctx context.Context, userID snowflake.JsonInt64) ([]*model.AuthMenuOutput, error) {
	if cached, ok := s.getCachedMenus(ctx, userID); ok {
		return cached, nil
	}

	menus, err := s.loadMenus(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.setCachedMenus(ctx, userID, menus)
	return menus, nil
}

func (s *sAuth) loadMenus(ctx context.Context, userID snowflake.JsonInt64) ([]*model.AuthMenuOutput, error) {
	// 查询用户角色
	var userRoles []struct {
		RoleId int64 `json:"roleId"`
	}
	err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userID).Scan(&userRoles)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}

	roleIDs := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleId)
	}
	roleIDs = compactInt64s(roleIDs)
	if len(roleIDs) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}

	// 检查是否有超级管理员角色
	isAdmin := false
	adminCount, _ := g.DB().Ctx(ctx).Model("system_role").
		Where("id", roleIDs).
		Where("deleted_at", nil).
		Where("status", 1).
		Where("is_admin", 1).
		Count()
	isAdmin = adminCount > 0

	if isAdmin {
		// 超级管理员获取所有菜单
		var list []*model.AuthMenuOutput
		err = g.DB().Ctx(ctx).Model("system_menu").
			Where("deleted_at", nil).
			Where("status", 1).
			OrderAsc("sort").
			Scan(&list)
		if err != nil {
			return nil, err
		}
		return buildMenuTree(list), nil
	}

	// 查询角色关联的菜单ID（去重）
	var roleMenus []struct {
		MenuId int64 `json:"menuId"`
	}
	err = dao.RoleMenu.Ctx(ctx).WhereIn(dao.RoleMenu.Columns().RoleId, roleIDs).Scan(&roleMenus)
	if err != nil {
		return nil, err
	}

	if len(roleMenus) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}

	menuIDs := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuId)
	}
	menuIDs = compactInt64s(menuIDs)
	if len(menuIDs) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}

	// 查询菜单详情
	var list []*model.AuthMenuOutput
	err = g.DB().Ctx(ctx).Model("system_menu").
		Where("id", menuIDs).
		Where("deleted_at", nil).
		Where("status", 1).
		OrderAsc("sort").
		Scan(&list)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(list), nil
}

func (s *sAuth) isLoginRateLimited(ctx context.Context, username string) bool {
	count, err := cache.GetInt64(ctx, s.loginFailKey(ctx, username))
	return err == nil && count >= authLoginFailLimit
}

func (s *sAuth) recordLoginFailure(ctx context.Context, username string) {
	_, _ = cache.IncrWithTTL(ctx, s.loginFailKey(ctx, username), authLoginFailWindow)
}

func (s *sAuth) clearLoginFailures(ctx context.Context, username string) {
	_ = cache.Delete(ctx, s.loginFailKey(ctx, username))
}

func (s *sAuth) getCachedInfo(ctx context.Context, userID snowflake.JsonInt64) (*model.AuthInfoOutput, bool) {
	var out model.AuthInfoOutput
	ok, err := cache.GetJSON(ctx, s.infoCacheKey(userID), &out)
	if err != nil || !ok {
		return nil, false
	}
	return &out, true
}

func (s *sAuth) setCachedInfo(ctx context.Context, userID snowflake.JsonInt64, out *model.AuthInfoOutput) {
	if out == nil {
		return
	}
	_ = cache.SetJSON(ctx, s.infoCacheKey(userID), out, authInfoCacheTTL)
}

func (s *sAuth) getCachedMenus(ctx context.Context, userID snowflake.JsonInt64) ([]*model.AuthMenuOutput, bool) {
	var menus []*model.AuthMenuOutput
	ok, err := cache.GetJSON(ctx, s.menusCacheKey(userID), &menus)
	if err != nil || !ok {
		return nil, false
	}
	return menus, true
}

func (s *sAuth) setCachedMenus(ctx context.Context, userID snowflake.JsonInt64, menus []*model.AuthMenuOutput) {
	if menus == nil {
		menus = make([]*model.AuthMenuOutput, 0)
	}
	_ = cache.SetJSON(ctx, s.menusCacheKey(userID), menus, authMenusCacheTTL)
}

func (s *sAuth) clearAuthCache(ctx context.Context, userID int64) {
	_ = cache.Delete(ctx, userCacheKeys(userID)...)
}

func ClearUserCaches(ctx context.Context, userIDs ...int64) {
	if len(userIDs) == 0 {
		return
	}
	keys := make([]string, 0, len(userIDs)*2)
	seen := make(map[int64]struct{}, len(userIDs))
	for _, userID := range userIDs {
		if userID <= 0 {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		keys = append(keys, userCacheKeys(userID)...)
	}
	_ = cache.Delete(ctx, keys...)
}

func ClearAllUserCaches(ctx context.Context) {
	var users []struct {
		Id int64 `json:"id"`
	}
	if err := dao.Users.Ctx(ctx).
		Fields(dao.Users.Columns().Id).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&users); err != nil {
		return
	}
	userIDs := make([]int64, 0, len(users))
	for _, item := range users {
		userIDs = append(userIDs, item.Id)
	}
	ClearUserCaches(ctx, userIDs...)
}

func (s *sAuth) loginFailKey(ctx context.Context, username string) string {
	ip := "unknown"
	if req := g.RequestFromCtx(ctx); req != nil {
		if clientIP := strings.TrimSpace(req.GetClientIp()); clientIP != "" {
			ip = clientIP
		}
	}
	return loginFailCacheKey(username, ip)
}

func (s *sAuth) infoCacheKey(userID snowflake.JsonInt64) string {
	return infoCacheKey(int64(userID))
}

func (s *sAuth) menusCacheKey(userID snowflake.JsonInt64) string {
	return menusCacheKey(int64(userID))
}

func userCacheKeys(userID int64) []string {
	if userID <= 0 {
		return nil
	}
	return []string{
		infoCacheKey(userID),
		menusCacheKey(userID),
	}
}

func infoCacheKey(userID int64) string {
	return fmt.Sprintf("system:auth:info:%d", userID)
}

func menusCacheKey(userID int64) string {
	return fmt.Sprintf("system:auth:menus:%d", userID)
}

func loginFailCacheKey(username, ip string) string {
	return fmt.Sprintf("system:auth:login_fail:%s:%s", normalizeAuthKeyPart(username), normalizeAuthKeyPart(ip))
}

func normalizeAuthKeyPart(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "unknown"
	}
	return normalized
}

func compactInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	normalized := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func compactPermissions(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

func collectPermissions(rows []permissionRow) []string {
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.Permission)
	}
	return values
}

func buildMenuTree(list []*model.AuthMenuOutput) []*model.AuthMenuOutput {
	if len(list) == 0 {
		return make([]*model.AuthMenuOutput, 0)
	}
	nodeMap := make(map[int64]*model.AuthMenuOutput, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		item.Children = make([]*model.AuthMenuOutput, 0)
		nodeMap[int64(item.ID)] = item
	}
	tree := make([]*model.AuthMenuOutput, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		if int64(item.ParentID) == 0 {
			tree = append(tree, item)
		} else if parent, ok := nodeMap[int64(item.ParentID)]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			tree = append(tree, item)
		}
	}
	return tree
}
