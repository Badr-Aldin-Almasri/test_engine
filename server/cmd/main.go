package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/visual-api-testing-platform/server/internal/engine"
	"github.com/visual-api-testing-platform/server/internal/handlers"
	"github.com/visual-api-testing-platform/server/internal/repository"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://visual_testing:visual_testing_pass@localhost:5433/visual_testing_db?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	flowRepo := repository.NewFlowRepository(pool)
	testRunRepo := repository.NewTestRunRepository(pool)

	// Initialize execution hub
	hub := engine.NewExecutionHub()
	go hub.Run()

	// Initialize flow runner
	flowRunner := engine.NewFlowRunner(hub)

	// Initialize handlers
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	authHandler := handlers.NewAuthHandler(userRepo, jwtSecret)
	flowHandler := handlers.NewFlowHandler(flowRepo)
	nodeHandler := handlers.NewNodeHandler(flowRepo)
	testRunHandler := handlers.NewTestRunHandler(testRunRepo, flowRepo, flowRunner)
	wsHandler := handlers.NewWebSocketHandler(hub)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authHandler.AuthMiddleware())
		{
			// Flows
			flows := protected.Group("/flows")
			{
				flows.POST("", flowHandler.CreateFlow)
				flows.GET("", flowHandler.ListFlows)
				flows.GET("/:id", flowHandler.GetFlow)
				flows.PUT("/:id", flowHandler.UpdateFlow)
				flows.DELETE("/:id", flowHandler.DeleteFlow)

				// Test runs
				flows.POST("/:id/run", testRunHandler.RunFlow)
				flows.GET("/:id/test-runs", testRunHandler.GetTestRunsByFlow)
			}

			// Nodes
			nodes := protected.Group("/nodes")
			{
				nodes.POST("/:flowId/:nodeId/execute", nodeHandler.ExecuteNode)
			}

			// Test runs
			testRuns := protected.Group("/test-runs")
			{
				testRuns.GET("/:id", testRunHandler.GetTestRun)
			}

			// WebSocket
			protected.GET("/ws", wsHandler.HandleWebSocket)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

