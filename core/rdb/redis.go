package rdb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Connect establishes a connection to a Redis server using the provided address,
// password, and database index. It returns a *redis.Client instance if successful,
// or an error if the connection fails or if the Redis server is unreachable.
func Connect(addr, password string, db int) (*redis.Client, error) {
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to ping redis server", "addr", addr, "err", err)
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	testKey := "test-key-" + uuid.New().String()
	testValue := "value"
	if err := rdb.Set(ctx, testKey, testValue, 1*time.Second).Err(); err != nil {
		slog.Error("redis client could not set value", "key", testKey, "err", err)
		return nil, fmt.Errorf("redis set failed: %w", err)
	}

	val, err := rdb.Get(ctx, testKey).Result()
	if err != nil {
		slog.Error("redis client could not get value", "key", testKey, "err", err)
		return nil, fmt.Errorf("redis get failed: %w", err)
	}

	if val != testValue {
		slog.Error("redis value mismatch", "key", testKey, "expected", testValue, "got", val)
		return nil, fmt.Errorf("redis value mismatch for key %s: expected %s, got %s", testKey, testValue, val)
	}

	slog.Info("Successfully connected to redis")
	return rdb, nil
}
