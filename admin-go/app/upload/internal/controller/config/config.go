package config

import (
	"context"

	v1 "gbaseadmin/app/upload/api/upload/v1"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
)

var Config = cConfig{}

type cConfig struct{}

// Create 创建上传配置
func (c *cConfig) Create(ctx context.Context, req *v1.ConfigCreateReq) (res *v1.ConfigCreateRes, err error) {
	err = service.Config().Create(ctx, &model.ConfigCreateInput{
		Name: req.Name,
		Storage: req.Storage,
		IsDefault: req.IsDefault,
		LocalPath: req.LocalPath,
		OssEndpoint: req.OssEndpoint,
		OssBucket: req.OssBucket,
		OssAccessKey: req.OssAccessKey,
		OssSecretKey: req.OssSecretKey,
		CosRegion: req.CosRegion,
		CosBucket: req.CosBucket,
		CosSecretID: req.CosSecretID,
		CosSecretKey: req.CosSecretKey,
		MaxSize: req.MaxSize,
		Status: req.Status,
	})
	return
}

// Update 更新上传配置
func (c *cConfig) Update(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error) {
	err = service.Config().Update(ctx, &model.ConfigUpdateInput{
		ID: req.ID,
		Name: req.Name,
		Storage: req.Storage,
		IsDefault: req.IsDefault,
		LocalPath: req.LocalPath,
		OssEndpoint: req.OssEndpoint,
		OssBucket: req.OssBucket,
		OssAccessKey: req.OssAccessKey,
		OssSecretKey: req.OssSecretKey,
		CosRegion: req.CosRegion,
		CosBucket: req.CosBucket,
		CosSecretID: req.CosSecretID,
		CosSecretKey: req.CosSecretKey,
		MaxSize: req.MaxSize,
		Status: req.Status,
	})
	return
}

// Delete 删除上传配置
func (c *cConfig) Delete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error) {
	err = service.Config().Delete(ctx, req.ID)
	return
}

// Detail 获取上传配置详情
func (c *cConfig) Detail(ctx context.Context, req *v1.ConfigDetailReq) (res *v1.ConfigDetailRes, err error) {
	res = &v1.ConfigDetailRes{}
	res.ConfigDetailOutput, err = service.Config().Detail(ctx, req.ID)
	return
}

// List 获取上传配置列表
func (c *cConfig) List(ctx context.Context, req *v1.ConfigListReq) (res *v1.ConfigListRes, err error) {
	res = &v1.ConfigListRes{}
	res.List, res.Total, err = service.Config().List(ctx, &model.ConfigListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Storage: req.Storage,
		IsDefault: req.IsDefault,
		Status: req.Status,
	})
	return
}

