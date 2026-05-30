package middleware

import (
	"net/http"
	"strings"
	"warehouse/pkg/api/auth"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
			})
			return
		}
		const BearerSchema = "Bearer "
		paths := strings.Split(authHeader, " ")
		if len(paths) != 2 || paths[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Token format should be: Bearer <token>",
				},
			})
			return

		}
		tokenStr := authHeader[:len(BearerSchema)]
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			return

		}
		if claims.UserId == 0 {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}

}
