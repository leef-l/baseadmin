package uploader

import (
	"context"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model/entity"
)

type uploadStorageConfig struct {
	MaxSize     int64
	StorageType int
	LocalPath   string
	OSS         ossConfig
	COS         cosConfig
}

func loadUploadStorageConfig(ctx context.Context) (uploadStorageConfig, error) {
	cfg := uploadStorageConfig{
		MaxSize:     10 * 1024 * 1024,
		StorageType: 1,
		LocalPath:   defaultLocalStoragePath,
	}

	var record *entity.UploadConfig
	err := dao.UploadConfig.Ctx(ctx).
		Where("is_default", 1).
		Where("status", 1).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&record)
	if err != nil || record == nil {
		cfg.LocalPath = normalizeLocalStoragePath(cfg.LocalPath)
		return cfg, err
	}

	if record.MaxSize > 0 {
		cfg.MaxSize = int64(record.MaxSize) * 1024 * 1024
	}
	if record.Storage > 0 {
		cfg.StorageType = record.Storage
	}
	if record.LocalPath != "" {
		cfg.LocalPath = record.LocalPath
	}
	cfg.LocalPath = normalizeLocalStoragePath(cfg.LocalPath)
	cfg.OSS = ossConfig{
		Endpoint:  record.OssEndpoint,
		Bucket:    record.OssBucket,
		AccessKey: record.OssAccessKey,
		SecretKey: record.OssSecretKey,
	}
	cfg.COS = cosConfig{
		Region:    record.CosRegion,
		Bucket:    record.CosBucket,
		SecretId:  record.CosSecretId,
		SecretKey: record.CosSecretKey,
	}
	return cfg, nil
}
