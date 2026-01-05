package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/receiv3/backend/internal/config"
	"github.com/receiv3/backend/internal/database"
	"github.com/receiv3/backend/internal/handlers"
	"github.com/receiv3/backend/internal/middleware"
	"github.com/receiv3/backend/internal/repository"
	"github.com/receiv3/backend/internal/services"
	"github.com/receiv3/backend/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Initialize database
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis (optional)
	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	kycRepo := repository.NewKYCRepository(db)
	buyerRepo := repository.NewBuyerRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	fundingRepo := repository.NewFundingRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours, cfg.JWTRefreshExpiryHours)

	// Initialize services
	pinataService := services.NewPinataService(cfg)
	authService := services.NewAuthService(userRepo, jwtManager)
	invoiceService := services.NewInvoiceService(invoiceRepo, buyerRepo, fundingRepo, pinataService, cfg)
	fundingService := services.NewFundingService(fundingRepo, invoiceRepo, txRepo, cfg)
	blockchainService, err := services.NewBlockchainService(cfg, invoiceRepo, pinataService)
	if err != nil {
		log.Printf("Warning: Blockchain service init failed: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo, kycRepo)
	buyerHandler := handlers.NewBuyerHandler(buyerRepo)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, blockchainService)
	fundingHandler := handlers.NewFundingHandler(fundingService)

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Rate limiter (if Redis is available)
	if redisClient != nil {
		rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
		router.Use(rateLimiter.Middleware())
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "receiv3-backend",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.PUT("/wallet", userHandler.UpdateWallet)
				user.POST("/kyc", userHandler.SubmitKYC)
				user.GET("/kyc", userHandler.GetKYCStatus)
			}

			// Buyer routes (exporter only)
			buyers := protected.Group("/buyers")
			buyers.Use(middleware.ExporterOnly())
			{
				buyers.POST("", buyerHandler.Create)
				buyers.GET("", buyerHandler.List)
				buyers.GET("/:id", buyerHandler.Get)
				buyers.PUT("/:id", buyerHandler.Update)
				buyers.DELETE("/:id", buyerHandler.Delete)
			}

			// Invoice routes (exporter only for CRUD)
			invoices := protected.Group("/invoices")
			{
				// Exporter routes
				invoices.POST("", middleware.ExporterOnly(), invoiceHandler.Create)
				invoices.GET("", middleware.ExporterOnly(), invoiceHandler.List)
				invoices.GET("/fundable", invoiceHandler.ListFundable) // Open to all
				invoices.GET("/:id", invoiceHandler.Get)
				invoices.PUT("/:id", middleware.ExporterOnly(), invoiceHandler.Update)
				invoices.DELETE("/:id", middleware.ExporterOnly(), invoiceHandler.Delete)
				invoices.POST("/:id/submit", middleware.ExporterOnly(), invoiceHandler.Submit)
				invoices.POST("/:id/documents", middleware.ExporterOnly(), invoiceHandler.UploadDocument)
				invoices.GET("/:id/documents", invoiceHandler.GetDocuments)
				invoices.POST("/:id/tokenize", middleware.AdminOnly(), invoiceHandler.Tokenize)
				invoices.POST("/:id/pool", middleware.AdminOnly(), fundingHandler.CreatePool)
			}

			// Funding/Investment routes
			pools := protected.Group("/pools")
			{
				pools.GET("", fundingHandler.ListPools)
				pools.GET("/:id", fundingHandler.GetPool)
			}

			investments := protected.Group("/investments")
			investments.Use(middleware.InvestorOnly())
			{
				investments.POST("", fundingHandler.Invest)
				investments.GET("", fundingHandler.GetMyInvestments)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/kyc/pending", userHandler.GetPendingKYC)
				admin.POST("/kyc/:id/approve", userHandler.ApproveKYC)
				admin.POST("/kyc/:id/reject", userHandler.RejectKYC)
				admin.POST("/invoices/:id/approve", invoiceHandler.Approve)
				admin.POST("/invoices/:id/reject", invoiceHandler.Reject)
				admin.POST("/pools/:id/disburse", fundingHandler.Disburse)
				admin.POST("/invoices/:id/repay", fundingHandler.ProcessRepayment)
			}
		}
	}

	// Start server
	log.Printf("Receiv3 Backend starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
