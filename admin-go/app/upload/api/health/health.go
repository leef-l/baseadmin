package health

import (
	"context"

	v1 "gbaseadmin/app/upload/api/health/v1"
)

type IHealthV1 interface {
	Check(ctx context.Context, req *v1.CheckReq) (res *v1.CheckRes, err error)
	Ready(ctx context.Context, req *v1.ReadyReq) (res *v1.ReadyRes, err error)
}
