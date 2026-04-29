// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadFile is the golang structure for table upload_file.
type UploadFile struct {
	Id         uint64      `orm:"id"          description:"ID"`                          // ID
	DirId      uint64      `orm:"dir_id"      description:"所属目录"`                        // 所属目录
	Name       string      `orm:"name"        description:"文件名称"`                        // 文件名称
	Url        string      `orm:"url"         description:"文件地址"`                        // 文件地址
	Ext        string      `orm:"ext"         description:"文件扩展名"`                       // 文件扩展名
	Size       uint64      `orm:"size"        description:"文件大小"`                        // 文件大小
	Mime       string      `orm:"mime"        description:"MIME类型"`                      // MIME类型
	Storage    int         `orm:"storage"     description:"存储类型:1=本地,2=阿里云OSS,3=腾讯云COS"` // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsImage    int         `orm:"is_image"    description:"是否图片:0=否,1=是"`                // 是否图片:0=否,1=是
	CreatedAt  *gtime.Time `orm:"created_at"  description:"创建时间"`                        // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"  description:"更新时间"`                        // 更新时间
	DeletedAt  *gtime.Time `orm:"deleted_at"  description:"删除时间"`                        // 删除时间
	CreatedBy  uint64      `orm:"created_by"  description:"创建人"`                         // 创建人
	DeptId     uint64      `orm:"dept_id"     description:"部门ID"`                        // 部门ID
	TenantId   uint64      `orm:"tenant_id"   description:"租户"`                          // 租户
	MerchantId uint64      `orm:"merchant_id" description:"商户"`                          // 商户
}
