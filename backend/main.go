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
	invoiceService.SetUserRepo(userRepo) // Set user repo for grade suggestion
	fundingService := services.NewFundingService(fundingRepo, invoiceRepo, txRepo, userRepo, buyerRepo, rqRepo, emailService, escrowService, cfg)
	paymentService := services.NewPaymentService(userRepo, txRepo, fundingRepo, invoiceRepo) // Updated with fundingRepo and invoiceRepo for Flow 3
	rqService := services.NewRiskQuestionnaireService(rqRepo)
	currencyService := services.NewCurrencyService(cfg)
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
	currencyHandler := handlers.NewCurrencyHandler(currencyService)

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
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/profile", userHandler.UpdateProfile)
				user.POST("/kyc", userHandler.SubmitKYC)
				user.GET("/kyc", userHandler.GetKYCStatus)
				user.GET("/balance", paymentHandler.GetBalance)

				// Profile Management (Flow: MANAGEMENT PROFIL USER)
				user.GET("/profile/data", userHandler.GetPersonalData)           // Data Diri (Read-only)
				user.GET("/profile/bank-account", userHandler.GetBankAccount)    // Rekening Bank
				user.PUT("/profile/bank-account", userHandler.ChangeBankAccount) // Ubah Rekening (OTP required)
				user.PUT("/profile/password", userHandler.ChangePassword)        // Keamanan - Ubah Password
				user.GET("/profile/banks", userHandler.GetSupportedBanks)        // List supported banks

				// MITRA application routes (Flow 2)
				mitra := user.Group("/mitra")
				{
					mitra.POST("/apply", mitraHandler.Apply)
					mitra.GET("/status", mitraHandler.GetStatus)
					mitra.POST("/documents", mitraHandler.UploadDocument)
				}
			}

			// Currency conversion routes (Flow 4 - BE-4)
			currency := protected.Group("/currency")
			{
				currency.POST("/convert", currencyHandler.GetLockedExchangeRate)
				currency.GET("/supported", currencyHandler.GetSupportedCurrencies)
				currency.GET("/disbursement-estimate", currencyHandler.CalculateEstimatedDisbursement)
			}

			// Payment routes (PROTOTYPE)
			payments := protected.Group("/payments")
			{
				payments.POST("/deposit", paymentHandler.Deposit)
				payments.POST("/withdraw", paymentHandler.Withdraw)
				payments.GET("/balance", paymentHandler.GetBalance)
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

			// Invoice routes (exporter/mitra for CRUD)
			invoices := protected.Group("/invoices")
			{
				// Mitra/Exporter routes
				invoices.POST("", middleware.ExporterOnly(), invoiceHandler.Create)
				invoices.POST("/funding-request", middleware.ExporterOnly(), invoiceHandler.CreateFundingRequest) // Flow 4
				invoices.POST("/check-repeat-buyer", middleware.ExporterOnly(), invoiceHandler.CheckRepeatBuyer)  // Flow 4 Pre-condition
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

			// Funding/Investment routes (Marketplace)
			pools := protected.Group("/pools")
			{
				pools.GET("", fundingHandler.ListPools)
				pools.GET("/:id", fundingHandler.GetPool)
			}

			// Marketplace routes (with filters) - Flow 6
			marketplace := protected.Group("/marketplace")
			{
				marketplace.GET("", fundingHandler.GetMarketplace)
				marketplace.GET("/:id/detail", fundingHandler.GetPoolDetail)       // Pool detail for investor
				marketplace.POST("/calculate", fundingHandler.CalculateInvestment) // Investment calculator
			}

			// Risk Questionnaire routes (for investors)
			riskQuestionnaire := protected.Group("/risk-questionnaire")
			riskQuestionnaire.Use(middleware.InvestorOnly())
			{
				riskQuestionnaire.GET("/questions", rqHandler.GetQuestions)
				riskQuestionnaire.POST("", rqHandler.Submit)
				riskQuestionnaire.GET("/status", rqHandler.GetStatus)
			}

			// Investment routes (Flow 6, 9, 10)
			investments := protected.Group("/investments")
			investments.Use(middleware.InvestorOnly())
			{
				investments.POST("", fundingHandler.Invest)
				investments.POST("/confirm", fundingHandler.ConfirmInvestment) // Flow 6 confirmation
				investments.GET("", fundingHandler.GetMyInvestments)
				investments.GET("/portfolio", fundingHandler.GetPortfolio)      // Flow 9
				investments.GET("/active", fundingHandler.GetActiveInvestments) // Flow 10
			}

			// Exporter/Mitra routes (Flow 8, 11)
			exporter := protected.Group("/exporter")
			exporter.Use(middleware.ExporterOnly())
			{
				exporter.POST("/disbursement", fundingHandler.ExporterDisbursement) // Flow 11
			}

			// Mitra Dashboard (Flow 8)
			mitraDashboard := protected.Group("/mitra")
			mitraDashboard.Use(middleware.ExporterOnly())
			{
				mitraDashboard.GET("/dashboard", fundingHandler.GetMitraDashboard)
				mitraDashboard.GET("/invoices", fundingHandler.GetMitraActiveInvoices)

				// Mitra Repayment (Flow: MITRA MEMBAYAR HUTANG)
				mitraDashboard.GET("/invoices/active", mitraHandler.GetActiveInvoices)                // Active invoices needing repayment
				mitraDashboard.GET("/pools/:id/breakdown", mitraHandler.GetRepaymentBreakdown)        // Repayment breakdown by tranche
				mitraDashboard.GET("/payment-methods", mitraHandler.GetVAPaymentMethods)              // Available VA banks
				mitraDashboard.POST("/repayment/va", mitraHandler.CreateVAPayment)                    // Create VA for payment
				mitraDashboard.GET("/repayment/va/:id", mitraHandler.GetVAPaymentStatus)              // VA payment page details
				mitraDashboard.POST("/repayment/va/:id/simulate-pay", mitraHandler.SimulateVAPayment) // MVP: Simulate payment
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/kyc/pending", userHandler.GetPendingKYC)
				admin.POST("/kyc/:id/approve", userHandler.ApproveKYC)
				admin.POST("/kyc/:id/reject", userHandler.RejectKYC)

				// Invoice approval routes (Flow 5)
				admin.GET("/invoices/pending", invoiceHandler.GetPendingInvoices)              // List pending for review
				admin.GET("/invoices/:id/grade-suggestion", invoiceHandler.GetGradeSuggestion) // BE-ADM-1 logic
				admin.GET("/invoices/:id/review", invoiceHandler.GetInvoiceReviewData)         // Split-screen data
				admin.POST("/invoices/:id/approve", invoiceHandler.Approve)                    // Approve with grade
				admin.POST("/invoices/:id/reject", invoiceHandler.Reject)

				// Pool management
				admin.POST("/pools/:id/disburse", fundingHandler.Disburse)
				admin.POST("/pools/:id/close", fundingHandler.ClosePoolAndNotify)
				admin.POST("/invoices/:id/repay", fundingHandler.ProcessRepayment)

				// Admin Mitra Application routes (Flow 2)
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
