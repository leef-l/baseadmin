package consts

// WalletLogWalletType 钱包类型
const (
	WalletLogWalletTypeV1 = 1 // 优惠券余额
	WalletLogWalletTypeV2 = 2 // 奖金余额
	WalletLogWalletTypeV3 = 3 // 推广奖余额
)

// WalletLogChangeType 变动类型
const (
	WalletLogChangeTypeRecharge = 1 // 充值
	WalletLogChangeTypeConsume = 2 // 消费
	WalletLogChangeTypeV3 = 3 // 推广奖
	WalletLogChangeTypeV4 = 4 // 仓库卖出收入
	WalletLogChangeTypeV5 = 5 // 平台扣除
	WalletLogChangeTypeV6 = 6 // 后台调整
)

