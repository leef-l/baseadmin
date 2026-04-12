// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadDirRule is the golang structure of table upload_dir_rule for DAO operations like Where/Data.
type UploadDirRule struct {
	g.Meta    `orm:"table:upload_dir_rule, do:true"`
	Id        any         // ID
	DirId     any         // 目录ID
	Category  any         // 类别:1=默认,2=类型,3=接口
	FileType  any         // 文件类型，多个用逗号分隔
	SavePath  any         // 保存目录
	Status    any         // 状态:0=禁用,1=启用
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
	CreatedBy any         // 创建人
	DeptId    any         // 部门ID
}
