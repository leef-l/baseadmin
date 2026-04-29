package consts

// CustomerGender 性别
const (
	CustomerGenderUnknown = 0 // 未知
	CustomerGenderMale = 1 // 男
	CustomerGenderFemale = 2 // 女
)

// CustomerLevel 等级
const (
	CustomerLevelRegular = 1 // 普通
	CustomerLevelVIP = 2 // VIP
	CustomerLevelCharged = 3 // 付费
	CustomerLevelFrozen = 4 // 冻结
)

// CustomerSourceType 来源
const (
	CustomerSourceTypeOfficial = 1 // 官网
	CustomerSourceTypeMiniApp = 2 // 小程序
	CustomerSourceTypeOffline = 3 // 线下
	CustomerSourceTypeImport = 4 // 导入
)

// CustomerIsVip 是否VIP
const (
	CustomerIsVipNo = 0 // 否
	CustomerIsVipYes = 1 // 是
)

// CustomerStatus 状态
const (
	CustomerStatusDisabled = 0 // 禁用
	CustomerStatusEnabled = 1 // 启用
)

