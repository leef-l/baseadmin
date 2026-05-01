package middleware

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// CurrentMemberID 从 context 中获取当前登录会员 ID（C 端 portal 路由）。
// 调用前确保已经过 PortalAuth 中间件，否则返回 0。
func CurrentMemberID(ctx context.Context) snowflake.JsonInt64 {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return 0
	}
	return snowflake.JsonInt64(req.GetCtxVar("member_id").Int64())
}

// CurrentMemberPhone 返回当前登录会员的手机号。
func CurrentMemberPhone(ctx context.Context) string {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ""
	}
	return req.GetCtxVar("member_phone").String()
}
