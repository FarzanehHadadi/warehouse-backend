package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitGroup creates a child router group with rate limiting applied.
func RateLimitGroup(parent *gin.RouterGroup, rdb *redis.Client, cfg RateLimiterConfig, relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	group := parent.Group(relativePath, handlers...)
	group.Use(RateLimiter(rdb, cfg))
	return group
}

// WithPrefix returns a copy of cfg with the given Redis key prefix.
func WithPrefix(cfg RateLimiterConfig, prefix string) RateLimiterConfig {
	cfg.Prefix = prefix
	return cfg
}

// Common rate limit configs
var (
	StrictRateLimit = RateLimiterConfig{
		Limit:  10,
		Window: time.Minute,
		Prefix: "strict",
	}

	MediumRateLimit = RateLimiterConfig{
		Limit:  80,
		Window: time.Minute,
		Prefix: "medium",
	}

	NormalRateLimit = RateLimiterConfig{
		Limit:  150,
		Window: time.Minute,
		Prefix: "normal",
	}

	GlobalRateLimit = RateLimiterConfig{
		Limit:  150,
		Window: time.Minute,
		Prefix: "global",
	}

	ProtectedRateLimit = RateLimiterConfig{
		Limit:  120,
		Window: time.Minute,
		Prefix: "protected",
	}

	RegisterRateLimit = RateLimiterConfig{
		Limit:  5,
		Window: time.Minute,
		Prefix: "register",
	}

	HighRateLimit = RateLimiterConfig{
		Limit:  300,
		Window: time.Minute,
		Prefix: "high",
	}
)
