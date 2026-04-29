// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberWalletLogDao is the data access object for the table member_wallet_log.
type MemberWalletLogDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  MemberWalletLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// MemberWalletLogColumns defines and stores column names for the table member_wallet_log.
type MemberWalletLogColumns struct {
	Id             string // ID（Snowflake）
	UserId         string // 会员|ref:member_user.nickname|search:select
	WalletType     string // 钱包类型:1=优惠券余额,2=奖金余额,3=推广奖余额|search:select
	ChangeType     string // 变动类型:1=充值,2=消费,3=推广奖,4=仓库卖出收入,5=平台扣除,6=后台调整|search:select
	ChangeAmount   string // 变动金额（分，正增负减）
	BeforeBalance  string // 变动前余额（分）
	AfterBalance   string // 变动后余额（分）
	RelatedOrderNo string // 关联单号|search:eq|keyword:on
	Remark         string // 备注说明|search:off
	TenantId       string // 租户
	MerchantId     string // 商户
	CreatedBy      string // 创建人ID
	DeptId         string // 所属部门ID
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
	DeletedAt      string // 软删除时间，非 NULL 表示已删除
}

// memberWalletLogColumns holds the columns for the table member_wallet_log.
var memberWalletLogColumns = MemberWalletLogColumns{
	Id:             "id",
	UserId:         "user_id",
	WalletType:     "wallet_type",
	ChangeType:     "change_type",
	ChangeAmount:   "change_amount",
	BeforeBalance:  "before_balance",
	AfterBalance:   "after_balance",
	RelatedOrderNo: "related_order_no",
	Remark:         "remark",
	TenantId:       "tenant_id",
	MerchantId:     "merchant_id",
	CreatedBy:      "created_by",
	DeptId:         "dept_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	DeletedAt:      "deleted_at",
}

// NewMemberWalletLogDao creates and returns a new DAO object for table data access.
func NewMemberWalletLogDao(handlers ...gdb.ModelHandler) *MemberWalletLogDao {
	return &MemberWalletLogDao{
		group:    "default",
		table:    "member_wallet_log",
		columns:  memberWalletLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberWalletLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberWalletLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberWalletLogDao) Columns() MemberWalletLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberWalletLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberWalletLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberWalletLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
