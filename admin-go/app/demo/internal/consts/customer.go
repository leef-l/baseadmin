package consts

// CustomerGender 性别
const (
	CustomerGenderV0 = 0 // 未知
	CustomerGenderMale = 1 // 男
	CustomerGenderFemale = 2 // 女
)

// CustomerLevel 等级
const (
	CustomerLevelRegular = 1 // 普通
	CustomerLevelVIP = 2 // VIP
	CustomerLevelPaid = 3 // 付费
	CustomerLevelFrozen = 4 // 冻结
)

// CustomerSourceType 来源
const (
	CustomerSourceTypeV1 = 1 // 官网
	CustomerSourceTypeV2 = 2 // 小程序
	CustomerSourceTypeV3 = 3 // 线下
	CustomerSourceTypeV4 = 4 // 导入
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

