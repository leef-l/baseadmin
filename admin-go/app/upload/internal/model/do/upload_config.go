// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadConfig is the golang structure of table upload_config for DAO operations like Where/Data.
type UploadConfig struct {
	g.Meta       `orm:"table:upload_config, do:true"`
	Id           any         // ID
	Name         any         // 配置名称
	Storage      any         // 存储类型:1=本地,2=阿里云OSS,3=腾讯云COS
	IsDefault    any         // 是否默认:0=否,1=是
	LocalPath    any         // 本地存储路径
	OssEndpoint  any         // OSS Endpoint
	OssBucket    any         // OSS Bucket
	OssAccessKey any         // OSS AccessKey
	OssSecretKey any         // OSS SecretKey
	CosRegion    any         // COS Region
	CosBucket    any         // COS Bucket
	CosSecretId  any         // COS SecretId
	CosSecretKey any         // COS SecretKey
	MaxSize      any         // 最大文件大小(MB)
	Status       any         // 状态:0=禁用,1=启用
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 删除时间
	CreatedBy    any         // 创建人
	DeptId       any         // 部门ID
	TenantId     any         // 租户
	MerchantId   any         // 商户
}
