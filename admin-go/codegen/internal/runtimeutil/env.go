package runtimeutil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const maxEnvLineBytes = 1024 * 1024

// LoadEnvFileIfExists 按 KEY=VALUE 格式加载环境变量文件。
// 文件不存在时静默跳过；已存在的环境变量不会被覆盖。
func LoadEnvFileIfExists(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("打开环境变量文件失败: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 4096), maxEnvLineBytes)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		value = strings.TrimSpace(value)
		if isQuotedValue(value) {
			value = strings.Trim(value, `"'`)
		} else {
			value = trimInlineComment(value)
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("设置环境变量 %s 失败: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取环境变量文件失败: %w", err)
	}
	return nil
}

func isQuotedValue(value string) bool {
	if len(value) < 2 {
		return false
	}
	if value[0] != value[len(value)-1] {
		return false
	}
	return value[0] == '"' || value[0] == '\''
}

func trimInlineComment(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] != '#' {
			continue
		}
		if i == 0 || unicode.IsSpace(rune(value[i-1])) {
			return strings.TrimSpace(value[:i])
		}
	}
	return value
}
