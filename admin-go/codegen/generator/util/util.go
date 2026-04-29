package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"gbaseadmin/codegen/parser"
)

// ReplacePlaceholders 将路径中的 {app} 和 {module} 占位符替换为实际名称
func ReplacePlaceholders(path, app, module string) string {
	result := strings.ReplaceAll(path, "{app}", app)
	result = strings.ReplaceAll(result, "{module}", module)
	return result
}

// TemplateMapping 模板文件名 → 输出相对路径
type TemplateMapping struct {
	TplFile    string
	OutputPath string
}

type FilePlanAction string

const (
	FilePlanActionCreate         FilePlanAction = "create"
	FilePlanActionUpdate         FilePlanAction = "update"
	FilePlanActionNoChange       FilePlanAction = "no_change"
	FilePlanActionSkipExisting   FilePlanAction = "skip_existing"
	FilePlanActionProtectEnhance FilePlanAction = "protect_enhance"
)

type PlannedFile struct {
	TemplateFile string         `json:"templateFile"`
	OutputPath   string         `json:"outputPath"`
	Action       FilePlanAction `json:"action"`
	Bytes        int            `json:"bytes"`
	Content      []byte         `json:"-"`
}

// SharedFuncMap 所有模板共享的自定义函数
var SharedFuncMap = template.FuncMap{
	"ModuleCamel": parser.SnakeToCamelSimple,
	"IsNumeric": func(s string) bool {
		if s == "" {
			return false
		}
		hasDigit := false
		hasDot := false
		for i, c := range s {
			if c == '-' && i == 0 {
				continue
			}
			if c == '.' && !hasDot {
				hasDot = true
				continue
			}
			if c < '0' || c > '9' {
				return false
			}
			hasDigit = true
		}
		return hasDigit
	},
}

// TemplateCache 模板缓存，避免重复解析
type TemplateCache struct {
	mu    sync.RWMutex
	cache map[string]*template.Template
}

// NewTemplateCache 创建模板缓存
func NewTemplateCache() *TemplateCache {
	return &TemplateCache{cache: make(map[string]*template.Template)}
}

func (p PlannedFile) MarshalJSON() ([]byte, error) {
	type alias PlannedFile
	return json.Marshal(alias{
		TemplateFile: p.TemplateFile,
		OutputPath:   p.OutputPath,
		Action:       p.Action,
		Bytes:        p.Bytes,
	})
}

// WriteFileIfChanged 仅在内容有变化时写入文件。
// 返回值表示是否实际落盘。
func WriteFileIfChanged(path string, content []byte) (bool, error) {
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return false, fmt.Errorf("目标路径是目录: %s", path)
		}
		existing, err := os.ReadFile(path)
		if err != nil {
			return false, err
		}
		if bytes.Equal(existing, content) {
			return false, nil
		}
	} else if !os.IsNotExist(err) {
		return false, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// Get 获取或解析模板（自动注入 SharedFuncMap）
func (tc *TemplateCache) Get(tplPath string) (*template.Template, error) {
	tc.mu.RLock()
	if t, ok := tc.cache[tplPath]; ok {
		tc.mu.RUnlock()
		return t, nil
	}
	tc.mu.RUnlock()

	tc.mu.Lock()
	defer tc.mu.Unlock()
	// 双重检查
	if t, ok := tc.cache[tplPath]; ok {
		return t, nil
	}
	t, err := template.New(filepath.Base(tplPath)).Funcs(SharedFuncMap).ParseFiles(tplPath)
	if err != nil {
		return nil, err
	}
	tc.cache[tplPath] = t
	return t, nil
}

func PlanFiles(mappings []TemplateMapping, tplDir, outDir, appName, moduleName string, force bool, data interface{}, cache ...*TemplateCache) ([]PlannedFile, error) {
	var planned []PlannedFile
	for _, m := range mappings {
		tplPath := filepath.Join(tplDir, m.TplFile)

		var tpl *template.Template
		var err error
		if len(cache) > 0 && cache[0] != nil {
			tpl, err = cache[0].Get(tplPath)
		} else {
			tpl, err = template.New(filepath.Base(tplPath)).Funcs(SharedFuncMap).ParseFiles(tplPath)
		}
		if err != nil {
			return nil, fmt.Errorf("解析模板 %s 失败: %w", m.TplFile, err)
		}

		relPath := ReplacePlaceholders(m.OutputPath, appName, moduleName)
		outPath := filepath.Join(outDir, relPath)
		entry := PlannedFile{
			TemplateFile: m.TplFile,
			OutputPath:   outPath,
		}

		// Enhance 文件保护：force 模式下跳过 *_enhance.* 文件
		if force && isEnhanceFile(outPath) {
			entry.Action = FilePlanActionProtectEnhance
			planned = append(planned, entry)
			continue
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("渲染模板 %s 失败: %w", m.TplFile, err)
		}
		entry.Content = buf.Bytes()
		entry.Bytes = len(entry.Content)

		info, err := os.Stat(outPath)
		if err == nil {
			if info.IsDir() {
				return nil, fmt.Errorf("目标路径是目录: %s", outPath)
			}
			existing, err := os.ReadFile(outPath)
			if err != nil {
				return nil, fmt.Errorf("读取现有文件 %s 失败: %w", outPath, err)
			}
			if bytes.Equal(existing, entry.Content) {
				entry.Action = FilePlanActionNoChange
			} else if !force {
				entry.Action = FilePlanActionSkipExisting
			} else {
				entry.Action = FilePlanActionUpdate
			}
		} else if os.IsNotExist(err) {
			entry.Action = FilePlanActionCreate
		} else {
			return nil, fmt.Errorf("检查文件 %s 状态失败: %w", outPath, err)
		}
		planned = append(planned, entry)
	}
	return planned, nil
}

// CommitPlannedFiles 将预渲染好的文件批量写入磁盘。
// 写入前会先把所有待提交内容落到临时文件，再统一替换目标文件。
func CommitPlannedFiles(plans []PlannedFile) ([]string, error) {
	type stagedFile struct {
		target string
		temp   string
	}

	staged := make([]stagedFile, 0, len(plans))
	cleanup := func() {
		for _, file := range staged {
			_ = os.Remove(file.temp)
		}
	}

	for _, plan := range plans {
		if plan.Action != FilePlanActionCreate && plan.Action != FilePlanActionUpdate {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(plan.OutputPath), 0o755); err != nil {
			cleanup()
			return nil, err
		}
		tempFile, err := os.CreateTemp(filepath.Dir(plan.OutputPath), "."+filepath.Base(plan.OutputPath)+".codegen-")
		if err != nil {
			cleanup()
			return nil, err
		}
		tempPath := tempFile.Name()
		if err := tempFile.Close(); err != nil {
			cleanup()
			return nil, err
		}
		staged = append(staged, stagedFile{target: plan.OutputPath, temp: tempPath})
		if err := os.WriteFile(tempPath, plan.Content, 0o644); err != nil {
			cleanup()
			return nil, err
		}
	}

	generated := make([]string, 0, len(staged))
	for _, file := range staged {
		if err := os.Rename(file.temp, file.target); err != nil {
			cleanup()
			return generated, err
		}
		generated = append(generated, file.target)
	}
	return generated, nil
}

// isEnhanceFile 判断是否是 enhance 文件
func isEnhanceFile(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(base, "_enhance.") || strings.Contains(base, "enhance.")
}
