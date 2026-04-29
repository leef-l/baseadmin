package uploader

import (
	"context"
	"fmt"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model/entity"
)

type uploadStorageConfig struct {
	ConfigID    int64
	MaxSize     int64
	StorageType int
	LocalPath   string
	OSS         ossConfig
	COS         cosConfig
}

func loadUploadStorageConfig(ctx context.Context, configIDs ...int64) (uploadStorageConfig, error) {
	cfg := uploadStorageConfig{
		MaxSize:     10 * 1024 * 1024,
		StorageType: 1,
		LocalPath:   defaultLocalStoragePath,
	}

	var configID int64
	if len(configIDs) > 0 {
		configID = configIDs[0]
	}

	var record *entity.UploadConfig
	query := dao.UploadConfig.Ctx(ctx).Where("status", 1)
	query = shared.ApplyTenantScopeToModel(ctx, query, shared.ColumnTenantID, shared.ColumnMerchantID)
	if configID > 0 {
		query = query.Where("id", configID)
	} else {
		query = query.Where("is_default", 1)
	}
	err := query.Scan(&record)
	if err != nil || record == nil {
		cfg.LocalPath = normalizeLocalStoragePath(cfg.LocalPath)
		if err != nil {
			return cfg, err
		}
		if configID > 0 {
			return cfg, fmt.Errorf("上传配置不存在或未启用")
		}
		return cfg, nil
	}

	cfg.ConfigID = int64(record.Id)
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
