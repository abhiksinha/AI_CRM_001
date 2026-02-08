package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// AllowRate increments the counter for key and returns whether the limit is respected.
func AllowRate(ctx context.Context, client *goredis.Client, key string, limit int64, window time.Duration) (bool, error) {
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		_, _ = client.Expire(ctx, key, window).Result()
	}
	return count <= limit, nil
}
