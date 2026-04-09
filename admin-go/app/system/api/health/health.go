package health

import (
	"context"

	v1 "gbaseadmin/app/system/api/health/v1"
)

type IHealthV1 interface {
	Check(ctx context.Context, req *v1.CheckReq) (res *v1.CheckRes, err error)
}
