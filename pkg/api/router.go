package api

import (
	"warehouse/docs"
	"warehouse/pkg/api/handlers"
	"warehouse/pkg/api/middleware"
	"warehouse/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
	handler *handlers.Handler
	redis   *redis.Client
}

func NewRouter(repo *repository.Repository, redisClient *redis.Client) *Router {
	r := &Router{
		Engine:  gin.Default(),
		handler: handlers.NewHandler(repo),
		redis:   redisClient,
	}

	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.PersistAuthorization(true)))
	docs.SwaggerInfo.BasePath = "/"

	v1 := r.Group("/v1")
	v1.Use(middleware.ApiKeyAuth())
	v1.Use(middleware.RateLimiter(r.redis, middleware.GlobalRateLimit))
	v1.GET("/ping", CheckHealthHandler)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", middleware.RateLimiter(r.redis, middleware.WithPrefix(middleware.StrictRateLimit, "login")),
			r.handler.HandleLogin)
		auth.POST("/refresh", middleware.RateLimiter(r.redis, middleware.WithPrefix(middleware.StrictRateLimit, "refresh")),
			r.handler.HandleRefreshToken)
		register := middleware.RateLimitGroup(auth, r.redis, middleware.RegisterRateLimit, "/register", middleware.AdminRegistrationKeyAuth())
		{
			register.POST("", r.handler.HandlePostRegister)
		}
	}

	categories := v1.Group("/categories")

	{
		categories.GET("/", r.handler.HandleGetListCategories)
		protectedCategories := middleware.RateLimitGroup(categories, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
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
		protectedUnits := middleware.RateLimitGroup(units, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
		{
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
		protectedDepartments := middleware.RateLimitGroup(departments, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
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
		protectedManagers := middleware.RateLimitGroup(managers, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
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
		protectedProducts := middleware.RateLimitGroup(products, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
		{
			protectedProducts.POST("/", r.handler.HandlePostProduct)
			protectedProducts.GET("/search", r.handler.HandleSearchProductList)
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
		protectedStores := middleware.RateLimitGroup(stores, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
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
	orders := v1.Group("/orders")
	{
		orders.GET("/", r.handler.HandleGetOrderList)
		protectedOrders := middleware.RateLimitGroup(orders, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
		{
			protectedOrders.POST("/", r.handler.HandlePostOrder)
			protectedOrders.GET("/export", r.handler.HandleExportOrder)
			withIdOrders := protectedOrders.Group("/:id", middleware.IDMiddleware())
			{
				withIdOrders.GET("/", r.handler.HandleGetOrder)
				withIdOrders.PATCH("/", r.handler.HandlePatchOrder)
				withIdOrders.DELETE("/", r.handler.HandleDeleteOrder)
			}
		}

	}

	reports := v1.Group("/reports")
	{
		protectedReports := middleware.RateLimitGroup(reports, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
		{
			protectedReports.GET("/threshold-proximity/", r.handler.HandleGetThresholdProximityReport)
			protectedReports.GET("/threshold-proximity/export", r.handler.HandleExportThresholdProximityReport)
			protectedReports.GET("/store-product-quantities/export", r.handler.HandleExportStoreProductQuantitiesReport)
			protectedReports.GET("/store-product-quantities/", r.handler.HandleGetStoreProductQuantitiesReport)
		}
	}
	dashboard := v1.Group("/dashboard")
	protectedDashboard := middleware.RateLimitGroup(dashboard, r.redis, middleware.ProtectedRateLimit, "/", middleware.JwtAuth())
	protectedDashboard.GET("/", r.handler.HandleGetDashboard)
	{
		protectedDashboard.GET("/activities", r.handler.HandleGetRecentActivities)
	}
}
