// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DemoOrder is the golang structure of table demo_order for DAO operations like Where/Data.
type DemoOrder struct {
	g.Meta        `orm:"table:demo_order, do:true"`
	Id            any         // 订单ID（Snowflake）
	OrderNo       any         // 订单号|search:eq|priority:100
	CustomerId    any         // 客户
	ProductId     any         // 商品
	Quantity      any         // 购买数量
	Amount        any         // 订单金额（分）
	PayStatus     any         // 支付状态:0=待支付,1=已支付,2=已退款
	DeliverStatus any         // 发货状态:0=待发货,1=已发货,2=已签收
	PaidAt        *gtime.Time // 支付时间
	DeliverAt     *gtime.Time // 发货时间
	ReceiverPhone any         // 收货电话
	Address       any         // 收货地址|keyword:only
	Remark        any         // 备注|keyword:only
	Status        any         // 状态:0=待确认,1=已确认,2=已取消
	TenantId      any         // 租户
	MerchantId    any         // 商户
	CreatedBy     any         // 创建人ID
	DeptId        any         // 所属部门ID
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 软删除时间，非 NULL 表示已删除
}
