// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberWarehouseListingDao is the data access object for the table member_warehouse_listing.
type MemberWarehouseListingDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  MemberWarehouseListingColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// MemberWarehouseListingColumns defines and stores column names for the table member_warehouse_listing.
type MemberWarehouseListingColumns struct {
	Id            string // ID（Snowflake）
	GoodsId       string // 仓库商品|ref:member_warehouse_goods.title|search:select
	SellerId      string // 卖家|ref:member_user.nickname|search:select
	ListingPrice  string // 挂卖价格（分，自动加价后）
	ListingStatus string // 挂卖状态:1=挂卖中,2=已售出,3=已取消|search:select
	ListedAt      string // 挂卖时间
	SoldAt        string // 售出时间
	Remark        string // 备注|search:off
	Status        string // 状态:0=关闭,1=开启|search:select
	TenantId      string // 租户
	MerchantId    string // 商户
	CreatedBy     string // 创建人ID
	DeptId        string // 所属部门ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 软删除时间，非 NULL 表示已删除
}

// memberWarehouseListingColumns holds the columns for the table member_warehouse_listing.
var memberWarehouseListingColumns = MemberWarehouseListingColumns{
	Id:            "id",
	GoodsId:       "goods_id",
	SellerId:      "seller_id",
	ListingPrice:  "listing_price",
	ListingStatus: "listing_status",
	ListedAt:      "listed_at",
	SoldAt:        "sold_at",
	Remark:        "remark",
	Status:        "status",
	TenantId:      "tenant_id",
	MerchantId:    "merchant_id",
	CreatedBy:     "created_by",
	DeptId:        "dept_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewMemberWarehouseListingDao creates and returns a new DAO object for table data access.
func NewMemberWarehouseListingDao(handlers ...gdb.ModelHandler) *MemberWarehouseListingDao {
	return &MemberWarehouseListingDao{
		group:    "default",
		table:    "member_warehouse_listing",
		columns:  memberWarehouseListingColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberWarehouseListingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberWarehouseListingDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberWarehouseListingDao) Columns() MemberWarehouseListingColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberWarehouseListingDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberWarehouseListingDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberWarehouseListingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
