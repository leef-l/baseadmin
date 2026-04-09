package uploader

import (
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
)

type cleanupHook func(context.Context) error

type storeRequest struct {
	DateDir       string
	LocalFilePath string
	ObjectKey     string
	UniqueName    string
}

type storeResult struct {
	FileURL    string
	OnCommit   cleanupHook
	OnRollback cleanupHook
}

type storageProvider interface {
	Store(ctx context.Context, req storeRequest) (*storeResult, error)
}

func newStorageProvider(cfg uploadStorageConfig) storageProvider {
	switch cfg.StorageType {
	case 2:
		return ossStorageProvider{cfg: cfg.OSS}
	case 3:
		return cosStorageProvider{cfg: cfg.COS}
	default:
		return localStorageProvider{}
	}
}

type localStorageProvider struct{}

func (p localStorageProvider) Store(ctx context.Context, req storeRequest) (*storeResult, error) {
	return &storeResult{
		FileURL: buildLocalFileURL(req.DateDir, req.UniqueName),
		OnRollback: func(context.Context) error {
			return removeLocalFile(req.LocalFilePath)
		},
	}, nil
}

type ossStorageProvider struct {
	cfg ossConfig
}

func (p ossStorageProvider) Store(ctx context.Context, req storeRequest) (*storeResult, error) {
	fileURL, err := uploadToOSS(p.cfg, req.LocalFilePath, req.ObjectKey)
	if err != nil {
		_ = removeLocalFile(req.LocalFilePath)
		return nil, fmt.Errorf("上传至OSS失败: %w", err)
	}
	return &storeResult{
		FileURL:  fileURL,
		OnCommit: func(context.Context) error { return removeLocalFile(req.LocalFilePath) },
		OnRollback: combineCleanupHooks(
			func(context.Context) error { return deleteFromOSS(p.cfg, req.ObjectKey) },
			func(context.Context) error { return removeLocalFile(req.LocalFilePath) },
		),
	}, nil
}

type cosStorageProvider struct {
	cfg cosConfig
}

func (p cosStorageProvider) Store(ctx context.Context, req storeRequest) (*storeResult, error) {
	fileURL, err := uploadToCOS(ctx, p.cfg, req.LocalFilePath, req.ObjectKey)
	if err != nil {
		_ = removeLocalFile(req.LocalFilePath)
		return nil, fmt.Errorf("上传至COS失败: %w", err)
	}
	return &storeResult{
		FileURL:  fileURL,
		OnCommit: func(context.Context) error { return removeLocalFile(req.LocalFilePath) },
		OnRollback: combineCleanupHooks(
			func(ctx context.Context) error { return deleteFromCOS(ctx, p.cfg, req.ObjectKey) },
			func(context.Context) error { return removeLocalFile(req.LocalFilePath) },
		),
	}, nil
}

func combineCleanupHooks(hooks ...cleanupHook) cleanupHook {
	return func(ctx context.Context) error {
		for _, hook := range hooks {
			if hook == nil {
				continue
			}
			if err := hook(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func runCleanupHook(ctx context.Context, stage string, hook cleanupHook) {
	if hook == nil {
		return
	}
	if err := hook(ctx); err != nil {
		g.Log().Warningf(ctx, "upload cleanup failed at stage=%s: %v", stage, err)
	}
}

func removeLocalFile(path string) error {
	if path == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
