// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoAppointment is the golang structure for table demo_appointment.
type DemoAppointment struct {
	Id            uint64      `orm:"id"             description:"预约ID（Snowflake）"`                         // 预约ID（Snowflake）
	AppointmentNo string      `orm:"appointment_no" description:"预约编号|search:eq|priority:100"`             // 预约编号|search:eq|priority:100
	CustomerId    uint64      `orm:"customer_id"    description:"客户"`                                      // 客户
	Subject       string      `orm:"subject"        description:"预约主题|search:like|keyword:on|priority:95"` // 预约主题|search:like|keyword:on|priority:95
	AppointmentAt *gtime.Time `orm:"appointment_at" description:"预约时间"`                                    // 预约时间
	ContactPhone  string      `orm:"contact_phone"  description:"联系电话|search:like|keyword:on|priority:90"` // 联系电话|search:like|keyword:on|priority:90
	Address       string      `orm:"address"        description:"预约地址|keyword:only"`                       // 预约地址|keyword:only
	Remark        string      `orm:"remark"         description:"备注|keyword:only"`                         // 备注|keyword:only
	Status        int         `orm:"status"         description:"状态:0=待确认,1=已确认,2=已完成,3=已取消"`              // 状态:0=待确认,1=已确认,2=已完成,3=已取消
	TenantId      uint64      `orm:"tenant_id"      description:"租户"`                                      // 租户
	MerchantId    uint64      `orm:"merchant_id"    description:"商户"`                                      // 商户
	CreatedBy     uint64      `orm:"created_by"     description:"创建人ID"`                                   // 创建人ID
	DeptId        uint64      `orm:"dept_id"        description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"     description:"创建时间"`                                    // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     description:"更新时间"`                                    // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"     description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
