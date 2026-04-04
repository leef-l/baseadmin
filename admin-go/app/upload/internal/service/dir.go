package service

import (
	"context"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IDir interface {
	Create(ctx context.Context, in *model.DirCreateInput) error
	Update(ctx context.Context, in *model.DirUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DirDetailOutput, err error)
	List(ctx context.Context, in *model.DirListInput) (list []*model.DirListOutput, total int, err error)
	Tree(ctx context.Context) (tree []*model.DirTreeOutput, err error)
}

var localDir IDir

func Dir() IDir {
	return localDir
}

func RegisterDir(i IDir) {
	localDir = i
}
