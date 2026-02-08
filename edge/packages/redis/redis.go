package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

// NewClient creates a redis client from config.
func NewClient(cfg Config) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// Ping verifies connectivity; caller can decide startup behavior.
func Ping(ctx context.Context, client *goredis.Client) error {
	return client.Ping(ctx).Err()
}
