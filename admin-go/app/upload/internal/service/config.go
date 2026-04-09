package service

import (
	"context"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IConfig interface {
	Create(ctx context.Context, in *model.ConfigCreateInput) error
	Update(ctx context.Context, in *model.ConfigUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) error
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ConfigDetailOutput, err error)
	List(ctx context.Context, in *model.ConfigListInput) (list []*model.ConfigListOutput, total int, err error)
}

var localConfig IConfig

func Config() IConfig {
	return localConfig
}

func RegisterConfig(i IConfig) {
	localConfig = i
}
