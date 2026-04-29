// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberTeamExportDao is the data access object for the table member_team_export.
type MemberTeamExportDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  MemberTeamExportColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// MemberTeamExportColumns defines and stores column names for the table member_team_export.
type MemberTeamExportColumns struct {
	Id              string // ID（Snowflake）
	UserId          string // 目标会员|ref:member_user.nickname|search:select
	TeamMemberCount string // 团队成员数
	ExportType      string // 导出类型:1=手动导出,2=自动升级导出|search:select
	FileUrl         string // 导出文件地址
	FileSize        string // 文件大小（字节）
	DeployStatus    string // 部署状态:0=未部署,1=部署中,2=已部署,3=部署失败|search:select
	DeployDomain    string // 部署域名|search:like
	DeployedAt      string // 部署完成时间
	Remark          string // 备注|search:off
	Status          string // 状态:0=关闭,1=开启|search:select
	TenantId        string // 租户
	MerchantId      string // 商户
	CreatedBy       string // 创建人ID
	DeptId          string // 所属部门ID
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	DeletedAt       string // 软删除时间，非 NULL 表示已删除
}

// memberTeamExportColumns holds the columns for the table member_team_export.
var memberTeamExportColumns = MemberTeamExportColumns{
	Id:              "id",
	UserId:          "user_id",
	TeamMemberCount: "team_member_count",
	ExportType:      "export_type",
	FileUrl:         "file_url",
	FileSize:        "file_size",
	DeployStatus:    "deploy_status",
	DeployDomain:    "deploy_domain",
	DeployedAt:      "deployed_at",
	Remark:          "remark",
	Status:          "status",
	TenantId:        "tenant_id",
	MerchantId:      "merchant_id",
	CreatedBy:       "created_by",
	DeptId:          "dept_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewMemberTeamExportDao creates and returns a new DAO object for table data access.
func NewMemberTeamExportDao(handlers ...gdb.ModelHandler) *MemberTeamExportDao {
	return &MemberTeamExportDao{
		group:    "default",
		table:    "member_team_export",
		columns:  memberTeamExportColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberTeamExportDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberTeamExportDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberTeamExportDao) Columns() MemberTeamExportColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberTeamExportDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberTeamExportDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *MemberTeamExportDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
