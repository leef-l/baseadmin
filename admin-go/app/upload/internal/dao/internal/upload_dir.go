// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UploadDirDao is the data access object for the table upload_dir.
type UploadDirDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UploadDirColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UploadDirColumns defines and stores column names for the table upload_dir.
type UploadDirColumns struct {
	Id         string // ID
	ParentId   string // 上级目录
	Name       string // 目录名称
	Path       string // 目录路径
	KeepName   string // 保留原文件名:0=否,1=是
	Sort       string // 排序
	Status     string // 状态:0=禁用,1=启用
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 删除时间
	CreatedBy  string // 创建人
	DeptId     string // 部门ID
	TenantId   string // 租户
	MerchantId string // 商户
}

// uploadDirColumns holds the columns for the table upload_dir.
var uploadDirColumns = UploadDirColumns{
	Id:         "id",
	ParentId:   "parent_id",
	Name:       "name",
	Path:       "path",
	KeepName:   "keep_name",
	Sort:       "sort",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
	CreatedBy:  "created_by",
	DeptId:     "dept_id",
	TenantId:   "tenant_id",
	MerchantId: "merchant_id",
}

// NewUploadDirDao creates and returns a new DAO object for table data access.
func NewUploadDirDao(handlers ...gdb.ModelHandler) *UploadDirDao {
	return &UploadDirDao{
		group:    "default",
		table:    "upload_dir",
		columns:  uploadDirColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UploadDirDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UploadDirDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UploadDirDao) Columns() UploadDirColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UploadDirDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UploadDirDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UploadDirDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
