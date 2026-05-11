package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gbaseadmin/codegen/parser"
)

func defaultAllowedApps() []string {
	return []string{"system", "upload"}
}

func isAllowedApp(cfg *Config, appName string) bool {
	if cfg == nil {
		return false
	}
	for _, item := range cfg.AllowedApps {
		if item == appName {
			return true
		}
	}
	return false
}

func hasFrontendDictModule(cfg *Config) bool {
	if cfg == nil || strings.TrimSpace(cfg.Frontend.Output) == "" {
		return false
	}

	dictAPIPath := filepath.Join(cfg.Frontend.Output, "api", "system", "dict", "index.ts")
	info, err := os.Stat(dictAPIPath)
	return err == nil && !info.IsDir()
}

func validateMetaScope(cfg *Config, meta *parser.TableMeta) error {
	if meta == nil {
		return fmt.Errorf("table meta is nil")
	}
	if !isAllowedApp(cfg, meta.AppName) {
		return fmt.Errorf("当前仓库只允许生成应用 %s，收到 %q", strings.Join(cfg.AllowedApps, ","), meta.AppName)
	}
	allowDictFields := cfg.AllowMissingDictModule || hasFrontendDictModule(cfg)
	for _, field := range meta.Fields {
		if field.DictType != "" && !allowDictFields {
			return fmt.Errorf("字段 %s 使用了 dict:%s，但当前精简仓库未保留 system/dict 模块", field.Name, field.DictType)
		}
	}
	return nil
}
