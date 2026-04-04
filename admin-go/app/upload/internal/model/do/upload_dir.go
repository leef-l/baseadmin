// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadDir is the golang structure of table upload_dir for DAO operations like Where/Data.
type UploadDir struct {
	g.Meta    `orm:"table:upload_dir, do:true"`
	Id        any         // ID
	ParentId  any         // 上级目录
	Name      any         // 目录名称
	Path      any         // 目录路径
	Sort      any         // 排序
	Status    any         // 状态:0=禁用,1=启用
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
	CreatedBy any         // 创建人
	DeptId    any         // 部门ID
}
