package shared

import (
	"net/url"
	"path"
	"strings"
)

// NormalizeUploadRuleSource standardizes upload source identifiers such as
// frontend routes or full URLs into a comparable absolute path.
func NormalizeUploadRuleSource(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	if strings.HasPrefix(value, "#") {
		value = strings.TrimPrefix(value, "#")
	}
	if parsed, err := url.Parse(value); err == nil {
		switch {
		case parsed.Path != "" && (parsed.Scheme != "" || parsed.Host != ""):
			value = parsed.Path
		case parsed.Fragment != "" && (parsed.Path == "" || parsed.Path == "/"):
			value = parsed.Fragment
		}
	}

	if idx := strings.Index(value, "#"); idx >= 0 {
		value = value[idx+1:]
	}
	if idx := strings.Index(value, "?"); idx >= 0 {
		value = value[:idx]
	}

	value = strings.ReplaceAll(strings.TrimSpace(value), `\`, "/")
	if value == "" {
		return ""
	}

	cleaned := path.Clean(value)
	switch cleaned {
	case "", ".":
		return ""
	case "/":
		return "/"
	}
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + strings.TrimPrefix(cleaned, "./")
	}
	return cleaned
}
