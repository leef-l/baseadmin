package file

import (
	"testing"

	"gbaseadmin/app/upload/internal/model"
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

func TestMatchOSSObjectKeySupportsPortAndQuery(t *testing.T) {
	config := &entity.UploadConfig{
		OssBucket:   "demo-bucket",
		OssEndpoint: "oss-cn-shanghai.aliyuncs.com",
	}

	objectKey, ok := matchOSSObjectKey("https://demo-bucket.oss-cn-shanghai.aliyuncs.com:443/2026-04-07/demo.png?versionId=1", config)
	if !ok {
		t.Fatal("expected OSS url with port/query to match config")
	}
	if objectKey != "2026-04-07/demo.png" {
		t.Fatalf("unexpected objectKey: %s", objectKey)
	}
}

func TestMatchOSSObjectKeyDecodesEscapedPathAndTrimmedConfig(t *testing.T) {
	config := &entity.UploadConfig{
		OssBucket:   " demo-bucket ",
		OssEndpoint: " oss-cn-shanghai.aliyuncs.com ",
	}

	objectKey, ok := matchOSSObjectKey("https://demo-bucket.oss-cn-shanghai.aliyuncs.com/2026-04-07/demo%20file.png", config)
	if !ok {
		t.Fatal("expected escaped OSS url to match config")
	}
	if objectKey != "2026-04-07/demo file.png" {
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

func TestMatchCOSObjectKeySupportsPortAndQuery(t *testing.T) {
	config := &entity.UploadConfig{
		CosBucket: "demo-1250000000",
		CosRegion: "ap-guangzhou",
	}

	objectKey, ok := matchCOSObjectKey("https://demo-1250000000.cos.ap-guangzhou.myqcloud.com:443/archive/file.pdf?sign=1", config)
	if !ok {
		t.Fatal("expected COS url with port/query to match config")
	}
	if objectKey != "archive/file.pdf" {
		t.Fatalf("unexpected objectKey: %s", objectKey)
	}
}

func TestMatchCOSObjectKeyDecodesEscapedPathAndTrimmedConfig(t *testing.T) {
	config := &entity.UploadConfig{
		CosBucket: " demo-1250000000 ",
		CosRegion: " ap-guangzhou ",
	}

	objectKey, ok := matchCOSObjectKey("https://demo-1250000000.cos.ap-guangzhou.myqcloud.com/archive/demo%20file.pdf", config)
	if !ok {
		t.Fatal("expected escaped COS url to match config")
	}
	if objectKey != "archive/demo file.pdf" {
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

func TestMatchUploadConfigByURLRejectsInvalidURL(t *testing.T) {
	configs := []*entity.UploadConfig{
		{
			Id:          1,
			Storage:     2,
			OssBucket:   "bucket-a",
			OssEndpoint: "oss-cn-beijing.aliyuncs.com",
		},
	}

	config, objectKey := matchUploadConfigByURL(configs, 2, "://bad-url")
	if config != nil || objectKey != "" {
		t.Fatalf("invalid URL should not match any config, got config=%v objectKey=%q", config, objectKey)
	}
}

func TestGetStrSupportsCommonValueTypes(t *testing.T) {
	values := map[string]interface{}{
		"name": []byte("demo"),
		"size": 123,
		"nil":  nil,
	}

	if got := getStr(values, "name"); got != "demo" {
		t.Fatalf("getStr bytes mismatch: %q", got)
	}
	if got := getStr(values, "size"); got != "123" {
		t.Fatalf("getStr fmt mismatch: %q", got)
	}
	if got := getStr(values, "nil"); got != "" {
		t.Fatalf("getStr nil mismatch: %q", got)
	}
}

func TestNormalizeFileInputs(t *testing.T) {
	createIn := &model.FileCreateInput{
		Name: " demo.png ",
		URL:  " /upload/demo.png ",
		Ext:  " .png ",
		Mime: " image/png ",
	}
	normalizeFileCreateInput(createIn)
	if createIn.Name != "demo.png" || createIn.URL != "/upload/demo.png" || createIn.Ext != ".png" || createIn.Mime != "image/png" {
		t.Fatalf("normalizeFileCreateInput mismatch: %+v", createIn)
	}

	updateIn := &model.FileUpdateInput{
		Name: " report.pdf ",
		URL:  " https://example.com/report.pdf ",
		Ext:  " .pdf ",
		Mime: " application/pdf ",
	}
	normalizeFileUpdateInput(updateIn)
	if updateIn.Name != "report.pdf" || updateIn.URL != "https://example.com/report.pdf" || updateIn.Ext != ".pdf" || updateIn.Mime != "application/pdf" {
		t.Fatalf("normalizeFileUpdateInput mismatch: %+v", updateIn)
	}

	listIn := &model.FileListInput{
		Keyword: " demo ",
		Name:    " logo ",
	}
	normalizeFileListInput(listIn)
	if listIn.Keyword != "demo" || listIn.Name != "logo" {
		t.Fatalf("normalizeFileListInput mismatch: %+v", listIn)
	}
}
