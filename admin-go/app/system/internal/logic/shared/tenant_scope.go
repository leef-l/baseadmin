package shared

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

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

func CurrentActorTenantID(ctx context.Context) int64 {
	scope := ResolveTenantAccessScope(ctx)
	return scope.TenantID
}

func CurrentActorMerchantID(ctx context.Context) int64 {
	scope := ResolveTenantAccessScope(ctx)
	return scope.MerchantID
}

func ResolveTenantAccessScope(ctx context.Context) TenantAccessScope {
	if ctx == nil {
		return TenantAccessScope{}
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return TenantAccessScope{}
	}
	tenantID := req.GetCtxVar("jwt_tenant_id").Int64()
	merchantID := req.GetCtxVar("jwt_merchant_id").Int64()
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
	if err := m.Fields(fields...).WhereIn(idColumn, ids).Where("deleted_at", nil).Scan(&rows); err != nil {
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
