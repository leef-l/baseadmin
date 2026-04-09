package uploader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gerror"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func CheckStorageReady(ctx context.Context) error {
	cfg, err := loadUploadStorageConfig(ctx)
	if err != nil {
		return err
	}
	return checkStorageReady(ctx, cfg)
}

func checkStorageReady(ctx context.Context, cfg uploadStorageConfig) error {
	switch cfg.StorageType {
	case 1:
		return checkLocalStorageReady(cfg.LocalPath)
	case 2:
		if strings.TrimSpace(cfg.OSS.Endpoint) == "" ||
			strings.TrimSpace(cfg.OSS.Bucket) == "" ||
			strings.TrimSpace(cfg.OSS.AccessKey) == "" ||
			strings.TrimSpace(cfg.OSS.SecretKey) == "" {
			return gerror.New("阿里云OSS配置不完整")
		}
		return checkOSSStorageReady(cfg.OSS)
	case 3:
		if strings.TrimSpace(cfg.COS.Region) == "" ||
			strings.TrimSpace(cfg.COS.Bucket) == "" ||
			strings.TrimSpace(cfg.COS.SecretId) == "" ||
			strings.TrimSpace(cfg.COS.SecretKey) == "" {
			return gerror.New("腾讯云COS配置不完整")
		}
		return checkCOSStorageReady(ctx, cfg.COS)
	default:
		return gerror.New("不支持的存储类型")
	}
}

func checkLocalStorageReady(localPath string) error {
	localPath = normalizeLocalStoragePath(localPath)
	if strings.TrimSpace(localPath) == "" {
		return gerror.New("本地存储路径不能为空")
	}
	if err := os.MkdirAll(localPath, 0o755); err != nil {
		return fmt.Errorf("本地存储目录不可创建: %w", err)
	}
	tempFile, err := os.CreateTemp(localPath, ".readyz-*")
	if err != nil {
		return fmt.Errorf("本地存储目录不可写: %w", err)
	}
	_ = tempFile.Close()
	_ = os.Remove(tempFile.Name())
	return nil
}

func checkOSSStorageReady(cfg ossConfig) error {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}
	exists, err := client.IsBucketExist(cfg.Bucket)
	if err != nil {
		return fmt.Errorf("访问OSS Bucket失败: %w", err)
	}
	if !exists {
		return gerror.New("OSS Bucket不存在或不可访问")
	}
	return nil
}

func checkCOSStorageReady(ctx context.Context, cfg cosConfig) error {
	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region)
	u, err := url.Parse(bucketURL)
	if err != nil {
		return fmt.Errorf("解析COS URL失败: %w", err)
	}
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	if _, err := client.Bucket.Head(ctx); err != nil {
		return fmt.Errorf("访问COS Bucket失败: %w", err)
	}
	return nil
}
