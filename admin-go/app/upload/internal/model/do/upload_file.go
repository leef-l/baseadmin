// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadFile is the golang structure of table upload_file for DAO operations like Where/Data.
type UploadFile struct {
	g.Meta     `orm:"table:upload_file, do:true"`
	Id         any         // ID
	DirId      any         // 所属目录
	Name       any         // 文件名称
	Url        any         // 文件地址
	Ext        any         // 文件扩展名
	Size       any         // 文件大小
	Mime       any         // MIME类型
	Storage    any         // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsImage    any         // 是否图片:0=否,1=是
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 删除时间
	CreatedBy  any         // 创建人
	DeptId     any         // 部门ID
	TenantId   any         // 租户
	MerchantId any         // 商户
}
