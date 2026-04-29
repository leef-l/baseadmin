package consts

// CampaignType 活动类型
const (
	CampaignTypeFree = 1 // 免费
	CampaignTypeCharged = 2 // 付费
	CampaignTypePublic = 3 // 公开
	CampaignTypePrivate = 4 // 私密
)

// CampaignChannel 投放渠道
const (
	CampaignChannelOfficial = 1 // 官网
	CampaignChannelMiniApp = 2 // 小程序
	CampaignChannelSMS = 3 // 短信
	CampaignChannelOffline = 4 // 线下
)

// CampaignIsPublic 是否公开
const (
	CampaignIsPublicNo = 0 // 否
	CampaignIsPublicYes = 1 // 是
)

// CampaignStatus 状态
const (
	CampaignStatusDraft = 0 // 草稿
	CampaignStatusPublished = 1 // 已发布
	CampaignStatusOffline = 2 // 已下架
)

