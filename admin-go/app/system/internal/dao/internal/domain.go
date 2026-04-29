// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DomainDao is the data access object for the table system_domain.
type DomainDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DomainColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DomainColumns defines and stores column names for the table system_domain.
type DomainColumns struct {
	Id           string // 域名ID（Snowflake）
	Domain       string // 绑定域名
	OwnerType    string // 主体类型:1=租户,2=商户
	TenantId     string // 租户
	MerchantId   string // 商户
	AppCode      string // 应用编码：admin/upload/shop
	VerifyToken  string // 域名校验令牌
	VerifyStatus string // 校验状态:0=未校验,1=已校验
	SslStatus    string // SSL状态:0=未配置,1=已配置
	NginxStatus  string // Nginx配置状态:0=未应用,1=已应用
	Status       string // 状态:0=关闭,1=开启
	Remark       string // 备注
	CreatedBy    string // 创建人ID
	DeptId       string // 所属部门ID
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 软删除时间，非 NULL 表示已删除
}

// domainColumns holds the columns for the table system_domain.
var domainColumns = DomainColumns{
	Id:           "id",
	Domain:       "domain",
	OwnerType:    "owner_type",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
	AppCode:      "app_code",
	VerifyToken:  "verify_token",
	VerifyStatus: "verify_status",
	SslStatus:    "ssl_status",
	NginxStatus:  "nginx_status",
	Status:       "status",
	Remark:       "remark",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewDomainDao creates and returns a new DAO object for table data access.
func NewDomainDao(handlers ...gdb.ModelHandler) *DomainDao {
	return &DomainDao{
		group:    "default",
		table:    "system_domain",
		columns:  domainColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DomainDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DomainDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DomainDao) Columns() DomainColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DomainDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DomainDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DomainDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
