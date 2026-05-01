package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/response"
)

// PortalAuth 是 C 端会员的鉴权中间件，挂在 /api/member-portal 受保护路由组上。
//
// 与后台 middleware.Auth 的差异：
//   - 解析的是 jwt.MemberClaims（会员独立 secret，避免和后台 token 互通）
//   - 写入的 ctx 变量前缀为 member_，与后台的 jwt_user_id 不冲突
//   - 不做后台权限码检查、不做域名归属校验（C 端没有租户概念，会员归属在 token 内）
func PortalAuth(r *ghttp.Request) {
	tokenStr := r.GetHeader("Authorization")
	if tokenStr == "" {
		response.Unauthorized(r)
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)
	if tokenStr == "" {
		response.Unauthorized(r)
		return
	}

	claims, err := jwt.ParseMemberToken(tokenStr)
	if err != nil {
		response.Unauthorized(r, "登录已过期，请重新登录")
		return
	}

	r.SetCtxVar("member_id", claims.MemberID)
	r.SetCtxVar("member_phone", claims.Phone)
	r.SetCtxVar("member_is_coach", claims.IsCoach)
	r.SetCtxVar("member_coach_id", claims.CoachID)
	r.SetCtxVar("member_current_role", claims.CurrentRole)
	r.SetCtxVar("member_claims", claims)

	r.Middleware.Next()
}
