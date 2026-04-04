package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

// 确保 gtime 被引用
var _ = gtime.New

// Config API

// ConfigCreateReq 创建上传配置请求
type ConfigCreateReq struct {
	g.Meta       `path:"/config/create" method:"post" tags:"上传配置" summary:"创建上传配置"`
	Name         string `json:"name" v:"required#配置名称不能为空" dc:"配置名称"`
	Storage      int    `json:"storage"  dc:"存储类型"`
	IsDefault    int    `json:"isDefault"  dc:"是否默认"`
	LocalPath    string `json:"localPath"  dc:"本地存储路径"`
	OssEndpoint  string `json:"ossEndpoint"  dc:"OSS Endpoint"`
	OssBucket    string `json:"ossBucket"  dc:"OSS Bucket"`
	OssAccessKey string `json:"ossAccessKey"  dc:"OSS AccessKey"`
	OssSecretKey string `json:"ossSecretKey"  dc:"OSS SecretKey"`
	CosRegion    string `json:"cosRegion"  dc:"COS Region"`
	CosBucket    string `json:"cosBucket"  dc:"COS Bucket"`
	CosSecretID  string `json:"cosSecretID"  dc:"COS SecretId"`
	CosSecretKey string `json:"cosSecretKey"  dc:"COS SecretKey"`
	MaxSize      int    `json:"maxSize"  dc:"最大文件大小(MB)"`
	Status       int    `json:"status"  dc:"状态"`
}

// ConfigCreateRes 创建上传配置响应
type ConfigCreateRes struct {
	g.Meta `mime:"application/json"`
}

// ConfigUpdateReq 更新上传配置请求
type ConfigUpdateReq struct {
	g.Meta       `path:"/config/update" method:"put" tags:"上传配置" summary:"更新上传配置"`
	ID           snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"上传配置ID"`
	Name         string              `json:"name" dc:"配置名称"`
	Storage      int                 `json:"storage" dc:"存储类型"`
	IsDefault    int                 `json:"isDefault" dc:"是否默认"`
	LocalPath    string              `json:"localPath" dc:"本地存储路径"`
	OssEndpoint  string              `json:"ossEndpoint" dc:"OSS Endpoint"`
	OssBucket    string              `json:"ossBucket" dc:"OSS Bucket"`
	OssAccessKey string              `json:"ossAccessKey" dc:"OSS AccessKey"`
	OssSecretKey string              `json:"ossSecretKey" dc:"OSS SecretKey"`
	CosRegion    string              `json:"cosRegion" dc:"COS Region"`
	CosBucket    string              `json:"cosBucket" dc:"COS Bucket"`
	CosSecretID  string              `json:"cosSecretID" dc:"COS SecretId"`
	CosSecretKey string              `json:"cosSecretKey" dc:"COS SecretKey"`
	MaxSize      int                 `json:"maxSize" dc:"最大文件大小(MB)"`
	Status       int                 `json:"status" dc:"状态"`
}

// ConfigUpdateRes 更新上传配置响应
type ConfigUpdateRes struct {
	g.Meta `mime:"application/json"`
}

// ConfigDeleteReq 删除上传配置请求
type ConfigDeleteReq struct {
	g.Meta `path:"/config/delete" method:"delete" tags:"上传配置" summary:"删除上传配置"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"上传配置ID"`
}

// ConfigDeleteRes 删除上传配置响应
type ConfigDeleteRes struct {
	g.Meta `mime:"application/json"`
}

// ConfigDetailReq 获取上传配置详情请求
type ConfigDetailReq struct {
	g.Meta `path:"/config/detail" method:"get" tags:"上传配置" summary:"获取上传配置详情"`
	ID     snowflake.JsonInt64 `json:"id" v:"required#ID不能为空" dc:"上传配置ID"`
}

// ConfigDetailRes 获取上传配置详情响应
type ConfigDetailRes struct {
	g.Meta `mime:"application/json"`
	*model.ConfigDetailOutput
}

// ConfigListReq 获取上传配置列表请求
type ConfigListReq struct {
	g.Meta    `path:"/config/list" method:"get" tags:"上传配置" summary:"获取上传配置列表"`
	PageNum   int  `json:"pageNum" d:"1" dc:"页码"`
	PageSize  int  `json:"pageSize" d:"10" dc:"每页数量"`
	Storage   *int `json:"storage" dc:"存储类型"`
	IsDefault *int `json:"isDefault" dc:"是否默认"`
	Status    *int `json:"status" dc:"状态"`
}

// ConfigListRes 获取上传配置列表响应
type ConfigListRes struct {
	g.Meta `mime:"application/json"`
	List   []*model.ConfigListOutput `json:"list" dc:"列表数据"`
	Total  int                       `json:"total" dc:"总数"`
}
