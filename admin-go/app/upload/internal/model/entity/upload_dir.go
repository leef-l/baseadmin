// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadDir is the golang structure for table upload_dir.
type UploadDir struct {
	Id        uint64      `orm:"id"         description:"ID"`           // ID
	ParentId  uint64      `orm:"parent_id"  description:"上级目录"`         // 上级目录
	Name      string      `orm:"name"       description:"目录名称"`         // 目录名称
	Path      string      `orm:"path"       description:"目录路径"`         // 目录路径
	Sort      int         `orm:"sort"       description:"排序"`           // 排序
	Status    int         `orm:"status"     description:"状态:0=禁用,1=启用"` // 状态:0=禁用,1=启用
	CreatedAt *gtime.Time `orm:"created_at" description:"创建时间"`         // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" description:"更新时间"`         // 更新时间
	DeletedAt *gtime.Time `orm:"deleted_at" description:"删除时间"`         // 删除时间
	CreatedBy uint64      `orm:"created_by" description:"创建人"`          // 创建人
	DeptId    uint64      `orm:"dept_id"    description:"部门ID"`         // 部门ID
}
