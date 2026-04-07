package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gbaseadmin/codegen/internal/runtimeutil"
	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type BackendConfig struct {
	Output string `yaml:"output"`
}

type FrontendConfig struct {
	Output string `yaml:"output"`
}

// MenuAppConfig 菜单应用目录配置
type MenuAppConfig struct {
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
	Sort  int    `yaml:"sort"` // 目录排序，默认 50
}

// MenuModuleConfig 模块级菜单配置
type MenuModuleConfig struct {
	Sort   int  `yaml:"sort"`    // 菜单排序
	IsShow *int `yaml:"is_show"` // 菜单是否显示（nil=默认1）
}

type Config struct {
	Database     DatabaseConfig              `yaml:"database"`
	Backend      BackendConfig               `yaml:"backend"`
	Frontend     FrontendConfig              `yaml:"frontend"`
	SkipFields   []string                    `yaml:"skip_fields"`
	MenuApps     map[string]MenuAppConfig    `yaml:"menu_apps"`
	MenuModules  map[string]MenuModuleConfig `yaml:"menu_modules"`  // key: "appName/moduleName"
	OperationLog bool                        `yaml:"operation_log"` // 全局操作日志开关
}

func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件路径失败: %w", err)
	}
	configDir := filepath.Dir(absPath)

	if err := runtimeutil.LoadEnvFileIfExists(filepath.Join(configDir, "..", ".env")); err != nil {
		return nil, fmt.Errorf("加载环境变量失败: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	expanded := expandEnvPlaceholders(string(data))
	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	cfg.Backend.Output = resolveConfigPath(configDir, cfg.Backend.Output)
	cfg.Frontend.Output = resolveConfigPath(configDir, cfg.Frontend.Output)
	return &cfg, nil
}

var envPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

func expandEnvPlaceholders(input string) string {
	return envPattern.ReplaceAllStringFunc(input, func(match string) string {
		sub := envPattern.FindStringSubmatch(match)
		if len(sub) != 2 {
			return match
		}
		if value, ok := os.LookupEnv(sub[1]); ok {
			return value
		}
		return match
	})
}

func resolveConfigPath(configDir, value string) string {
	if value == "" || filepath.IsAbs(value) {
		return value
	}
	return filepath.Clean(filepath.Join(configDir, value))
}

// DSN 生成数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// DSNForHack 生成 hack/config.yaml 中 gf gen dao 使用的 link 格式
// 格式: mysql:user:password@tcp(host:port)/dbname?charset=utf8mb4
func (c *DatabaseConfig) DSNForHack() string {
	return fmt.Sprintf("mysql:%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}
