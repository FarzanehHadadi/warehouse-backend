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
		protectedCategories := categories.Group("/", middleware.JwtAuth())
		{
			protectedCategories.POST("/", r.handler.HandlePostCategory)
			categoriesWithId := protectedCategories.Group("/:id", middleware.IDMiddleware())
			{
				categoriesWithId.GET("/", r.handler.HandleGetCategory)
				categoriesWithId.DELETE("/", r.handler.HandleDeleteCategory)
				categoriesWithId.PATCH("/", r.handler.HandlePatchCategory)
			}
		}
	}
	units := v1.Group("/units")
	{
		units.GET("/", r.handler.HandleGetUnitList)
		protectedUnits := units.Group("/", middleware.JwtAuth())
		{
			protectedUnits.Use()
			protectedUnits.POST("/", r.handler.HandlePostUnit)
			unitsWithId := protectedUnits.Group("/:id", middleware.IDMiddleware())
			{
				unitsWithId.GET("/", r.handler.HandleGetUnitById)
				unitsWithId.PATCH("/", r.handler.HandlePatchUnit)
				unitsWithId.DELETE("/", r.handler.HandleDeleteUnit)
			}
		}
	}
	departments := v1.Group("/departments")
	{
		departments.GET("/", r.handler.HandleGetDepartmentList)
		protectedDepartments := departments.Group("/", middleware.JwtAuth())
		{
			protectedDepartments.POST("/", r.handler.HandlePostDepartment)
			withIdDepartments := protectedDepartments.Group("/:id", middleware.IDMiddleware())
			{
				withIdDepartments.GET("/", r.handler.HandleGetDepartment)
				withIdDepartments.PATCH("/", r.handler.HandlePatchDepartment)
				withIdDepartments.DELETE("/", r.handler.HandleDeleteDepartment)

			}
		}
	}
}
