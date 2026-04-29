// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberLevelLogDao is the data access object for the table member_level_log.
type MemberLevelLogDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  MemberLevelLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// MemberLevelLogColumns defines and stores column names for the table member_level_log.
type MemberLevelLogColumns struct {
	Id         string // ID（Snowflake）
	UserId     string // 会员|ref:member_user.nickname|search:select
	OldLevelId string // 变更前等级|ref:member_level.name
	NewLevelId string // 变更后等级|ref:member_level.name
	ChangeType string // 变更类型:1=自动升级,2=后台调整,3=过期降级|search:select
	ExpireAt   string // 新等级到期时间
	Remark     string // 变更说明|search:off
	TenantId   string // 租户
	MerchantId string // 商户
	CreatedBy  string // 创建人ID
	DeptId     string // 所属部门ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间，非 NULL 表示已删除
}

// memberLevelLogColumns holds the columns for the table member_level_log.
var memberLevelLogColumns = MemberLevelLogColumns{
	Id:         "id",
	UserId:     "user_id",
	OldLevelId: "old_level_id",
	NewLevelId: "new_level_id",
	ChangeType: "change_type",
	ExpireAt:   "expire_at",
	Remark:     "remark",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewMemberLevelLogDao creates and returns a new DAO object for table data access.
func NewMemberLevelLogDao(handlers ...gdb.ModelHandler) *MemberLevelLogDao {
	return &MemberLevelLogDao{
		group:    "default",
		table:    "member_level_log",
		columns:  memberLevelLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberLevelLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberLevelLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberLevelLogDao) Columns() MemberLevelLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberLevelLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberLevelLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberLevelLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
