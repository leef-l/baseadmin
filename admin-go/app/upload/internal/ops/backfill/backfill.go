package backfill

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"gbaseadmin/app/upload/internal/consts"
	"gbaseadmin/app/upload/internal/dao"
	"gbaseadmin/app/upload/internal/model/do"
	"gbaseadmin/app/upload/internal/model/entity"
)

const defaultBatchSize = 500

type Options struct {
	BatchSize int
	DryRun    bool
}

type Result struct {
	Scanned       int
	Updated       int
	NoRelativeDir int
	NoMatch       int
	Ambiguous     int
}

type fileRow struct {
	ID      int64  `json:"id"`
	URL     string `json:"url"`
	Ext     string `json:"ext"`
	Storage int    `json:"storage"`
}

type ruleMatcher struct {
	dirID        int64
	category     int
	fileType     string
	storageTypes string
	regex        *regexp.Regexp
	score        int
}

func Run(ctx context.Context, opts Options) (*Result, error) {
	opts = normalizeOptions(opts)

	rules, err := loadActiveRules(ctx)
	if err != nil {
		return nil, err
	}
	matchers, err := buildRuleMatchers(rules)
	if err != nil {
		return nil, err
	}

	result := &Result{}
	if len(matchers) == 0 {
		return result, nil
	}

	var lastID int64
	for {
		rows, err := loadFileBatch(ctx, lastID, opts.BatchSize)
		if err != nil {
			return nil, err
		}
		if len(rows) == 0 {
			return result, nil
		}

		for _, row := range rows {
			lastID = row.ID
			result.Scanned++

			relativeDir := extractRelativeDirFromURL(row.URL)
			if relativeDir == "" {
				result.NoRelativeDir++
				continue
			}

			dirID, matched, ambiguous := selectDirIDForBackfill(matchers, relativeDir, row.Ext, row.Storage)
			switch {
			case ambiguous:
				result.Ambiguous++
				continue
			case !matched || dirID <= 0:
				result.NoMatch++
				continue
			}

			if opts.DryRun {
				result.Updated++
				continue
			}
			if _, err := dao.UploadFile.Ctx(ctx).
				Where(dao.UploadFile.Columns().Id, row.ID).
				Data(do.UploadFile{DirId: dirID}).
				Update(); err != nil {
				return nil, err
			}
			result.Updated++
		}
	}
}

func normalizeOptions(opts Options) Options {
	if opts.BatchSize <= 0 {
		opts.BatchSize = defaultBatchSize
	}
	return opts
}

func loadActiveRules(ctx context.Context) ([]*entity.UploadDirRule, error) {
	var rules []*entity.UploadDirRule
	if err := dao.UploadDirRule.Ctx(ctx).
		Where(dao.UploadDirRule.Columns().Status, consts.DirRuleStatus1).
		Where(dao.UploadDirRule.Columns().DeletedAt, nil).
		OrderAsc(dao.UploadDirRule.Columns().Id).
		Scan(&rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func buildRuleMatchers(rules []*entity.UploadDirRule) ([]ruleMatcher, error) {
	matchers := make([]ruleMatcher, 0, len(rules))
	for _, rule := range rules {
		if rule == nil || rule.DirId == 0 {
			continue
		}
		regex, score, ok, err := compileRulePathPattern(rule.SavePath)
		if err != nil {
			return nil, gerror.Wrapf(err, "compile dir rule %d save_path", rule.Id)
		}
		if !ok {
			continue
		}
		matchers = append(matchers, ruleMatcher{
			dirID:        int64(rule.DirId),
			category:     rule.Category,
			fileType:     normalizeRuleFileTypes(rule.FileType),
			storageTypes: normalizeStorageTypes(rule.StorageTypes),
			regex:        regex,
			score:        score + categoryBackfillWeight(rule.Category),
		})
	}
	return matchers, nil
}

func loadFileBatch(ctx context.Context, lastID int64, batchSize int) ([]fileRow, error) {
	model := dao.UploadFile.Ctx(ctx).
		Fields(
			dao.UploadFile.Columns().Id,
			dao.UploadFile.Columns().Url,
			dao.UploadFile.Columns().Ext,
			dao.UploadFile.Columns().Storage,
		).
		Where(dao.UploadFile.Columns().DirId, 0).
		Where(dao.UploadFile.Columns().DeletedAt, nil).
		OrderAsc(dao.UploadFile.Columns().Id).
		Limit(batchSize)
	if lastID > 0 {
		model = model.Where(fmt.Sprintf("%s > ?", dao.UploadFile.Columns().Id), lastID)
	}

	var rows []fileRow
	if err := model.Scan(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func selectDirIDForBackfill(matchers []ruleMatcher, relativeDir, fileExt string, storageType int) (int64, bool, bool) {
	var (
		bestDirID int64
		bestScore = -1
		matched   bool
		ambiguous bool
	)

	normalizedDir := normalizeRelativeDir(relativeDir)
	normalizedExt := normalizeExt(fileExt)
	for _, matcher := range matchers {
		if !ruleSupportsStorageType(matcher.storageTypes, storageType) {
			continue
		}
		if matcher.category == consts.DirRuleCategory2 && !ruleFileTypeMatches(matcher.fileType, normalizedExt) {
			continue
		}
		if matcher.regex == nil || !matcher.regex.MatchString(normalizedDir) {
			continue
		}

		matched = true
		switch {
		case matcher.score > bestScore:
			bestDirID = matcher.dirID
			bestScore = matcher.score
			ambiguous = false
		case matcher.score == bestScore && matcher.dirID != bestDirID:
			ambiguous = true
		}
	}

	if ambiguous {
		return 0, false, true
	}
	if !matched || bestDirID <= 0 {
		return 0, false, false
	}
	return bestDirID, true, false
}

func extractRelativeDirFromURL(fileURL string) string {
	fileURL = strings.TrimSpace(fileURL)
	if fileURL == "" {
		return ""
	}

	parsed, err := url.Parse(fileURL)
	switch {
	case err == nil && strings.TrimSpace(parsed.Path) != "":
		fileURL = parsed.Path
	case strings.Contains(fileURL, "://"):
		return ""
	}

	fileURL = strings.ReplaceAll(strings.TrimSpace(fileURL), `\`, "/")
	switch {
	case strings.HasPrefix(fileURL, "/upload/"):
		fileURL = strings.TrimPrefix(fileURL, "/upload/")
	case strings.HasPrefix(fileURL, "/resource/"):
		fileURL = strings.TrimPrefix(fileURL, "/resource/")
	default:
		fileURL = strings.TrimPrefix(fileURL, "/")
	}

	cleaned := path.Clean("/" + strings.TrimPrefix(fileURL, "/"))
	if cleaned == "/" {
		return ""
	}
	dir := path.Dir(strings.TrimPrefix(cleaned, "/"))
	if dir == "." || dir == "/" {
		return ""
	}
	return normalizeRelativeDir(dir)
}

func compileRulePathPattern(savePath string) (*regexp.Regexp, int, bool, error) {
	template := normalizeRelativeDir(savePath)
	if template == "" {
		return nil, 0, false, nil
	}

	var (
		builder strings.Builder
		score   int
	)
	builder.WriteString("^")

	for i := 0; i < len(template); {
		token, replacement := matchRuleToken(template[i:])
		if token != "" {
			builder.WriteString(replacement)
			i += len(token)
			continue
		}

		builder.WriteString(regexp.QuoteMeta(string(template[i])))
		if template[i] != '/' {
			score++
		}
		i++
	}

	builder.WriteString("$")
	regex, err := regexp.Compile(builder.String())
	if err != nil {
		return nil, 0, false, err
	}
	return regex, score, true, nil
}

func matchRuleToken(value string) (string, string) {
	switch {
	case strings.HasPrefix(value, "{Y-m-d}"):
		return "{Y-m-d}", `\d{4}-\d{2}-\d{2}`
	case strings.HasPrefix(value, "{Y-m}"):
		return "{Y-m}", `\d{4}-\d{2}`
	case strings.HasPrefix(value, "{Y}"):
		return "{Y}", `\d{4}`
	case strings.HasPrefix(value, "{m}"):
		return "{m}", `\d{2}`
	case strings.HasPrefix(value, "{d}"):
		return "{d}", `\d{2}`
	case strings.HasPrefix(value, "{H}"):
		return "{H}", `\d{2}`
	case strings.HasPrefix(value, "{i}"):
		return "{i}", `\d{2}`
	case strings.HasPrefix(value, "{s}"):
		return "{s}", `\d{2}`
	case strings.HasPrefix(value, "{ext}"):
		return "{ext}", `[^/]+`
	default:
		return "", ""
	}
}

func categoryBackfillWeight(category int) int {
	switch category {
	case consts.DirRuleCategory2:
		return 1000
	default:
		return 0
	}
}

func ruleSupportsStorageType(storageTypes string, storageType int) bool {
	target := strings.TrimSpace(fmt.Sprintf("%d", storageType))
	if target == "" || target == "0" {
		return false
	}
	if storageTypes == "" {
		return true
	}
	for _, item := range strings.Split(storageTypes, ",") {
		if strings.TrimSpace(item) == target {
			return true
		}
	}
	return false
}

func ruleFileTypeMatches(fileTypes, fileExt string) bool {
	fileExt = normalizeExt(fileExt)
	if fileExt == "" {
		return false
	}
	for _, item := range strings.Split(normalizeRuleFileTypes(fileTypes), ",") {
		if item == fileExt {
			return true
		}
	}
	return false
}

func normalizeStorageTypes(value string) string {
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

func normalizeExt(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimPrefix(value, ".")
	return value
}

func normalizeRuleFileTypes(value string) string {
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

func normalizeRelativeDir(value string) string {
	value = strings.ReplaceAll(strings.TrimSpace(value), `\`, "/")
	value = strings.Trim(value, "/")
	if value == "" {
		return ""
	}
	cleaned := path.Clean("/" + value)
	switch cleaned {
	case "", ".", "/":
		return ""
	default:
		return strings.TrimPrefix(cleaned, "/")
	}
}
