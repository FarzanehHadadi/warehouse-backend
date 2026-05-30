package api

import (
	"warehouse/docs"
	"warehouse/pkg/api/handlers"
	"warehouse/pkg/api/middleware"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/"

	v1 := r.Group("/v1")
	v1.Use(middleware.ApiKeyAuth())
	v1.GET("/ping", CheckHealthHandler)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", r.handler.HandleLogin)

	}

	categories := v1.Group("/categories")

	{
		categories.GET("/", r.handler.HandleGetListCategories)

		categories.POST("/", middleware.JwtAuth(), r.handler.HandlePostCategory)
		categoriesWithId := categories.Group("/:id", middleware.IDMiddleware())
		{
			categoriesWithId.GET("/", r.handler.HandleGetCategory)
			categoriesWithId.DELETE("/", r.handler.HandleDeleteCategory)
			categoriesWithId.PATCH("/", r.handler.HandlePatchCategory)
		}

	}
}
