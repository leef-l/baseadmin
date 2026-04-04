package config

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/snowflake"
)

func init() {
	service.RegisterConfig(New())
}

func New() *sConfig {
	return &sConfig{}
}

type sConfig struct{}

// Create 创建上传配置
func (s *sConfig) Create(ctx context.Context, in *model.ConfigCreateInput) error {
	id := snowflake.Generate()
	_, err := dao.UploadConfig.Ctx(ctx).Data(g.Map{
		dao.UploadConfig.Columns().Id:        id,
		dao.UploadConfig.Columns().Name: in.Name,
		dao.UploadConfig.Columns().Storage: in.Storage,
		dao.UploadConfig.Columns().IsDefault: in.IsDefault,
		dao.UploadConfig.Columns().LocalPath: in.LocalPath,
		dao.UploadConfig.Columns().OssEndpoint: in.OssEndpoint,
		dao.UploadConfig.Columns().OssBucket: in.OssBucket,
		dao.UploadConfig.Columns().OssAccessKey: in.OssAccessKey,
		dao.UploadConfig.Columns().OssSecretKey: in.OssSecretKey,
		dao.UploadConfig.Columns().CosRegion: in.CosRegion,
		dao.UploadConfig.Columns().CosBucket: in.CosBucket,
		dao.UploadConfig.Columns().CosSecretId: in.CosSecretID,
		dao.UploadConfig.Columns().CosSecretKey: in.CosSecretKey,
		dao.UploadConfig.Columns().MaxSize: in.MaxSize,
		dao.UploadConfig.Columns().Status: in.Status,
		dao.UploadConfig.Columns().CreatedAt: gtime.Now(),
		dao.UploadConfig.Columns().UpdatedAt: gtime.Now(),
	}).Insert()
	return err
}

// Update 更新上传配置
func (s *sConfig) Update(ctx context.Context, in *model.ConfigUpdateInput) error {
	data := g.Map{
		dao.UploadConfig.Columns().Name: in.Name,
		dao.UploadConfig.Columns().Storage: in.Storage,
		dao.UploadConfig.Columns().IsDefault: in.IsDefault,
		dao.UploadConfig.Columns().LocalPath: in.LocalPath,
		dao.UploadConfig.Columns().OssEndpoint: in.OssEndpoint,
		dao.UploadConfig.Columns().OssBucket: in.OssBucket,
		dao.UploadConfig.Columns().OssAccessKey: in.OssAccessKey,
		dao.UploadConfig.Columns().OssSecretKey: in.OssSecretKey,
		dao.UploadConfig.Columns().CosRegion: in.CosRegion,
		dao.UploadConfig.Columns().CosBucket: in.CosBucket,
		dao.UploadConfig.Columns().CosSecretId: in.CosSecretID,
		dao.UploadConfig.Columns().CosSecretKey: in.CosSecretKey,
		dao.UploadConfig.Columns().MaxSize: in.MaxSize,
		dao.UploadConfig.Columns().Status: in.Status,
		dao.UploadConfig.Columns().UpdatedAt: gtime.Now(),
	}
	_, err := dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().Id, in.ID).Data(data).Update()
	return err
}

// Delete 软删除上传配置
func (s *sConfig) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	_, err := dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().Id, id).Data(g.Map{
		dao.UploadConfig.Columns().DeletedAt: gtime.Now(),
	}).Update()
	return err
}

// Detail 获取上传配置详情
func (s *sConfig) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ConfigDetailOutput, err error) {
	out = &model.ConfigDetailOutput{}
	err = dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().Id, id).Where(dao.UploadConfig.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	return
}

// List 获取上传配置列表
func (s *sConfig) List(ctx context.Context, in *model.ConfigListInput) (list []*model.ConfigListOutput, total int, err error) {
	m := dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().DeletedAt, nil)
	if in.Storage != nil {
		m = m.Where(dao.UploadConfig.Columns().Storage, *in.Storage)
	}
	if in.IsDefault != nil {
		m = m.Where(dao.UploadConfig.Columns().IsDefault, *in.IsDefault)
	}
	if in.Status != nil {
		m = m.Where(dao.UploadConfig.Columns().Status, *in.Status)
	}
	total, err = m.Count()
	if err != nil {
		return
	}
	err = m.Page(in.PageNum, in.PageSize).OrderAsc(dao.UploadConfig.Columns().Id).Scan(&list)
	if err != nil {
		return
	}
	return
}

