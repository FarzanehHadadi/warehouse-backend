package api

import (
	"warehouse/pkg/api/handlers"
	"warehouse/pkg/api/middleware"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	handler *handlers.Handler
}

func NewRouter(repo *repository.Repository) *Router {
	r := &Router{
		Engine:  gin.Default(),
		handler: handlers.NewHandler(repo),
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	v1 := r.Group("/v1")

	v1.GET("/ping", CheckHealthHandler)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", r.handler.HandleLogin)

	}
	categories := v1.Group("/categories")
	{
		categories.POST("/", r.handler.HandlePostCategory)
		categories.GET("/", r.handler.HandleGetListCategories)
		categoriesWithId := categories.Group("/:id", middleware.IDMiddleware())
		{
			categoriesWithId.GET("/", r.handler.HandleGetCategory)
			categoriesWithId.DELETE("/", r.handler.HandleDeleteCategory)
			categoriesWithId.PATCH("/", r.handler.HandlePatchCategory)
		}

	}
}
