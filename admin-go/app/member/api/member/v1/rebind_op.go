package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// RebindLogRebindParentReq 后台触发换绑上级。
type RebindLogRebindParentReq struct {
	g.Meta      `path:"/rebind_log/rebind-parent" method:"post" tags:"换绑上级日志" summary:"执行换绑上级"`
	UserID      snowflake.JsonInt64 `json:"userId" v:"required#会员 ID 不能为空"`
	NewParentID snowflake.JsonInt64 `json:"newParentId" dc:"新上级 ID（0 表示置顶无上级）"`
	Reason      string              `json:"reason" v:"max-length:500" dc:"换绑原因"`
}

// RebindLogRebindParentRes 响应。
type RebindLogRebindParentRes struct {
	g.Meta `mime:"application/json"`
}
