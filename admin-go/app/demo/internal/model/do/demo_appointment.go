// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoAppointment is the golang structure of table demo_appointment for DAO operations like Where/Data.
type DemoAppointment struct {
	g.Meta        `orm:"table:demo_appointment, do:true"`
	Id            any         // 预约ID（Snowflake）
	AppointmentNo any         // 预约编号|search:eq|priority:100
	CustomerId    any         // 客户
	Subject       any         // 预约主题|search:like|keyword:on|priority:95
	AppointmentAt *gtime.Time // 预约时间
	ContactPhone  any         // 联系电话|search:like|keyword:on|priority:90
	Address       any         // 预约地址|keyword:only
	Remark        any         // 备注|keyword:only
	Status        any         // 状态:0=待确认,1=已确认,2=已完成,3=已取消
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
