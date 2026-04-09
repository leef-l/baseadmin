package uploader

import (
	"context"

	"gbaseadmin/app/upload/internal/dao"
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

	var record map[string]interface{}
	err := dao.UploadConfig.Ctx(ctx).
		Where("is_default", 1).
		Where("status", 1).
		Where(dao.UploadConfig.Columns().DeletedAt, nil).
		Scan(&record)
	if err != nil || record == nil {
		cfg.LocalPath = normalizeLocalStoragePath(cfg.LocalPath)
		return cfg, err
	}

	if v := getInt64(record, "max_size"); v > 0 {
		cfg.MaxSize = v * 1024 * 1024
	}
	if v := getInt64(record, "storage"); v > 0 {
		cfg.StorageType = int(v)
	}
	if v := getString(record, "local_path"); v != "" {
		cfg.LocalPath = v
	}
	cfg.LocalPath = normalizeLocalStoragePath(cfg.LocalPath)
	cfg.OSS = ossConfig{
		Endpoint:  getString(record, "oss_endpoint"),
		Bucket:    getString(record, "oss_bucket"),
		AccessKey: getString(record, "oss_access_key"),
		SecretKey: getString(record, "oss_secret_key"),
	}
	cfg.COS = cosConfig{
		Region:    getString(record, "cos_region"),
		Bucket:    getString(record, "cos_bucket"),
		SecretId:  getString(record, "cos_secret_id"),
		SecretKey: getString(record, "cos_secret_key"),
	}
	return cfg, nil
}
