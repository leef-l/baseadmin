package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"gbaseadmin/app/system/internal/dao"
	"gbaseadmin/app/system/internal/logic/shared"
)

const domainAppAdmin = "admin"

type domainScopeRow struct {
	Domain     string `json:"domain"`
	OwnerType  int    `json:"ownerType"`
	TenantID   int64  `json:"tenantId"`
	MerchantID int64  `json:"merchantId"`
	AppCode    string `json:"appCode"`
}

func DomainContext(r *ghttp.Request) {
	host := shared.NormalizeDomainHost(r.GetHeader("X-Forwarded-Host"))
	if host == "" {
		host = shared.NormalizeDomainHost(r.Host)
	}
	if host != "" {
		var row domainScopeRow
		_ = dao.Domain.Ctx(r.Context()).
			Fields(
				dao.Domain.Columns().Domain,
				dao.Domain.Columns().OwnerType+" AS ownerType",
				dao.Domain.Columns().TenantId+" AS tenantId",
				dao.Domain.Columns().MerchantId+" AS merchantId",
				dao.Domain.Columns().AppCode+" AS appCode",
			).
			Where(dao.Domain.Columns().Domain, host).
			Where(dao.Domain.Columns().AppCode, domainAppAdmin).
			Where(dao.Domain.Columns().VerifyStatus, 1).
			Where(dao.Domain.Columns().Status, 1).
			Where(dao.Domain.Columns().DeletedAt, nil).
			Scan(&row)
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
