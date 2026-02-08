package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// PersistSession stores session + token mappings with TTL.
func PersistSession(ctx context.Context, client *goredis.Client, userID, sessionID, apiToken string, ttl time.Duration) error {
	sessionSetKey := fmt.Sprintf("auth:sessions:%s", userID)
	sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
	tokenHash := HashValue(apiToken)
	tokenKey := fmt.Sprintf("auth:token:%s", tokenHash)

	pipe := client.TxPipeline()
	pipe.SAdd(ctx, sessionSetKey, sessionID)
	pipe.Set(ctx, sessionKey, tokenHash, ttl)
	pipe.HSet(ctx, tokenKey, map[string]string{
		"user_id":    userID,
		"session_id": sessionID,
	})
	pipe.Expire(ctx, tokenKey, ttl)
	pipe.Expire(ctx, sessionSetKey, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

// ClearSession removes session + token mappings.
func ClearSession(ctx context.Context, client *goredis.Client, userID, sessionID, apiToken string) error {
	sessionSetKey := fmt.Sprintf("auth:sessions:%s", userID)
	sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
	tokenKey := fmt.Sprintf("auth:token:%s", HashValue(apiToken))

	pipe := client.TxPipeline()
	pipe.Del(ctx, tokenKey)
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, sessionSetKey, sessionID)
	_, err := pipe.Exec(ctx)
	return err
}

// ActiveSessions returns active session IDs and prunes stale ones.
func ActiveSessions(ctx context.Context, client *goredis.Client, userID string) ([]string, error) {
	sessionSetKey := fmt.Sprintf("auth:sessions:%s", userID)
	members, err := client.SMembers(ctx, sessionSetKey).Result()
	if err != nil {
		return nil, err
	}

	active := make([]string, 0, len(members))
	for _, sessionID := range members {
		sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
		exists, err := client.Exists(ctx, sessionKey).Result()
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			_ = client.SRem(ctx, sessionSetKey, sessionID).Err()
			continue
		}
		active = append(active, sessionID)
	}
	return active, nil
}

// ExpireSessionByID removes a session by ID if it exists.
func ExpireSessionByID(ctx context.Context, client *goredis.Client, sessionID string) error {
	sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
	tokenHash, err := client.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil
		}
		return err
	}

	tokenKey := fmt.Sprintf("auth:token:%s", tokenHash)
	data, err := client.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		return err
	}
	userID := data["user_id"]
	if userID == "" {
		return nil
	}

	pipe := client.TxPipeline()
	pipe.Del(ctx, tokenKey)
	pipe.Del(ctx, sessionKey)
	pipe.SRem(ctx, fmt.Sprintf("auth:sessions:%s", userID), sessionID)
	_, err = pipe.Exec(ctx)
	return err
}

// RefreshSessionTTL refreshes session-related keys to the provided TTL.
func RefreshSessionTTL(ctx context.Context, client *goredis.Client, userID, sessionID, apiToken string, ttl time.Duration) {
	sessionKey := fmt.Sprintf("auth:session:%s", sessionID)
	tokenKey := fmt.Sprintf("auth:token:%s", HashValue(apiToken))
	sessionSetKey := fmt.Sprintf("auth:sessions:%s", userID)

	pipe := client.TxPipeline()
	pipe.Expire(ctx, tokenKey, ttl)
	pipe.Expire(ctx, sessionKey, ttl)
	pipe.Expire(ctx, sessionSetKey, ttl)
	_, _ = pipe.Exec(ctx)
}

// MigrateRawToken removed: only hashed tokens are supported.
