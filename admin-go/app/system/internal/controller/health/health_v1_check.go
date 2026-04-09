package health

import (
	"context"

	v1 "gbaseadmin/app/system/api/health/v1"
)

func (c *ControllerV1) Check(ctx context.Context, req *v1.CheckReq) (res *v1.CheckRes, err error) {
	return &v1.CheckRes{
		Status:  "ok",
		Service: "system",
	}, nil
}
