package service

import (
	"context"

	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/utility/snowflake"
)

type IDaemon interface {
	Create(ctx context.Context, in *model.DaemonCreateInput) error
	Update(ctx context.Context, in *model.DaemonUpdateInput) error
	Delete(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error)
	BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error)
	Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.DaemonDetailOutput, err error)
	List(ctx context.Context, in *model.DaemonListInput) (list []*model.DaemonListOutput, total int, err error)
	Restart(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error)
	BatchRestart(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error)
	Stop(ctx context.Context, id snowflake.JsonInt64) (*model.DaemonOperationOutput, error)
	BatchStop(ctx context.Context, ids []snowflake.JsonInt64) (*model.DaemonBatchOperationOutput, error)
	Log(ctx context.Context, id snowflake.JsonInt64, logType string, lines int) (*model.DaemonLogOutput, error)
}

var localDaemon IDaemon

func Daemon() IDaemon {
	return localDaemon
}

func RegisterDaemon(i IDaemon) {
	localDaemon = i
}
