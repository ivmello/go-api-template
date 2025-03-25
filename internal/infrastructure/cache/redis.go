package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/ivmello/go-api-template/internal/config"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

// Set stores a value in Redis with expiration
func Set(ctx context.Context, client *redis.Client, key string, value interface{}, expiration int) error {
	return client.Set(ctx, key, value, 0).Err()
}

// Get retrieves a value from Redis
func Get(ctx context.Context, client *redis.Client, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Delete removes a key from Redis
func Delete(ctx context.Context, client *redis.Client, key string) error {
	return client.Del(ctx, key).Err()
}