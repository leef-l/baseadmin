package daemon

import (
	"context"

	v1 "gbaseadmin/app/system/api/system/v1"
	"gbaseadmin/app/system/internal/model"
	"gbaseadmin/app/system/internal/service"
)

var Daemon = cDaemon{}

type cDaemon struct{}

func (c *cDaemon) Create(ctx context.Context, req *v1.DaemonCreateReq) (res *v1.DaemonCreateRes, err error) {
	err = service.Daemon().Create(ctx, &model.DaemonCreateInput{
		Name:         req.Name,
		Program:      req.Program,
		Command:      req.Command,
		Directory:    req.Directory,
		RunUser:      req.RunUser,
		Numprocs:     req.Numprocs,
		Priority:     req.Priority,
		Autostart:    req.Autostart,
		Autorestart:  req.Autorestart,
		Startsecs:    req.Startsecs,
		Startretries: req.Startretries,
		StopSignal:   req.StopSignal,
		Environment:  req.Environment,
		Remark:       req.Remark,
	})
	return
}

func (c *cDaemon) Update(ctx context.Context, req *v1.DaemonUpdateReq) (res *v1.DaemonUpdateRes, err error) {
	err = service.Daemon().Update(ctx, &model.DaemonUpdateInput{
		ID:           req.ID,
		Name:         req.Name,
		Program:      req.Program,
		Command:      req.Command,
		Directory:    req.Directory,
		RunUser:      req.RunUser,
		Numprocs:     req.Numprocs,
		Priority:     req.Priority,
		Autostart:    req.Autostart,
		Autorestart:  req.Autorestart,
		Startsecs:    req.Startsecs,
		Startretries: req.Startretries,
		StopSignal:   req.StopSignal,
		Environment:  req.Environment,
		Remark:       req.Remark,
	})
	return
}

func (c *cDaemon) Delete(ctx context.Context, req *v1.DaemonDeleteReq) (res *v1.DaemonDeleteRes, err error) {
	_, err = service.Daemon().Delete(ctx, req.ID)
	return
}

func (c *cDaemon) BatchDelete(ctx context.Context, req *v1.DaemonBatchDeleteReq) (res *v1.DaemonBatchDeleteRes, err error) {
	res = &v1.DaemonBatchDeleteRes{}
	res.DaemonBatchOperationOutput, err = service.Daemon().BatchDelete(ctx, req.IDs)
	return
}

func (c *cDaemon) Detail(ctx context.Context, req *v1.DaemonDetailReq) (res *v1.DaemonDetailRes, err error) {
	res = &v1.DaemonDetailRes{}
	res.DaemonDetailOutput, err = service.Daemon().Detail(ctx, req.ID)
	return
}

func (c *cDaemon) List(ctx context.Context, req *v1.DaemonListReq) (res *v1.DaemonListRes, err error) {
	res = &v1.DaemonListRes{}
	res.List, res.Total, err = service.Daemon().List(ctx, &model.DaemonListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Program:  req.Program,
	})
	return
}

func (c *cDaemon) Restart(ctx context.Context, req *v1.DaemonRestartReq) (res *v1.DaemonRestartRes, err error) {
	res = &v1.DaemonRestartRes{}
	res.DaemonOperationOutput, err = service.Daemon().Restart(ctx, req.ID)
	return
}

func (c *cDaemon) BatchRestart(ctx context.Context, req *v1.DaemonBatchRestartReq) (res *v1.DaemonBatchRestartRes, err error) {
	res = &v1.DaemonBatchRestartRes{}
	res.DaemonBatchOperationOutput, err = service.Daemon().BatchRestart(ctx, req.IDs)
	return
}

func (c *cDaemon) Stop(ctx context.Context, req *v1.DaemonStopReq) (res *v1.DaemonStopRes, err error) {
	res = &v1.DaemonStopRes{}
	res.DaemonOperationOutput, err = service.Daemon().Stop(ctx, req.ID)
	return
}

func (c *cDaemon) BatchStop(ctx context.Context, req *v1.DaemonBatchStopReq) (res *v1.DaemonBatchStopRes, err error) {
	res = &v1.DaemonBatchStopRes{}
	res.DaemonBatchOperationOutput, err = service.Daemon().BatchStop(ctx, req.IDs)
	return
}

func (c *cDaemon) Log(ctx context.Context, req *v1.DaemonLogReq) (res *v1.DaemonLogRes, err error) {
	res = &v1.DaemonLogRes{}
	res.DaemonLogOutput, err = service.Daemon().Log(ctx, req.ID, req.LogType, req.Lines)
	return
}
