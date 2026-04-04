// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadDirRule is the golang structure for table upload_dir_rule.
type UploadDirRule struct {
	Id        uint64      `orm:"id"         description:"ID"`                // ID
	DirId     uint64      `orm:"dir_id"     description:"目录ID"`              // 目录ID
	Category  int         `orm:"category"   description:"类别:1=默认,2=类型,3=接口"` // 类别:1=默认,2=类型,3=接口
	SavePath  string      `orm:"save_path"  description:"保存目录"`              // 保存目录
	Status    int         `orm:"status"     description:"状态:0=禁用,1=启用"`      // 状态:0=禁用,1=启用
	CreatedAt *gtime.Time `orm:"created_at" description:"创建时间"`              // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" description:"更新时间"`              // 更新时间
	DeletedAt *gtime.Time `orm:"deleted_at" description:"删除时间"`              // 删除时间
	CreatedBy uint64      `orm:"created_by" description:"创建人"`               // 创建人
	DeptId    uint64      `orm:"dept_id"    description:"部门ID"`              // 部门ID
}
