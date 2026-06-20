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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))
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
	managers := v1.Group("/managers")
	{
		managers.GET("/", r.handler.HandleGetManagerList)
		protectedManagers := managers.Group("/", middleware.JwtAuth())
		{
			protectedManagers.POST("/", r.handler.HandlePostManager)
			withIdManagers := protectedManagers.Group("/:id", middleware.IDMiddleware())
			{
				withIdManagers.GET("/", r.handler.HandleGetManager)
				withIdManagers.PATCH("/", r.handler.HandlePatchManager)
				withIdManagers.DELETE("/", r.handler.HandleDeleteManager)
			}
		}

	}
	products := v1.Group("/products")
	{
		products.GET("/", r.handler.HandleGetProductList)
		protectedProducts := products.Group("/", middleware.JwtAuth())
		{
			protectedProducts.POST("/", r.handler.HandlePostProduct)
			withIdProducts := protectedProducts.Group("/:id", middleware.IDMiddleware())
			{
				withIdProducts.GET("/", r.handler.HandleGetProduct)
				withIdProducts.PATCH("/", r.handler.HandlePatchProduct)
				withIdProducts.DELETE("/", r.handler.HandleDeleteProduct)
			}
		}

	}
	stores := v1.Group("/stores")
	{
		stores.GET("/", r.handler.HandleGetStoreList)
		protectedStores := stores.Group("/", middleware.JwtAuth())
		{
			protectedStores.POST("/", r.handler.HandlePostStore)
			withIdStores := protectedStores.Group("/:id", middleware.IDMiddleware())
			{
				withIdStores.GET("/", r.handler.HandleGetStore)
				withIdStores.PATCH("/", r.handler.HandlePatchStore)
				withIdStores.DELETE("/", r.handler.HandleDeleteStore)
			}
		}

	}

}
