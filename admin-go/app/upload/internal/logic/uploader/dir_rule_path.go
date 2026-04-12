package uploader

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model/entity"
)

func resolveUploadSavePath(
	ctx context.Context,
	cfg uploadStorageConfig,
	dirID int64,
	fileExt string,
	now time.Time,
) (relativeDir string, physicalDir string, err error) {
	relativeDir, err = resolveUploadRelativeDir(ctx, dirID, fileExt, now)
	if err != nil {
		return "", "", err
	}
	return relativeDir, joinLocalStorageDir(cfg.LocalPath, relativeDir), nil
}

func resolveUploadRelativeDir(ctx context.Context, dirID int64, fileExt string, now time.Time) (string, error) {
	defaultDir := now.Format("2006-01-02")
	if dirID <= 0 {
		return defaultDir, nil
	}

	var rules []*entity.UploadDirRule
	err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().DirId, dirID).
		Where(dao.UploadDirRule.Columns().Status, 1).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		OrderAsc(dao.UploadDirRule.Columns().Id).
		Scan(&rules)
	if err != nil {
		return "", fmt.Errorf("读取目录规则失败: %w", err)
	}

	renderedDir := renderUploadRulePath(selectUploadRulePath(rules, fileExt), now, fileExt)
	if renderedDir == "" {
		return defaultDir, nil
	}
	return renderedDir, nil
}

func selectUploadRulePath(rules []*entity.UploadDirRule, fileExt string) string {
	var defaultPath string
	normalizedExt := normalizeExt(fileExt)
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		switch rule.Category {
		case 2:
			if dirRuleFileTypeMatches(rule.FileType, normalizedExt) {
				return rule.SavePath
			}
		case 1:
			if defaultPath == "" {
				defaultPath = rule.SavePath
			}
		}
	}
	return defaultPath
}

func dirRuleFileTypeMatches(fileTypes, fileExt string) bool {
	fileExt = normalizeExt(fileExt)
	if fileExt == "" {
		return false
	}
	for _, item := range strings.Split(normalizeUploadRuleFileTypes(fileTypes), ",") {
		if item == fileExt {
			return true
		}
	}
	return false
}

func normalizeUploadRuleFileTypes(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || r == ' ' || r == '\n' || r == '\r' || r == '\t'
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = normalizeExt(part)
		if part == "" {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		normalized = append(normalized, part)
	}
	return strings.Join(normalized, ",")
}

func renderUploadRulePath(savePath string, now time.Time, fileExt string) string {
	savePath = strings.TrimSpace(savePath)
	if savePath == "" {
		return ""
	}
	rendered := strings.NewReplacer(
		"{Y-m-d}", now.Format("2006-01-02"),
		"{Y-m}", now.Format("2006-01"),
		"{Y}", now.Format("2006"),
		"{m}", now.Format("01"),
		"{d}", now.Format("02"),
		"{H}", now.Format("15"),
		"{i}", now.Format("04"),
		"{s}", now.Format("05"),
		"{ext}", normalizeExt(fileExt),
	).Replace(savePath)
	return strings.Join(splitUploadStoragePath(rendered), "/")
}

func buildObjectKey(relativeDir, uniqueName string) string {
	parts := splitUploadStoragePath(relativeDir)
	if uniqueName != "" {
		parts = append(parts, uniqueName)
	}
	return strings.Join(parts, "/")
}

func joinLocalStorageDir(baseDir, relativeDir string) string {
	parts := append([]string{baseDir}, splitUploadStoragePath(relativeDir)...)
	return filepath.Join(parts...)
}

func splitUploadStoragePath(path string) []string {
	path = strings.Trim(strings.TrimSpace(path), `/\`)
	if path == "" {
		return nil
	}
	rawSegments := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	segments := make([]string, 0, len(rawSegments))
	for _, segment := range rawSegments {
		segment = strings.TrimSpace(segment)
		if segment == "" || segment == "." || segment == ".." {
			continue
		}
		segments = append(segments, segment)
	}
	return segments
}
