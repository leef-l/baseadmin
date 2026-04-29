// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberLevelDao is the data access object for the table member_level.
type MemberLevelDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MemberLevelColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MemberLevelColumns defines and stores column names for the table member_level.
type MemberLevelColumns struct {
	Id               string // 等级ID（Snowflake）
	Name             string // 等级名称|search:like|keyword:on|priority:100
	LevelNo          string // 等级编号（越大越高）|search:eq
	Icon             string // 等级图标
	DurationDays     string // 有效天数（0=永久）
	NeedActiveCount  string // 升级所需有效用户数
	NeedTeamTurnover string // 升级所需团队营业额（分）
	IsTop            string // 是否最高等级:0=否,1=是|search:select
	AutoDeploy       string // 到达后自动部署站点:0=否,1=是
	Remark           string // 等级说明|search:off
	Sort             string // 排序（升序）
	Status           string // 状态:0=关闭,1=开启|search:select
	TenantId         string // 租户
	MerchantId       string // 商户
	CreatedBy        string // 创建人ID
	DeptId           string // 所属部门ID
	CreatedAt        string // 创建时间
	UpdatedAt        string // 更新时间
	DeletedAt        string // 软删除时间，非 NULL 表示已删除
}

// memberLevelColumns holds the columns for the table member_level.
var memberLevelColumns = MemberLevelColumns{
	Id:               "id",
	Name:             "name",
	LevelNo:          "level_no",
	Icon:             "icon",
	DurationDays:     "duration_days",
	NeedActiveCount:  "need_active_count",
	NeedTeamTurnover: "need_team_turnover",
	IsTop:            "is_top",
	AutoDeploy:       "auto_deploy",
	Remark:           "remark",
	Sort:             "sort",
	Status:           "status",
	TenantId:         "tenant_id",
	MerchantId:       "merchant_id",
	CreatedBy:        "created_by",
	DeptId:           "dept_id",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewMemberLevelDao creates and returns a new DAO object for table data access.
func NewMemberLevelDao(handlers ...gdb.ModelHandler) *MemberLevelDao {
	return &MemberLevelDao{
		group:    "default",
		table:    "member_level",
		columns:  memberLevelColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberLevelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberLevelDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberLevelDao) Columns() MemberLevelColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberLevelDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberLevelDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberLevelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
