package api

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		r.GET("/ping", CheckHealthHandler)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", HandleLogin)
		}
	}
	return r
}
