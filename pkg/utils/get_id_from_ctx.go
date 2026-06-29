package utils

import "github.com/gin-gonic/gin"

// GetIDFromContext returns the entity ID set by IDMiddleware from the URL path.
func GetIDFromContext(c *gin.Context) uint {
	if id, exists := c.Get("id"); exists {
		return id.(uint)
	}
	return 0
}

// GetUserIDFromContext returns the authenticated user ID set by JwtAuth middleware.
func GetUserIDFromContext(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		return id.(uint)
	}
	return 0
}
