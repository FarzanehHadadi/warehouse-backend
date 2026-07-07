package database

import (
	"context"
	"fmt"
	"time"
	"warehouse/pkg/utils"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	var redisClient *redis.Client
	addr := utils.GetEnv("REDIS_ADDR", "localhost:6379")
	password := utils.GetEnv("REDIS_PASSWORD", "")

	redisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password, // leave empty if no password
		DB:           0,        // default DB
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		return nil, err
	}

	fmt.Println("✅ Redis connected successfully")
	return redisClient, nil
}
