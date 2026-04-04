package middleware

import (
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/response"
)

// Auth JWT 鉴权中间件，只验证 token 合法性，不区分身份
func Auth(r *ghttp.Request) {
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

	if !jwt.VerifyAnyToken(tokenStr) {
		response.Unauthorized(r, "Token无效或已过期")
		return
	}

	r.Middleware.Next()
}
