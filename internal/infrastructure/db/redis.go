package db

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func GetRedisInstance(uri string, dbName string) *redis.Client {
	redisOnce.Do(func() {
		opts, err := redis.ParseURL(uri)
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		db, _ := strconv.ParseInt(dbName, 10, 64)
		opts.DB = int(db)
		redisClient = redis.NewClient(opts)

		// Check the Redis connection
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}

		log.Println("Connected to Redis successfully")
	})

	return redisClient
}
