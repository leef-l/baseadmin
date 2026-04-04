package model

import (
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/utility/snowflake"
)

// File DTO 模型

// FileCreateInput 创建文件记录输入
type FileCreateInput struct {
	DirID   snowflake.JsonInt64 `json:"dirID"`
	Name    string              `json:"name"`
	URL     string              `json:"url"`
	Ext     string              `json:"ext"`
	Size    int64               `json:"size"`
	Mime    string              `json:"mime"`
	Storage int                 `json:"storage"`
	IsImage int                 `json:"isImage"`
}

// FileUpdateInput 更新文件记录输入
type FileUpdateInput struct {
	ID      snowflake.JsonInt64 `json:"id"`
	DirID   snowflake.JsonInt64 `json:"dirID"`
	Name    string              `json:"name"`
	URL     string              `json:"url"`
	Ext     string              `json:"ext"`
	Size    int64               `json:"size"`
	Mime    string              `json:"mime"`
	Storage int                 `json:"storage"`
	IsImage int                 `json:"isImage"`
}

// FileDetailOutput 文件记录详情输出
type FileDetailOutput struct {
	ID        snowflake.JsonInt64 `json:"id"`
	DirID     snowflake.JsonInt64 `json:"dirID"`
	DirName   string              `json:"dirName"`
	Name      string              `json:"name"`
	URL       string              `json:"url"`
	Ext       string              `json:"ext"`
	Size      int64               `json:"size"`
	Mime      string              `json:"mime"`
	Storage   int                 `json:"storage"`
	IsImage   int                 `json:"isImage"`
	CreatedAt *gtime.Time         `json:"createdAt"`
	UpdatedAt *gtime.Time         `json:"updatedAt"`
}

// FileListOutput 文件记录列表输出
type FileListOutput struct {
	ID        snowflake.JsonInt64 `json:"id"`
	DirID     snowflake.JsonInt64 `json:"dirID"`
	DirName   string              `json:"dirName"`
	Name      string              `json:"name"`
	URL       string              `json:"url"`
	Ext       string              `json:"ext"`
	Size      int64               `json:"size"`
	Mime      string              `json:"mime"`
	Storage   int                 `json:"storage"`
	IsImage   int                 `json:"isImage"`
	CreatedAt *gtime.Time         `json:"createdAt"`
	UpdatedAt *gtime.Time         `json:"updatedAt"`
}

// FileListInput 文件记录列表查询输入
type FileListInput struct {
	PageNum  int                 `json:"pageNum"`
	PageSize int                 `json:"pageSize"`
	DirID    snowflake.JsonInt64 `json:"dirID"`
	Name     string              `json:"name"`
	Storage  int                 `json:"storage"`
	IsImage  int                 `json:"isImage"`
}

// UploadOutput 上传文件输出
type UploadOutput struct {
	ID      snowflake.JsonInt64 `json:"id"`
	URL     string              `json:"url"`
	Name    string              `json:"name"`
	Size    int64               `json:"size"`
	Ext     string              `json:"ext"`
	Mime    string              `json:"mime"`
	IsImage int                 `json:"isImage"`
}
