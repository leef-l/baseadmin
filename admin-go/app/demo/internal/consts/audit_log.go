package consts

// AuditLogAction 动作
const (
	AuditLogActionCreate = 1 // 创建
	AuditLogActionUpdate = 2 // 修改
	AuditLogActionDelete = 3 // 删除
	AuditLogActionExport = 4 // 导出
	AuditLogActionImport = 5 // 导入
)

// AuditLogTargetType 对象类型
const (
	AuditLogTargetTypeCustomer = 1 // 客户
	AuditLogTargetTypeProduct = 2 // 商品
	AuditLogTargetTypeOrder = 3 // 订单
	AuditLogTargetTypeWorkOrder = 4 // 工单
)

// AuditLogResult 结果
const (
	AuditLogResultFailed = 0 // 失败
	AuditLogResultSuccess = 1 // 成功
)

