package consts

// WorkOrderPriority 优先级
const (
	WorkOrderPriorityLow = 1 // 低
	WorkOrderPriorityRegular = 2 // 普通
	WorkOrderPriorityHigh = 3 // 高
	WorkOrderPriorityUrgent = 4 // 紧急
)

// WorkOrderSourceType 来源
const (
	WorkOrderSourceTypeOfficial = 1 // 官网
	WorkOrderSourceTypePhone = 2 // 电话
	WorkOrderSourceTypeWechat = 3 // 微信
	WorkOrderSourceTypeBackend = 4 // 后台
)

// WorkOrderStatus 状态
const (
	WorkOrderStatusPending = 0 // 待处理
	WorkOrderStatusInProgress = 1 // 进行中
	WorkOrderStatusDone = 2 // 已完成
	WorkOrderStatusCancelled = 3 // 已取消
)

