package shared

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// LoadCurrentActorAssignableRoleIDs returns the enabled role IDs the current actor
// is allowed to assign. Super admins may assign every enabled role. Non-admin
// actors may only assign their own enabled roles.
func LoadCurrentActorAssignableRoleIDs(ctx context.Context) ([]int64, error) {
	return LoadUserAssignableRoleIDs(ctx, CurrentActorUserID(ctx))
}

func LoadUserAssignableRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	if userID <= 0 {
		return nil, nil
	}
	actorHasAdmin, err := HasUserAdminRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	var rows []struct {
		RoleID int64 `json:"roleId"`
	}
	m := g.DB().Ctx(ctx).
		Model("system_role r").
		Fields("r.id AS roleId").
		Where("r.deleted_at", nil).
		Where("r.status", 1)
	if !actorHasAdmin {
		m = m.LeftJoin("system_user_role ur", "ur.role_id = r.id").
			Where("ur.user_id", userID)
	}
	if err := m.Scan(&rows); err != nil {
		return nil, err
	}

	roleIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		roleIDs = append(roleIDs, row.RoleID)
	}
	return compactPositiveIDs(roleIDs), nil
}

func RoleIDsWithinScope(targetIDs, allowedIDs []int64) bool {
	targetIDs = compactPositiveIDs(targetIDs)
	if len(targetIDs) == 0 {
		return true
	}
	allowedIDs = compactPositiveIDs(allowedIDs)
	if len(allowedIDs) == 0 {
		return false
	}

	allowedSet := make(map[int64]struct{}, len(allowedIDs))
	for _, id := range allowedIDs {
		allowedSet[id] = struct{}{}
	}
	for _, id := range targetIDs {
		if _, ok := allowedSet[id]; !ok {
			return false
		}
	}
	return true
}
