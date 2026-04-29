// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoAuditLog is the golang structure of table demo_audit_log for DAO operations like Where/Data.
type DemoAuditLog struct {
	g.Meta      `orm:"table:demo_audit_log, do:true"`
	Id          any         // 审计日志ID（Snowflake）
	LogNo       any         // 日志编号|search:eq|priority:100
	OperatorId  any         // 操作人|ref:system_users.username
	Action      any         // 动作:1=创建,2=修改,3=删除,4=导出,5=导入
	TargetType  any         // 对象类型:1=客户,2=商品,3=订单,4=工单
	TargetCode  any         // 对象编号|search:eq|priority:88
	RequestJson any         // 请求JSON
	Result      any         // 结果:0=失败,1=成功
	ClientIp    any         // 客户端IP|search:eq|priority:80
	OccurredAt  *gtime.Time // 发生时间
	Remark      any         // 备注|keyword:only
	TenantId    any         // 租户
	MerchantId  any         // 商户
	CreatedBy   any         // 创建人ID
	DeptId      any         // 所属部门ID
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 软删除时间，非 NULL 表示已删除
}
