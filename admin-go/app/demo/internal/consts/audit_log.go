package consts

// AuditLogAction 动作
const (
	AuditLogActionV1 = 1 // 创建
	AuditLogActionV2 = 2 // 修改
	AuditLogActionV3 = 3 // 删除
	AuditLogActionV4 = 4 // 导出
	AuditLogActionV5 = 5 // 导入
)

// AuditLogTargetType 对象类型
const (
	AuditLogTargetTypeV1 = 1 // 客户
	AuditLogTargetTypeV2 = 2 // 商品
	AuditLogTargetTypeV3 = 3 // 订单
	AuditLogTargetTypeV4 = 4 // 工单
)

// AuditLogResult 结果
const (
	AuditLogResultFailed = 0 // 失败
	AuditLogResultSuccess = 1 // 成功
)

