package cmdutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UseConfigFile points GoFrame config loading at the provided config file.
func UseConfigFile(configPath string) error {
	configPath = strings.TrimSpace(configPath)
	if configPath == "" {
		return nil
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("stat config path: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("config path must be a file: %s", absPath)
	}

	if err := os.Setenv("GF_GCFG_PATH", filepath.Dir(absPath)); err != nil {
		return fmt.Errorf("set GF_GCFG_PATH: %w", err)
	}
	if err := os.Setenv("GF_GCFG_FILE", filepath.Base(absPath)); err != nil {
		return fmt.Errorf("set GF_GCFG_FILE: %w", err)
	}
	return nil
}
