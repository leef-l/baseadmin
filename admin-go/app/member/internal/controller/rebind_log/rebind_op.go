package rebind_log

import (
	"context"

	v1 "gbaseadmin/app/member/api/member/v1"
	"gbaseadmin/app/member/internal/logic/rebind_log"
	"gbaseadmin/app/member/internal/middleware"
)

// RebindParent 后台执行换绑上级。
//
// 路由权限通过现有后台 system:role 的菜单权限控制，调用前提是已登录后台账号。
func (c *cRebindLog) RebindParent(ctx context.Context, req *v1.RebindLogRebindParentReq) (res *v1.RebindLogRebindParentRes, err error) {
	operatorID := int64(middleware.GetUserID(ctx))
	if err = rebind_log.RebindParent(ctx, &rebind_log.RebindParentInput{
		UserID:      int64(req.UserID),
		NewParentID: int64(req.NewParentID),
		Reason:      req.Reason,
		OperatorID:  operatorID,
	}); err != nil {
		return nil, err
	}
	return &v1.RebindLogRebindParentRes{}, nil
}
