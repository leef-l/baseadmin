package uploader

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// ossConfig 阿里云 OSS 配置
type ossConfig struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
}

// uploadToOSS 上传文件到阿里云 OSS，返回可访问的完整 URL
func uploadToOSS(cfg ossConfig, localFilePath, objectKey string) (string, error) {
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return "", fmt.Errorf("阿里云OSS配置不完整，请检查 endpoint/bucket/access_key/secret_key")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return "", fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return "", fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	f, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("打开本地文件失败: %w", err)
	}
	defer f.Close()

	if err = bucket.PutObject(objectKey, f); err != nil {
		return "", fmt.Errorf("OSS上传失败: %w", err)
	}

	// 拼接 CDN/公网 URL：https://{bucket}.{endpoint}/{objectKey}
	url := fmt.Sprintf("https://%s.%s/%s", cfg.Bucket, cfg.Endpoint, objectKey)
	return url, nil
}

// deleteFromOSS 从阿里云 OSS 删除文件
func deleteFromOSS(cfg ossConfig, objectKey string) error {
	if cfg.Endpoint == "" || cfg.Bucket == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return fmt.Errorf("阿里云OSS配置不完整")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	if err = bucket.DeleteObject(objectKey); err != nil {
		return fmt.Errorf("OSS删除失败: %w", err)
	}
	return nil
}
