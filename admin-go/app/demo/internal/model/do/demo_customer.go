// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoCustomer is the golang structure of table demo_customer for DAO operations like Where/Data.
type DemoCustomer struct {
	g.Meta       `orm:"table:demo_customer, do:true"`
	Id           any         // 客户ID（Snowflake）
	Avatar       any         // 头像
	Name         any         // 客户名称|search:like|keyword:on|priority:95
	CustomerNo   any         // 客户编号|search:eq|priority:100
	Phone        any         // 联系电话|search:like|keyword:on|priority:90
	Email        any         // 邮箱|search:like|keyword:on|priority:90
	Gender       any         // 性别:0=未知,1=男,2=女
	Level        any         // 等级:1=普通,2=VIP,3=付费,4=冻结
	SourceType   any         // 来源:1=官网,2=小程序,3=线下,4=导入
	IsVip        any         // 是否VIP:0=否,1=是
	RegisteredAt *gtime.Time // 注册时间
	Remark       any         // 备注|search:like|keyword:only
	Status       any         // 状态:0=禁用,1=启用
	TenantId     any         // 租户
	MerchantId   any         // 商户
	CreatedBy    any         // 创建人ID
	DeptId       any         // 所属部门ID
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 软删除时间，非 NULL 表示已删除
}
