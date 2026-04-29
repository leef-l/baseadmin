package middleware

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"

	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/utility/snowflake"
)

type TenantAccessScope = shared.TenantAccessScope
type TenantScopedRow = shared.TenantScopedRow

func GetUserID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(shared.CurrentActorUserID(ctx))
}

func GetDeptID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(shared.CurrentActorDeptID(ctx))
}

func GetTenantID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(shared.CurrentActorTenantID(ctx))
}

func GetMerchantID(ctx context.Context) snowflake.JsonInt64 {
	return snowflake.JsonInt64(shared.CurrentActorMerchantID(ctx))
}

func ResolveTenantAccessScope(ctx context.Context) TenantAccessScope {
	return shared.ResolveTenantAccessScope(ctx)
}

func ApplyTenantScopeToModel(ctx context.Context, m *gdb.Model, tenantIDColumn, merchantIDColumn string) *gdb.Model {
	return shared.ApplyTenantScopeToModel(ctx, m, tenantIDColumn, merchantIDColumn)
}

func ApplyTenantScopeToWrite(ctx context.Context, tenantID, merchantID *snowflake.JsonInt64) {
	shared.ApplyTenantScopeToWrite(ctx, tenantID, merchantID)
}

func CanAccessTenantMerchant(ctx context.Context, tenantID, merchantID int64) bool {
	return shared.CanAccessTenantMerchant(ctx, tenantID, merchantID)
}

func EnsureTenantMerchantAccessible(ctx context.Context, tenantID, merchantID snowflake.JsonInt64) error {
	return shared.EnsureTenantMerchantAccessible(ctx, tenantID, merchantID)
}

func EnsureTenantScopedRowAccessible(ctx context.Context, m *gdb.Model, id snowflake.JsonInt64, idColumn, tenantIDColumn, merchantIDColumn, resourceName string) error {
	return shared.EnsureTenantScopedRowAccessible(ctx, m, id, idColumn, tenantIDColumn, merchantIDColumn, resourceName)
}

func EnsureTenantScopedRowsAccessible(ctx context.Context, m *gdb.Model, ids []snowflake.JsonInt64, idColumn, tenantIDColumn, merchantIDColumn, resourceName string) error {
	return shared.EnsureTenantScopedRowsAccessible(ctx, m, ids, idColumn, tenantIDColumn, merchantIDColumn, resourceName)
}

func ApplyDataScope(ctx context.Context, m *gdb.Model, cols ...string) *gdb.Model {
	if m == nil {
		return nil
	}
	scope, err := shared.ResolveDataAccessScope(ctx)
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
