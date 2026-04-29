// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoOrder is the golang structure for table demo_order.
type DemoOrder struct {
	Id            uint64      `orm:"id"             description:"订单ID（Snowflake）"`            // 订单ID（Snowflake）
	OrderNo       string      `orm:"order_no"       description:"订单号|search:eq|priority:100"` // 订单号|search:eq|priority:100
	CustomerId    uint64      `orm:"customer_id"    description:"客户"`                         // 客户
	ProductId     uint64      `orm:"product_id"     description:"商品"`                         // 商品
	Quantity      int         `orm:"quantity"       description:"购买数量"`                       // 购买数量
	Amount        int         `orm:"amount"         description:"订单金额（分）"`                    // 订单金额（分）
	PayStatus     int         `orm:"pay_status"     description:"支付状态:0=待支付,1=已支付,2=已退款"`     // 支付状态:0=待支付,1=已支付,2=已退款
	DeliverStatus int         `orm:"deliver_status" description:"发货状态:0=待发货,1=已发货,2=已签收"`     // 发货状态:0=待发货,1=已发货,2=已签收
	PaidAt        *gtime.Time `orm:"paid_at"        description:"支付时间"`                       // 支付时间
	DeliverAt     *gtime.Time `orm:"deliver_at"     description:"发货时间"`                       // 发货时间
	ReceiverPhone string      `orm:"receiver_phone" description:"收货电话"`                       // 收货电话
	Address       string      `orm:"address"        description:"收货地址|keyword:only"`          // 收货地址|keyword:only
	Remark        string      `orm:"remark"         description:"备注|keyword:only"`            // 备注|keyword:only
	Status        int         `orm:"status"         description:"状态:0=待确认,1=已确认,2=已取消"`       // 状态:0=待确认,1=已确认,2=已取消
	TenantId      uint64      `orm:"tenant_id"      description:"租户"`                         // 租户
	MerchantId    uint64      `orm:"merchant_id"    description:"商户"`                         // 商户
	CreatedBy     uint64      `orm:"created_by"     description:"创建人ID"`                      // 创建人ID
	DeptId        uint64      `orm:"dept_id"        description:"所属部门ID"`                     // 所属部门ID
	CreatedAt     *gtime.Time `orm:"created_at"     description:"创建时间"`                       // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     description:"更新时间"`                       // 更新时间
	DeletedAt     *gtime.Time `orm:"deleted_at"     description:"软删除时间，非 NULL 表示已删除"`         // 软删除时间，非 NULL 表示已删除
}
