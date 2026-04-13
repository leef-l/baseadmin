package middleware

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/utility/authz"
	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/response"
)

const denyPermission = "__deny__"

// Auth JWT 鉴权中间件，仅接受管理端 token，并校验 upload 管理权限。
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

	claims, err := jwt.ParseToken(tokenStr)
	if err != nil {
		response.Unauthorized(r, "Token无效或已过期")
		return
	}

	r.SetCtxVar("jwt_user_id", claims.UserID)
	r.SetCtxVar("jwt_username", claims.Username)
	r.SetCtxVar("jwt_dept_id", claims.DeptID)
	r.SetCtxVar("jwt_claims", claims)

	if permission := resolveUploadPermission(r.Method, r.URL.Path); permission != "" {
		if permission == denyPermission {
			response.Forbidden(r, "未配置的权限动作")
			return
		}
		allowed, err := authz.HasPermission(r.Context(), claims.UserID, permission)
		if err != nil {
			response.Forbidden(r, "权限校验失败")
			return
		}
		if !allowed {
			response.Forbidden(r, fmt.Sprintf("缺少权限: %s", permission))
			return
		}
	}

	r.Middleware.Next()
}

func resolveUploadPermission(method, path string) string {
	module, action := splitUploadRouteAction("/api/upload/", path)
	if module == "" {
		return ""
	}
	module = normalizeUploadModule(module)
	if module == "" {
		_ = method
		return denyPermission
	}
	if resolved := resolveUploadAction(action); resolved != "" {
		return "upload:" + module + ":" + resolved
	}
	return denyPermission
}

func resolveUploadAction(action string) string {
	switch action {
	case "create":
		return "create"
	case "update":
		return "update"
	case "delete":
		return "delete"
	case "detail", "list", "tree":
		return "list"
	case "batch-delete":
		return "batch-delete"
	case "upload":
		return "create"
	default:
		return ""
	}
}

func normalizeUploadModule(module string) string {
	switch module {
	case "config", "dir", "dir_rule", "file":
		return module
	default:
		return ""
	}
}

func splitUploadRouteAction(prefix, path string) (string, string) {
	path = strings.TrimSpace(path)
	if !strings.HasPrefix(path, prefix) {
		return "", ""
	}
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return "", ""
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
