package shared

import (
	"context"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type DomainScope struct {
	Matched    bool
	Domain     string
	AppCode    string
	OwnerType  int
	TenantID   int64
	MerchantID int64
}

var (
	domainStrictOnce sync.Once
	domainStrictMode bool
)

func isDomainStrictMode() bool {
	domainStrictOnce.Do(func() {
		ctx := gctx.New()
		val, _ := g.Cfg().Get(ctx, "domain.strictMode", false)
		domainStrictMode = val.Bool()
	})
	return domainStrictMode
}

var (
	domainExistsCacheMu   sync.Mutex
	domainExistsCacheTime time.Time
	domainExistsCacheVal  bool
)

func hasDomainRecords(ctx context.Context) bool {
	domainExistsCacheMu.Lock()
	defer domainExistsCacheMu.Unlock()
	if time.Since(domainExistsCacheTime) < time.Minute {
		return domainExistsCacheVal
	}
	count, err := g.DB().Ctx(ctx).Model("system_domain").
		Where("app_code", "admin").
		Where("verify_status", 1).
		Where("status", 1).
		Where("deleted_at", nil).
		Count()
	if err != nil {
		return false
	}
	domainExistsCacheVal = count > 0
	domainExistsCacheTime = time.Now()
	return domainExistsCacheVal
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
		if isDomainStrictMode() && hasDomainRecords(ctx) {
			return false
		}
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
