// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UploadConfigDao is the data access object for the table upload_config.
type UploadConfigDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  UploadConfigColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// UploadConfigColumns defines and stores column names for the table upload_config.
type UploadConfigColumns struct {
	Id           string // ID
	Name         string // 配置名称
	Storage      string // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsDefault    string // 是否默认:0=否,1=是
	LocalPath    string // 本地存储路径
	OssEndpoint  string // OSS Endpoint
	OssBucket    string // OSS Bucket
	OssAccessKey string // OSS AccessKey
	OssSecretKey string // OSS SecretKey
	CosRegion    string // COS Region
	CosBucket    string // COS Bucket
	CosSecretId  string // COS SecretId
	CosSecretKey string // COS SecretKey
	MaxSize      string // 最大文件大小(MB)
	Status       string // 状态:0=禁用,1=启用
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 删除时间
	CreatedBy    string // 创建人
	DeptId       string // 部门ID
	TenantId     string // 租户
	MerchantId   string // 商户
}

// uploadConfigColumns holds the columns for the table upload_config.
var uploadConfigColumns = UploadConfigColumns{
	Id:           "id",
	Name:         "name",
	Storage:      "storage",
	IsDefault:    "is_default",
	LocalPath:    "local_path",
	OssEndpoint:  "oss_endpoint",
	OssBucket:    "oss_bucket",
	OssAccessKey: "oss_access_key",
	OssSecretKey: "oss_secret_key",
	CosRegion:    "cos_region",
	CosBucket:    "cos_bucket",
	CosSecretId:  "cos_secret_id",
	CosSecretKey: "cos_secret_key",
	MaxSize:      "max_size",
	Status:       "status",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
	TenantId:     "tenant_id",
	MerchantId:   "merchant_id",
}

// NewUploadConfigDao creates and returns a new DAO object for table data access.
func NewUploadConfigDao(handlers ...gdb.ModelHandler) *UploadConfigDao {
	return &UploadConfigDao{
		group:    "default",
		table:    "upload_config",
		columns:  uploadConfigColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UploadConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UploadConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UploadConfigDao) Columns() UploadConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UploadConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UploadConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UploadConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
