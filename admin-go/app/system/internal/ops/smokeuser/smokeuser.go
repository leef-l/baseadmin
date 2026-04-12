package smokeuser

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"unicode"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/system/internal/dao"
	authlogic "gbaseadmin/app/system/internal/logic/auth"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/snowflake"
)

const (
	smokeRoleTitle       = "CI冒烟测试"
	defaultSmokeNickname = "CI 冒烟"
	smokeRoleDataScope   = 4
	smokeRoleSort        = 999
)

var smokePermissions = []string{
	"upload:dir:list",
	"upload:dir_rule:list",
	"upload:dir_rule:create",
	"upload:dir_rule:delete",
	"upload:file:list",
	"upload:file:create",
	"upload:file:delete",
}

type EnsureOptions struct {
	Username string
	Password string
	Nickname string
}

type Result struct {
	RoleID      int64
	UserID      int64
	RoleCreated bool
	RoleUpdated bool
	UserCreated bool
	UserUpdated bool
	MenuCount   int
}

type roleState struct {
	id      int64
	created bool
	updated bool
}

type userState struct {
	id      int64
	created bool
	updated bool
}

func Ensure(ctx context.Context, opts EnsureOptions) (*Result, error) {
	opts = normalizeOptions(opts)
	if err := validateOptions(opts); err != nil {
		return nil, err
	}
	if err := password.ValidatePolicy(opts.Password); err != nil {
		return nil, gerror.New(err.Error())
	}

	deptID, err := loadDefaultDeptID(ctx)
	if err != nil {
		return nil, err
	}
	menuIDs, err := loadSmokeMenuIDs(ctx)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash(opts.Password)
	if err != nil {
		return nil, gerror.Wrap(err, "hash smoke password")
	}

	var (
		result    = &Result{MenuCount: len(menuIDs)}
		roleAfter roleState
		userAfter userState
	)

	err = dao.Role.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		roleAfter, err = ensureRoleTx(ctx, tx, deptID, menuIDs)
		if err != nil {
			return err
		}
		userAfter, err = ensureUserTx(ctx, tx, opts, hashedPassword, deptID, roleAfter.id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result.RoleID = roleAfter.id
	result.UserID = userAfter.id
	result.RoleCreated = roleAfter.created
	result.RoleUpdated = roleAfter.updated
	result.UserCreated = userAfter.created
	result.UserUpdated = userAfter.updated

	if !userAfter.created || roleAfter.updated {
		authlogic.ClearUserCaches(ctx, userAfter.id)
	}
	return result, nil
}

func normalizeOptions(opts EnsureOptions) EnsureOptions {
	opts.Username = strings.TrimSpace(opts.Username)
	opts.Password = strings.TrimSpace(opts.Password)
	opts.Nickname = strings.TrimSpace(opts.Nickname)
	if opts.Nickname == "" {
		opts.Nickname = defaultSmokeNickname
	}
	return opts
}

func validateOptions(opts EnsureOptions) error {
	if opts.Username == "" {
		return gerror.New("烟雾账号用户名不能为空")
	}
	if opts.Password == "" {
		return gerror.New("烟雾账号密码不能为空")
	}
	if containsWhitespace(opts.Username) {
		return gerror.New("烟雾账号用户名不能包含空白字符")
	}
	return nil
}

func containsWhitespace(value string) bool {
	for _, item := range value {
		if unicode.IsSpace(item) {
			return true
		}
	}
	return false
}

func loadDefaultDeptID(ctx context.Context) (int64, error) {
	var row struct {
		ID int64 `json:"id"`
	}
	err := dao.Dept.Ctx(ctx).
		Fields(dao.Dept.Columns().Id).
		Where(dao.Dept.Columns().Status, 1).
		OrderAsc(dao.Dept.Columns().Sort).
		OrderAsc(dao.Dept.Columns().Id).
		Limit(1).
		Scan(&row)
	if err != nil {
		return 0, err
	}
	if row.ID == 0 {
		return 0, gerror.New("未找到可用部门，无法创建 CI 冒烟账号")
	}
	return row.ID, nil
}

func loadSmokeMenuIDs(ctx context.Context) ([]int64, error) {
	type menuRow struct {
		ID         int64  `json:"id"`
		Permission string `json:"permission"`
	}

	var rows []menuRow
	if err := dao.Menu.Ctx(ctx).
		Fields(dao.Menu.Columns().Id, dao.Menu.Columns().Permission).
		Where(dao.Menu.Columns().Status, 1).
		WhereIn(dao.Menu.Columns().Permission, smokePermissions).
		OrderAsc(dao.Menu.Columns().Id).
		Scan(&rows); err != nil {
		return nil, err
	}

	menuIDs := make([]int64, 0, len(rows))
	seenPerms := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if row.ID <= 0 || strings.TrimSpace(row.Permission) == "" {
			continue
		}
		menuIDs = append(menuIDs, row.ID)
		seenPerms[row.Permission] = struct{}{}
	}

	missing := make([]string, 0, len(smokePermissions))
	for _, permission := range smokePermissions {
		if _, ok := seenPerms[permission]; ok {
			continue
		}
		missing = append(missing, permission)
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		return nil, gerror.Newf("CI 冒烟角色缺少权限菜单: %s", strings.Join(missing, ", "))
	}
	return menuIDs, nil
}

func ensureRoleTx(ctx context.Context, tx gdb.TX, deptID int64, menuIDs []int64) (roleState, error) {
	type roleRow struct {
		ID        int64       `json:"id"`
		DeletedAt *gtime.Time `json:"deletedAt"`
	}

	var (
		state roleState
		row   roleRow
	)

	err := tx.Model(dao.Role.Table()).Ctx(ctx).
		Unscoped().
		Fields(dao.Role.Columns().Id, dao.Role.Columns().DeletedAt).
		Where(dao.Role.Columns().Title, smokeRoleTitle).
		Limit(1).
		Scan(&row)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return state, err
	}

	roleData := do.Role{
		ParentId:  0,
		Title:     smokeRoleTitle,
		DataScope: smokeRoleDataScope,
		Sort:      smokeRoleSort,
		Status:    1,
		IsAdmin:   0,
		CreatedBy: 0,
		DeptId:    deptID,
	}

	if row.ID == 0 {
		state.id = int64(snowflake.Generate())
		roleData.Id = state.id
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).Data(roleData).Insert(); err != nil {
			return state, err
		}
		state.created = true
	} else {
		state.id = row.ID
		if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
			Unscoped().
			Where(dao.Role.Columns().Id, row.ID).
			Data(roleData).
			Update(); err != nil {
			return state, err
		}
		if row.DeletedAt != nil {
			if _, err := tx.Model(dao.Role.Table()).Ctx(ctx).
				Unscoped().
				Where(dao.Role.Columns().Id, row.ID).
				Data(g.Map{dao.Role.Columns().DeletedAt: nil}).
				Update(); err != nil {
				return state, err
			}
		}
		state.updated = true
	}

	if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).
		Where(dao.RoleMenu.Columns().RoleId, state.id).
		Delete(); err != nil {
		return state, err
	}

	roleMenus := make([]do.RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, do.RoleMenu{
			RoleId: state.id,
			MenuId: menuID,
		})
	}
	if len(roleMenus) > 0 {
		if _, err := tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Data(roleMenus).Insert(); err != nil {
			return state, err
		}
	}

	return state, nil
}

func ensureUserTx(
	ctx context.Context,
	tx gdb.TX,
	opts EnsureOptions,
	hashedPassword string,
	deptID int64,
	roleID int64,
) (userState, error) {
	type userRow struct {
		ID        int64       `json:"id"`
		DeletedAt *gtime.Time `json:"deletedAt"`
	}

	var (
		state userState
		row   userRow
	)

	err := tx.Model(dao.Users.Table()).Ctx(ctx).
		Unscoped().
		Fields(dao.Users.Columns().Id, dao.Users.Columns().DeletedAt).
		Where(dao.Users.Columns().Username, opts.Username).
		Limit(1).
		Scan(&row)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return state, err
	}

	userData := do.Users{
		Username:  opts.Username,
		Password:  hashedPassword,
		Nickname:  opts.Nickname,
		Status:    1,
		CreatedBy: 0,
		DeptId:    deptID,
	}

	if row.ID == 0 {
		state.id = int64(snowflake.Generate())
		userData.Id = state.id
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).Data(userData).Insert(); err != nil {
			return state, err
		}
		state.created = true
	} else {
		state.id = row.ID
		if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).
			Unscoped().
			Where(dao.Users.Columns().Id, row.ID).
			Data(userData).
			Update(); err != nil {
			return state, err
		}
		if row.DeletedAt != nil {
			if _, err := tx.Model(dao.Users.Table()).Ctx(ctx).
				Unscoped().
				Where(dao.Users.Columns().Id, row.ID).
				Data(g.Map{dao.Users.Columns().DeletedAt: nil}).
				Update(); err != nil {
				return state, err
			}
		}
		state.updated = true
	}

	if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).
		Where(dao.UserRole.Columns().UserId, state.id).
		Delete(); err != nil {
		return state, err
	}
	if _, err := tx.Model(dao.UserRole.Table()).Ctx(ctx).Data(do.UserRole{
		UserId: state.id,
		RoleId: roleID,
	}).Insert(); err != nil {
		return state, err
	}

	return state, nil
}
