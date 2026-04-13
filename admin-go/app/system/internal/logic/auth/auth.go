package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/app/system/internal/service"
	"gbaseadmin/utility/appticket"
	"gbaseadmin/utility/authz"
	"gbaseadmin/utility/cache"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/snowflake"
	"gbaseadmin/utility/treeutil"
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

type roleSnapshot struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	IsAdmin int    `json:"isAdmin"`
}

type authUserRecord struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	DeptId   int64  `json:"deptId"`
	Status   int    `json:"status"`
}

const (
	authLoginFailLimit  = 5
	authLoginFailWindow = 10 * time.Minute
	authInfoCacheTTL    = time.Minute
	authMenusCacheTTL   = time.Minute
)

var (
	localTicketReplayMu sync.Mutex
	localTicketReplay   = make(map[string]time.Time)
)

func normalizeAuthLoginInput(in *model.AuthLoginInput) {
	if in == nil {
		return
	}
	in.Username = strings.TrimSpace(in.Username)
}

func normalizeAuthChangePasswordInput(in *model.AuthChangePasswordInput) {
	if in == nil {
		return
	}
	in.OldPassword = strings.TrimSpace(in.OldPassword)
	in.NewPassword = strings.TrimSpace(in.NewPassword)
}

func normalizeAuthTicketLoginInput(in *model.AuthTicketLoginInput) {
	if in == nil {
		return
	}
	in.Ticket = strings.TrimSpace(in.Ticket)
}

func normalizeAuthIssueTicketInput(in *model.AuthIssueTicketInput) {
	if in == nil {
		return
	}
	in.TargetApp = strings.TrimSpace(in.TargetApp)
}

// Login 用户登录
func (s *sAuth) Login(ctx context.Context, in *model.AuthLoginInput) (out *model.AuthLoginOutput, err error) {
	if err := inpututil.Require(in); err != nil {
		return nil, err
	}
	normalizeAuthLoginInput(in)
	if in.Username == "" {
		return nil, gerror.New("用户名不能为空")
	}
	if s.isLoginRateLimited(ctx, in.Username) {
		return nil, gerror.New("登录失败次数过多，请10分钟后再试")
	}

	user, err := s.loadUserByUsername(ctx, in.Username)
	if err != nil {
		return nil, gerror.New("用户名或密码错误")
	}
	if user == nil || user.Id == 0 {
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
				Data(do.Users{
					Password: upgraded,
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

// TicketLogin 应用间票据登录
func (s *sAuth) TicketLogin(ctx context.Context, in *model.AuthTicketLoginInput) (out *model.AuthLoginOutput, err error) {
	if err := inpututil.Require(in); err != nil {
		return nil, err
	}
	normalizeAuthTicketLoginInput(in)
	if in.Ticket == "" {
		return nil, gerror.New("票据不能为空")
	}

	claims, err := appticket.Parse(ctx, in.Ticket)
	if err != nil {
		return nil, gerror.New("票据无效或已过期")
	}
	if strings.TrimSpace(claims.Username) == "" {
		return nil, gerror.New("票据缺少用户名")
	}
	if err := appticket.ValidateTarget(ctx, claims); err != nil {
		return nil, gerror.New("票据目标应用不匹配")
	}
	replayed, err := s.markTicketUsed(ctx, claims, in.Ticket)
	if err != nil {
		return nil, err
	}
	if replayed {
		return nil, gerror.New("票据已使用或已失效")
	}

	user, err := s.loadUserByUsername(ctx, claims.Username)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Id == 0 {
		return nil, gerror.New("票据用户不存在或已删除")
	}
	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

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
	return out, nil
}

// IssueTicket 为当前登录用户生成应用间票据
func (s *sAuth) IssueTicket(ctx context.Context, in *model.AuthIssueTicketInput) (out *model.AuthIssueTicketOutput, err error) {
	if err := inpututil.Require(in); err != nil {
		return nil, err
	}
	normalizeAuthIssueTicketInput(in)
	if in.UserID <= 0 {
		return nil, gerror.New("用户ID不能为空")
	}
	if in.TargetApp == "" {
		return nil, gerror.New("目标应用不能为空")
	}

	user, err := s.loadUserByID(ctx, int64(in.UserID))
	if err != nil {
		return nil, err
	}
	if user == nil || user.Id == 0 {
		return nil, gerror.New("用户不存在或已删除")
	}
	if user.Status == 0 {
		return nil, gerror.New("账号已被禁用")
	}

	ticket, err := appticket.Generate(ctx, user.Username, "", in.TargetApp)
	if err != nil {
		return nil, gerror.New("生成票据失败")
	}

	out = &model.AuthIssueTicketOutput{
		Ticket:    ticket,
		SourceApp: appticket.CurrentAppID(ctx),
		TargetApp: in.TargetApp,
		ExpiresIn: appticket.ExpireSeconds(ctx),
	}
	return out, nil
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
	if user.Id == 0 {
		return nil, gerror.New("用户不存在或已删除")
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

	roles, err := loadUserRoles(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return out, nil
	}
	roleIDs := collectRoleIDs(roles)
	out.Roles = collectRoleTitles(roles)
	if hasAdminRole(roles) {
		out.Perms, err = loadActiveMenuPermissions(ctx, nil)
		return out, err
	}
	menuIDs, err := loadRoleMenuIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	out.Perms, err = loadActiveMenuPermissions(ctx, menuIDs)
	if err != nil {
		return nil, err
	}

	return
}

// ChangePassword 修改密码
func (s *sAuth) ChangePassword(ctx context.Context, in *model.AuthChangePasswordInput) error {
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeAuthChangePasswordInput(in)
	if in.OldPassword == "" {
		return gerror.New("旧密码不能为空")
	}
	if in.NewPassword == "" {
		return gerror.New("新密码不能为空")
	}
	if in.NewPassword == in.OldPassword {
		return gerror.New("新密码不能与旧密码相同")
	}
	if err := password.ValidatePolicy(in.NewPassword); err != nil {
		return gerror.New(err.Error())
	}
	// 查询当前密码
	currentPassword, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, in.UserID).
		Where(dao.Users.Columns().DeletedAt, nil).
		Value(dao.Users.Columns().Password)
	if err != nil {
		return err
	}
	if currentPassword.IsEmpty() {
		return gerror.New("用户不存在或已删除")
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
		Where(dao.Users.Columns().DeletedAt, nil).
		Data(do.Users{
			Password: hashedNew,
		}).
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
	roles, err := loadUserRoles(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}
	roleIDs := collectRoleIDs(roles)
	if len(roleIDs) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}
	if hasAdminRole(roles) {
		list, err := loadActiveMenus(ctx, nil)
		if err != nil {
			return nil, err
		}
		return buildMenuTree(list), nil
	}

	// 查询角色关联的菜单ID（去重）
	menuIDs, err := loadRoleMenuIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	if len(menuIDs) == 0 {
		return make([]*model.AuthMenuOutput, 0), nil
	}

	// 查询菜单详情
	list, err := loadActiveMenus(ctx, menuIDs)
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
		authz.PermissionCacheKey(userID),
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

func (s *sAuth) loadUserByUsername(ctx context.Context, username string) (*authUserRecord, error) {
	var user authUserRecord
	err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Username, username).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *sAuth) loadUserByID(ctx context.Context, userID int64) (*authUserRecord, error) {
	if userID <= 0 {
		return nil, nil
	}
	var user authUserRecord
	err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userID).
		Where(dao.Users.Columns().DeletedAt, nil).
		Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *sAuth) markTicketUsed(ctx context.Context, claims *appticket.Claims, rawTicket string) (bool, error) {
	key := appticket.ReplayCacheKey(claims, rawTicket)
	ttl := appticket.ReplayTTL(claims)
	count, err := cache.IncrWithTTL(ctx, key, ttl)
	if err != nil {
		g.Log().Warningf(ctx, "redis replay protection unavailable, fallback to local ticket cache: %v", err)
		return markTicketUsedLocal(key, ttl), nil
	}
	return count > 1, nil
}

func markTicketUsedLocal(key string, ttl time.Duration) bool {
	key = strings.TrimSpace(key)
	if key == "" {
		return false
	}
	if ttl < time.Second {
		ttl = time.Second
	}
	now := time.Now()

	localTicketReplayMu.Lock()
	defer localTicketReplayMu.Unlock()

	for cacheKey, expiresAt := range localTicketReplay {
		if !expiresAt.After(now) {
			delete(localTicketReplay, cacheKey)
		}
	}
	if expiresAt, ok := localTicketReplay[key]; ok && expiresAt.After(now) {
		return true
	}
	localTicketReplay[key] = now.Add(ttl)
	return false
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

func loadUserRoles(ctx context.Context, userID int64) ([]roleSnapshot, error) {
	if userID <= 0 {
		return nil, nil
	}
	var userRoles []struct {
		RoleId int64 `json:"roleId"`
	}
	if err := dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userID).Scan(&userRoles); err != nil {
		return nil, err
	}
	roleIDs := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleId)
	}
	roleIDs = compactInt64s(roleIDs)
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var roles []roleSnapshot
	if err := g.DB().Ctx(ctx).Model("system_role").
		Fields("id,title,is_admin").
		Where("id", roleIDs).
		Where("deleted_at", nil).
		Where("status", 1).
		Scan(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func collectRoleIDs(roles []roleSnapshot) []int64 {
	ids := make([]int64, 0, len(roles))
	for _, role := range roles {
		ids = append(ids, role.ID)
	}
	return compactInt64s(ids)
}

func collectRoleTitles(roles []roleSnapshot) []string {
	titles := make([]string, 0, len(roles))
	for _, role := range roles {
		titles = append(titles, role.Title)
	}
	return compactPermissions(titles)
}

func hasAdminRole(roles []roleSnapshot) bool {
	for _, role := range roles {
		if role.IsAdmin == 1 {
			return true
		}
	}
	return false
}

func loadRoleMenuIDs(ctx context.Context, roleIDs []int64) ([]int64, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var roleMenus []struct {
		MenuId int64 `json:"menuId"`
	}
	if err := dao.RoleMenu.Ctx(ctx).WhereIn(dao.RoleMenu.Columns().RoleId, roleIDs).Scan(&roleMenus); err != nil {
		return nil, err
	}
	menuIDs := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuId)
	}
	return compactInt64s(menuIDs), nil
}

func activeMenuModel(ctx context.Context) *gdb.Model {
	return g.DB().Ctx(ctx).Model("system_menu").
		Where("deleted_at", nil).
		Where("status", 1)
}

func loadActiveMenuPermissions(ctx context.Context, menuIDs []int64) ([]string, error) {
	var perms []permissionRow
	model := activeMenuModel(ctx)
	if len(menuIDs) > 0 {
		model = model.Where("id", menuIDs)
	}
	if err := model.WhereNot("permission", "").Scan(&perms); err != nil {
		return nil, err
	}
	return compactPermissions(collectPermissions(perms)), nil
}

func loadActiveMenus(ctx context.Context, menuIDs []int64) ([]*model.AuthMenuOutput, error) {
	var list []*model.AuthMenuOutput
	model := activeMenuModel(ctx).WhereNot("type", 3)
	if len(menuIDs) > 0 {
		model = model.Where("id", menuIDs)
	}
	if err := model.OrderAsc("sort").Scan(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func collectPermissions(rows []permissionRow) []string {
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.Permission)
	}
	return values
}

func buildMenuTree(list []*model.AuthMenuOutput) []*model.AuthMenuOutput {
	return treeutil.BuildForest(list, treeutil.TreeNodeAccessor[*model.AuthMenuOutput]{
		ID: func(item *model.AuthMenuOutput) int64 {
			if item == nil {
				return 0
			}
			return int64(item.ID)
		},
		ParentID: func(item *model.AuthMenuOutput) int64 {
			if item == nil {
				return 0
			}
			return int64(item.ParentID)
		},
		Init: func(item *model.AuthMenuOutput) {
			if item != nil {
				item.Children = make([]*model.AuthMenuOutput, 0)
			}
		},
		Append: func(parent *model.AuthMenuOutput, child *model.AuthMenuOutput) {
			if parent != nil && child != nil {
				parent.Children = append(parent.Children, child)
			}
		},
	})
}
