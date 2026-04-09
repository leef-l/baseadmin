package shared

import (
	"fmt"
	"net/url"
	"strings"

	"gbaseadmin/app/upload/internal/model/entity"
)

func MatchUploadConfigByURL(configs []*entity.UploadConfig, storage int, fileURL string) (*entity.UploadConfig, string) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return nil, ""
	}
	for _, config := range configs {
		if config == nil {
			continue
		}
		switch storage {
		case 2:
			if objectKey, ok := MatchOSSObjectKeyParsed(parsedURL, config); ok {
				return config, objectKey
			}
		case 3:
			if objectKey, ok := MatchCOSObjectKeyParsed(parsedURL, config); ok {
				return config, objectKey
			}
		}
	}
	return nil, ""
}

func MatchOSSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", false
	}
	return MatchOSSObjectKeyParsed(parsedURL, config)
}

func MatchOSSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	if config == nil || config.OssBucket == "" || config.OssEndpoint == "" {
		return "", false
	}
	if parsedURL == nil {
		return "", false
	}
	expectedHost := strings.ToLower(fmt.Sprintf("%s.%s", normalizeHostPart(config.OssBucket), normalizeHostPart(config.OssEndpoint)))
	if strings.ToLower(parsedURL.Hostname()) != expectedHost {
		return "", false
	}
	objectKey := objectKeyFromPath(parsedURL.Path)
	if objectKey == "" {
		return "", false
	}
	return objectKey, true
}

func MatchCOSObjectKey(fileURL string, config *entity.UploadConfig) (string, bool) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", false
	}
	return MatchCOSObjectKeyParsed(parsedURL, config)
}

func MatchCOSObjectKeyParsed(parsedURL *url.URL, config *entity.UploadConfig) (string, bool) {
	if config == nil || config.CosBucket == "" || config.CosRegion == "" {
		return "", false
	}
	if parsedURL == nil {
		return "", false
	}
	expectedHost := strings.ToLower(fmt.Sprintf("%s.cos.%s.myqcloud.com", normalizeHostPart(config.CosBucket), normalizeHostPart(config.CosRegion)))
	if strings.ToLower(parsedURL.Hostname()) != expectedHost {
		return "", false
	}
	objectKey := objectKeyFromPath(parsedURL.Path)
	if objectKey == "" {
		return "", false
	}
	return objectKey, true
}

func normalizeHostPart(value string) string {
	return strings.TrimSpace(value)
}

func objectKeyFromPath(path string) string {
	objectKey := strings.TrimPrefix(path, "/")
	if objectKey == "" {
		return ""
	}
	if decoded, err := url.PathUnescape(objectKey); err == nil && decoded != "" {
		return decoded
	}
	return objectKey
}
