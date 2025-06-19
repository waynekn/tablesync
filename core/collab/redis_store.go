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

// InitRedisSheet initializes a collaborative editing session in Redis for the given sheet ID.
// It sets the sheet data and an expiration time based on the provided deadline.
// The expiration time is set to 5 minutes after the deadline to allow time for processing and
// storage of the data in the database.
func (s *Store) InitRedisSheet(sheetID string, sheetDeadline time.Time, sheetData *[][]string) error {
	ttl := time.Until(sheetDeadline.Add(5 * time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := s.rdb.TxPipeline()

	// flatten HSETs into one big call per hash
	cells := make(map[string]string)
	for i, row := range *sheetData {
		for j, cell := range row {
			key := fmt.Sprintf("%d:%d", i, j)
			cells[key] = cell
		}
	}

	pipe.HSet(ctx, sheetID, cells)
	pipe.Expire(ctx, sheetID, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Error("failed to initialize sheet in redis", "err", err)
		return fmt.Errorf("could not initialize sheet in redis: %w", err)
	}

	return nil
}
