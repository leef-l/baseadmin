// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoWorkOrder is the golang structure for table demo_work_order.
type DemoWorkOrder struct {
	Id             uint64      `orm:"id"              description:"工单ID（Snowflake）"`                         // 工单ID（Snowflake）
	TicketNo       string      `orm:"ticket_no"       description:"工单号|search:eq|priority:100"`              // 工单号|search:eq|priority:100
	CustomerId     uint64      `orm:"customer_id"     description:"客户"`                                      // 客户
	ProductId      uint64      `orm:"product_id"      description:"商品"`                                      // 商品
	OrderId        uint64      `orm:"order_id"        description:"订单"`                                      // 订单
	Title          string      `orm:"title"           description:"工单标题|search:like|keyword:on|priority:95"` // 工单标题|search:like|keyword:on|priority:95
	Priority       int         `orm:"priority"        description:"优先级:1=低,2=普通,3=高,4=紧急"`                   // 优先级:1=低,2=普通,3=高,4=紧急
	SourceType     int         `orm:"source_type"     description:"来源:1=官网,2=电话,3=微信,4=后台"`                  // 来源:1=官网,2=电话,3=微信,4=后台
	Description    string      `orm:"description"     description:"问题描述|search:like|keyword:only"`           // 问题描述|search:like|keyword:only
	AttachmentFile string      `orm:"attachment_file" description:"附件"`                                      // 附件
	DueAt          *gtime.Time `orm:"due_at"          description:"截止时间"`                                    // 截止时间
	Status         int         `orm:"status"          description:"状态:0=待处理,1=进行中,2=已完成,3=已取消"`              // 状态:0=待处理,1=进行中,2=已完成,3=已取消
	TenantId       uint64      `orm:"tenant_id"       description:"租户"`                                      // 租户
	MerchantId     uint64      `orm:"merchant_id"     description:"商户"`                                      // 商户
	CreatedBy      uint64      `orm:"created_by"      description:"创建人ID"`                                   // 创建人ID
	DeptId         uint64      `orm:"dept_id"         description:"所属部门ID"`                                  // 所属部门ID
	CreatedAt      *gtime.Time `orm:"created_at"      description:"创建时间"`                                    // 创建时间
	UpdatedAt      *gtime.Time `orm:"updated_at"      description:"更新时间"`                                    // 更新时间
	DeletedAt      *gtime.Time `orm:"deleted_at"      description:"软删除时间，非 NULL 表示已删除"`                      // 软删除时间，非 NULL 表示已删除
}
