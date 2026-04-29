package consts

// OrderPayStatus 支付状态
const (
	OrderPayStatusUnpaid = 0 // 待支付
	OrderPayStatusPaid = 1 // 已支付
	OrderPayStatusRefunded = 2 // 已退款
)

// OrderDeliverStatus 发货状态
const (
	OrderDeliverStatusUnshipped = 0 // 待发货
	OrderDeliverStatusShipped = 1 // 已发货
	OrderDeliverStatusReceived = 2 // 已签收
)

// OrderStatus 状态
const (
	OrderStatusUnconfirmed = 0 // 待确认
	OrderStatusConfirmed = 1 // 已确认
	OrderStatusCancelled = 2 // 已取消
)

