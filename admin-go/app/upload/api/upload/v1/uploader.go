package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"gbaseadmin/utility/snowflake"
)

type UploaderUploadReq struct {
	g.Meta   `path:"/uploader/upload" method:"post" mime:"multipart/form-data" tags:"文件上传" summary:"上传文件"`
	DirID    string `json:"dirId" dc:"目录ID"`
	ConfigID string `json:"configId" dc:"上传配置ID"`
	Scene    string `json:"scene" dc:"上传场景/接口标识"`
}

type UploaderUploadByTicketReq struct {
	g.Meta `path:"/public/uploader/upload_by_ticket" method:"post" mime:"multipart/form-data" tags:"文件上传" summary:"通过上传票据上传文件"`
	Ticket string `json:"ticket" v:"required#上传票据不能为空" dc:"上传票据"`
}

type UploaderUploadRes struct {
	g.Meta  `mime:"application/json"`
	ID      snowflake.JsonInt64 `json:"id"`
	URL     string              `json:"url"`
	Name    string              `json:"name"`
	Size    int64               `json:"size"`
	Ext     string              `json:"ext"`
	Mime    string              `json:"mime"`
	IsImage int                 `json:"isImage"`
}

type UploaderUploadByTicketRes = UploaderUploadRes
