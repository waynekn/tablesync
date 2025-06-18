package collab

import (
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
