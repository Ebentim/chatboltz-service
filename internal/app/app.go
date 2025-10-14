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
	"github.com/alpinesboltltd/boltz-ai/internal/handler"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
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

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, firebaseService)
	agentUsecase := usecase.NewAgentUseCase(agentRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUsecase)
	agentHandler := handler.NewAgentHandler(agentUsecase)

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
		{
			agent.POST("/create", agentHandler.CreateAgent)
			agent.PATCH("/update/:id", agentHandler.UpdateAgent)
			agent.GET("/:agentId", agentHandler.GetAgent)
			agent.GET("/agents/:userId", agentHandler.GetAgentByUser)
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
