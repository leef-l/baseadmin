// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberUserDao is the data access object for the table member_user.
type MemberUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MemberUserColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MemberUserColumns defines and stores column names for the table member_user.
type MemberUserColumns struct {
	Id                 string // 会员ID（Snowflake）
	ParentId           string // 上级会员
	Username           string // 用户名（登录账号）|search:eq|keyword:on|priority:100
	Password           string // 密码（bcrypt加密）
	Nickname           string // 昵称|search:like|keyword:on|priority:95
	Phone              string // 手机号|search:eq|keyword:on|priority:90
	Avatar             string // 头像
	RealName           string // 真实姓名|search:like|keyword:on
	LevelId            string // 当前等级|ref:member_level.name|search:select
	LevelExpireAt      string // 等级到期时间
	TeamCount          string // 团队总人数
	DirectCount        string // 直推人数
	ActiveCount        string // 有效用户数
	TeamTurnover       string // 团队总营业额（分）
	IsActive           string // 是否激活:0=未激活,1=已激活|search:select
	IsQualified        string // 仓库资格:0=已失效,1=有效|search:select
	DailyPurchaseLimit string // 本会员每日限购单数（按等级初始化，可单独调整）|search:eq
	TodayPurchaseCount string // 今日已购单数|search:off
	LastPurchaseDate   string // 最近购买日期（跨日重置 today_purchase_count）|search:off
	TotalPurchaseCount string // 历史总购单数（用于阶梯返佣判断）|search:off
	InviteCode         string // 邀请码|search:eq
	RegisterIp         string // 注册IP
	LastLoginAt        string // 最后登录时间
	Remark             string // 备注|search:off
	Sort               string // 排序（升序）
	Status             string // 状态:0=冻结,1=正常|search:select
	TenantId           string // 租户
	MerchantId         string // 商户
	CreatedBy          string // 创建人ID
	DeptId             string // 所属部门ID
	CreatedAt          string // 创建时间
	UpdatedAt          string // 更新时间
	DeletedAt          string // 软删除时间，非 NULL 表示已删除
}

// memberUserColumns holds the columns for the table member_user.
var memberUserColumns = MemberUserColumns{
	Id:                 "id",
	ParentId:           "parent_id",
	Username:           "username",
	Password:           "password",
	Nickname:           "nickname",
	Phone:              "phone",
	Avatar:             "avatar",
	RealName:           "real_name",
	LevelId:            "level_id",
	LevelExpireAt:      "level_expire_at",
	TeamCount:          "team_count",
	DirectCount:        "direct_count",
	ActiveCount:        "active_count",
	TeamTurnover:       "team_turnover",
	IsActive:           "is_active",
	IsQualified:        "is_qualified",
	DailyPurchaseLimit: "daily_purchase_limit",
	TodayPurchaseCount: "today_purchase_count",
	LastPurchaseDate:   "last_purchase_date",
	TotalPurchaseCount: "total_purchase_count",
	InviteCode:         "invite_code",
	RegisterIp:         "register_ip",
	LastLoginAt:        "last_login_at",
	Remark:             "remark",
	Sort:               "sort",
	Status:             "status",
	TenantId:           "tenant_id",
	MerchantId:         "merchant_id",
	CreatedBy:          "created_by",
	DeptId:             "dept_id",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewMemberUserDao creates and returns a new DAO object for table data access.
func NewMemberUserDao(handlers ...gdb.ModelHandler) *MemberUserDao {
	return &MemberUserDao{
		group:    "default",
		table:    "member_user",
		columns:  memberUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberUserDao) Columns() MemberUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
