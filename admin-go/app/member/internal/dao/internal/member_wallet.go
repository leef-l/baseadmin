// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberWalletDao is the data access object for the table member_wallet.
type MemberWalletDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  MemberWalletColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// MemberWalletColumns defines and stores column names for the table member_wallet.
type MemberWalletColumns struct {
	Id           string // ID（Snowflake）
	UserId       string // 会员|ref:member_user.nickname|search:select
	WalletType   string // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	Balance      string // 当前余额（分）
	TotalIncome  string // 累计收入（分）
	TotalExpense string // 累计支出（分）
	FrozenAmount string // 冻结金额（分）
	Status       string // 状态:0=冻结,1=正常|search:select
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// memberWalletColumns holds the columns for the table member_wallet.
var memberWalletColumns = MemberWalletColumns{
	Id:           "id",
	UserId:       "user_id",
	WalletType:   "wallet_type",
	Balance:      "balance",
	TotalIncome:  "total_income",
	TotalExpense: "total_expense",
	FrozenAmount: "frozen_amount",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewMemberWalletDao creates and returns a new DAO object for table data access.
func NewMemberWalletDao(handlers ...gdb.ModelHandler) *MemberWalletDao {
	return &MemberWalletDao{
		group:    "default",
		table:    "member_wallet",
		columns:  memberWalletColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberWalletDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberWalletDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberWalletDao) Columns() MemberWalletColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberWalletDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberWalletDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberWalletDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
