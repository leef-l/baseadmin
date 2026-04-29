// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberTeamExport is the golang structure of table member_team_export for DAO operations like Where/Data.
type MemberTeamExport struct {
	g.Meta          `orm:"table:member_team_export, do:true"`
	Id              any         // ID（Snowflake）
	UserId          any         // 目标会员|ref:member_user.nickname|search:select
	TeamMemberCount any         // 团队成员数
	ExportType      any         // 导出类型:1=手动导出,2=自动升级导出|search:select
	FileUrl         any         // 导出文件地址
	FileSize        any         // 文件大小（字节）
	DeployStatus    any         // 部署状态:0=未部署,1=部署中,2=已部署,3=部署失败|search:select
	DeployDomain    any         // 部署域名|search:like
	DeployedAt      *gtime.Time // 部署完成时间
	Remark          any         // 备注|search:off
	Status          any         // 状态:0=关闭,1=开启|search:select
	TenantId        any         // 租户
	MerchantId      any         // 商户
	CreatedBy       any         // 创建人ID
	DeptId          any         // 所属部门ID
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
	DeletedAt       *gtime.Time // 软删除时间，非 NULL 表示已删除
}
