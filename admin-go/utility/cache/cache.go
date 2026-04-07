package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

func GetJSON(ctx context.Context, key string, dst any) (bool, error) {
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
		_, _ = client.Del(ctx, key)
		return false, err
	}
	return true, nil
}

func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	client := safeRedis(ctx)
	if client == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.SetEX(ctx, key, string(data), ttlSeconds(ttl))
}

func Delete(ctx context.Context, keys ...string) error {
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
	return strconv.ParseInt(raw, 10, 64)
}

func IncrWithTTL(ctx context.Context, key string, ttl time.Duration) (count int64, err error) {
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
