package middleware

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/system/internal/logic/shared"
	"gbaseadmin/utility/authz"
	"gbaseadmin/utility/jwt"
	"gbaseadmin/utility/response"
)

const denyPermission = "__deny__"

// Auth JWT 鉴权中间件
func Auth(r *ghttp.Request) {
	tokenStr := r.GetHeader("Authorization")
	if tokenStr == "" {
		response.Unauthorized(r)
		return
	}

	// 支持 Bearer token 格式
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

	// 将用户信息写入 context
	r.SetCtxVar("jwt_user_id", claims.UserID)
	r.SetCtxVar("jwt_username", claims.Username)
	r.SetCtxVar("jwt_dept_id", claims.DeptID)
	r.SetCtxVar("jwt_tenant_id", claims.TenantID)
	r.SetCtxVar("jwt_merchant_id", claims.MerchantID)
	r.SetCtxVar("jwt_claims", claims)

	if !shared.DomainScopeAllows(r.Context(), claims.TenantID, claims.MerchantID) {
		response.Forbidden(r, "当前账号不属于该访问域名")
		return
	}

	if permission := resolveSystemPermission(r.Method, r.URL.Path); permission != "" {
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

func resolveSystemPermission(method, path string) string {
	module, action := splitRouteAction("/api/system/", path)
	if module == "" {
		return ""
	}
	switch module {
	case "auth":
		if isAllowedAuthAction(action) {
			return ""
		}
		return denyPermission
	case "users":
		if resolved := resolveSystemAction(action); resolved != "" {
			return "system:user:" + resolved
		}
		return denyPermission
	case "dept", "domain", "menu", "tenant", "merchant", "daemon":
		if module == "daemon" && action == "detail" {
			return "system:daemon:view"
		}
		if resolved := resolveSystemAction(action); resolved != "" {
			return "system:" + module + ":" + resolved
		}
		if module == "domain" && action == "apply-nginx" {
			return "system:domain:apply"
		}
		if module == "domain" && action == "apply-ssl" {
			return "system:domain:ssl"
		}
		if module == "daemon" && (action == "restart" || action == "batch-restart") {
			return "system:daemon:restart"
		}
		if module == "daemon" && (action == "stop" || action == "batch-stop") {
			return "system:daemon:stop"
		}
		if module == "daemon" && action == "log" {
			return "system:daemon:view"
		}
		return denyPermission
	case "role":
		switch action {
		case "grant-menu", "menu-ids":
			return "system:role:grant:menu"
		case "grant-dept", "dept-ids":
			return "system:role:grant:dept"
		default:
			if resolved := resolveSystemAction(action); resolved != "" {
				return "system:role:" + resolved
			}
			return denyPermission
		}
	default:
		_ = method
		return denyPermission
	}
}

func isAllowedAuthAction(action string) bool {
	switch action {
	case "login", "ticket-login", "ticket", "info", "menus", "change-password":
		return true
	default:
		return false
	}
}

func resolveSystemAction(action string) string {
	switch action {
	case "create":
		return "create"
	case "update", "reset-password":
		return "update"
	case "delete":
		return "delete"
	case "detail", "list", "tree":
		return "list"
	case "batch-delete":
		return "batch-delete"
	case "batch-update":
		return "batch-update"
	case "export":
		return "export"
	case "import", "import-template":
		return "import"
	default:
		return ""
	}
}

func splitRouteAction(prefix, path string) (string, string) {
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
