package service

import (
	"context"

	"gbaseadmin/app/upload/internal/model"
)

type IUploader interface {
	Upload(ctx context.Context) (*model.UploadOutput, error)
	UploadByTicket(ctx context.Context, ticket string) (*model.UploadOutput, error)
}

var localUploader IUploader

func RegisterUploader(s IUploader) {
	localUploader = s
}

func Uploader() IUploader {
	return localUploader
}
