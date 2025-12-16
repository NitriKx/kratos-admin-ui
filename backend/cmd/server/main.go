package main

import (
	"log"
	"os"

	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/auth"
	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/config"
	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/handlers"
	"github.com/benoit-sauvere/kratos-admin-ui/backend/internal/kratos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Kratos client
	kratosClient := kratos.NewClient(cfg.KratosAdminURL)
	kratosClient.SetPublicURL(cfg.KratosPublicURL)

	// Initialize handlers
	authHandler := auth.NewHandler(cfg)
	identitiesHandler := handlers.NewIdentitiesHandler(kratosClient)
	sessionsHandler := handlers.NewSessionsHandler(kratosClient)
	schemasHandler := handlers.NewSchemasHandler(kratosClient)
	statsHandler := handlers.NewStatsHandler(kratosClient)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	log.Printf("CORS allowed origins: %v", cfg.CORSOrigins)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint (for Kubernetes probes)
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(auth.JWTMiddleware(cfg.JWTSecret))
	{
		// Identities
		protected.GET("/identities", identitiesHandler.List)
		protected.GET("/identities/:id", identitiesHandler.Get)
		protected.POST("/identities", identitiesHandler.Create)
		protected.PUT("/identities/:id", identitiesHandler.Update)
		protected.DELETE("/identities/:id", identitiesHandler.Delete)
		protected.GET("/identities/:id/sessions", identitiesHandler.GetSessions)

		// Sessions
		protected.GET("/sessions", sessionsHandler.List)
		protected.DELETE("/sessions/:id", sessionsHandler.Revoke)

		// Schemas
		protected.GET("/schemas", schemasHandler.List)

		// Stats
		protected.GET("/stats", statsHandler.Get)
	}

	// Get port from config or default
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
