// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberShopGoodsDao is the data access object for the table member_shop_goods.
type MemberShopGoodsDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  MemberShopGoodsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// MemberShopGoodsColumns defines and stores column names for the table member_shop_goods.
type MemberShopGoodsColumns struct {
	Id            string // ID（Snowflake）
	CategoryId    string // 商品分类|ref:member_shop_category.name|search:select
	Title         string // 商品名称|search:like|keyword:on|priority:100
	Cover         string // 封面图
	Images        string // 商品图片（JSON数组）|search:off
	Price         string // 售价（分，优惠券余额支付）
	OriginalPrice string // 原价（分）
	Stock         string // 库存
	Sales         string // 销量
	Content       string // 商品详情|search:off
	Sort          string // 排序（升序）
	IsRecommend   string // 是否推荐:0=否,1=是|search:select
	Status        string // 状态:0=下架,1=上架|search:select
	TenantId      string // 租户
	MerchantId    string // 商户
	CreatedBy     string // 创建人ID
	DeptId        string // 所属部门ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 软删除时间，非 NULL 表示已删除
}

// memberShopGoodsColumns holds the columns for the table member_shop_goods.
var memberShopGoodsColumns = MemberShopGoodsColumns{
	Id:            "id",
	CategoryId:    "category_id",
	Title:         "title",
	Cover:         "cover",
	Images:        "images",
	Price:         "price",
	OriginalPrice: "original_price",
	Stock:         "stock",
	Sales:         "sales",
	Content:       "content",
	Sort:          "sort",
	IsRecommend:   "is_recommend",
	Status:        "status",
	TenantId:      "tenant_id",
	MerchantId:    "merchant_id",
	CreatedBy:     "created_by",
	DeptId:        "dept_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewMemberShopGoodsDao creates and returns a new DAO object for table data access.
func NewMemberShopGoodsDao(handlers ...gdb.ModelHandler) *MemberShopGoodsDao {
	return &MemberShopGoodsDao{
		group:    "default",
		table:    "member_shop_goods",
		columns:  memberShopGoodsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberShopGoodsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberShopGoodsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberShopGoodsDao) Columns() MemberShopGoodsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberShopGoodsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberShopGoodsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberShopGoodsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
