package middleware

import "github.com/gin-gonic/gin"

type IDUri struct {
	ID uint `uri:"id" binding:"required,numeric,min=1"`
}

func IDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri IDUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error":   "invalid id",
				"message": "ID must be a positive integer",
			})
			return
		}
		c.Set("id", uri.ID)
		c.Next()
	}
}
