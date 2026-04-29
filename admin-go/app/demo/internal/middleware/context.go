package middleware

import (
	"context"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
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

type TenantAccessScope struct {
	TenantID   int64
	MerchantID int64
	All        bool
}

type TenantScopedRow struct {
	ID         int64 `json:"id"`
	TenantID   int64 `json:"tenantId"`
	MerchantID int64 `json:"merchantId"`
}

// GetTenantID 从 context 中获取当前登录用户租户 ID
func GetTenantID(ctx context.Context) snowflake.JsonInt64 {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0
	}
	val := req.GetCtxVar("jwt_tenant_id")
	return snowflake.JsonInt64(val.Int64())
}

// GetMerchantID 从 context 中获取当前登录用户商户 ID
func GetMerchantID(ctx context.Context) snowflake.JsonInt64 {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0
	}
	val := req.GetCtxVar("jwt_merchant_id")
	return snowflake.JsonInt64(val.Int64())
}

func ResolveTenantAccessScope(ctx context.Context) TenantAccessScope {
	tenantID := int64(GetTenantID(ctx))
	merchantID := int64(GetMerchantID(ctx))
	return TenantAccessScope{
		TenantID:   tenantID,
		MerchantID: merchantID,
		All:        tenantID <= 0,
	}
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

func ApplyTenantScopeToWrite(ctx context.Context, tenantID, merchantID *snowflake.JsonInt64) {
	scope := ResolveTenantAccessScope(ctx)
	if scope.All {
		return
	}
	if tenantID != nil {
		*tenantID = snowflake.JsonInt64(scope.TenantID)
	}
	if merchantID != nil && scope.MerchantID > 0 {
		*merchantID = snowflake.JsonInt64(scope.MerchantID)
	}
}

func CanAccessTenantMerchant(ctx context.Context, tenantID, merchantID int64) bool {
	scope := ResolveTenantAccessScope(ctx)
	if scope.All {
		return true
	}
	if tenantID != scope.TenantID {
		return false
	}
	if scope.MerchantID > 0 && merchantID != scope.MerchantID {
		return false
	}
	return true
}

func EnsureTenantMerchantAccessible(ctx context.Context, tenantID, merchantID snowflake.JsonInt64) error {
	targetTenantID := int64(tenantID)
	targetMerchantID := int64(merchantID)
	if targetTenantID <= 0 {
		if targetMerchantID > 0 {
			return gerror.New("商户必须归属于租户")
		}
		if !CanAccessTenantMerchant(ctx, targetTenantID, targetMerchantID) {
			return gerror.New("无权操作平台级数据")
		}
		return nil
	}
	if !CanAccessTenantMerchant(ctx, targetTenantID, targetMerchantID) {
		return gerror.New("无权操作该租户或商户数据")
	}
	if err := ensureTenantExists(ctx, targetTenantID); err != nil {
		return err
	}
	if targetMerchantID > 0 {
		if err := ensureMerchantBelongsToTenant(ctx, targetTenantID, targetMerchantID); err != nil {
			return err
		}
	}
	return nil
}

func EnsureTenantScopedRowAccessible(ctx context.Context, m *gdb.Model, id snowflake.JsonInt64, idColumn, tenantIDColumn, merchantIDColumn, resourceName string) error {
	return EnsureTenantScopedRowsAccessible(ctx, m, []snowflake.JsonInt64{id}, idColumn, tenantIDColumn, merchantIDColumn, resourceName)
}

func EnsureTenantScopedRowsAccessible(ctx context.Context, m *gdb.Model, ids []snowflake.JsonInt64, idColumn, tenantIDColumn, merchantIDColumn, resourceName string) error {
	if m == nil {
		return gerror.New("数据模型不能为空")
	}
	if idColumn == "" || tenantIDColumn == "" {
		return gerror.New("租户数据权限字段未配置")
	}
	ids = compactTenantScopedIDs(ids)
	if len(ids) == 0 {
		return nil
	}
	resourceName = normalizeTenantScopedResourceName(resourceName)
	fields := []interface{}{
		idColumn + " AS id",
		tenantIDColumn + " AS tenantId",
	}
	if merchantIDColumn != "" {
		fields = append(fields, merchantIDColumn+" AS merchantId")
	}
	var rows []TenantScopedRow
	if err := m.Fields(fields...).WhereIn(idColumn, ids).Scan(&rows); err != nil {
		return err
	}
	if len(rows) != len(ids) {
		return gerror.New(resourceName + "不存在或已删除")
	}
	for _, row := range rows {
		if !CanAccessTenantMerchant(ctx, row.TenantID, row.MerchantID) {
			return gerror.New("无权操作该" + resourceName)
		}
	}
	return nil
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

func ensureTenantExists(ctx context.Context, tenantID int64) error {
	count, err := g.DB().Ctx(ctx).
		Model("system_tenant").
		Where("id", tenantID).
		Where("deleted_at", nil).
		Where("status", 1).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("租户不存在或已禁用")
	}
	return nil
}

func ensureMerchantBelongsToTenant(ctx context.Context, tenantID, merchantID int64) error {
	count, err := g.DB().Ctx(ctx).
		Model("system_merchant").
		Where("id", merchantID).
		Where("tenant_id", tenantID).
		Where("deleted_at", nil).
		Where("status", 1).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.New("商户不存在、已禁用或不属于该租户")
	}
	return nil
}

func compactTenantScopedIDs(ids []snowflake.JsonInt64) []snowflake.JsonInt64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]snowflake.JsonInt64, 0, len(ids))
	for _, id := range ids {
		value := int64(id)
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, id)
	}
	return out
}

func normalizeTenantScopedResourceName(resourceName string) string {
	if resourceName == "" {
		return "数据"
	}
	return resourceName
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
