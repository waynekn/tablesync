package collab

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store is a struct that holds a Redis client for managing collaborative editing sessions.
type Store struct {
	rdb *redis.Client
}

// NewStore creates a new Store instance with the provided Redis client.
// It initializes the Store with the Redis client for managing collaborative editing sessions.
func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

// SheetExists checks if a collaborative editing session for the given sheet ID exists in Redis.
func (s *Store) SheetExists(sheetID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := s.rdb.HGetAll(ctx, sheetID).Result()
	if err != nil {
		slog.Error("failed to check if sheet HSET exists in redis", "err", err)
		return false, fmt.Errorf("could not get sheet from redis: %w", err)
	}

	if len(result) == 0 {
		return false, nil
	}

	return true, nil
}
