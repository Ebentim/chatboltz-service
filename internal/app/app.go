package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/config"
	"github.com/alpinesboltltd/boltz-ai/internal/crypto"
	"github.com/alpinesboltltd/boltz-ai/internal/handler"
	"github.com/alpinesboltltd/boltz-ai/internal/middleware"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// initialize Encryption service
	crypto.NewEncryptionKey([]byte(cfg.GCM_KEY))

	// Initialize database
	db, err := repository.InitDB(cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database handle:", err)
	}
	defer sqlDB.Close()

	// Initialize Firebase
	// firebaseService, err := usecase.NewFirebaseService(cfg.FIREBASE_PROJECT_ID, cfg.FIREBASE_CREDENTIALS)
	firebaseService, err := usecase.NewFirebaseService(cfg.FIREBASE_SERVICE_ACCOUNT)
	if err != nil {
		log.Fatal("Failed to initialize Firebase:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	systemRepo := repository.NewSystemRepository(db)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, firebaseService)
	agentUsecase := usecase.NewAgentUseCase(agentRepo)
	systemUsecase := usecase.NewSystemUsecase(systemRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUsecase, []byte(cfg.JWT_SECRET))
	agentHandler := handler.NewAgentHandler(agentUsecase)
	systemHandler := handler.NewSystemHandler(systemUsecase)

	// Setup routes
	r := gin.Default()

	// Shutdown middleware
	shuttingDown := false
	r.Use(func(c *gin.Context) {
		if shuttingDown {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Service Unavailable",
				"message": "The server is currently shutting down. Please try again later.",
				"code":    503,
			})
			c.Abort()
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignupWithEmail)
			auth.POST("/login", authHandler.LoginWithEmail)
			auth.POST("/verify", authHandler.AuthenticateWithToken)
		}

		agent := api.Group("/agent")
		agent.Use(middleware.AuthMiddleware([]byte(cfg.JWT_SECRET)))
		{
			// Core agent operations
			agent.POST("/create", agentHandler.CreateAgent)
			agent.PATCH("/update/:agentId", agentHandler.UpdateAgent)
			agent.GET("/:agentId", agentHandler.GetAgent)
			agent.GET("/agents/:userId", agentHandler.GetAgentByUser)
			agent.DELETE("/:agentId", agentHandler.DeleteAgent)

			// Agent appearance
			agent.POST("/create/appearance", agentHandler.CreateAgentAppearance)
			agent.GET("/:agentId/appearance", agentHandler.GetAgentAppearance)
			agent.PATCH("/:agentId/appearance", agentHandler.UpdateAgentAppearance)
			agent.DELETE("/:agentId/appearance", agentHandler.DeleteAgentAppearance)

			// Agent behavior
			agent.POST("/create/behavior", agentHandler.CreateAgentBehavior)
			agent.GET("/:agentId/behavior", agentHandler.GetAgentBehavior)
			agent.PATCH("/:agentId/behavior", agentHandler.UpdateAgentBehavior)
			agent.DELETE("/:agentId/behavior", agentHandler.DeleteAgentBehavior)

			// Agent channel
			agent.POST("/create/channel", agentHandler.CreateAgentChannel)
			agent.GET("/:agentId/channel", agentHandler.GetAgentChannel)
			agent.PATCH("/:agentId/channel", agentHandler.UpdateAgentChannel)
			agent.DELETE("/:agentId/channel", agentHandler.DeleteAgentChannel)

			// Agent stats
			agent.GET("/:agentId/stats", agentHandler.GetAgentStats)
			agent.DELETE("/:agentId/stats", agentHandler.DeleteAgentStats)

			// Agent integration
			agent.POST("/create/integration", agentHandler.CreateAgentIntegration)
			agent.GET("/:agentId/integration", agentHandler.GetAgentIntegration)
			agent.PATCH("/:agentId/integration", agentHandler.UpdateAgentIntegration)
			agent.DELETE("/:agentId/integration", agentHandler.DeleteAgentIntegration)
		}

		system := api.Group("/system")
		system.Use(middleware.AuthMiddleware([]byte(cfg.JWT_SECRET)))
		{
			// System instructions (SuperAdmin only)
			system.POST("/instructions", systemHandler.CreateSystemInstruction)
			system.GET("/instructions/:id", systemHandler.GetSystemInstruction)
			system.PATCH("/instructions/:id", systemHandler.UpdateSystemInstruction)
			system.DELETE("/instructions/:id", systemHandler.DeleteSystemInstruction)
			system.GET("/instructions", systemHandler.ListSystemInstructions)

			// Prompt templates (SuperAdmin only)
			system.POST("/templates", systemHandler.CreatePromptTemplate)
			system.GET("/templates/:id", systemHandler.GetPromptTemplate)
			system.GET("/templates", systemHandler.ListPromptTemplates)
		}
	}
	ws := r.Group("/ws/v1")
	{
		chat := ws.Group("/chat")
		{
			chat.GET("/", agentHandler.GetAgent)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	shuttingDown = true

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
