package uploader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// cosConfig 腾讯云 COS 配置
type cosConfig struct {
	Region    string
	Bucket    string
	SecretId  string
	SecretKey string
}

// newCOSClient 创建 COS 客户端
func newCOSClient(cfg cosConfig) (*cos.Client, error) {
	if cfg.Region == "" || cfg.Bucket == "" || cfg.SecretId == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("腾讯云COS配置不完整，请检查 region/bucket/secret_id/secret_key")
	}

	// 格式：https://{bucket}.cos.{region}.myqcloud.com
	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region)
	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, fmt.Errorf("解析COS URL失败: %w", err)
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	return client, nil
}

// uploadToCOS 上传文件到腾讯云 COS，返回可访问的完整 URL
func uploadToCOS(ctx context.Context, cfg cosConfig, localFilePath, objectKey string) (string, error) {
	client, err := newCOSClient(cfg)
	if err != nil {
		return "", err
	}

	f, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("打开本地文件失败: %w", err)
	}
	defer f.Close()

	_, err = client.Object.Put(ctx, objectKey, f, nil)
	if err != nil {
		return "", fmt.Errorf("COS上传失败: %w", err)
	}

	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", cfg.Bucket, cfg.Region, objectKey)
	return url, nil
}

// deleteFromCOS 从腾讯云 COS 删除文件
func deleteFromCOS(ctx context.Context, cfg cosConfig, objectKey string) error {
	client, err := newCOSClient(cfg)
	if err != nil {
		return err
	}

	_, err = client.Object.Delete(ctx, objectKey, nil)
	if err != nil {
		return fmt.Errorf("COS删除失败: %w", err)
	}
	return nil
}
