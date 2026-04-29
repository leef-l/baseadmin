package shared

import (
	"context"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/model/do"
	"gbaseadmin/utility/password"
	"gbaseadmin/utility/snowflake"
)

type AdminBootstrapMenuProfile int

const (
	AdminBootstrapMenuProfileTenant AdminBootstrapMenuProfile = iota + 1
	AdminBootstrapMenuProfileMerchant
)

type AdminBootstrapInput struct {
	TenantID        snowflake.JsonInt64
	MerchantID      snowflake.JsonInt64
	DeptParentID    snowflake.JsonInt64
	DeptTitle       string
	DeptManagerName string
	RoleTitle       string
	AdminUsername   string
	AdminPassword   string
	AdminNickname   string
	AdminEmail      string
	MenuProfile     AdminBootstrapMenuProfile
	CreatedBy       int64
}

type adminBootstrapMenuRow struct {
	ID         int64  `json:"id"`
	ParentID   int64  `json:"parentId"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
}

func BootstrapScopedAdmin(ctx context.Context, tx gdb.TX, in AdminBootstrapInput) error {
	normalizeAdminBootstrapInput(&in)
	if err := validateAdminBootstrapInput(in); err != nil {
		return err
	}
	if err := ensureBootstrapUserUnique(ctx, tx, in.AdminUsername, in.AdminEmail); err != nil {
		return err
	}
	if err := ensureBootstrapDeptTitleUnique(ctx, tx, in.DeptParentID, in.DeptTitle, in.TenantID, in.MerchantID); err != nil {
		return err
	}
	if err := ensureBootstrapRoleTitleUnique(ctx, tx, in.RoleTitle, in.TenantID, in.MerchantID); err != nil {
		return err
	}

	hashedPassword, err := password.Hash(in.AdminPassword)
	if err != nil {
		return err
	}

	deptID := snowflake.Generate()
	roleID := snowflake.Generate()
	userID := snowflake.Generate()

	if _, err = tx.Model(dao.Dept.Table()).Ctx(ctx).Data(do.Dept{
		Id:         deptID,
		ParentId:   in.DeptParentID,
		Title:      in.DeptTitle,
		Username:   in.DeptManagerName,
		Email:      in.AdminEmail,
		Sort:       0,
		Status:     1,
		CreatedBy:  in.CreatedBy,
		DeptId:     deptID,
		TenantId:   in.TenantID,
		MerchantId: in.MerchantID,
	}).Insert(); err != nil {
		return err
	}

	if _, err = tx.Model(dao.Role.Table()).Ctx(ctx).Data(do.Role{
		Id:         roleID,
		ParentId:   0,
		Title:      in.RoleTitle,
		DataScope:  1,
		Sort:       0,
		Status:     1,
		IsAdmin:    1,
		CreatedBy:  in.CreatedBy,
		DeptId:     deptID,
		TenantId:   in.TenantID,
		MerchantId: in.MerchantID,
	}).Insert(); err != nil {
		return err
	}

	menuIDs, err := loadBootstrapMenuIDs(ctx, tx, in.MenuProfile)
	if err != nil {
		return err
	}
	if len(menuIDs) > 0 {
		roleMenus := make([]do.RoleMenu, 0, len(menuIDs))
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, do.RoleMenu{
				RoleId: roleID,
				MenuId: menuID,
			})
		}
		if _, err = tx.Model(dao.RoleMenu.Table()).Ctx(ctx).Data(roleMenus).Insert(); err != nil {
			return err
		}
	}

	if _, err = tx.Model(dao.Users.Table()).Ctx(ctx).Data(do.Users{
		Id:         userID,
		Username:   in.AdminUsername,
		Password:   hashedPassword,
		Nickname:   in.AdminNickname,
		Email:      in.AdminEmail,
		Status:     1,
		CreatedBy:  in.CreatedBy,
		DeptId:     deptID,
		TenantId:   in.TenantID,
		MerchantId: in.MerchantID,
	}).Insert(); err != nil {
		return err
	}

	_, err = tx.Model(dao.UserRole.Table()).Ctx(ctx).Data(do.UserRole{
		UserId: userID,
		RoleId: roleID,
	}).Insert()
	return err
}

func normalizeAdminBootstrapInput(in *AdminBootstrapInput) {
	if in == nil {
		return
	}
	in.DeptTitle = strings.TrimSpace(in.DeptTitle)
	in.DeptManagerName = strings.TrimSpace(in.DeptManagerName)
	in.RoleTitle = strings.TrimSpace(in.RoleTitle)
	in.AdminUsername = strings.TrimSpace(in.AdminUsername)
	in.AdminPassword = strings.TrimSpace(in.AdminPassword)
	in.AdminNickname = strings.TrimSpace(in.AdminNickname)
	in.AdminEmail = strings.TrimSpace(in.AdminEmail)
}

func validateAdminBootstrapInput(in AdminBootstrapInput) error {
	if in.TenantID <= 0 {
		return gerror.New("租户不能为空")
	}
	if in.DeptTitle == "" {
		return gerror.New("默认部门名称不能为空")
	}
	if in.RoleTitle == "" {
		return gerror.New("管理员角色名称不能为空")
	}
	if in.AdminUsername == "" {
		return gerror.New("管理员用户名不能为空")
	}
	if err := password.ValidatePolicy(in.AdminPassword); err != nil {
		return gerror.New(err.Error())
	}
	return nil
}

func ensureBootstrapUserUnique(ctx context.Context, tx gdb.TX, username, email string) error {
	count, err := tx.Model(dao.Users.Table()).Ctx(ctx).
		Where(dao.Users.Columns().Username, username).
		Where(dao.Users.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("管理员用户名已存在")
	}
	if email == "" {
		return nil
	}
	count, err = tx.Model(dao.Users.Table()).Ctx(ctx).
		Where(dao.Users.Columns().Email, email).
		Where(dao.Users.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("管理员邮箱已存在")
	}
	return nil
}

func ensureBootstrapDeptTitleUnique(
	ctx context.Context,
	tx gdb.TX,
	parentID snowflake.JsonInt64,
	title string,
	tenantID snowflake.JsonInt64,
	merchantID snowflake.JsonInt64,
) error {
	count, err := tx.Model(dao.Dept.Table()).Ctx(ctx).
		Where(dao.Dept.Columns().ParentId, parentID).
		Where(dao.Dept.Columns().Title, title).
		Where(dao.Dept.Columns().TenantId, tenantID).
		Where(dao.Dept.Columns().MerchantId, merchantID).
		Where(dao.Dept.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("默认部门名称已存在")
	}
	return nil
}

func ensureBootstrapRoleTitleUnique(
	ctx context.Context,
	tx gdb.TX,
	title string,
	tenantID snowflake.JsonInt64,
	merchantID snowflake.JsonInt64,
) error {
	count, err := tx.Model(dao.Role.Table()).Ctx(ctx).
		Where(dao.Role.Columns().Title, title).
		Where(dao.Role.Columns().TenantId, tenantID).
		Where(dao.Role.Columns().MerchantId, merchantID).
		Where(dao.Role.Columns().DeletedAt, nil).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("管理员角色名称已存在")
	}
	return nil
}

func loadBootstrapMenuIDs(ctx context.Context, tx gdb.TX, profile AdminBootstrapMenuProfile) ([]int64, error) {
	var rows []adminBootstrapMenuRow
	err := tx.Model(dao.Menu.Table()).Ctx(ctx).
		Fields(
			dao.Menu.Columns().Id+" AS id",
			dao.Menu.Columns().ParentId+" AS parentId",
			dao.Menu.Columns().Path,
			dao.Menu.Columns().Permission,
		).
		Where(dao.Menu.Columns().DeletedAt, nil).
		Where(dao.Menu.Columns().Status, 1).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	return selectBootstrapMenuIDs(rows, profile), nil
}

func selectBootstrapMenuIDs(rows []adminBootstrapMenuRow, profile AdminBootstrapMenuProfile) []int64 {
	if len(rows) == 0 {
		return nil
	}
	byID := make(map[int64]adminBootstrapMenuRow, len(rows))
	selected := make(map[int64]struct{}, len(rows))
	pathSet := map[string]struct{}{
		"/dashboard": {},
		"/system":    {},
		"/workspace": {},
	}
	prefixes := []string{
		"system:dept:",
		"system:domain:",
		"system:role:",
		"system:user:",
	}
	if profile == AdminBootstrapMenuProfileTenant {
		prefixes = append(prefixes, "system:merchant:")
	}

	for _, row := range rows {
		if row.ID <= 0 {
			continue
		}
		byID[row.ID] = row
		if _, ok := pathSet[strings.TrimSpace(row.Path)]; ok {
			selected[row.ID] = struct{}{}
			continue
		}
		permission := strings.TrimSpace(row.Permission)
		if permission == "system:domain:apply" || permission == "system:domain:ssl" {
			continue
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(permission, prefix) {
				selected[row.ID] = struct{}{}
				break
			}
		}
	}

	for id := range selected {
		seenParents := make(map[int64]struct{})
		for parentID := byID[id].ParentID; parentID > 0; {
			parent, ok := byID[parentID]
			if !ok {
				break
			}
			if _, seen := seenParents[parentID]; seen {
				break
			}
			seenParents[parentID] = struct{}{}
			selected[parentID] = struct{}{}
			parentID = parent.ParentID
		}
	}

	menuIDs := make([]int64, 0, len(selected))
	for id := range selected {
		if id > 0 {
			menuIDs = append(menuIDs, id)
		}
	}
	sort.Slice(menuIDs, func(i, j int) bool {
		return menuIDs[i] < menuIDs[j]
	})
	return menuIDs
}
