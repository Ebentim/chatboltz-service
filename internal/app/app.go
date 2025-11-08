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
	"github.com/alpinesboltltd/boltz-ai/internal/provider/smtp"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/scraper"
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
	smtpConfig := smtp.Config{Host: cfg.SMTP_HOST, Port: cfg.SMTP_PORT, User: cfg.SMTP_USER, Pass: cfg.SMTP_PASS}
	smtpClient := smtp.NewClient(smtpConfig)
	emailService, err := smtp.NewEmailService(smtpConfig)
	if err != nil {
		log.Fatal("Failed to initialize email service:", err)
	}
	userUsecase := usecase.NewUserUsecase(userRepo, firebaseService, smtpClient)
	agentUsecase := usecase.NewAgentUseCase(agentRepo)
	systemUsecase := usecase.NewSystemUsecase(systemRepo)
	otpUsecase := usecase.NewOTPUsecase(repository.NewUserToken(db), userRepo, 10*time.Minute)
	trainingUsecase, err := usecase.NewTrainingUseCase(cfg.COHERE_API_KEY, cfg.OPENAI_API_KEY, cfg.GOOGLE_API_KEY, db, agentRepo)
	if err != nil {
		log.Fatal("Failed to initialize training usecase:", err)
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userUsecase, []byte(cfg.JWT_SECRET))
	agentHandler := handler.NewAgentHandler(agentUsecase)
	systemHandler := handler.NewSystemHandler(systemUsecase)
	otpHandler := handler.NewOTPHandler(otpUsecase, emailService)
	trainingHandler := handler.NewTrainingHandler(trainingUsecase)

	// Initialize scraper service + handler
	scraperService := scraper.NewService(nil)
	scraperHandler := handler.NewScraperHandler(scraperService)

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

		r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status":"ok",
		})
	})

	api := r.Group("/api/v1")
	
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignupWithEmail)
			auth.POST("/login", authHandler.LoginWithEmail)
			auth.POST("/verify", authHandler.AuthenticateWithToken)
		}

		// OTP routes
		otp := api.Group("/otp")
		{
			// Public routes
			otp.POST("/request", otpHandler.RequestOTP)
			otp.POST("/verify", otpHandler.VerifyOTP)
			otp.POST("/password-reset/complete", otpHandler.CompletePasswordReset)
			otp.POST("/login/complete", otpHandler.CompleteOTPLogin)

			// Protected routes
			protected := otp.Group("/")
			protected.Use(middleware.AuthMiddleware([]byte(cfg.JWT_SECRET)))
			{
				protected.POST("/enable", authHandler.EnableOTP)
				protected.POST("/disable", authHandler.DisableOTP)
				protected.POST("/2fa/enable", otpHandler.Enable2FA)
				protected.POST("/2fa/disable", otpHandler.Disable2FA)
			}
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

			// Agent training
			agent.POST("/:agentId/train/text", trainingHandler.TrainWithText)
			agent.POST("/:agentId/train/file", trainingHandler.TrainWithFile)
			agent.POST("/:agentId/train/url", trainingHandler.TrainWithURL)
			agent.GET("/:agentId/training/documents", trainingHandler.GetTrainingDocuments)
			agent.GET("/:agentId/training/stats", trainingHandler.GetTrainingStats)
			agent.POST("/:agentId/training/query", trainingHandler.QueryKnowledgeBase)
			agent.DELETE("/:agentId/training", trainingHandler.DeleteTrainingData)
			agent.POST("/:agentId/training/migrate", trainingHandler.MigrateLegacyTraining)
		}

		// Scraper endpoint (public). Accepts JSON {url, trace, exclude, max_pages}
		api.POST("/scrape", scraperHandler.Scrape)

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
