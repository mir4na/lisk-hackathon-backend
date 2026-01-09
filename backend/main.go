package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/database"
	"github.com/vessel/backend/internal/handlers"
	"github.com/vessel/backend/internal/middleware"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
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
	otpRepo := repository.NewOTPRepository(db)
	mitraRepo := repository.NewMitraRepository(db)
	importerPaymentRepo := repository.NewImporterPaymentRepository(db)
	rqRepo := repository.NewRiskQuestionnaireRepository(db)

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours, cfg.JWTRefreshExpiryHours)

	// Initialize services
	pinataService := services.NewPinataService(cfg)
	emailService := services.NewEmailService(cfg)
	escrowService := services.NewEscrowService()
	otpService := services.NewOTPService(otpRepo, emailService, cfg)
	authService := services.NewAuthService(userRepo, jwtManager, otpService)
	mitraService := services.NewMitraService(mitraRepo, userRepo, emailService, pinataService)
	invoiceService := services.NewInvoiceService(invoiceRepo, buyerRepo, fundingRepo, pinataService, cfg)
	fundingService := services.NewFundingService(fundingRepo, invoiceRepo, txRepo, userRepo, buyerRepo, rqRepo, emailService, escrowService, cfg)
	paymentService := services.NewPaymentService(userRepo, txRepo)
	rqService := services.NewRiskQuestionnaireService(rqRepo)
	blockchainService, err := services.NewBlockchainService(cfg, invoiceRepo, pinataService)
	if err != nil {
		log.Printf("Warning: Blockchain service init failed: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, otpService)
	userHandler := handlers.NewUserHandler(userRepo, kycRepo)
	buyerHandler := handlers.NewBuyerHandler(buyerRepo)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceService, blockchainService)
	fundingHandler := handlers.NewFundingHandler(fundingService)
	mitraHandler := handlers.NewMitraHandler(mitraService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	importerHandler := handlers.NewImporterHandler(importerPaymentRepo, fundingService, fundingRepo, invoiceRepo)
	rqHandler := handlers.NewRiskQuestionnaireHandler(rqService)

	// Initialize wallet middleware
	walletMiddleware := middleware.NewWalletMiddleware(userRepo)

	// Initialize profile middleware
	profileMiddleware := middleware.NewProfileMiddleware(userRepo)

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
			"service": "vessel-backend",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/send-otp", authHandler.SendOTP)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Public routes (no auth required) - for importers to pay
		public := v1.Group("/public")
		{
			public.GET("/payments/:payment_id", importerHandler.GetPaymentInfo)
			public.POST("/payments/:payment_id/pay", importerHandler.Pay)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// User routes (no profile completion required for profile endpoints)
			user := protected.Group("/user")
			{
				// Profile endpoints - don't require profile completion (allow user to complete profile)
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)

				// These require profile completion
				userWithProfile := user.Group("")
				userWithProfile.Use(profileMiddleware.RequireProfileComplete())
				{
					userWithProfile.PUT("/wallet", userHandler.UpdateWallet)
					userWithProfile.POST("/kyc", userHandler.SubmitKYC)
					userWithProfile.GET("/kyc", userHandler.GetKYCStatus)
					userWithProfile.GET("/balance", paymentHandler.GetBalance)

					// MITRA application routes
					mitra := userWithProfile.Group("/mitra")
					{
						mitra.POST("/apply", mitraHandler.Apply)
						mitra.GET("/status", mitraHandler.GetStatus)
						mitra.POST("/documents", mitraHandler.UploadDocument)
					}
				}
			}

			// Payment routes (PROTOTYPE) - require profile completion
			payments := protected.Group("/payments")
			payments.Use(profileMiddleware.RequireProfileComplete())
			{
				payments.POST("/deposit", paymentHandler.Deposit)
				payments.POST("/withdraw", paymentHandler.Withdraw)
				payments.GET("/balance", paymentHandler.GetBalance)
			}

			// Buyer routes (exporter only, read-only - buyer creation removed as buyers don't use the app)
			buyers := protected.Group("/buyers")
			buyers.Use(middleware.ExporterOnly(), profileMiddleware.RequireProfileComplete())
			{
				buyers.GET("", buyerHandler.List)
				buyers.GET("/:id", buyerHandler.Get)
			}

			// Invoice routes (exporter only for CRUD) - require profile completion
			invoices := protected.Group("/invoices")
			invoices.Use(profileMiddleware.RequireProfileComplete())
			{
				// Exporter routes - require wallet for blockchain operations
				invoices.POST("", middleware.ExporterOnly(), walletMiddleware.RequireWallet(), invoiceHandler.Create)
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

			// Funding/Investment routes (Marketplace) - require profile completion
			pools := protected.Group("/pools")
			pools.Use(profileMiddleware.RequireProfileComplete())
			{
				pools.GET("", fundingHandler.ListPools)
				pools.GET("/:id", fundingHandler.GetPool)
			}

			// Marketplace routes (with filters) - require profile completion
			marketplace := protected.Group("/marketplace")
			marketplace.Use(profileMiddleware.RequireProfileComplete())
			{
				marketplace.GET("", fundingHandler.GetMarketplace)
			}

			// Risk Questionnaire routes (for investors) - require profile completion
			riskQuestionnaire := protected.Group("/risk-questionnaire")
			riskQuestionnaire.Use(middleware.InvestorOnly(), profileMiddleware.RequireProfileComplete())
			{
				riskQuestionnaire.GET("/questions", rqHandler.GetQuestions)
				riskQuestionnaire.POST("", rqHandler.Submit)
				riskQuestionnaire.GET("/status", rqHandler.GetStatus)
			}

			// Investment routes - require wallet and profile completion
			investments := protected.Group("/investments")
			investments.Use(middleware.InvestorOnly(), profileMiddleware.RequireProfileComplete(), walletMiddleware.RequireWallet())
			{
				investments.POST("", fundingHandler.Invest)
				investments.GET("", fundingHandler.GetMyInvestments)
				investments.GET("/portfolio", fundingHandler.GetPortfolio)
			}

			// Exporter routes - require profile completion
			exporter := protected.Group("/exporter")
			exporter.Use(middleware.ExporterOnly(), profileMiddleware.RequireProfileComplete())
			{
				exporter.POST("/disbursement", fundingHandler.ExporterDisbursement)
			}

			// Mitra Dashboard - require profile completion
			mitraDashboard := protected.Group("/mitra")
			mitraDashboard.Use(middleware.ExporterOnly(), profileMiddleware.RequireProfileComplete())
			{
				mitraDashboard.GET("/dashboard", fundingHandler.GetMitraDashboard)
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
				admin.POST("/pools/:id/close", fundingHandler.ClosePoolAndNotify)
				admin.POST("/invoices/:id/repay", fundingHandler.ProcessRepayment)

				// Admin Mitra Application routes
				admin.GET("/mitra/pending", mitraHandler.GetPendingApplications)
				admin.POST("/mitra/:id/approve", mitraHandler.Approve)
				admin.POST("/mitra/:id/reject", mitraHandler.Reject)

				// Admin Balance Management (MVP)
				admin.POST("/balance/grant", paymentHandler.AdminGrantBalance)
			}
		}
	}

	// Start server
	log.Printf("VESSEL Backend starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
