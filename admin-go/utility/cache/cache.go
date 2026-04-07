package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

func GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	key = normalizeKey(key)
	if key == "" {
		return false, nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return false, nil
	}
	value, err := client.Get(ctx, key)
	if err != nil || value == nil {
		return false, err
	}
	raw := value.String()
	if raw == "" {
		return false, nil
	}
	if err := json.Unmarshal([]byte(raw), dst); err != nil {
		if shouldDeleteInvalidJSON(err) {
			_, _ = client.Del(ctx, key)
		}
		return false, err
	}
	return true, nil
}

func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	key = normalizeKey(key)
	if key == "" {
		return nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return nil
	}
	if value == nil {
		_, err := client.Del(ctx, key)
		return err
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.SetEX(ctx, key, string(data), ttlSeconds(ttl))
}

func Delete(ctx context.Context, keys ...string) error {
	keys = normalizeKeys(keys)
	if len(keys) == 0 {
		return nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return nil
	}
	_, err := client.Del(ctx, keys...)
	return err
}

func GetInt64(ctx context.Context, key string) (int64, error) {
	key = normalizeKey(key)
	if key == "" {
		return 0, nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return 0, nil
	}
	value, err := client.Get(ctx, key)
	if err != nil || value == nil {
		return 0, err
	}
	raw := value.String()
	if raw == "" {
		return 0, nil
	}
	parsed, err := parseCachedInt64(raw)
	if err != nil {
		_, _ = client.Del(ctx, key)
		return 0, err
	}
	return parsed, nil
}

func IncrWithTTL(ctx context.Context, key string, ttl time.Duration) (count int64, err error) {
	key = normalizeKey(key)
	if key == "" {
		return 0, nil
	}
	client := safeRedis(ctx)
	if client == nil {
		return 0, nil
	}
	count, err = client.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	if count == 1 {
		_, _ = client.Expire(ctx, key, ttlSeconds(ttl))
	}
	return count, nil
}

func safeRedis(ctx context.Context) (client *gredis.Redis) {
	defer func() {
		if recovered := recover(); recovered != nil {
			client = nil
			g.Log().Warningf(ctx, "redis unavailable, fallback to direct execution: %v", recovered)
		}
	}()
	return g.Redis()
}

func ttlSeconds(ttl time.Duration) int64 {
	if ttl <= 0 {
		return 1
	}
	seconds := int64(ttl / time.Second)
	if ttl%time.Second != 0 {
		seconds++
	}
	if seconds <= 0 {
		return 1
	}
	return seconds
}

func normalizeKeys(keys []string) []string {
	if len(keys) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(keys))
	normalized := make([]string, 0, len(keys))
	for _, key := range keys {
		key = normalizeKey(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, key)
	}
	return normalized
}

func normalizeKey(key string) string {
	return strings.TrimSpace(key)
}

func shouldDeleteInvalidJSON(err error) bool {
	if err == nil {
		return false
	}
	var invalidUnmarshalErr *json.InvalidUnmarshalError
	return !errors.As(err, &invalidUnmarshalErr)
}

func parseCachedInt64(raw string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
}
