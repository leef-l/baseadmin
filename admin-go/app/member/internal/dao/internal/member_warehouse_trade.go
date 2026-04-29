// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberWarehouseTradeDao is the data access object for the table member_warehouse_trade.
type MemberWarehouseTradeDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  MemberWarehouseTradeColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// MemberWarehouseTradeColumns defines and stores column names for the table member_warehouse_trade.
type MemberWarehouseTradeColumns struct {
	Id           string // ID（Snowflake）
	TradeNo      string // 交易编号|search:eq|keyword:on|priority:100
	GoodsId      string // 仓库商品|ref:member_warehouse_goods.title|search:select
	ListingId    string // 挂卖记录|ref:member_warehouse_listing.id
	SellerId     string // 卖家|ref:member_user.nickname|search:select
	BuyerId      string // 买家|ref:member_user.nickname|search:select
	TradePrice   string // 成交价格（分）
	PlatformFee  string // 平台扣除费用（分）
	SellerIncome string // 卖家实收（分）
	TradeStatus  string // 交易状态:1=待卖家确认,2=已确认完成,3=已取消|search:select
	ConfirmedAt  string // 确认时间
	Remark       string // 备注|search:off
	Status       string // 状态:0=关闭,1=开启|search:select
	TenantId     string // 租户
	MerchantId   string // 商户
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// memberWarehouseTradeColumns holds the columns for the table member_warehouse_trade.
var memberWarehouseTradeColumns = MemberWarehouseTradeColumns{
	Id:           "id",
	TradeNo:      "trade_no",
	GoodsId:      "goods_id",
	ListingId:    "listing_id",
	SellerId:     "seller_id",
	BuyerId:      "buyer_id",
	TradePrice:   "trade_price",
	PlatformFee:  "platform_fee",
	SellerIncome: "seller_income",
	TradeStatus:  "trade_status",
	ConfirmedAt:  "confirmed_at",
	Remark:       "remark",
	Status:       "status",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewMemberWarehouseTradeDao creates and returns a new DAO object for table data access.
func NewMemberWarehouseTradeDao(handlers ...gdb.ModelHandler) *MemberWarehouseTradeDao {
	return &MemberWarehouseTradeDao{
		group:    "default",
		table:    "member_warehouse_trade",
		columns:  memberWarehouseTradeColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberWarehouseTradeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberWarehouseTradeDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberWarehouseTradeDao) Columns() MemberWarehouseTradeColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberWarehouseTradeDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberWarehouseTradeDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberWarehouseTradeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
