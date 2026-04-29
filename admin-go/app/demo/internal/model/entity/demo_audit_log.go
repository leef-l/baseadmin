// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoAuditLog is the golang structure for table demo_audit_log.
type DemoAuditLog struct {
	Id          uint64      `orm:"id"           description:"审计日志ID（Snowflake）"`             // 审计日志ID（Snowflake）
	LogNo       string      `orm:"log_no"       description:"日志编号|search:eq|priority:100"`   // 日志编号|search:eq|priority:100
	OperatorId  uint64      `orm:"operator_id"  description:"操作人|ref:system_users.username"` // 操作人|ref:system_users.username
	Action      int         `orm:"action"       description:"动作:1=创建,2=修改,3=删除,4=导出,5=导入"`   // 动作:1=创建,2=修改,3=删除,4=导出,5=导入
	TargetType  int         `orm:"target_type"  description:"对象类型:1=客户,2=商品,3=订单,4=工单"`      // 对象类型:1=客户,2=商品,3=订单,4=工单
	TargetCode  string      `orm:"target_code"  description:"对象编号|search:eq|priority:88"`    // 对象编号|search:eq|priority:88
	RequestJson string      `orm:"request_json" description:"请求JSON"`                        // 请求JSON
	Result      int         `orm:"result"       description:"结果:0=失败,1=成功"`                  // 结果:0=失败,1=成功
	ClientIp    string      `orm:"client_ip"    description:"客户端IP|search:eq|priority:80"`   // 客户端IP|search:eq|priority:80
	OccurredAt  *gtime.Time `orm:"occurred_at"  description:"发生时间"`                          // 发生时间
	Remark      string      `orm:"remark"       description:"备注|keyword:only"`               // 备注|keyword:only
	TenantId    uint64      `orm:"tenant_id"    description:"租户"`                            // 租户
	MerchantId  uint64      `orm:"merchant_id"  description:"商户"`                            // 商户
	CreatedBy   uint64      `orm:"created_by"   description:"创建人ID"`                         // 创建人ID
	DeptId      uint64      `orm:"dept_id"      description:"所属部门ID"`                        // 所属部门ID
	CreatedAt   *gtime.Time `orm:"created_at"   description:"创建时间"`                          // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"   description:"更新时间"`                          // 更新时间
	DeletedAt   *gtime.Time `orm:"deleted_at"   description:"软删除时间，非 NULL 表示已删除"`            // 软删除时间，非 NULL 表示已删除
}
