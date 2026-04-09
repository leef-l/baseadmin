package middleware

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// GetUserID 从 context 中获取当前登录用户 ID
func GetUserID(ctx context.Context) snowflake.JsonInt64 {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0
	}
	val := req.GetCtxVar("jwt_user_id")
	return snowflake.JsonInt64(val.Int64())
}

// GetDeptID 从 context 中获取当前登录用户部门 ID
func GetDeptID(ctx context.Context) snowflake.JsonInt64 {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0
	}
	val := req.GetCtxVar("jwt_dept_id")
	return snowflake.JsonInt64(val.Int64())
}

// GetIsAdmin 从 context 中获取当前用户是否是超级管理员
func GetIsAdmin(ctx context.Context) bool {
	userID := int64(GetUserID(ctx))
	if userID <= 0 {
		return false
	}
	var row struct {
		IsAdmin int `json:"isAdmin"`
	}
	if err := g.DB().Ctx(ctx).Raw(`
		SELECT MAX(r.is_admin) AS isAdmin
		FROM system_role r
		INNER JOIN system_user_role ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.status = 1 AND r.deleted_at IS NULL
	`, userID).Scan(&row); err != nil {
		return false
	}
	return row.IsAdmin == 1
}

// ApplyDataScope 应用数据权限过滤
// 根据当前用户的角色数据范围，自动添加 WHERE 条件
// cols 参数为可选的字段名：第一个为 created_by 字段名，第二个为 dept_id 字段名
func ApplyDataScope(ctx context.Context, m *gdb.Model, cols ...string) *gdb.Model {
	if m == nil {
		return nil
	}

	scope, err := resolveDataScope(ctx)
	if err != nil {
		return m.Where("1 = 0")
	}
	if scope.All {
		return m
	}

	var (
		createdByColumn string
		deptIDColumn    string
	)
	if len(cols) > 0 {
		createdByColumn = cols[0]
	}
	if len(cols) > 1 {
		deptIDColumn = cols[1]
	}

	builder := m.Builder()
	hasFilter := false
	if scope.IncludeSelf && scope.UserID > 0 && createdByColumn != "" {
		builder = builder.Where(createdByColumn, scope.UserID)
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

type dataAccessScope struct {
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

type dataScopeModes struct {
	AllowAll               bool
	IncludeDeptAndChildren bool
	IncludeCurrentDept     bool
	IncludeSelf            bool
	IncludeCustomDept      bool
}

func resolveDataScope(ctx context.Context) (dataAccessScope, error) {
	scope := dataAccessScope{
		UserID: int64(GetUserID(ctx)),
		DeptID: int64(GetDeptID(ctx)),
	}
	if scope.UserID <= 0 {
		scope.IncludeSelf = true
		return scope, nil
	}

	roleRows, err := loadDataScopeRoles(ctx, scope.UserID)
	if err != nil {
		return dataAccessScope{}, err
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
		loadedDeptIDs, loadErr := loadDeptAndChildrenIDs(ctx, scope.DeptID)
		if loadErr != nil {
			return dataAccessScope{}, loadErr
		}
		deptIDs = append(deptIDs, loadedDeptIDs...)
	}
	if modes.IncludeCurrentDept || modes.IncludeSelf {
		deptIDs = append(deptIDs, scope.DeptID)
	}
	if modes.IncludeCustomDept {
		loadedDeptIDs, loadErr := loadCustomScopeDeptIDs(ctx, roleRows)
		if loadErr != nil {
			return dataAccessScope{}, loadErr
		}
		deptIDs = append(deptIDs, loadedDeptIDs...)
	}
	scope.IncludeSelf = modes.IncludeSelf
	scope.DeptIDs = compactScopeIDs(deptIDs)
	if !scope.IncludeSelf && len(scope.DeptIDs) == 0 {
		scope.IncludeSelf = true
	}
	return scope, nil
}

func loadDataScopeRoles(ctx context.Context, userID int64) ([]dataScopeRoleRow, error) {
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

func hasAdminDataScope(rows []dataScopeRoleRow) bool {
	for _, row := range rows {
		if row.IsAdmin == 1 {
			return true
		}
	}
	return false
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

func loadDeptAndChildrenIDs(ctx context.Context, rootDeptID int64) ([]int64, error) {
	rootIDs := compactScopeIDs([]int64{rootDeptID})
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
	return compactScopeIDs(deptIDs), nil
}

func customScopeRoleIDs(rows []dataScopeRoleRow) []int64 {
	roleIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if normalizeDataScope(row.DataScope) == 5 {
			roleIDs = append(roleIDs, row.RoleID)
		}
	}
	return compactScopeIDs(roleIDs)
}

func normalizeDataScope(scope int) int {
	switch scope {
	case 1, 2, 3, 4, 5:
		return scope
	default:
		return 4
	}
}

func expandDeptTree(rootIDs []int64, tree map[int64][]int64) []int64 {
	rootIDs = compactScopeIDs(rootIDs)
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

func compactScopeIDs(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		ids = append(ids, value)
	}
	if len(ids) == 0 {
		return nil
	}
	return ids
}
