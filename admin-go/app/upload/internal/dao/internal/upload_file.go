// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UploadFileDao is the data access object for the table upload_file.
type UploadFileDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UploadFileColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UploadFileColumns defines and stores column names for the table upload_file.
type UploadFileColumns struct {
	Id        string // ID
	DirId     string // 所属目录
	Name      string // 文件名称
	Url       string // 文件地址
	Ext       string // 文件扩展名
	Size      string // 文件大小
	Mime      string // MIME类型
	Storage   string // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsImage   string // 是否图片:0=否,1=是
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
	DeletedAt string // 删除时间
	CreatedBy string // 创建人
	DeptId    string // 部门ID
}

// uploadFileColumns holds the columns for the table upload_file.
var uploadFileColumns = UploadFileColumns{
	Id:        "id",
	DirId:     "dir_id",
	Name:      "name",
	Url:       "url",
	Ext:       "ext",
	Size:      "size",
	Mime:      "mime",
	Storage:   "storage",
	IsImage:   "is_image",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
	CreatedBy: "created_by",
	DeptId:    "dept_id",
}

// NewUploadFileDao creates and returns a new DAO object for table data access.
func NewUploadFileDao(handlers ...gdb.ModelHandler) *UploadFileDao {
	return &UploadFileDao{
		group:    "default",
		table:    "upload_file",
		columns:  uploadFileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UploadFileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UploadFileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UploadFileDao) Columns() UploadFileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UploadFileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UploadFileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UploadFileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
