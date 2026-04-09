package authz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/cache"
)

const permissionCacheTTL = time.Minute

type permissionState struct {
	IsAdmin bool     `json:"isAdmin"`
	Perms   []string `json:"perms"`
}

type roleRow struct {
	RoleID  int64 `json:"roleId"`
	IsAdmin int   `json:"isAdmin"`
}

type permissionRow struct {
	Permission string `json:"permission"`
}

func PermissionCacheKey(userID int64) string {
	if userID <= 0 {
		return ""
	}
	return fmt.Sprintf("system:authz:perms:%d", userID)
}

func HasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	permission = strings.TrimSpace(permission)
	if permission == "" {
		return true, nil
	}
	state, err := loadPermissionState(ctx, userID)
	if err != nil {
		return false, err
	}
	if state.IsAdmin {
		return true, nil
	}
	for _, item := range state.Perms {
		if item == permission {
			return true, nil
		}
	}
	return false, nil
}

func loadPermissionState(ctx context.Context, userID int64) (permissionState, error) {
	if userID <= 0 {
		return permissionState{}, nil
	}
	cacheKey := PermissionCacheKey(userID)
	var state permissionState
	if ok, err := cache.GetJSON(ctx, cacheKey, &state); err == nil && ok {
		return state, nil
	}

	state, err := queryPermissionState(ctx, userID)
	if err != nil {
		return permissionState{}, err
	}
	_ = cache.SetJSON(ctx, cacheKey, state, permissionCacheTTL)
	return state, nil
}

func queryPermissionState(ctx context.Context, userID int64) (permissionState, error) {
	var roles []roleRow
	err := g.DB().Ctx(ctx).
		Model("system_user_role ur").
		LeftJoin("system_role r", "r.id = ur.role_id").
		Fields("ur.role_id AS roleId", "r.is_admin AS isAdmin").
		Where("ur.user_id", userID).
		Where("r.deleted_at", nil).
		Where("r.status", 1).
		Scan(&roles)
	if err != nil {
		return permissionState{}, err
	}

	if hasAdminRole(roles) {
		return permissionState{
			IsAdmin: true,
			Perms:   nil,
		}, nil
	}

	roleIDs := collectRoleIDs(roles)
	if len(roleIDs) == 0 {
		return permissionState{}, nil
	}

	var rows []permissionRow
	err = g.DB().Ctx(ctx).
		Model("system_role_menu rm").
		LeftJoin("system_menu m", "m.id = rm.menu_id").
		Fields("m.permission AS permission").
		WhereIn("rm.role_id", roleIDs).
		Where("m.deleted_at", nil).
		Where("m.status", 1).
		WhereNot("m.permission", "").
		Scan(&rows)
	if err != nil {
		return permissionState{}, err
	}

	return permissionState{
		IsAdmin: false,
		Perms:   compactPermissions(rows),
	}, nil
}

func hasAdminRole(rows []roleRow) bool {
	for _, row := range rows {
		if row.IsAdmin == 1 {
			return true
		}
	}
	return false
}

func collectRoleIDs(rows []roleRow) []int64 {
	if len(rows) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(rows))
	roleIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row.RoleID <= 0 {
			continue
		}
		if _, ok := seen[row.RoleID]; ok {
			continue
		}
		seen[row.RoleID] = struct{}{}
		roleIDs = append(roleIDs, row.RoleID)
	}
	if len(roleIDs) == 0 {
		return nil
	}
	return roleIDs
}

func compactPermissions(rows []permissionRow) []string {
	if len(rows) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(rows))
	perms := make([]string, 0, len(rows))
	for _, row := range rows {
		value := strings.TrimSpace(row.Permission)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		perms = append(perms, value)
	}
	if len(perms) == 0 {
		return nil
	}
	return perms
}
