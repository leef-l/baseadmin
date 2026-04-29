// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCustomer is the golang structure for table demo_customer.
type DemoCustomer struct {
	Id           uint64      `orm:"id"            description:"客户ID（Snowflake）"`                         // 客户ID（Snowflake）
	Avatar       string      `orm:"avatar"        description:"头像"`                                      // 头像
	Name         string      `orm:"name"          description:"客户名称|search:like|keyword:on|priority:95"` // 客户名称|search:like|keyword:on|priority:95
	CustomerNo   string      `orm:"customer_no"   description:"客户编号|search:eq|priority:100"`             // 客户编号|search:eq|priority:100
	Phone        string      `orm:"phone"         description:"联系电话|search:like|keyword:on|priority:90"` // 联系电话|search:like|keyword:on|priority:90
	Email        string      `orm:"email"         description:"邮箱|search:like|keyword:on|priority:90"`   // 邮箱|search:like|keyword:on|priority:90
	Gender       int         `orm:"gender"        description:"性别:0=未知,1=男,2=女"`                         // 性别:0=未知,1=男,2=女
	Level        int         `orm:"level"         description:"等级:1=普通,2=VIP,3=付费,4=冻结"`                 // 等级:1=普通,2=VIP,3=付费,4=冻结
	SourceType   int         `orm:"source_type"   description:"来源:1=官网,2=小程序,3=线下,4=导入"`                 // 来源:1=官网,2=小程序,3=线下,4=导入
	IsVip        int         `orm:"is_vip"        description:"是否VIP:0=否,1=是"`                           // 是否VIP:0=否,1=是
	RegisteredAt *gtime.Time `orm:"registered_at" description:"注册时间"`                                    // 注册时间
	Remark       string      `orm:"remark"        description:"备注|search:like|keyword:only"`             // 备注|search:like|keyword:only
	Status       int         `orm:"status"        description:"状态:0=禁用,1=启用"`                            // 状态:0=禁用,1=启用
	TenantId     uint64      `orm:"tenant_id"     description:"租户"`                                      // 租户
	MerchantId   uint64      `orm:"merchant_id"   description:"商户"`                                      // 商户
	CreatedBy    uint64      `orm:"created_by"    description:"创建人ID"`                                   // 创建人ID
	DeptId       uint64      `orm:"dept_id"       description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt    *gtime.Time `orm:"created_at"    description:"创建时间"`                                    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"    description:"更新时间"`                                    // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"    description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
