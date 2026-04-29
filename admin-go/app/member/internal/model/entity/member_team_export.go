// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// MemberTeamExport is the golang structure for table member_team_export.
type MemberTeamExport struct {
	Id              uint64      `orm:"id"                description:"ID（Snowflake）"`                               // ID（Snowflake）
	UserId          uint64      `orm:"user_id"           description:"目标会员|ref:member_user.nickname|search:select"` // 目标会员|ref:member_user.nickname|search:select
	TeamMemberCount uint        `orm:"team_member_count" description:"团队成员数"`                                       // 团队成员数
	ExportType      int         `orm:"export_type"       description:"导出类型:1=手动导出,2=自动升级导出|search:select"`          // 导出类型:1=手动导出,2=自动升级导出|search:select
	FileUrl         string      `orm:"file_url"          description:"导出文件地址"`                                      // 导出文件地址
	FileSize        uint64      `orm:"file_size"         description:"文件大小（字节）"`                                    // 文件大小（字节）
	DeployStatus    int         `orm:"deploy_status"     description:"部署状态:0=未部署,1=部署中,2=已部署,3=部署失败|search:select"` // 部署状态:0=未部署,1=部署中,2=已部署,3=部署失败|search:select
	DeployDomain    string      `orm:"deploy_domain"     description:"部署域名|search:like"`                            // 部署域名|search:like
	DeployedAt      *gtime.Time `orm:"deployed_at"       description:"部署完成时间"`                                      // 部署完成时间
	Remark          string      `orm:"remark"            description:"备注|search:off"`                               // 备注|search:off
	Status          int         `orm:"status"            description:"状态:0=关闭,1=开启|search:select"`                  // 状态:0=关闭,1=开启|search:select
	TenantId        uint64      `orm:"tenant_id"         description:"租户"`                                          // 租户
	MerchantId      uint64      `orm:"merchant_id"       description:"商户"`                                          // 商户
	CreatedBy       uint64      `orm:"created_by"        description:"创建人ID"`                                       // 创建人ID
	DeptId          uint64      `orm:"dept_id"           description:"所属部门ID"`                                      // 所属部门ID
	CreatedAt       *gtime.Time `orm:"created_at"        description:"创建时间"`                                        // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"        description:"更新时间"`                                        // 更新时间
	DeletedAt       *gtime.Time `orm:"deleted_at"        description:"软删除时间，非 NULL 表示已删除"`                          // 软删除时间，非 NULL 表示已删除
}
