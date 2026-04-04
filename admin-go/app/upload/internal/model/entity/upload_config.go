// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadConfig is the golang structure for table upload_config.
type UploadConfig struct {
	Id           uint64      `orm:"id"             description:"ID"`                          // ID
	Name         string      `orm:"name"           description:"配置名称"`                        // 配置名称
	Storage      int         `orm:"storage"        description:"存储类型:1=本地,2=阿里云OSS,3=腾讯云COS"` // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsDefault    int         `orm:"is_default"     description:"是否默认:0=否,1=是"`                // 是否默认:0=否,1=是
	LocalPath    string      `orm:"local_path"     description:"本地存储路径"`                      // 本地存储路径
	OssEndpoint  string      `orm:"oss_endpoint"   description:"OSS Endpoint"`                // OSS Endpoint
	OssBucket    string      `orm:"oss_bucket"     description:"OSS Bucket"`                  // OSS Bucket
	OssAccessKey string      `orm:"oss_access_key" description:"OSS AccessKey"`               // OSS AccessKey
	OssSecretKey string      `orm:"oss_secret_key" description:"OSS SecretKey"`               // OSS SecretKey
	CosRegion    string      `orm:"cos_region"     description:"COS Region"`                  // COS Region
	CosBucket    string      `orm:"cos_bucket"     description:"COS Bucket"`                  // COS Bucket
	CosSecretId  string      `orm:"cos_secret_id"  description:"COS SecretId"`                // COS SecretId
	CosSecretKey string      `orm:"cos_secret_key" description:"COS SecretKey"`               // COS SecretKey
	MaxSize      int         `orm:"max_size"       description:"最大文件大小(MB)"`                  // 最大文件大小(MB)
	Status       int         `orm:"status"         description:"状态:0=禁用,1=启用"`                // 状态:0=禁用,1=启用
	CreatedAt    *gtime.Time `orm:"created_at"     description:"创建时间"`                        // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"     description:"更新时间"`                        // 更新时间
	DeletedAt    *gtime.Time `orm:"deleted_at"     description:"删除时间"`                        // 删除时间
	CreatedBy    uint64      `orm:"created_by"     description:"创建人"`                         // 创建人
	DeptId       uint64      `orm:"dept_id"        description:"部门ID"`                        // 部门ID
}
