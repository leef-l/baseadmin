package v1

import "github.com/gogf/gf/v2/frame/g"

type CheckReq struct {
	g.Meta `path:"/healthz" tags:"Health" method:"get" summary:"Upload health check"`
}

type CheckRes struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type ReadyReq struct {
	g.Meta `path:"/readyz" tags:"Health" method:"get" summary:"Upload readiness check"`
}

type ReadyRes struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Checks  map[string]string `json:"checks"`
}
