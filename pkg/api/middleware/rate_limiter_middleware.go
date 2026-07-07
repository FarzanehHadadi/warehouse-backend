package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"warehouse/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimiterConfig struct {
	Limit  int
	Window time.Duration
	Prefix string
}

func RateLimiter(rdb *redis.Client, cfg RateLimiterConfig) gin.HandlerFunc {
	if cfg.Limit <= 0 {
		cfg.Limit = 10
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "rate"
	}

	return func(c *gin.Context) {
		ctx := context.Background()

		ip := c.ClientIP()
		key := fmt.Sprintf("%s:%s", cfg.Prefix, ip)

		// INCR counter
		count, err := rdb.Incr(ctx, key).Result()
		logger.Log.Info("count", zap.Int64("count", count), zap.String("key", key))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "rate limiter error",
			})
			return
		}

		// Set expiration only on first request
		if count == 1 {
			rdb.Expire(ctx, key, cfg.Window)
		}

		remaining := cfg.Limit - int(count)

		// Headers (important for production APIs)
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, remaining)))

		if count > int64(cfg.Limit) {
			ttl, _ := rdb.TTL(ctx, key).Result()

			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
			return
		}

		c.Next()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
