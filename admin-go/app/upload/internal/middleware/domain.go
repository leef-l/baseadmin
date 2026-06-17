package middleware

import (
	"database/sql"
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// DomainContext 从请求域名解析租户/商户归属，写入上下文
func DomainContext(r *ghttp.Request) {
	host := normalizeDomainHost(r.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = normalizeDomainHost(r.Host)
	}
	if host != "" {
		var row domainScopeRow
		err := g.DB().Ctx(r.Context()).Model("system_domain").
			Fields("domain", "owner_type AS ownerType", "tenant_id AS tenantId", "merchant_id AS merchantId", "app_code AS appCode").
			Where("domain", host).
			Where("verify_status", 1).
			Where("status", 1).
			Where("deleted_at", nil).
			Scan(&row)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				g.Log().Debugf(r.Context(), "DomainContext: no domain configured for host %s", host)
			} else {
				g.Log().Errorf(r.Context(), "DomainContext: query failed for host %s: %v", host, err)
			}
		}
		if row.Domain != "" {
			r.SetCtxVar("domain_scope_matched", true)
			r.SetCtxVar("domain_scope_domain", row.Domain)
			r.SetCtxVar("domain_scope_app_code", row.AppCode)
			r.SetCtxVar("domain_scope_owner_type", row.OwnerType)
			r.SetCtxVar("domain_scope_tenant_id", row.TenantID)
			r.SetCtxVar("domain_scope_merchant_id", row.MerchantID)
		}
	}
	r.Middleware.Next()
}

func normalizeDomainHost(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	if strings.Contains(value, ",") {
		value = strings.TrimSpace(strings.Split(value, ",")[0])
	}
	if strings.Contains(value, "://") {
		if parsed, err := url.Parse(value); err == nil {
			value = parsed.Host
		}
	}
	if slash := strings.Index(value, "/"); slash >= 0 {
		value = value[:slash]
	}
	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	} else if strings.Count(value, ":") == 1 {
		value = strings.Split(value, ":")[0]
	}
	return strings.Trim(value, ".")
}

type domainScopeRow struct {
	Domain     string `json:"domain"`
	OwnerType  int    `json:"ownerType"`
	TenantID   int64  `json:"tenantId"`
	MerchantID int64  `json:"merchantId"`
	AppCode    string `json:"appCode"`
}
