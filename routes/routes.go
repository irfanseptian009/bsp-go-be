package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/irfanseptian/fims-backend/config"
	"github.com/irfanseptian/fims-backend/handlers"
	"github.com/irfanseptian/fims-backend/middleware"
	"github.com/irfanseptian/fims-backend/models"
	"github.com/irfanseptian/fims-backend/services"
)

// Setup configures all API routes with their handlers and middleware.
func Setup(router *gin.Engine, cfg *config.Config) {
	// ─── Initialize Services ───
	authService := services.NewAuthService(cfg)
	userService := services.NewUserService()
	branchService := services.NewBranchService()
	occupationTypeService := services.NewOccupationTypeService()
	insuranceRequestService := services.NewInsuranceRequestService()
	policyService := services.NewPolicyService()

	// ─── Initialize Handlers ───
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	branchHandler := handlers.NewBranchHandler(branchService)
	occupationTypeHandler := handlers.NewOccupationTypeHandler(occupationTypeService)
	insuranceRequestHandler := handlers.NewInsuranceRequestHandler(insuranceRequestService)
	policyHandler := handlers.NewPolicyHandler(policyService)

	// ─── Health Check ───
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ─── API Group ───
	api := router.Group("/api")

	// ─── Public Routes (No Auth Required) ───
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ─── Protected Routes (JWT Required) ───
	protected := api.Group("")
	protected.Use(middleware.Auth(cfg))
	{
		// ── Users ──
		users := protected.Group("/users")
		{
			users.GET("/me", userHandler.GetProfile)
			users.PATCH("/me", userHandler.UpdateProfile)
		}

		// ── Branches ──
		branches := protected.Group("/branches")
		{
			branches.GET("", branchHandler.FindAll)
			branches.GET("/:id", branchHandler.FindByID)
			branches.POST("", middleware.RequireRole(models.RoleAdmin), branchHandler.Create)
			branches.PATCH("/:id", middleware.RequireRole(models.RoleAdmin), branchHandler.Update)
			branches.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), branchHandler.Delete)
		}

		// ── Occupation Types ──
		occupationTypes := protected.Group("/occupation-types")
		{
			occupationTypes.GET("", occupationTypeHandler.FindAll)
			occupationTypes.GET("/:id", occupationTypeHandler.FindByID)
			occupationTypes.POST("", middleware.RequireRole(models.RoleAdmin), occupationTypeHandler.Create)
			occupationTypes.PATCH("/:id", middleware.RequireRole(models.RoleAdmin), occupationTypeHandler.Update)
			occupationTypes.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), occupationTypeHandler.Delete)
		}

		// ── Insurance Requests ──
		insuranceRequests := protected.Group("/insurance-requests")
		{
			insuranceRequests.POST("", middleware.RequireRole(models.RoleCustomer), insuranceRequestHandler.Create)
			insuranceRequests.GET("/my-requests", middleware.RequireRole(models.RoleCustomer), insuranceRequestHandler.FindMyRequests)
			insuranceRequests.GET("/invoice/:invoiceNumber", insuranceRequestHandler.FindByInvoiceNumber)
			insuranceRequests.GET("", middleware.RequireRole(models.RoleAdmin), insuranceRequestHandler.FindAll)
			insuranceRequests.GET("/:id", insuranceRequestHandler.FindByID)
			insuranceRequests.PATCH("/:id/approve", middleware.RequireRole(models.RoleAdmin), insuranceRequestHandler.Approve)
			insuranceRequests.PATCH("/:id/reject", middleware.RequireRole(models.RoleAdmin), insuranceRequestHandler.Reject)
		}

		// ── Policies ──
		policies := protected.Group("/policies")
		{
			policies.GET("", policyHandler.FindAll)
			policies.GET("/:id", policyHandler.FindByID)
			policies.POST("", middleware.RequireRole(models.RoleAdmin), policyHandler.Create)
			policies.PATCH("/:id", middleware.RequireRole(models.RoleAdmin), policyHandler.Update)
			policies.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), policyHandler.Delete)
		}
	}
}
