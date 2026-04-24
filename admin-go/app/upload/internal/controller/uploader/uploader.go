package uploader

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

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
		ID:      out.ID,
		URL:     out.URL,
		Name:    out.Name,
		Size:    out.Size,
		Ext:     out.Ext,
		Mime:    out.Mime,
		IsImage: out.IsImage,
	}
	return
}

var PublicUploader = cPublicUploader{}

type cPublicUploader struct{}

func (c *cPublicUploader) UploadByTicket(ctx context.Context, req *v1.UploaderUploadByTicketReq) (res *v1.UploaderUploadByTicketRes, err error) {
	ticket := req.Ticket
	if ticket == "" {
		if request := g.RequestFromCtx(ctx); request != nil {
			ticket = request.Get("ticket").String()
		}
	}
	out, err := service.Uploader().UploadByTicket(ctx, ticket)
	if err != nil {
		return nil, err
	}
	res = &v1.UploaderUploadByTicketRes{
		ID:      out.ID,
		URL:     out.URL,
		Name:    out.Name,
		Size:    out.Size,
		Ext:     out.Ext,
		Mime:    out.Mime,
		IsImage: out.IsImage,
	}
	return
}
