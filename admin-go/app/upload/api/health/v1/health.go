package v1

import "github.com/gogf/gf/v2/frame/g"

type CheckReq struct {
	g.Meta `path:"/healthz" tags:"Health" method:"get" summary:"Upload health check"`
}

type CheckRes struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
