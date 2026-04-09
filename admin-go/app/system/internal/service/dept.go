package service

import (
	"context"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IDept interface {
	Create(ctx context.Context, in *model.DeptCreateInput) error
	Update(ctx context.Context, in *model.DeptUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DeptDetailOutput, err error)
	List(ctx context.Context, in *model.DeptListInput) (list []*model.DeptListOutput, total int, err error)
	Tree(ctx context.Context, in *model.DeptTreeInput) (tree []*model.DeptTreeOutput, err error)
}

var localDept IDept

func Dept() IDept {
	return localDept
}

func RegisterDept(i IDept) {
	localDept = i
}
