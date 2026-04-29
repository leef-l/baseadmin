// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberRebindLogDao is the data access object for the table member_rebind_log.
type MemberRebindLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  MemberRebindLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// MemberRebindLogColumns defines and stores column names for the table member_rebind_log.
type MemberRebindLogColumns struct {
	Id          string // ID（Snowflake）
	UserId      string // 会员|ref:member_user.nickname|search:select
	OldParentId string // 原上级|ref:member_user.nickname
	NewParentId string // 新上级|ref:member_user.nickname
	Reason      string // 换绑原因|search:off
	OperatorId  string // 操作人|ref:system_users.username
	TenantId    string // 租户
	MerchantId  string // 商户
	CreatedBy   string // 创建人ID
	DeptId      string // 所属部门ID
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 软删除时间，非 NULL 表示已删除
}

// memberRebindLogColumns holds the columns for the table member_rebind_log.
var memberRebindLogColumns = MemberRebindLogColumns{
	Id:          "id",
	UserId:      "user_id",
	OldParentId: "old_parent_id",
	NewParentId: "new_parent_id",
	Reason:      "reason",
	OperatorId:  "operator_id",
	TenantId:    "tenant_id",
	MerchantId:  "merchant_id",
	CreatedBy:   "created_by",
	DeptId:      "dept_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewMemberRebindLogDao creates and returns a new DAO object for table data access.
func NewMemberRebindLogDao(handlers ...gdb.ModelHandler) *MemberRebindLogDao {
	return &MemberRebindLogDao{
		group:    "default",
		table:    "member_rebind_log",
		columns:  memberRebindLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberRebindLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberRebindLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberRebindLogDao) Columns() MemberRebindLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberRebindLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberRebindLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberRebindLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
