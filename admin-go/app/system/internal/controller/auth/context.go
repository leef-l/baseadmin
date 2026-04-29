package auth

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/jwt"
)

// GetClaims 从 context 中获取 JWT Claims
func GetClaims(ctx context.Context) *jwt.Claims {
	val := g.RequestFromCtx(ctx).GetCtxVar("jwt_claims")
	if val.IsNil() {
		return nil
	}
	if claims, ok := val.Val().(*jwt.Claims); ok {
		return claims
	}
	return nil
}

// GetUserID 从 context 中获取当前用户 ID
func GetUserID(ctx context.Context) int64 {
	claims := GetClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.UserID
}

// GetDeptID 从 context 中获取当前用户部门 ID
func GetDeptID(ctx context.Context) int64 {
	claims := GetClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.DeptID
}

// GetTenantID 从 context 中获取当前用户租户 ID
func GetTenantID(ctx context.Context) int64 {
	claims := GetClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.TenantID
}

// GetMerchantID 从 context 中获取当前用户商户 ID
func GetMerchantID(ctx context.Context) int64 {
	claims := GetClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.MerchantID
}
