package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// TeamExportRunReq 后台触发团队数据导出。
type TeamExportRunReq struct {
	g.Meta `path:"/team_export/run" method:"post" tags:"团队数据导出" summary:"导出指定会员团队数据"`
	UserID snowflake.JsonInt64 `json:"userId" v:"required#目标会员 ID 不能为空"`
	Remark string              `json:"remark" v:"max-length:500"`
}

// TeamExportRunRes 导出结果。
type TeamExportRunRes struct {
	g.Meta      `mime:"application/json"`
	ExportID    string `json:"exportId"`
	FileURL     string `json:"fileUrl" dc:"导出文件相对路径"`
	FileSize    int64  `json:"fileSize" dc:"压缩后字节数"`
	MemberCount int    `json:"memberCount" dc:"导出的会员数"`
}

// TeamExportDeployReq 后台触发站点裂变部署。
type TeamExportDeployReq struct {
	g.Meta   `path:"/team_export/deploy" method:"post" tags:"团队数据导出" summary:"基于导出包部署独立站点"`
	ExportID snowflake.JsonInt64 `json:"exportId" v:"required#导出记录 ID 不能为空"`
	Domain   string              `json:"domain" v:"required|max-length:200#部署域名不能为空|域名长度不能超过200位"`
}

// TeamExportDeployRes 部署响应。
type TeamExportDeployRes struct {
	g.Meta       `mime:"application/json"`
	DeployStatus int    `json:"deployStatus"`
	DeployDomain string `json:"deployDomain"`
}
