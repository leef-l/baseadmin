package shared

import (
	"context"
	"net"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

type DomainScope struct {
	Matched    bool
	Domain     string
	AppCode    string
	OwnerType  int
	TenantID   int64
	MerchantID int64
}

func CurrentDomainScope(ctx context.Context) DomainScope {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return DomainScope{}
	}
	scope := DomainScope{
		Matched:    req.GetCtxVar("domain_scope_matched").Bool(),
		Domain:     req.GetCtxVar("domain_scope_domain").String(),
		AppCode:    req.GetCtxVar("domain_scope_app_code").String(),
		OwnerType:  req.GetCtxVar("domain_scope_owner_type").Int(),
		TenantID:   req.GetCtxVar("domain_scope_tenant_id").Int64(),
		MerchantID: req.GetCtxVar("domain_scope_merchant_id").Int64(),
	}
	return scope
}

func CurrentRequestHost(ctx context.Context) string {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ""
	}
	host := req.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}
	return NormalizeDomainHost(host)
}

func DomainScopeAllows(ctx context.Context, tenantID, merchantID int64) bool {
	scope := CurrentDomainScope(ctx)
	if !scope.Matched {
		return true
	}
	if scope.TenantID > 0 && tenantID != scope.TenantID {
		return false
	}
	if scope.MerchantID > 0 && merchantID != scope.MerchantID {
		return false
	}
	return true
}

func NormalizeDomainHost(value string) string {
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
