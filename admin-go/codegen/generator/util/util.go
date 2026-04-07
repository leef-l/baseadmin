package util

import (
	"bytes"
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

// SharedFuncMap 所有模板共享的自定义函数
var SharedFuncMap = template.FuncMap{
	"ModuleCamel": parser.SnakeToCamelSimple,
	"IsNumeric": func(s string) bool {
		if s == "" {
			return false
		}
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
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

// GenerateFiles 通用文件生成逻辑
func GenerateFiles(mappings []TemplateMapping, tplDir, outDir, appName, moduleName string, force bool, data interface{}, cache ...*TemplateCache) ([]string, error) {
	var generated []string

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
			return generated, fmt.Errorf("解析模板 %s 失败: %w", m.TplFile, err)
		}

		relPath := ReplacePlaceholders(m.OutputPath, appName, moduleName)
		outPath := filepath.Join(outDir, relPath)

		// Enhance 文件保护：force 模式下跳过 *_enhance.* 文件
		if force && isEnhanceFile(outPath) {
			fmt.Printf("  [保护] %s（enhance 文件不覆盖）\n", outPath)
			continue
		}

		if !force {
			if _, err := os.Stat(outPath); err == nil {
				fmt.Printf("  [跳过] %s（已存在，使用 --force 覆盖）\n", outPath)
				continue
			}
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			return generated, fmt.Errorf("渲染模板 %s 失败: %w", m.TplFile, err)
		}

		written, err := WriteFileIfChanged(outPath, buf.Bytes())
		if err != nil {
			return generated, fmt.Errorf("写入文件 %s 失败: %w", outPath, err)
		}
		if !written {
			fmt.Printf("  [无变化] %s\n", outPath)
			continue
		}

		generated = append(generated, outPath)
		fmt.Printf("  [生成] %s\n", outPath)
	}

	return generated, nil
}

// GenerateToMemory 生成到内存（用于 dry-run diff 预览）
func GenerateToMemory(mappings []TemplateMapping, tplDir, outDir, appName, moduleName string, data interface{}, cache ...*TemplateCache) (map[string][]byte, error) {
	result := make(map[string][]byte)
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
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("渲染模板 %s 失败: %w", m.TplFile, err)
		}
		result[outPath] = buf.Bytes()
	}
	return result, nil
}

// isEnhanceFile 判断是否是 enhance 文件
func isEnhanceFile(path string) bool {
	base := filepath.Base(path)
	return strings.Contains(base, "_enhance.") || strings.Contains(base, "enhance.")
}
