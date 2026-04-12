package shared

import (
	"path/filepath"
	"strings"
)

const DefaultLocalStoragePath = "resource/upload"
const DefaultLocalStoragePublicPrefix = "/upload"
const ResourceStoragePublicPrefix = "/resource"

func NormalizeLocalStoragePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return DefaultLocalStoragePath
	}
	slashed := filepath.ToSlash(strings.ReplaceAll(path, "\\", "/"))
	slashed = strings.TrimPrefix(strings.TrimSpace(slashed), "./")
	slashed = strings.TrimPrefix(slashed, "/")
	switch {
	case slashed == "", slashed == ".":
		return DefaultLocalStoragePath
	case slashed == "upload":
		return DefaultLocalStoragePath
	case strings.HasPrefix(slashed, "upload/"):
		return filepath.Clean(filepath.Join("resource", filepath.FromSlash(slashed)))
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return DefaultLocalStoragePath
	}
	return cleaned
}

func LocalStoragePublicRoot(baseLocalPath string) string {
	baseLocalPath = NormalizeLocalStoragePath(baseLocalPath)
	return filepath.Clean(filepath.Dir(baseLocalPath))
}

func BuildLocalFileURL(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		filtered = append(filtered, sanitizeStoragePathSegment(part)...)
	}
	if len(filtered) == 0 {
		return DefaultLocalStoragePublicPrefix
	}
	return DefaultLocalStoragePublicPrefix + "/" + strings.Join(filtered, "/")
}

func BuildLocalFileURLWithBase(baseLocalPath, relativeDir, fileName string) string {
	baseLocalPath = NormalizeLocalStoragePath(baseLocalPath)
	publicRoot := LocalStoragePublicRoot(baseLocalPath)
	targetPath := resolveLocalStoragePath(baseLocalPath, relativeDir, fileName)
	if relToBase, ok := relativePathWithin(baseLocalPath, targetPath); ok {
		return joinPublicURL(DefaultLocalStoragePublicPrefix, relToBase)
	}
	if relToPublicRoot, ok := relativePathWithin(publicRoot, targetPath); ok {
		return joinPublicURL(ResourceStoragePublicPrefix, relToPublicRoot)
	}
	return joinPublicURL(DefaultLocalStoragePublicPrefix, fileName)
}

func LocalStoragePhysicalPath(fileURL string) string {
	fileURL = strings.TrimSpace(fileURL)
	if fileURL == "" {
		return DefaultLocalStoragePath
	}
	basePath := DefaultLocalStoragePath
	switch {
	case fileURL == DefaultLocalStoragePublicPrefix || strings.HasPrefix(fileURL, DefaultLocalStoragePublicPrefix+"/"):
		fileURL = strings.TrimPrefix(fileURL, DefaultLocalStoragePublicPrefix)
	case fileURL == ResourceStoragePublicPrefix || strings.HasPrefix(fileURL, ResourceStoragePublicPrefix+"/"):
		basePath = LocalStoragePublicRoot(DefaultLocalStoragePath)
		fileURL = strings.TrimPrefix(fileURL, ResourceStoragePublicPrefix)
	default:
		basePath = LocalStoragePublicRoot(DefaultLocalStoragePath)
	}
	return resolveLocalStoragePath(basePath, fileURL)
}

func ResolveLocalStorageDir(baseLocalPath, relativeDir string) (string, bool) {
	baseLocalPath = NormalizeLocalStoragePath(baseLocalPath)
	publicRoot := LocalStoragePublicRoot(baseLocalPath)
	targetPath := resolveLocalStoragePath(baseLocalPath, relativeDir)
	_, ok := relativePathWithin(publicRoot, targetPath)
	return targetPath, ok
}

func sanitizeStoragePathSegment(part string) []string {
	part = strings.Trim(strings.TrimSpace(part), "/\\")
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

func resolveLocalStoragePath(basePath string, parts ...string) string {
	all := make([]string, 0, len(parts)+1)
	all = append(all, NormalizeLocalStoragePath(basePath))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		all = append(all, filepath.FromSlash(strings.ReplaceAll(part, "\\", "/")))
	}
	return filepath.Clean(filepath.Join(all...))
}

func relativePathWithin(rootPath, targetPath string) (string, bool) {
	rootPath = filepath.Clean(rootPath)
	targetPath = filepath.Clean(targetPath)
	relPath, err := filepath.Rel(rootPath, targetPath)
	if err != nil {
		return "", false
	}
	if relPath == "." {
		return "", true
	}
	if relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.ToSlash(relPath), true
}

func joinPublicURL(prefix, relPath string) string {
	prefix = strings.TrimRight(strings.TrimSpace(prefix), "/")
	relPath = strings.Trim(strings.TrimSpace(filepath.ToSlash(relPath)), "/")
	if relPath == "" {
		if prefix == "" {
			return "/"
		}
		return prefix
	}
	if prefix == "" {
		return "/" + relPath
	}
	return prefix + "/" + relPath
}
