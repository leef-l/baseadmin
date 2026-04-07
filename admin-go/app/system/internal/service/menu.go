package service

import (
	"context"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IMenu interface {
	Create(ctx context.Context, in *model.MenuCreateInput) error
	Update(ctx context.Context, in *model.MenuUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.MenuDetailOutput, err error)
	List(ctx context.Context, in *model.MenuListInput) (list []*model.MenuListOutput, total int, err error)
	Tree(ctx context.Context, in *model.MenuTreeInput) (tree []*model.MenuTreeOutput, err error)
}

var localMenu IMenu

func Menu() IMenu {
	return localMenu
}

func RegisterMenu(i IMenu) {
	localMenu = i
}
