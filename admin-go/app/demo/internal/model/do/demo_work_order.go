// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoWorkOrder is the golang structure of table demo_work_order for DAO operations like Where/Data.
type DemoWorkOrder struct {
	g.Meta         `orm:"table:demo_work_order, do:true"`
	Id             any         // 工单ID（Snowflake）
	TicketNo       any         // 工单号|search:eq|priority:100
	CustomerId     any         // 客户
	ProductId      any         // 商品
	OrderId        any         // 订单
	Title          any         // 工单标题|search:like|keyword:on|priority:95
	Priority       any         // 优先级:1=低,2=普通,3=高,4=紧急
	SourceType     any         // 来源:1=官网,2=电话,3=微信,4=后台
	Description    any         // 问题描述|search:like|keyword:only
	AttachmentFile any         // 附件
	DueAt          *gtime.Time // 截止时间
	Status         any         // 状态:0=待处理,1=进行中,2=已完成,3=已取消
	TenantId       any         // 租户
	MerchantId     any         // 商户
	CreatedBy      any         // 创建人ID
	DeptId         any         // 所属部门ID
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	DeletedAt      *gtime.Time // 软删除时间，非 NULL 表示已删除
}
