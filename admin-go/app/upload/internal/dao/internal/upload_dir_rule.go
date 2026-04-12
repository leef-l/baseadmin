// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UploadDirRuleDao is the data access object for the table upload_dir_rule.
type UploadDirRuleDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  UploadDirRuleColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// UploadDirRuleColumns defines and stores column names for the table upload_dir_rule.
type UploadDirRuleColumns struct {
	Id           string // ID
	DirId        string // 目录ID
	Category     string // 类别:1=默认,2=类型,3=接口
	FileType     string // 文件类型，多个用逗号分隔
	StorageTypes string // 适用存储类型，多个用逗号分隔
	SavePath     string // 保存目录
	Status       string // 状态:0=禁用,1=启用
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
	DeletedAt    string // 删除时间
	CreatedBy    string // 创建人
	DeptId       string // 部门ID
}

// uploadDirRuleColumns holds the columns for the table upload_dir_rule.
var uploadDirRuleColumns = UploadDirRuleColumns{
	Id:           "id",
	DirId:        "dir_id",
	Category:     "category",
	FileType:     "file_type",
	StorageTypes: "storage_types",
	SavePath:     "save_path",
	Status:       "status",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
	CreatedBy:    "created_by",
	DeptId:       "dept_id",
}

// NewUploadDirRuleDao creates and returns a new DAO object for table data access.
func NewUploadDirRuleDao(handlers ...gdb.ModelHandler) *UploadDirRuleDao {
	return &UploadDirRuleDao{
		group:    "default",
		table:    "upload_dir_rule",
		columns:  uploadDirRuleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UploadDirRuleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UploadDirRuleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UploadDirRuleDao) Columns() UploadDirRuleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UploadDirRuleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UploadDirRuleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UploadDirRuleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
