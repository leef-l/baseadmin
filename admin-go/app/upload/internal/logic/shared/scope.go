package shared

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

const (
	ColumnTenantID   = "tenant_id"
	ColumnMerchantID = "merchant_id"
	ColumnCreatedBy  = "created_by"
	ColumnDeptID     = "dept_id"
)

type TenantAccessScope struct {
	TenantID   int64
	MerchantID int64
	All        bool
}

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

func CurrentActorUserID(ctx context.Context) snowflake.JsonInt64 {
	userID, _ := currentActor(ctx)
	return snowflake.JsonInt64(userID)
}

func CurrentActorDeptID(ctx context.Context) snowflake.JsonInt64 {
	_, deptID := currentActor(ctx)
	return snowflake.JsonInt64(deptID)
}

func CurrentActorTenantID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(ResolveTenantAccessScope(ctx).TenantID)
}

func CurrentActorMerchantID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(ResolveTenantAccessScope(ctx).MerchantID)
}

func ApplyWriteScope(ctx context.Context, tenantID, merchantID, createdBy, deptID *snowflake.JsonInt64) {
	scope := ResolveTenantAccessScope(ctx)
	if tenantID != nil {
		*tenantID = snowflake.JsonInt64(scope.TenantID)
	}
	if merchantID != nil {
		*merchantID = snowflake.JsonInt64(scope.MerchantID)
	}
	if createdBy != nil {
		*createdBy = CurrentActorUserID(ctx)
	}
	if deptID != nil {
		*deptID = CurrentActorDeptID(ctx)
	}
}

func ResolveTenantAccessScope(ctx context.Context) TenantAccessScope {
	if ctx == nil {
		return TenantAccessScope{All: true}
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return TenantAccessScope{All: true}
	}
	tenantID := req.GetCtxVar("jwt_tenant_id").Int64()
	merchantID := req.GetCtxVar("jwt_merchant_id").Int64()
	return TenantAccessScope{
		TenantID:   tenantID,
		MerchantID: merchantID,
		All:        tenantID <= 0,
	}
}

func ApplyAccessScope(ctx context.Context, m *gdb.Model) *gdb.Model {
	m = ApplyTenantScopeToModel(ctx, m, ColumnTenantID, ColumnMerchantID)
	return ApplyDataScope(ctx, m, ColumnCreatedBy, ColumnDeptID)
}

func ApplyTenantScopeToModel(ctx context.Context, m *gdb.Model, tenantIDColumn, merchantIDColumn string) *gdb.Model {
	if m == nil {
		return nil
	}
	scope := ResolveTenantAccessScope(ctx)
	if scope.All {
		return m
	}
	if tenantIDColumn == "" {
		return m.Where("1 = 0")
	}
	m = m.Where(tenantIDColumn, scope.TenantID)
	if scope.MerchantID > 0 {
		if merchantIDColumn == "" {
			return m.Where("1 = 0")
		}
		m = m.Where(merchantIDColumn, scope.MerchantID)
	}
	return m
}

func ApplyDataScope(ctx context.Context, m *gdb.Model, createdByColumn, deptIDColumn string) *gdb.Model {
	if m == nil {
		return nil
	}
	scope, err := ResolveDataAccessScope(ctx)
	if err != nil {
		return m.Where("1 = 0")
	}
	if scope.All {
		return m
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

func ResolveDataAccessScope(ctx context.Context) (DataAccessScope, error) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return DataAccessScope{All: true}, nil
	}
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
	if !scope.IncludeSelf && len(scope.DeptIDs) == 0 {
		scope.IncludeSelf = true
	}
	return scope, nil
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
	return rows, err
}

func loadDeptAndChildrenIDs(ctx context.Context, rootDeptID int64) ([]int64, error) {
	rootIDs := compactPositiveIDs([]int64{rootDeptID})
	if len(rootIDs) == 0 {
		return nil, nil
	}
	var rows []dataScopeDeptRow
	m := g.DB().Ctx(ctx).
		Model("system_dept").
		Fields("id", "parent_id AS parentId").
		Where("deleted_at", nil)
	tenantScope := ResolveTenantAccessScope(ctx)
	if !tenantScope.All {
		m = m.Where(ColumnTenantID, tenantScope.TenantID)
		if tenantScope.MerchantID > 0 {
			m = m.WhereIn(ColumnMerchantID, []int64{0, tenantScope.MerchantID})
		}
	}
	if err := m.Scan(&rows); err != nil {
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
	err := g.DB().Ctx(ctx).
		Model("system_role_dept").
		Fields("dept_id AS deptId").
		WhereIn("role_id", roleIDs).
		Scan(&deptRows)
	if err != nil {
		return nil, err
	}
	deptIDs := make([]int64, 0, len(deptRows))
	for _, row := range deptRows {
		deptIDs = append(deptIDs, row.DeptID)
	}
	return compactPositiveIDs(deptIDs), nil
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
	modes := dataScopeModes{}
	if len(rows) == 0 {
		modes.IncludeSelf = true
		return modes
	}
	for _, row := range rows {
		switch row.DataScope {
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
	return modes
}

func customScopeRoleIDs(rows []dataScopeRoleRow) []int64 {
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row.DataScope == 5 {
			ids = append(ids, row.RoleID)
		}
	}
	return compactPositiveIDs(ids)
}

func compactPositiveIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func expandDeptTree(rootIDs []int64, children map[int64][]int64) []int64 {
	rootIDs = compactPositiveIDs(rootIDs)
	if len(rootIDs) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(rootIDs))
	queue := append([]int64{}, rootIDs...)
	out := make([]int64, 0, len(rootIDs))
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		if _, ok := seen[id]; ok && len(out) > 0 {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
		for _, childID := range children[id] {
			if _, ok := seen[childID]; ok {
				continue
			}
			queue = append(queue, childID)
		}
	}
	return compactPositiveIDs(out)
}
