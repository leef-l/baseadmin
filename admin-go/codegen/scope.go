package main

import (
	"fmt"
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

func validateMetaScope(cfg *Config, meta *parser.TableMeta) error {
	if meta == nil {
		return fmt.Errorf("table meta is nil")
	}
	if !isAllowedApp(cfg, meta.AppName) {
		return fmt.Errorf("当前仓库只允许生成应用 %s，收到 %q", strings.Join(cfg.AllowedApps, ","), meta.AppName)
	}
	for _, field := range meta.Fields {
		if field.DictType != "" {
			return fmt.Errorf("字段 %s 使用了 dict:%s，但当前精简仓库未保留 system/dict 模块", field.Name, field.DictType)
		}
	}
	return nil
}
