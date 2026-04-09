package shared

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type DataAccessScope struct {
	UserID      int64
	DeptID      int64
	All         bool
	IncludeSelf bool
	DeptIDs     []int64
}

type dataScopeRoleRow struct {
	RoleID    int64 `json:"roleId"`
	DataScope int   `json:"dataScope"`
	IsAdmin   int   `json:"isAdmin"`
}

type dataScopeDeptRow struct {
	ID       int64 `json:"id"`
	ParentID int64 `json:"parentId"`
}

func CurrentActorUserID(ctx context.Context) int64 {
	userID, _ := currentActor(ctx)
	return userID
}

func CurrentActorDeptID(ctx context.Context) int64 {
	_, deptID := currentActor(ctx)
	return deptID
}

func HasCurrentActorAdminRole(ctx context.Context) (bool, error) {
	return HasUserAdminRole(ctx, CurrentActorUserID(ctx))
}

func HasUserAdminRole(ctx context.Context, userID int64) (bool, error) {
	roleRows, err := loadUserDataScopeRoles(ctx, userID)
	if err != nil {
		return false, err
	}
	return hasAdminDataScope(roleRows), nil
}

func ResolveDataAccessScope(ctx context.Context) (DataAccessScope, error) {
	userID, deptID := currentActor(ctx)
	scope := DataAccessScope{
		UserID: userID,
		DeptID: deptID,
	}
	if userID <= 0 {
		scope.IncludeSelf = true
		return scope, nil
	}

	roleRows, err := loadUserDataScopeRoles(ctx, userID)
	if err != nil {
		return DataAccessScope{}, err
	}
	if hasAdminDataScope(roleRows) {
		scope.All = true
		return scope, nil
	}

	modes := selectDataScopeModes(roleRows)
	if modes.AllowAll {
		scope.All = true
		return scope, nil
	}

	var deptIDs []int64
	if modes.IncludeDeptAndChildren {
		deptIDs, err = loadDeptAndChildrenIDs(ctx, deptID)
		if err != nil {
			return DataAccessScope{}, err
		}
	}
	if modes.IncludeCurrentDept || modes.IncludeSelf {
		deptIDs = append(deptIDs, deptID)
	}
	if modes.IncludeCustomDept {
		customDeptIDs, customErr := loadCustomScopeDeptIDs(ctx, roleRows)
		if customErr != nil {
			return DataAccessScope{}, customErr
		}
		deptIDs = append(deptIDs, customDeptIDs...)
	}
	scope.IncludeSelf = modes.IncludeSelf
	scope.DeptIDs = compactPositiveIDs(deptIDs)
	if err != nil {
		return DataAccessScope{}, err
	}
	if !scope.IncludeSelf && len(scope.DeptIDs) == 0 {
		scope.IncludeSelf = true
	}
	return scope, nil
}

func ApplyDeptScopeToUserModel(ctx context.Context, m *gdb.Model, userIDColumn, deptIDColumn string) (*gdb.Model, error) {
	scope, err := ResolveDataAccessScope(ctx)
	if err != nil {
		return nil, err
	}
	if scope.All {
		return m, nil
	}
	return applyUserScopeFilter(m, userIDColumn, deptIDColumn, scope), nil
}

func ApplyDeptScopeToDeptModel(ctx context.Context, m *gdb.Model, deptIDColumn string) (*gdb.Model, error) {
	scope, err := ResolveDataAccessScope(ctx)
	if err != nil {
		return nil, err
	}
	if scope.All {
		return m, nil
	}
	return applyScopeDeptFilter(m, deptIDColumn, scope.DeptIDs), nil
}

func CanAccessUser(ctx context.Context, userID, deptID int64) (bool, error) {
	scope, err := ResolveDataAccessScope(ctx)
	if err != nil {
		return false, err
	}
	if scope.All {
		return true, nil
	}
	if scope.IncludeSelf && scope.UserID > 0 && scope.UserID == userID {
		return true, nil
	}
	return containsInt64(scope.DeptIDs, deptID), nil
}

func CanAccessDept(ctx context.Context, deptID int64) (bool, error) {
	scope, err := ResolveDataAccessScope(ctx)
	if err != nil {
		return false, err
	}
	if scope.All {
		return true, nil
	}
	return containsInt64(scope.DeptIDs, deptID), nil
}

func loadUserDataScopeRoles(ctx context.Context, userID int64) ([]dataScopeRoleRow, error) {
	if userID <= 0 {
		return nil, nil
	}
	var rows []dataScopeRoleRow
	err := g.DB().Ctx(ctx).
		Model("system_user_role ur").
		LeftJoin("system_role r", "r.id = ur.role_id").
		Fields("ur.role_id AS roleId", "r.data_scope AS dataScope", "r.is_admin AS isAdmin").
		Where("ur.user_id", userID).
		Where("r.deleted_at", nil).
		Where("r.status", 1).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func loadDeptAndChildrenIDs(ctx context.Context, rootDeptID int64) ([]int64, error) {
	rootIDs := compactPositiveIDs([]int64{rootDeptID})
	if len(rootIDs) == 0 {
		return nil, nil
	}
	var rows []dataScopeDeptRow
	if err := g.DB().Ctx(ctx).
		Model("system_dept").
		Fields("id", "parent_id AS parentId").
		Where("deleted_at", nil).
		Scan(&rows); err != nil {
		return nil, err
	}
	tree := make(map[int64][]int64, len(rows))
	for _, row := range rows {
		tree[row.ParentID] = append(tree[row.ParentID], row.ID)
	}
	return expandDeptTree(rootIDs, tree), nil
}

func loadCustomScopeDeptIDs(ctx context.Context, rows []dataScopeRoleRow) ([]int64, error) {
	roleIDs := customScopeRoleIDs(rows)
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var deptRows []struct {
		DeptID int64 `json:"deptId"`
	}
	if err := g.DB().Ctx(ctx).
		Model("system_role_dept").
		Fields("dept_id AS deptId").
		WhereIn("role_id", roleIDs).
		Scan(&deptRows); err != nil {
		return nil, err
	}
	deptIDs := make([]int64, 0, len(deptRows))
	for _, row := range deptRows {
		deptIDs = append(deptIDs, row.DeptID)
	}
	return compactPositiveIDs(deptIDs), nil
}

func applyScopeDeptFilter(m *gdb.Model, deptIDColumn string, deptIDs []int64) *gdb.Model {
	if m == nil {
		return nil
	}
	if deptIDColumn == "" {
		return m.Where("1 = 0")
	}
	deptIDs = compactPositiveIDs(deptIDs)
	if len(deptIDs) == 0 {
		return m.Where("1 = 0")
	}
	return m.WhereIn(deptIDColumn, deptIDs)
}

func applyUserScopeFilter(m *gdb.Model, userIDColumn, deptIDColumn string, scope DataAccessScope) *gdb.Model {
	if m == nil {
		return nil
	}
	builder := m.Builder()
	hasFilter := false
	if scope.IncludeSelf && scope.UserID > 0 && userIDColumn != "" {
		builder = builder.Where(userIDColumn, scope.UserID)
		hasFilter = true
	}
	if deptIDColumn != "" && len(scope.DeptIDs) > 0 {
		if hasFilter {
			builder = builder.WhereOrIn(deptIDColumn, scope.DeptIDs)
		} else {
			builder = builder.WhereIn(deptIDColumn, scope.DeptIDs)
			hasFilter = true
		}
	}
	if !hasFilter {
		return m.Where("1 = 0")
	}
	return m.Where(builder)
}

func currentActor(ctx context.Context) (userID, deptID int64) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0, 0
	}
	return req.GetCtxVar("jwt_user_id").Int64(), req.GetCtxVar("jwt_dept_id").Int64()
}

func hasAdminDataScope(rows []dataScopeRoleRow) bool {
	for _, row := range rows {
		if row.IsAdmin == 1 {
			return true
		}
	}
	return false
}

type dataScopeModes struct {
	AllowAll               bool
	IncludeDeptAndChildren bool
	IncludeCurrentDept     bool
	IncludeSelf            bool
	IncludeCustomDept      bool
}

func selectDataScopeModes(rows []dataScopeRoleRow) dataScopeModes {
	if len(rows) == 0 {
		return dataScopeModes{IncludeSelf: true}
	}
	var modes dataScopeModes
	for _, row := range rows {
		switch normalizeDataScope(row.DataScope) {
		case 1:
			modes.AllowAll = true
		case 2:
			modes.IncludeDeptAndChildren = true
		case 3:
			modes.IncludeCurrentDept = true
		case 4:
			modes.IncludeSelf = true
		case 5:
			modes.IncludeCustomDept = true
		}
	}
	if !modes.AllowAll && !modes.IncludeDeptAndChildren && !modes.IncludeCurrentDept && !modes.IncludeSelf && !modes.IncludeCustomDept {
		modes.IncludeSelf = true
	}
	return modes
}

func normalizeDataScope(scope int) int {
	switch scope {
	case 1, 2, 3, 4, 5:
		return scope
	default:
		return 4
	}
}

func customScopeRoleIDs(rows []dataScopeRoleRow) []int64 {
	roleIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if normalizeDataScope(row.DataScope) == 5 {
			roleIDs = append(roleIDs, row.RoleID)
		}
	}
	return compactPositiveIDs(roleIDs)
}

func expandDeptTree(rootIDs []int64, tree map[int64][]int64) []int64 {
	rootIDs = compactPositiveIDs(rootIDs)
	if len(rootIDs) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(rootIDs))
	queue := append([]int64(nil), rootIDs...)
	expanded := make([]int64, 0, len(rootIDs))
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current <= 0 {
			continue
		}
		if _, ok := seen[current]; ok {
			continue
		}
		seen[current] = struct{}{}
		expanded = append(expanded, current)
		queue = append(queue, tree[current]...)
	}
	sort.Slice(expanded, func(i, j int) bool {
		return expanded[i] < expanded[j]
	})
	return expanded
}

func containsInt64(values []int64, target int64) bool {
	if target <= 0 {
		return false
	}
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
