package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	adminRegistrationKey []byte
	onceLoadAdminKey     sync.Once
)

func loadAdminRegistrationKey() {
	envSecret := os.Getenv("ADMIN_REGISTRATION_KEY")
	if len(envSecret) < 32 {
		panic("Length of admin registration key is too short")
	}
	adminRegistrationKey = []byte(envSecret)
}

func AdminRegistrationKeyAuth() gin.HandlerFunc {
	onceLoadAdminKey.Do(loadAdminRegistrationKey)
	return func(ctx *gin.Context) {

		provided := ctx.GetHeader("X-Admin-Key")
		if subtle.ConstantTimeCompare([]byte(provided), adminRegistrationKey) == 1 {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")

	}
}
