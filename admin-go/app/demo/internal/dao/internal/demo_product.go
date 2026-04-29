// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DemoProductDao is the data access object for the table demo_product.
type DemoProductDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DemoProductColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DemoProductColumns defines and stores column names for the table demo_product.
type DemoProductColumns struct {
	Id            string // 商品ID（Snowflake）
	CategoryId    string // 商品分类
	SkuNo         string // SKU编号|search:eq|priority:100
	Name          string // 商品名称|search:like|keyword:on|priority:95
	Cover         string // 封面
	ManualFile    string // 说明书文件
	DetailContent string // 详情内容
	SpecJson      string // 规格JSON
	WebsiteUrl    string // 官网URL
	Type          string // 类型:1=普通,2=置顶,3=推荐,4=热门
	IsRecommend   string // 是否推荐:0=否,1=是
	SalePrice     string // 销售价（分）
	StockNum      string // 库存数量
	WeightNum     string // 重量（克）
	Sort          string // 排序（升序）
	Icon          string // 图标
	Status        string // 状态:0=草稿,1=上架,2=下架
	TenantId      string // 租户
	MerchantId    string // 商户
	CreatedBy     string // 创建人ID
	DeptId        string // 所属部门ID
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 软删除时间，非 NULL 表示已删除
}

// demoProductColumns holds the columns for the table demo_product.
var demoProductColumns = DemoProductColumns{
	Id:            "id",
	CategoryId:    "category_id",
	SkuNo:         "sku_no",
	Name:          "name",
	Cover:         "cover",
	ManualFile:    "manual_file",
	DetailContent: "detail_content",
	SpecJson:      "spec_json",
	WebsiteUrl:    "website_url",
	Type:          "type",
	IsRecommend:   "is_recommend",
	SalePrice:     "sale_price",
	StockNum:      "stock_num",
	WeightNum:     "weight_num",
	Sort:          "sort",
	Icon:          "icon",
	Status:        "status",
	TenantId:      "tenant_id",
	MerchantId:    "merchant_id",
	CreatedBy:     "created_by",
	DeptId:        "dept_id",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewDemoProductDao creates and returns a new DAO object for table data access.
func NewDemoProductDao(handlers ...gdb.ModelHandler) *DemoProductDao {
	return &DemoProductDao{
		group:    "default",
		table:    "demo_product",
		columns:  demoProductColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DemoProductDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DemoProductDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DemoProductDao) Columns() DemoProductColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DemoProductDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DemoProductDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DemoProductDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
