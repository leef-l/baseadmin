package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

// UploaderUploadReq 上传文件请求
type UploaderUploadReq struct {
	g.Meta `path:"/uploader/upload" method:"post" mime:"multipart/form-data" tags:"文件上传" summary:"上传文件"`
	DirId  string `json:"dirId" dc:"目录ID"`
}

// UploaderUploadRes 上传文件响应
type UploaderUploadRes struct {
	Id      snowflake.JsonInt64 `json:"id"`
	Url     string              `json:"url"`
	Name    string              `json:"name"`
	Size    int64               `json:"size"`
	Ext     string              `json:"ext"`
	Mime    string              `json:"mime"`
	IsImage int                 `json:"isImage"`
}
