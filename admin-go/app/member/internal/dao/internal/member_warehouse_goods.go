// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberWarehouseGoodsDao is the data access object for the table member_warehouse_goods.
type MemberWarehouseGoodsDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  MemberWarehouseGoodsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// MemberWarehouseGoodsColumns defines and stores column names for the table member_warehouse_goods.
type MemberWarehouseGoodsColumns struct {
	Id              string // ID（Snowflake）
	GoodsNo         string // 商品编号|search:eq|keyword:on|priority:100
	Title           string // 商品名称|search:like|keyword:on|priority:95
	Cover           string // 商品封面
	InitPrice       string // 初始价格（分）
	CurrentPrice    string // 当前价格（分）
	PriceRiseRate   string // 每次加价比例（百分比，如10=10%）
	PlatformFeeRate string // 平台扣除比例（百分比，如5=5%）
	OwnerId         string // 当前持有人|ref:member_user.nickname|search:select
	TradeCount      string // 流转次数
	GoodsStatus     string // 商品状态:1=持有中,2=挂卖中,3=交易中|search:select
	Remark          string // 备注|search:off
	Sort            string // 排序（升序）
	Status          string // 状态:0=关闭,1=开启|search:select
	TenantId        string // 租户
	MerchantId      string // 商户
	CreatedBy       string // 创建人ID
	DeptId          string // 所属部门ID
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
	DeletedAt       string // 软删除时间，非 NULL 表示已删除
}

// memberWarehouseGoodsColumns holds the columns for the table member_warehouse_goods.
var memberWarehouseGoodsColumns = MemberWarehouseGoodsColumns{
	Id:              "id",
	GoodsNo:         "goods_no",
	Title:           "title",
	Cover:           "cover",
	InitPrice:       "init_price",
	CurrentPrice:    "current_price",
	PriceRiseRate:   "price_rise_rate",
	PlatformFeeRate: "platform_fee_rate",
	OwnerId:         "owner_id",
	TradeCount:      "trade_count",
	GoodsStatus:     "goods_status",
	Remark:          "remark",
	Sort:            "sort",
	Status:          "status",
	TenantId:        "tenant_id",
	MerchantId:      "merchant_id",
	CreatedBy:       "created_by",
	DeptId:          "dept_id",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewMemberWarehouseGoodsDao creates and returns a new DAO object for table data access.
func NewMemberWarehouseGoodsDao(handlers ...gdb.ModelHandler) *MemberWarehouseGoodsDao {
	return &MemberWarehouseGoodsDao{
		group:    "default",
		table:    "member_warehouse_goods",
		columns:  memberWarehouseGoodsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberWarehouseGoodsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberWarehouseGoodsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberWarehouseGoodsDao) Columns() MemberWarehouseGoodsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberWarehouseGoodsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberWarehouseGoodsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberWarehouseGoodsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
