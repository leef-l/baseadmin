package config

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
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
	normalizeConfigCreateInput(in)
	if err := validateConfigFields(
		in.Storage,
		in.LocalPath,
		in.OssEndpoint,
		in.OssBucket,
		in.OssAccessKey,
		in.OssSecretKey,
		in.CosRegion,
		in.CosBucket,
		in.CosSecretID,
		in.CosSecretKey,
	); err != nil {
		return err
	}
	id := snowflake.Generate()
	now := gtime.Now()
	return dao.UploadConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if in.IsDefault == 1 {
			if _, err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).
				Where(dao.UploadConfig.Columns().DeletedAt, nil).
				Data(g.Map{
					dao.UploadConfig.Columns().IsDefault: 0,
					dao.UploadConfig.Columns().UpdatedAt: now,
				}).
				Update(); err != nil {
				return err
			}
		}
		_, err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).Data(g.Map{
			dao.UploadConfig.Columns().Id:           id,
			dao.UploadConfig.Columns().Name:         in.Name,
			dao.UploadConfig.Columns().Storage:      in.Storage,
			dao.UploadConfig.Columns().IsDefault:    in.IsDefault,
			dao.UploadConfig.Columns().LocalPath:    in.LocalPath,
			dao.UploadConfig.Columns().OssEndpoint:  in.OssEndpoint,
			dao.UploadConfig.Columns().OssBucket:    in.OssBucket,
			dao.UploadConfig.Columns().OssAccessKey: in.OssAccessKey,
			dao.UploadConfig.Columns().OssSecretKey: in.OssSecretKey,
			dao.UploadConfig.Columns().CosRegion:    in.CosRegion,
			dao.UploadConfig.Columns().CosBucket:    in.CosBucket,
			dao.UploadConfig.Columns().CosSecretId:  in.CosSecretID,
			dao.UploadConfig.Columns().CosSecretKey: in.CosSecretKey,
			dao.UploadConfig.Columns().MaxSize:      in.MaxSize,
			dao.UploadConfig.Columns().Status:       in.Status,
			dao.UploadConfig.Columns().CreatedAt:    now,
			dao.UploadConfig.Columns().UpdatedAt:    now,
		}).Insert()
		return err
	})
}

// Update 更新上传配置
func (s *sConfig) Update(ctx context.Context, in *model.ConfigUpdateInput) error {
	normalizeConfigUpdateInput(in)
	now := gtime.Now()
	return dao.UploadConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var current struct {
			IsDefault    int    `json:"isDefault"`
			OssAccessKey string `json:"ossAccessKey"`
			OssSecretKey string `json:"ossSecretKey"`
			CosSecretID  string `json:"cosSecretID"`
			CosSecretKey string `json:"cosSecretKey"`
		}
		if err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).
			Where(dao.UploadConfig.Columns().Id, in.ID).
			Where(dao.UploadConfig.Columns().DeletedAt, nil).
			Scan(&current); err != nil {
			return err
		}
		if current.IsDefault == 1 && in.IsDefault != 1 {
			return gerror.New("默认上传配置不能直接取消默认，请先设置其他配置为默认")
		}
		if in.IsDefault == 1 {
			if _, err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).
				Where(dao.UploadConfig.Columns().DeletedAt, nil).
				WhereNot(dao.UploadConfig.Columns().Id, in.ID).
				Data(g.Map{
					dao.UploadConfig.Columns().IsDefault: 0,
					dao.UploadConfig.Columns().UpdatedAt: now,
				}).
				Update(); err != nil {
				return err
			}
		}
		ossAccessKey := pickSensitiveValue(in.OssAccessKey, current.OssAccessKey)
		ossSecretKey := pickSensitiveValue(in.OssSecretKey, current.OssSecretKey)
		cosSecretID := pickSensitiveValue(in.CosSecretID, current.CosSecretID)
		cosSecretKey := pickSensitiveValue(in.CosSecretKey, current.CosSecretKey)
		if err := validateConfigFields(
			in.Storage,
			in.LocalPath,
			in.OssEndpoint,
			in.OssBucket,
			ossAccessKey,
			ossSecretKey,
			in.CosRegion,
			in.CosBucket,
			cosSecretID,
			cosSecretKey,
		); err != nil {
			return err
		}
		data := g.Map{
			dao.UploadConfig.Columns().Name:         in.Name,
			dao.UploadConfig.Columns().Storage:      in.Storage,
			dao.UploadConfig.Columns().IsDefault:    in.IsDefault,
			dao.UploadConfig.Columns().LocalPath:    in.LocalPath,
			dao.UploadConfig.Columns().OssEndpoint:  in.OssEndpoint,
			dao.UploadConfig.Columns().OssBucket:    in.OssBucket,
			dao.UploadConfig.Columns().OssAccessKey: ossAccessKey,
			dao.UploadConfig.Columns().OssSecretKey: ossSecretKey,
			dao.UploadConfig.Columns().CosRegion:    in.CosRegion,
			dao.UploadConfig.Columns().CosBucket:    in.CosBucket,
			dao.UploadConfig.Columns().CosSecretId:  cosSecretID,
			dao.UploadConfig.Columns().CosSecretKey: cosSecretKey,
			dao.UploadConfig.Columns().MaxSize:      in.MaxSize,
			dao.UploadConfig.Columns().Status:       in.Status,
			dao.UploadConfig.Columns().UpdatedAt:    now,
		}
		_, err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).
			Where(dao.UploadConfig.Columns().Id, in.ID).
			Data(data).
			Update()
		return err
	})
}

// Delete 软删除上传配置
func (s *sConfig) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	var current struct {
		IsDefault int `json:"isDefault"`
	}
	if err := dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Id, id).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&current); err != nil {
		return err
	}
	if current.IsDefault == 1 {
		return gerror.New("默认上传配置不能删除，请先设置其他配置为默认")
	}
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
	sanitizeConfigOutput(out)
	return
}

// List 获取上传配置列表
func (s *sConfig) List(ctx context.Context, in *model.ConfigListInput) (list []*model.ConfigListOutput, total int, err error) {
	if in == nil {
		in = &model.ConfigListInput{}
	}
	normalizeConfigListInput(in)
	m := dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().DeletedAt, nil)
	if in.Keyword != "" {
		keywordBuilder := m.Builder().
			WhereLike(dao.UploadConfig.Columns().Name, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadConfig.Columns().LocalPath, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadConfig.Columns().OssEndpoint, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadConfig.Columns().OssBucket, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadConfig.Columns().CosRegion, "%"+in.Keyword+"%").
			WhereOrLike(dao.UploadConfig.Columns().CosBucket, "%"+in.Keyword+"%")
		m = m.Where(keywordBuilder)
	}
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
	for _, item := range list {
		sanitizeConfigOutput(item)
	}
	return
}

func pickSensitiveValue(input, fallback string) string {
	if strings.TrimSpace(input) == "" {
		return fallback
	}
	return strings.TrimSpace(input)
}

func sanitizeConfigOutput(v any) {
	switch out := v.(type) {
	case *model.ConfigDetailOutput:
		if out == nil {
			return
		}
		out.OssAccessKey = ""
		out.OssSecretKey = ""
		out.CosSecretID = ""
		out.CosSecretKey = ""
	case *model.ConfigListOutput:
		if out == nil {
			return
		}
		out.OssAccessKey = ""
		out.OssSecretKey = ""
		out.CosSecretID = ""
		out.CosSecretKey = ""
	}
}

func validateConfigFields(storage int, localPath, ossEndpoint, ossBucket, ossAccessKey, ossSecretKey, cosRegion, cosBucket, cosSecretID, cosSecretKey string) error {
	switch storage {
	case 1:
		if strings.TrimSpace(localPath) == "" {
			return gerror.New("本地存储必须填写存储路径")
		}
	case 2:
		if strings.TrimSpace(ossEndpoint) == "" ||
			strings.TrimSpace(ossBucket) == "" ||
			strings.TrimSpace(ossAccessKey) == "" ||
			strings.TrimSpace(ossSecretKey) == "" {
			return gerror.New("阿里云OSS配置不完整")
		}
	case 3:
		if strings.TrimSpace(cosRegion) == "" ||
			strings.TrimSpace(cosBucket) == "" ||
			strings.TrimSpace(cosSecretID) == "" ||
			strings.TrimSpace(cosSecretKey) == "" {
			return gerror.New("腾讯云COS配置不完整")
		}
	default:
		return gerror.New("不支持的存储类型")
	}
	return nil
}

func normalizeConfigCreateInput(in *model.ConfigCreateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.LocalPath = strings.TrimSpace(in.LocalPath)
	in.OssEndpoint = strings.TrimSpace(in.OssEndpoint)
	in.OssBucket = strings.TrimSpace(in.OssBucket)
	in.OssAccessKey = strings.TrimSpace(in.OssAccessKey)
	in.OssSecretKey = strings.TrimSpace(in.OssSecretKey)
	in.CosRegion = strings.TrimSpace(in.CosRegion)
	in.CosBucket = strings.TrimSpace(in.CosBucket)
	in.CosSecretID = strings.TrimSpace(in.CosSecretID)
	in.CosSecretKey = strings.TrimSpace(in.CosSecretKey)
}

func normalizeConfigUpdateInput(in *model.ConfigUpdateInput) {
	if in == nil {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.LocalPath = strings.TrimSpace(in.LocalPath)
	in.OssEndpoint = strings.TrimSpace(in.OssEndpoint)
	in.OssBucket = strings.TrimSpace(in.OssBucket)
	in.OssAccessKey = strings.TrimSpace(in.OssAccessKey)
	in.OssSecretKey = strings.TrimSpace(in.OssSecretKey)
	in.CosRegion = strings.TrimSpace(in.CosRegion)
	in.CosBucket = strings.TrimSpace(in.CosBucket)
	in.CosSecretID = strings.TrimSpace(in.CosSecretID)
	in.CosSecretKey = strings.TrimSpace(in.CosSecretKey)
}

func normalizeConfigListInput(in *model.ConfigListInput) {
	if in == nil {
		return
	}
	in.Keyword = strings.TrimSpace(in.Keyword)
}
