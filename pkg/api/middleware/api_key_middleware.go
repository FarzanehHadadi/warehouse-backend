package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	apiSecretKey []byte
	onceLoad     sync.Once
)

func loadSecretKey() {
	envSecret := os.Getenv("API_SECRET_KEY")
	if len(envSecret) < 32 {
		panic("Length of api secret key is too short")
	}
	apiSecretKey = []byte(envSecret)
}

func ApiKeyAuth() gin.HandlerFunc {
	onceLoad.Do(loadSecretKey)
	return func(ctx *gin.Context) {

		provided := ctx.GetHeader("X-API-Key")
		if subtle.ConstantTimeCompare([]byte(provided), apiSecretKey) == 1 {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")

	}
}
