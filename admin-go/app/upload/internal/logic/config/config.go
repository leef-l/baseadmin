package config

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model"
	"gbaseadmin/app/upload/internal/model/entity"
	"gbaseadmin/app/upload/internal/service"
	"gbaseadmin/utility/batchutil"
	"gbaseadmin/utility/fieldvalid"
	"gbaseadmin/utility/inpututil"
	"gbaseadmin/utility/pageutil"
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
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeConfigCreateInput(in)
	if err := validateConfigName(in.Name); err != nil {
		return err
	}
	if err := validateConfigMeta(in.Storage, in.IsDefault, in.MaxSize, in.Status); err != nil {
		return err
	}
	if err := s.ensureConfigNameUnique(ctx, 0, in.Name); err != nil {
		return err
	}
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
	if err := inpututil.Require(in); err != nil {
		return err
	}
	normalizeConfigUpdateInput(in)
	if err := validateConfigName(in.Name); err != nil {
		return err
	}
	if err := validateConfigMeta(in.Storage, in.IsDefault, in.MaxSize, in.Status); err != nil {
		return err
	}
	current, err := s.getConfigByID(ctx, in.ID)
	if err != nil {
		return err
	}
	if err := s.ensureConfigNameUnique(ctx, in.ID, in.Name); err != nil {
		return err
	}
	now := gtime.Now()
	return dao.UploadConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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
		cosSecretID := pickSensitiveValue(in.CosSecretID, current.CosSecretId)
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
			Where(dao.UploadConfig.Columns().DeletedAt, nil).
			Data(data).
			Update()
		return err
	})
}

// Delete 软删除上传配置
func (s *sConfig) Delete(ctx context.Context, id snowflake.JsonInt64) error {
	current, err := s.getConfigByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.ensureConfigDeletable(ctx, current); err != nil {
		return err
	}
	_, err = dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Id, id).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Data(g.Map{
			dao.UploadConfig.Columns().DeletedAt: gtime.Now(),
		}).
		Update()
	return err
}

// BatchDelete 批量删除上传配置
func (s *sConfig) BatchDelete(ctx context.Context, ids []snowflake.JsonInt64) error {
	ids = batchutil.CompactIDs(ids)
	if len(ids) == 0 {
		return gerror.New("请选择要删除的上传配置")
	}
	deleteIDs := batchutil.ToInt64s(ids)
	configs, err := s.listConfigsByIDs(ctx, deleteIDs)
	if err != nil {
		return err
	}
	if len(configs) != len(deleteIDs) {
		return gerror.New("包含不存在或已删除的上传配置")
	}
	for _, cfg := range configs {
		if err := s.ensureConfigDeletable(ctx, cfg); err != nil {
			return err
		}
	}
	return dao.UploadConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.UploadConfig.Table()).Ctx(ctx).
			WhereIn(dao.UploadConfig.Columns().Id, deleteIDs).
			Where(dao.UploadConfig.Columns().DeletedAt, nil).
			Data(g.Map{
				dao.UploadConfig.Columns().DeletedAt: gtime.Now(),
			}).
			Update()
		return err
	})
}

func (s *sConfig) getConfigByID(ctx context.Context, id snowflake.JsonInt64) (*entity.UploadConfig, error) {
	if id <= 0 {
		return nil, gerror.New("上传配置不存在或已删除")
	}
	var cfg *entity.UploadConfig
	if err := dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Id, id).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&cfg); err != nil {
		return nil, err
	}
	if cfg == nil || cfg.Id == 0 {
		return nil, gerror.New("上传配置不存在或已删除")
	}
	return cfg, nil
}

func (s *sConfig) listConfigsByIDs(ctx context.Context, ids []int64) ([]*entity.UploadConfig, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var configs []*entity.UploadConfig
	if err := dao.UploadConfig.Ctx(ctx).
		WhereIn(dao.UploadConfig.Columns().Id, ids).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *sConfig) ensureConfigDeletable(ctx context.Context, cfg *entity.UploadConfig) error {
	if cfg == nil || cfg.Id == 0 {
		return gerror.New("上传配置不存在或已删除")
	}
	if cfg.IsDefault == 1 {
		return gerror.New("默认上传配置不能删除，请先设置其他配置为默认")
	}
	refCount, err := s.countActiveFileReferences(ctx, cfg)
	if err != nil {
		return err
	}
	if refCount > 0 {
		return gerror.New("当前上传配置仍有关联文件，不能直接删除")
	}
	return nil
}

func (s *sConfig) ensureConfigNameUnique(ctx context.Context, currentID snowflake.JsonInt64, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	m := dao.UploadConfig.Ctx(ctx).
		Where(dao.UploadConfig.Columns().Name, name).
		Where(dao.UploadConfig.Columns().DeletedAt, nil)
	if currentID > 0 {
		m = m.WhereNot(dao.UploadConfig.Columns().Id, currentID)
	}
	count, err := m.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.New("配置名称已存在")
	}
	return nil
}

func (s *sConfig) countActiveFileReferences(ctx context.Context, cfg *entity.UploadConfig) (int, error) {
	if cfg == nil || cfg.Id == 0 {
		return 0, nil
	}
	if cfg.Storage != 2 && cfg.Storage != 3 {
		return 0, nil
	}
	var files []struct {
		Url string `json:"url"`
	}
	if err := dao.UploadFile.Ctx(ctx).
		Fields(dao.UploadFile.Columns().Url).
		Where(dao.UploadFile.Columns().Storage, cfg.Storage).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		Scan(&files); err != nil {
		return 0, err
	}
	refCount := 0
	for _, file := range files {
		if configMatchesFileURL(cfg, file.Url) {
			refCount++
		}
	}
	return refCount, nil
}

func configMatchesFileURL(cfg *entity.UploadConfig, fileURL string) bool {
	if cfg == nil || strings.TrimSpace(fileURL) == "" {
		return false
	}
	switch cfg.Storage {
	case 2:
		_, ok := shared.MatchOSSObjectKey(fileURL, cfg)
		return ok
	case 3:
		_, ok := shared.MatchCOSObjectKey(fileURL, cfg)
		return ok
	default:
		return false
	}
}

// Detail 获取上传配置详情
func (s *sConfig) Detail(ctx context.Context, id snowflake.JsonInt64) (out *model.ConfigDetailOutput, err error) {
	if id <= 0 {
		return nil, gerror.New("上传配置不存在或已删除")
	}
	out = &model.ConfigDetailOutput{}
	err = dao.UploadConfig.Ctx(ctx).Where(dao.UploadConfig.Columns().Id, id).Where(dao.UploadConfig.Columns().DeletedAt, nil).Scan(out)
	if err != nil {
		return nil, err
	}
	if out == nil || out.ID == 0 {
		return nil, gerror.New("上传配置不存在或已删除")
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
	in.PageNum, in.PageSize = pageutil.Normalize(in.PageNum, in.PageSize)
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
	in.LocalPath = normalizeOptionalLocalPath(in.LocalPath)
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
	in.LocalPath = normalizeOptionalLocalPath(in.LocalPath)
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

func validateConfigName(name string) error {
	if name == "" {
		return gerror.New("配置名称不能为空")
	}
	return nil
}

func validateConfigMeta(storage, isDefault, maxSize, status int) error {
	if err := fieldvalid.Enum("存储类型", storage, 1, 2, 3); err != nil {
		return err
	}
	if err := fieldvalid.Binary("是否默认", isDefault); err != nil {
		return err
	}
	if maxSize <= 0 {
		return gerror.New("最大文件大小必须大于0")
	}
	if err := fieldvalid.Binary("状态", status); err != nil {
		return err
	}
	return nil
}

func normalizeOptionalLocalPath(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	return shared.NormalizeLocalStoragePath(path)
}
