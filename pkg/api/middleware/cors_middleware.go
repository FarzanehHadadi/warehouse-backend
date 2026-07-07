package middleware

import (
	"strings"
	"time"
	"warehouse/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// allowedCORSRequestHeaders are headers clients may send on API requests.
var allowedCORSRequestHeaders = []string{
	"Authorization",
	"Content-Type",
	"Accept",
	"Origin",
	"X-Requested-With",
	"X-API-Key",
	"X-Admin-Key",
}

// exposedCORSResponseHeaders are headers browser clients are allowed to read.
var exposedCORSResponseHeaders = []string{
	"Content-Disposition",
	"Content-Type",
	"Content-Length",
	"X-RateLimit-Limit",
	"X-RateLimit-Remaining",
	"Retry-After",
}

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     parseAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     allowedCORSRequestHeaders,
		ExposeHeaders:    exposedCORSResponseHeaders,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func parseAllowedOrigins() []string {
	defaultOrigins := "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173"
	raw := utils.GetEnv("CORS_ALLOWED_ORIGINS", defaultOrigins)

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
