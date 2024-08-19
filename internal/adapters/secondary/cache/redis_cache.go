package cache

import (
	"context"
	"time"

	"log"

	"github.com/go-redis/redis/v8"
)

// RedisCache struct wraps the Redis client and provides caching functionalities
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new instance of RedisCache
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Set stores a key-value pair in Redis with an expiration time
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := r.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	log.Printf("Key %s set successfully in Redis", key)
	return nil
}

// Get retrieves a value from Redis based on the given key
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Key %s not found in Redis", key)
			return "", nil
		}
		return "", err
	}
	log.Printf("Key %s retrieved successfully from Redis", key)
	return result, nil
}

// Delete removes a key-value pair from Redis
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	log.Printf("Key %s deleted successfully from Redis", key)
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	log.Printf("Key %s existence check: %v", key, result > 0)
	return result > 0, nil
}
