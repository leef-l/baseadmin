package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// Config DTO 模型

// ConfigCreateInput 创建上传配置输入
type ConfigCreateInput struct {
	Name         string `json:"name"`
	Storage      int    `json:"storage"`
	IsDefault    int    `json:"isDefault"`
	LocalPath    string `json:"localPath"`
	OssEndpoint  string `json:"ossEndpoint"`
	OssBucket    string `json:"ossBucket"`
	OssAccessKey string `json:"ossAccessKey"`
	OssSecretKey string `json:"ossSecretKey"`
	CosRegion    string `json:"cosRegion"`
	CosBucket    string `json:"cosBucket"`
	CosSecretID  string `json:"cosSecretID"`
	CosSecretKey string `json:"cosSecretKey"`
	MaxSize      int    `json:"maxSize"`
	Status       int    `json:"status"`
}

// ConfigUpdateInput 更新上传配置输入
type ConfigUpdateInput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Storage      int                 `json:"storage"`
	IsDefault    int                 `json:"isDefault"`
	LocalPath    string              `json:"localPath"`
	OssEndpoint  string              `json:"ossEndpoint"`
	OssBucket    string              `json:"ossBucket"`
	OssAccessKey string              `json:"ossAccessKey"`
	OssSecretKey string              `json:"ossSecretKey"`
	CosRegion    string              `json:"cosRegion"`
	CosBucket    string              `json:"cosBucket"`
	CosSecretID  string              `json:"cosSecretID"`
	CosSecretKey string              `json:"cosSecretKey"`
	MaxSize      int                 `json:"maxSize"`
	Status       int                 `json:"status"`
}

// ConfigDetailOutput 上传配置详情输出
type ConfigDetailOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Storage      int                 `json:"storage"`
	IsDefault    int                 `json:"isDefault"`
	LocalPath    string              `json:"localPath"`
	OssEndpoint  string              `json:"ossEndpoint"`
	OssBucket    string              `json:"ossBucket"`
	OssAccessKey string              `json:"ossAccessKey"`
	OssSecretKey string              `json:"ossSecretKey"`
	CosRegion    string              `json:"cosRegion"`
	CosBucket    string              `json:"cosBucket"`
	CosSecretID  string              `json:"cosSecretID"`
	CosSecretKey string              `json:"cosSecretKey"`
	MaxSize      int                 `json:"maxSize"`
	Status       int                 `json:"status"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// ConfigListOutput 上传配置列表输出
type ConfigListOutput struct {
	ID           snowflake.JsonInt64 `json:"id"`
	Name         string              `json:"name"`
	Storage      int                 `json:"storage"`
	IsDefault    int                 `json:"isDefault"`
	LocalPath    string              `json:"localPath"`
	OssEndpoint  string              `json:"ossEndpoint"`
	OssBucket    string              `json:"ossBucket"`
	OssAccessKey string              `json:"ossAccessKey"`
	OssSecretKey string              `json:"ossSecretKey"`
	CosRegion    string              `json:"cosRegion"`
	CosBucket    string              `json:"cosBucket"`
	CosSecretID  string              `json:"cosSecretID"`
	CosSecretKey string              `json:"cosSecretKey"`
	MaxSize      int                 `json:"maxSize"`
	Status       int                 `json:"status"`
	CreatedAt    *gtime.Time         `json:"createdAt"`
	UpdatedAt    *gtime.Time         `json:"updatedAt"`
}

// ConfigListInput 上传配置列表查询输入
type ConfigListInput struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	Keyword   string `json:"keyword"`
	Storage   *int   `json:"storage"`
	IsDefault *int   `json:"isDefault"`
	Status    *int   `json:"status"`
}
