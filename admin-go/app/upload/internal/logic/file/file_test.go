package file

import (
	"testing"

	"gbaseadmin/app/upload/internal/model/entity"
)

func TestMatchOSSObjectKey(t *testing.T) {
	config := &entity.UploadConfig{
		OssBucket:   "demo-bucket",
		OssEndpoint: "oss-cn-shanghai.aliyuncs.com",
	}

	objectKey, ok := matchOSSObjectKey("https://demo-bucket.oss-cn-shanghai.aliyuncs.com/2026-04-07/demo.png", config)
	if !ok {
		t.Fatal("expected OSS url to match config")
	}
	if objectKey != "2026-04-07/demo.png" {
		t.Fatalf("unexpected objectKey: %s", objectKey)
	}
}

func TestMatchCOSObjectKey(t *testing.T) {
	config := &entity.UploadConfig{
		CosBucket: "demo-1250000000",
		CosRegion: "ap-guangzhou",
	}

	objectKey, ok := matchCOSObjectKey("https://demo-1250000000.cos.ap-guangzhou.myqcloud.com/archive/file.pdf", config)
	if !ok {
		t.Fatal("expected COS url to match config")
	}
	if objectKey != "archive/file.pdf" {
		t.Fatalf("unexpected objectKey: %s", objectKey)
	}
}

func TestMatchUploadConfigByURL(t *testing.T) {
	configs := []*entity.UploadConfig{
		{
			Id:          1,
			Storage:     2,
			OssBucket:   "bucket-a",
			OssEndpoint: "oss-cn-beijing.aliyuncs.com",
		},
		{
			Id:          2,
			Storage:     2,
			OssBucket:   "bucket-b",
			OssEndpoint: "oss-cn-hangzhou.aliyuncs.com",
		},
	}

	config, objectKey := matchUploadConfigByURL(configs, 2, "https://bucket-b.oss-cn-hangzhou.aliyuncs.com/path/to/file.txt")
	if config == nil {
		t.Fatal("expected to find matching config")
	}
	if config.Id != 2 {
		t.Fatalf("expected config 2, got %d", config.Id)
	}
	if objectKey != "path/to/file.txt" {
		t.Fatalf("unexpected objectKey: %s", objectKey)
	}
}
