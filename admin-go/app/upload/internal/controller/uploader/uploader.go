package uploader

import (
	"context"

	v1 "gbaseadmin/app/upload/api/upload/v1"
	"gbaseadmin/app/upload/internal/service"
)

var Uploader = cUploader{}

type cUploader struct{}

func (c *cUploader) Upload(ctx context.Context, req *v1.UploaderUploadReq) (res *v1.UploaderUploadRes, err error) {
	out, err := service.Uploader().Upload(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.UploaderUploadRes{
		Id:      out.ID,
		Url:     out.URL,
		Name:    out.Name,
		Size:    out.Size,
		Ext:     out.Ext,
		Mime:    out.Mime,
		IsImage: out.IsImage,
	}
	return
}
