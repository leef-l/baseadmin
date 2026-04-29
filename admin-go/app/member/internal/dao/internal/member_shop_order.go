// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberShopOrderDao is the data access object for the table member_shop_order.
type MemberShopOrderDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  MemberShopOrderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// MemberShopOrderColumns defines and stores column names for the table member_shop_order.
type MemberShopOrderColumns struct {
	Id          string // ID（Snowflake）
	OrderNo     string // 订单号|search:eq|keyword:on|priority:100
	UserId      string // 购买会员|ref:member_user.nickname|search:select
	GoodsId     string // 商品|ref:member_shop_goods.title|search:select
	GoodsTitle  string // 商品名称（快照）
	GoodsCover  string // 商品封面（快照）
	Quantity    string // 购买数量
	TotalPrice  string // 订单总价（分）
	PayWallet   string // 支付钱包:1=优惠券余额
	OrderStatus string // 订单状态:1=已完成,2=已取消|search:select
	Remark      string // 订单备注|search:off
	Status      string // 状态:0=关闭,1=开启|search:select
	TenantId    string // 租户
	MerchantId  string // 商户
	CreatedBy   string // 创建人ID
	DeptId      string // 所属部门ID
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 软删除时间，非 NULL 表示已删除
}

// memberShopOrderColumns holds the columns for the table member_shop_order.
var memberShopOrderColumns = MemberShopOrderColumns{
	Id:          "id",
	OrderNo:     "order_no",
	UserId:      "user_id",
	GoodsId:     "goods_id",
	GoodsTitle:  "goods_title",
	GoodsCover:  "goods_cover",
	Quantity:    "quantity",
	TotalPrice:  "total_price",
	PayWallet:   "pay_wallet",
	OrderStatus: "order_status",
	Remark:      "remark",
	Status:      "status",
	TenantId:    "tenant_id",
	MerchantId:  "merchant_id",
	CreatedBy:   "created_by",
	DeptId:      "dept_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewMemberShopOrderDao creates and returns a new DAO object for table data access.
func NewMemberShopOrderDao(handlers ...gdb.ModelHandler) *MemberShopOrderDao {
	return &MemberShopOrderDao{
		group:    "default",
		table:    "member_shop_order",
		columns:  memberShopOrderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberShopOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberShopOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberShopOrderDao) Columns() MemberShopOrderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberShopOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberShopOrderDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberShopOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
