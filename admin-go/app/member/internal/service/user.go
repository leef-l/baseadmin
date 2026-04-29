package service

import (
	"context"
	"gbaseadmin/app/member/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IUser interface {
	Create(ctx context.Context, in *model.UserCreateInput) error
	Update(ctx context.Context, in *model.UserUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.UserDetailOutput, err error)
	List(ctx context.Context, in *model.UserListInput) (list []*model.UserListOutput, total int, err error)
	Export(ctx context.Context, in *model.UserListInput) (list []*model.UserListOutput, err error)
	Tree(ctx context.Context, in *model.UserTreeInput) (tree []*model.UserTreeOutput, err error)
	BatchUpdate(ctx context.Context, in *model.UserBatchUpdateInput) error
}

var localUser IUser

func User() IUser {
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
