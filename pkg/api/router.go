package api

import (
	"warehouse/pkg/api/handlers"
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
		// auth.POST("/register", r.handler.Register)
		// Add more auth routes here
	}
	categories := v1.Group("/category")
	{
		categories.GET("/{categoryId}", r.handler.HandleGetCategory)
		categories.POST("/", r.handler.HandlePostCategory)
		categories.GET("/", r.handler.HandleGetListCategories)
		categories.DELETE("/{categoryId}", r.handler.HandleDeleteCategory)
		categories.PATCH("/{categoryId}", r.handler.HandlePatchCategory)

	}
	// Example: users group
	// users := v1.Group("/users")
	// {
	//     users.GET("/:id", r.handler.GetUser)
	// }
}
