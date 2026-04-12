package uploader

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
	"gbaseadmin/app/upload/internal/model/entity"
)

func resolveUploadSavePath(
	ctx context.Context,
	cfg uploadStorageConfig,
	dirID int64,
	fileExt string,
	source string,
	now time.Time,
) (relativeDir string, physicalDir string, err error) {
	relativeDir, err = resolveUploadRelativeDir(ctx, dirID, cfg.StorageType, fileExt, source, now)
	if err != nil {
		return "", "", err
	}
	physicalDir, ok := shared.ResolveLocalStorageDir(cfg.LocalPath, relativeDir)
	if !ok {
		return "", "", fmt.Errorf("保存目录超出允许范围")
	}
	if cfg.StorageType != 1 && hasParentRelativeDir(relativeDir) {
		return "", "", fmt.Errorf("云存储不支持父级目录规则")
	}
	return relativeDir, physicalDir, nil
}

func resolveUploadRelativeDir(ctx context.Context, dirID int64, storageType int, fileExt, source string, now time.Time) (string, error) {
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

	renderedDir := renderUploadRulePath(selectUploadRulePath(rules, storageType, fileExt, source), now, fileExt)
	if renderedDir == "" {
		return defaultDir, nil
	}
	return renderedDir, nil
}

func selectUploadRulePath(rules []*entity.UploadDirRule, storageType int, fileExt, source string) string {
	var sourcePath string
	var typePath string
	var defaultPath string
	normalizedExt := normalizeExt(fileExt)
	normalizedSource := shared.NormalizeUploadRuleSource(source)
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		if !dirRuleSupportsStorageType(rule.StorageTypes, storageType) {
			continue
		}
		switch rule.Category {
		case 3:
			if sourcePath == "" && dirRuleSourceMatches(rule.FileType, normalizedSource) {
				sourcePath = rule.SavePath
			}
		case 2:
			if typePath == "" && dirRuleFileTypeMatches(rule.FileType, normalizedExt) {
				typePath = rule.SavePath
			}
		case 1:
			if defaultPath == "" {
				defaultPath = rule.SavePath
			}
		}
	}
	if sourcePath != "" {
		return sourcePath
	}
	if typePath != "" {
		return typePath
	}
	return defaultPath
}

func dirRuleSupportsStorageType(storageTypes string, storageType int) bool {
	target := normalizeUploadRuleStorageType(storageType)
	if target == "" {
		return false
	}
	if strings.TrimSpace(storageTypes) == "" {
		return true
	}
	for _, item := range strings.Split(normalizeUploadRuleStorageTypes(storageTypes), ",") {
		if item == target {
			return true
		}
	}
	return false
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

func dirRuleSourceMatches(matchers, source string) bool {
	source = shared.NormalizeUploadRuleSource(source)
	if source == "" {
		return false
	}
	for _, item := range strings.Split(normalizeUploadRuleSources(matchers), ",") {
		if item == "" {
			continue
		}
		if strings.HasSuffix(item, "/*") {
			prefix := strings.TrimSuffix(item, "/*")
			if source == prefix || strings.HasPrefix(source, prefix+"/") {
				return true
			}
			continue
		}
		if source == item || strings.HasPrefix(source, item+"/") {
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

func normalizeUploadRuleSources(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || r == ' ' || r == '\n' || r == '\r' || r == '\t'
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = normalizeUploadRuleSourceMatcher(part)
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

func normalizeUploadRuleSourceMatcher(value string) string {
	value = strings.TrimSpace(value)
	hasWildcard := strings.HasSuffix(value, "/*")
	if hasWildcard {
		value = strings.TrimSuffix(value, "/*")
	}
	value = shared.NormalizeUploadRuleSource(value)
	if value == "" {
		return ""
	}
	if hasWildcard && value != "/" {
		return value + "/*"
	}
	return value
}

func normalizeUploadRuleStorageTypes(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || r == ' ' || r == '\n' || r == '\r' || r == '\t'
	})
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "1" && part != "2" && part != "3" {
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

func normalizeUploadRuleStorageType(storageType int) string {
	switch storageType {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	default:
		return ""
	}
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
	rendered = strings.ReplaceAll(rendered, `\`, "/")
	rendered = path.Clean(strings.TrimSpace(rendered))
	switch rendered {
	case "", ".", "/":
		return ""
	}
	return strings.TrimPrefix(rendered, "/")
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

func hasParentRelativeDir(relativeDir string) bool {
	relativeDir = strings.TrimSpace(strings.ReplaceAll(relativeDir, `\`, "/"))
	return relativeDir == ".." || strings.HasPrefix(relativeDir, "../")
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
