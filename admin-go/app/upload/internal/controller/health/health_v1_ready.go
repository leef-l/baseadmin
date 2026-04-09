package health

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "gbaseadmin/app/upload/api/health/v1"
	"gbaseadmin/app/upload/internal/logic/uploader"
	"gbaseadmin/utility/readyutil"
)

func (c *ControllerV1) Ready(ctx context.Context, req *v1.ReadyReq) (res *v1.ReadyRes, err error) {
	res = &v1.ReadyRes{
		Status:  "ok",
		Service: "upload",
		Checks: map[string]string{
			"db":      "ok",
			"storage": "ok",
		},
	}
	if err := readyutil.CheckDatabase(); err != nil {
		res.Status = "not_ready"
		res.Checks["db"] = err.Error()
	}
	if err := uploader.CheckStorageReady(ctx); err != nil {
		res.Status = "not_ready"
		res.Checks["storage"] = err.Error()
	}
	if res.Status != "ok" {
		g.RequestFromCtx(ctx).Response.WriteStatus(503)
	}
	return res, nil
}
