package consts

// WorkOrderPriority 优先级
const (
	WorkOrderPriorityV1 = 1 // 低
	WorkOrderPriorityRegular = 2 // 普通
	WorkOrderPriorityV3 = 3 // 高
	WorkOrderPriorityV4 = 4 // 紧急
)

// WorkOrderSourceType 来源
const (
	WorkOrderSourceTypeV1 = 1 // 官网
	WorkOrderSourceTypeV2 = 2 // 电话
	WorkOrderSourceTypeV3 = 3 // 微信
	WorkOrderSourceTypeV4 = 4 // 后台
)

// WorkOrderStatus 状态
const (
	WorkOrderStatusPending = 0 // 待处理
	WorkOrderStatusInProgress = 1 // 进行中
	WorkOrderStatusDone = 2 // 已完成
	WorkOrderStatusCancelled = 3 // 已取消
)

