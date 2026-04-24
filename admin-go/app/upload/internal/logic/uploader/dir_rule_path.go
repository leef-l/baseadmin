package uploader

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/logic/shared"
)

type uploadDirRuleRecord struct {
	Id           uint64 `orm:"id"`
	DirId        uint64 `orm:"dir_id"`
	Category     int    `orm:"category"`
	FileType     string `orm:"file_type"`
	StorageTypes string `orm:"storage_types"`
	SavePath     string `orm:"save_path"`
	KeepName     int    `orm:"keep_name"`
}

func resolveUploadSavePath(
	ctx context.Context,
	cfg uploadStorageConfig,
	dirID int64,
	fileExt string,
	mimeType string,
	source string,
	now time.Time,
	fallbackRelativeDir string,
	systemUserID int64,
) (resolvedDirID int64, relativeDir string, physicalDir string, keepName bool, err error) {
	resolvedDirID, relativeDir, keepName, err = resolveUploadRelativeDir(ctx, dirID, cfg.StorageType, fileExt, mimeType, source, now, fallbackRelativeDir, systemUserID)
	if err != nil {
		return 0, "", "", false, err
	}
	physicalDir, ok := shared.ResolveLocalStorageDir(cfg.LocalPath, relativeDir)
	if !ok {
		return 0, "", "", false, fmt.Errorf("保存目录超出允许范围")
	}
	if cfg.StorageType != 1 && hasParentRelativeDir(relativeDir) {
		return 0, "", "", false, fmt.Errorf("云存储不支持父级目录规则")
	}
	return resolvedDirID, relativeDir, physicalDir, keepName, nil
}

func resolveUploadRelativeDir(
	ctx context.Context,
	dirID int64,
	storageType int,
	fileExt string,
	mimeType string,
	source string,
	now time.Time,
	fallbackRelativeDir string,
	systemUserID int64,
) (int64, string, bool, error) {
	defaultDir := now.Format("2006-01-02")
	rules, err := loadUploadRules(ctx, dirID)
	if err != nil {
		return 0, "", false, err
	}
	selectedDirKeepName := false
	if dirID > 0 {
		selectedDirKeepName, err = loadUploadDirKeepName(ctx, dirID)
		if err != nil {
			return 0, "", false, err
		}
	}

	selectedRule := selectUploadRule(rules, storageType, fileExt, mimeType, source)
	renderedDir := renderUploadRulePath(selectedRuleSavePath(selectedRule), now, fileExt, systemUserID)
	if selectedRule != nil && selectedRule.DirId > 0 {
		if renderedDir == "" {
			renderedDir = defaultDir
		}
		resolvedDirID := int64(selectedRule.DirId)
		dirKeepName := selectedDirKeepName
		if resolvedDirID != dirID {
			dirKeepName, err = loadUploadDirKeepName(ctx, resolvedDirID)
			if err != nil {
				return 0, "", false, err
			}
		}
		return resolvedDirID, renderedDir, selectedRule.KeepName == 1 || dirKeepName, nil
	}
	if selectedRule != nil {
		if renderedDir == "" {
			renderedDir = defaultDir
		}
		return dirID, renderedDir, selectedRule.KeepName == 1 || selectedDirKeepName, nil
	}
	fallbackRelativeDir = sanitizeUploadSubdir(fallbackRelativeDir)
	if fallbackRelativeDir != "" {
		return dirID, path.Join(fallbackRelativeDir, defaultDir), selectedDirKeepName, nil
	}
	return dirID, defaultDir, selectedDirKeepName, nil
}

func loadUploadRules(ctx context.Context, dirID int64) ([]*uploadDirRuleRecord, error) {
	m := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Status, 1).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		OrderAsc(dao.UploadDirRule.Columns().Id)
	if dirID > 0 {
		m = m.Where(dao.UploadDirRule.Columns().DirId, dirID)
	}
	var rules []*uploadDirRuleRecord
	if err := m.Scan(&rules); err != nil {
		return nil, fmt.Errorf("读取目录规则失败: %w", err)
	}
	return rules, nil
}

func loadUploadDirKeepName(ctx context.Context, dirID int64) (bool, error) {
	if dirID <= 0 {
		return false, nil
	}
	value, err := dao.UploadDir.Ctx(ctx).
		Fields("keep_name").
		Where(dao.UploadDir.Columns().Id, dirID).
		Where(dao.UploadDir.Columns().DeletedAt, nil).
		Value("keep_name")
	if err != nil {
		return false, fmt.Errorf("读取目录配置失败: %w", err)
	}
	return value.Int() == 1, nil
}

func selectUploadRulePath(rules []*uploadDirRuleRecord, storageType int, fileExt, source string) string {
	return selectedRuleSavePath(selectUploadRule(rules, storageType, fileExt, "", source))
}

func selectedRuleSavePath(rule *uploadDirRuleRecord) string {
	if rule == nil {
		return ""
	}
	return rule.SavePath
}

func selectUploadRule(rules []*uploadDirRuleRecord, storageType int, fileExt, mimeType, source string) *uploadDirRuleRecord {
	var sourceRule *uploadDirRuleRecord
	var typeRule *uploadDirRuleRecord
	var defaultRule *uploadDirRuleRecord
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
			if sourceRule == nil && dirRuleSourceMatches(rule.FileType, normalizedSource) {
				sourceRule = rule
			}
		case 2:
			if typeRule == nil && dirRuleFileTypeMatches(rule.FileType, normalizedExt, mimeType) {
				typeRule = rule
			}
		case 1:
			if defaultRule == nil {
				defaultRule = rule
			}
		}
	}
	if sourceRule != nil {
		return sourceRule
	}
	if typeRule != nil {
		return typeRule
	}
	if defaultRule != nil {
		return defaultRule
	}
	return nil
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

func dirRuleFileTypeMatches(fileTypes, fileExt, mimeType string) bool {
	fileExt = normalizeExt(fileExt)
	mimeType = normalizeMimeForRuleMatch(mimeType)
	if fileExt == "" && mimeType == "" {
		return false
	}
	for _, item := range strings.Split(normalizeUploadRuleFileTypes(fileTypes), ",") {
		if matchUploadRuleFileTypeToken(item, fileExt, mimeType) {
			return true
		}
	}
	return false
}

func matchUploadRuleFileTypeToken(token, fileExt, mimeType string) bool {
	token = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(token)), ".")
	if token == "" {
		return false
	}
	switch token {
	case "image", "img":
		return strings.HasPrefix(mimeType, "image/")
	case "video":
		return strings.HasPrefix(mimeType, "video/")
	case "audio":
		return strings.HasPrefix(mimeType, "audio/")
	case "pdf":
		return fileExt == "pdf" || mimeType == "application/pdf"
	}
	if strings.HasSuffix(token, "/*") {
		return strings.HasPrefix(mimeType, strings.TrimSuffix(token, "*"))
	}
	if strings.Contains(token, "/") {
		return mimeType == token
	}
	return fileExt == token
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
		part = normalizeUploadRuleFileType(part)
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

func renderUploadRulePath(savePath string, now time.Time, fileExt string, systemUserID int64) string {
	savePath = strings.TrimSpace(savePath)
	if savePath == "" {
		return ""
	}
	systemUserIDText := "0"
	if systemUserID > 0 {
		systemUserIDText = strconv.FormatInt(systemUserID, 10)
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
		"{systemUserId}", systemUserIDText,
	).Replace(savePath)
	rendered = strings.ReplaceAll(rendered, `\`, "/")
	rendered = normalizeUploadRuleSavePathAlias(rendered)
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
	relativeDir = normalizeUploadRuleSavePathAlias(relativeDir)
	return relativeDir == ".." || strings.HasPrefix(relativeDir, "../")
}

func normalizeUploadRuleFileType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "/") {
		value = strings.TrimPrefix(value, ".")
	}
	return value
}

func normalizeMimeForRuleMatch(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if index := strings.Index(value, ";"); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return value
}

func normalizeUploadRuleSavePathAlias(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, `\`, "/"))
	switch {
	case value == "@up":
		return ".."
	case strings.HasPrefix(value, "@up/"):
		return "../" + strings.TrimPrefix(value, "@up/")
	default:
		return value
	}
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
