package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/irfanseptian/fims-backend/docs"
	"github.com/irfanseptian/fims-backend/config"
	"github.com/irfanseptian/fims-backend/database"
	"github.com/irfanseptian/fims-backend/middleware"
	"github.com/irfanseptian/fims-backend/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title FIMS Backend API (Go)
// @version 1.0
// @description API documentation for Fire Insurance Management System (Go backend).
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer <token>

func main() {
	// Load configuration
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set")
	}

	// Connect to database & run migrations
	database.Connect(cfg)

	// Run seed if SEED=true
	if os.Getenv("SEED") == "true" {
		database.Seed()
	}

	// Setup Gin router
	router := gin.Default()

	// Global middleware
	router.Use(middleware.CORS())

	// Register routes
	routes.Setup(router, cfg)
	if shouldEnableSwagger() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start server
	port := cfg.Port
	log.Printf("🔥 FIMS Backend (Go) running on http://localhost:%s", port)
	if shouldEnableSwagger() {
		log.Printf("📘 Swagger Dashboard: http://localhost:%s/swagger/index.html", port)
	} else {
		log.Println("🔒 Swagger disabled for production")
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func shouldEnableSwagger() bool {
	enableSwagger := strings.TrimSpace(strings.ToLower(os.Getenv("ENABLE_SWAGGER")))
	if enableSwagger != "" {
		return enableSwagger == "true" || enableSwagger == "1" || enableSwagger == "yes"
	}

	appEnv := strings.TrimSpace(strings.ToLower(os.Getenv("APP_ENV")))
	return appEnv == "" || appEnv == "local" || appEnv == "development"
}
