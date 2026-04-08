package shared

import (
	"path/filepath"
	"strings"
)

const DefaultLocalStoragePath = "resource/upload"

func NormalizeLocalStoragePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return DefaultLocalStoragePath
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return DefaultLocalStoragePath
	}
	return cleaned
}

func BuildLocalFileURL(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		filtered = append(filtered, sanitizeStoragePathSegment(part)...)
	}
	if len(filtered) == 0 {
		return "/upload"
	}
	return "/upload/" + strings.Join(filtered, "/")
}

func LocalStoragePhysicalPath(fileURL string) string {
	fileURL = strings.TrimSpace(fileURL)
	if fileURL == "" {
		return DefaultLocalStoragePath
	}
	trimmed := strings.TrimPrefix(fileURL, "/upload")
	trimmed = strings.TrimPrefix(trimmed, "/")
	if trimmed == "" {
		return DefaultLocalStoragePath
	}
	parts := strings.Split(trimmed, "/")
	all := make([]string, 0, len(parts)+1)
	all = append(all, DefaultLocalStoragePath)
	for _, part := range parts {
		all = append(all, sanitizeStoragePathSegment(part)...)
	}
	return filepath.Join(all...)
}

func sanitizeStoragePathSegment(part string) []string {
	part = strings.Trim(strings.TrimSpace(part), `/\`)
	if part == "" {
		return nil
	}

	rawSegments := strings.FieldsFunc(part, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	filtered := make([]string, 0, len(rawSegments))
	for _, segment := range rawSegments {
		segment = strings.TrimSpace(segment)
		if segment == "" || segment == "." || segment == ".." {
			continue
		}
		filtered = append(filtered, segment)
	}
	return filtered
}
