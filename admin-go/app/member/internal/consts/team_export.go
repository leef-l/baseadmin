package consts

// TeamExportExportType 导出类型
const (
	TeamExportExportTypeV1 = 1 // 手动导出
	TeamExportExportTypeV2 = 2 // 自动升级导出
)

// TeamExportDeployStatus 部署状态
const (
	TeamExportDeployStatusV0 = 0 // 未部署
	TeamExportDeployStatusV1 = 1 // 部署中
	TeamExportDeployStatusV2 = 2 // 已部署
	TeamExportDeployStatusV3 = 3 // 部署失败
)

// TeamExportStatus 状态
const (
	TeamExportStatusOff = 0 // 关闭
	TeamExportStatusOn = 1 // 开启
)

