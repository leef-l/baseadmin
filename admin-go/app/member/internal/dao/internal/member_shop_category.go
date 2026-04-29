// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberShopCategoryDao is the data access object for the table member_shop_category.
type MemberShopCategoryDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  MemberShopCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// MemberShopCategoryColumns defines and stores column names for the table member_shop_category.
type MemberShopCategoryColumns struct {
	Id         string // ID（Snowflake）
	ParentId   string // 上级分类
	Name       string // 分类名称|search:like|keyword:on|priority:100
	Icon       string // 分类图标
	Sort       string // 排序（升序）
	Status     string // 状态:0=关闭,1=开启|search:select
	TenantId   string // 租户
	MerchantId string // 商户
	CreatedBy  string // 创建人ID
	DeptId     string // 所属部门ID
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 软删除时间，非 NULL 表示已删除
}

// memberShopCategoryColumns holds the columns for the table member_shop_category.
var memberShopCategoryColumns = MemberShopCategoryColumns{
	Id:         "id",
	ParentId:   "parent_id",
	Name:       "name",
	Icon:       "icon",
	Sort:       "sort",
	Status:     "status",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewMemberShopCategoryDao creates and returns a new DAO object for table data access.
func NewMemberShopCategoryDao(handlers ...gdb.ModelHandler) *MemberShopCategoryDao {
	return &MemberShopCategoryDao{
		group:    "default",
		table:    "member_shop_category",
		columns:  memberShopCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MemberShopCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MemberShopCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MemberShopCategoryDao) Columns() MemberShopCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MemberShopCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MemberShopCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MemberShopCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
