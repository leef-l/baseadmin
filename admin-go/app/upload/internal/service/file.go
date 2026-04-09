package service

import (
	"context"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IFile interface {
	Create(ctx context.Context, in *model.FileCreateInput) error
	Update(ctx context.Context, in *model.FileUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.FileDetailOutput, err error)
	List(ctx context.Context, in *model.FileListInput) (list []*model.FileListOutput, total int, err error)
}

var localFile IFile

func File() IFile {
	return localFile
}

func RegisterFile(i IFile) {
	localFile = i
}
