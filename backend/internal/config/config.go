package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port    string
	GinMode string

	// Database
	DatabaseURL  string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret             string
	JWTExpiryHours        int
	JWTRefreshExpiryHours int

	// Blockchain
	PrivateKey              string
	BlockchainRPCURL        string
	ChainID                 int64
	InvoiceNFTContractAddr  string
	InvoicePoolContractAddr string

	// Pinata (IPFS)
	PinataAPIKey     string
	PinataSecretKey  string
	PinataJWT        string
	PinataGatewayURL string

	// File Upload
	MaxFileSizeMB    int
	AllowedFileTypes string

	// Logistics API

	// Platform Settings
	PlatformFeePercentage    float64
	DefaultAdvancePercentage float64
	MinInvoiceAmount         float64
	MaxInvoiceAmount         float64

	// CORS
	CORSAllowedOrigins string

	// Frontend URL for payment links
	FrontendURL string

	// SMTP Settings for OTP (VESSEL)
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// OTP Settings
	OTPExpiryMinutes int
	OTPMaxAttempts   int

	// Currency Conversion Settings
	DefaultBufferRate float64 // Default 1.5% buffer for currency conversion
}

func Load() (*Config, error) {
	wd, _ := os.Getwd()
	fmt.Printf("[CONFIG] Starting up... Current Working Directory: %s\n", wd)

	// Try loading .env from current directory first
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("[CONFIG] Note: .env not found in %s (Error: %v), checking other locations...\n", wd, err)

		// Try backend/.env (if running from root)
		if err := godotenv.Load("backend/.env"); err != nil {
			// Try parent directory
			if err := godotenv.Load("../.env"); err != nil {
				fmt.Printf("[CONFIG] Warning: Could not load .env from current, backend/, or parent directory\n")
			} else {
				fmt.Println("[CONFIG] Loaded .env from parent directory (../.env)")
			}
		} else {
			fmt.Println("[CONFIG] Loaded .env from backend/.env")
		}
	} else {
		fmt.Println("[CONFIG] Loaded .env from current directory")
	}

	// Verify critical config
	smtpUser := getEnv("SMTP_USERNAME", getEnv("SMTP_USER", ""))
	fmt.Printf("[CONFIG] Debug: Found SMTP user='%s' (Len: %d)\n", smtpUser, len(smtpUser))
	if smtpUser == "" {
		fmt.Println("[CONFIG] CRITICAL WARNING: SMTP username is empty! Email sending will fail.")
	}

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))
	chainID, _ := strconv.ParseInt(getEnv("CHAIN_ID", "4202"), 10, 64)
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	maxFileSize, _ := strconv.Atoi(getEnv("MAX_FILE_SIZE_MB", "10"))
	platformFee, _ := strconv.ParseFloat(getEnv("PLATFORM_FEE_PERCENTAGE", "2.0"), 64)
	defaultAdvance, _ := strconv.ParseFloat(getEnv("DEFAULT_ADVANCE_PERCENTAGE", "80.0"), 64)
	minInvoice, _ := strconv.ParseFloat(getEnv("MIN_INVOICE_AMOUNT", "1000"), 64)
	maxInvoice, _ := strconv.ParseFloat(getEnv("MAX_INVOICE_AMOUNT", "1000000"), 64)
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	otpExpiry, _ := strconv.Atoi(getEnv("OTP_EXPIRY_MINUTES", "5"))
	otpMaxAttempts, _ := strconv.Atoi(getEnv("OTP_MAX_ATTEMPTS", "5"))
	bufferRate, _ := strconv.ParseFloat(getEnv("DEFAULT_BUFFER_RATE", "0.015"), 64)

	return &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DatabaseURL:  getEnv("DATABASE_URL", ""),
		PostgresHost: getEnv("POSTGRES_HOST", ""),
		PostgresPort: getEnv("POSTGRES_PORT", ""),
		PostgresUser: getEnv("POSTGRES_USER", ""),
		PostgresPass: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:   getEnv("POSTGRES_DB", ""),

		RedisHost:     getEnv("REDIS_HOST", ""),
		RedisPort:     getEnv("REDIS_PORT", ""),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTSecret:             getEnv("JWT_SECRET", ""),
		JWTExpiryHours:        jwtExpiry,
		JWTRefreshExpiryHours: jwtRefreshExpiry,

		PrivateKey:              getEnv("PRIVATE_KEY", ""),
		BlockchainRPCURL:        getEnv("BLOCKCHAIN_RPC_URL", ""),
		ChainID:                 chainID,
		InvoiceNFTContractAddr:  getEnv("INVOICE_NFT_CONTRACT_ADDRESS", ""),
		InvoicePoolContractAddr: getEnv("INVOICE_POOL_CONTRACT_ADDRESS", ""),

		PinataAPIKey:     getEnv("PINATA_API_KEY", ""),
		PinataSecretKey:  getEnv("PINATA_SECRET_KEY", ""),
		PinataJWT:        getEnv("PINATA_JWT", ""),
		PinataGatewayURL: getEnv("PINATA_GATEWAY_URL", ""),

		MaxFileSizeMB:    maxFileSize,
		AllowedFileTypes: getEnv("ALLOWED_FILE_TYPES", ""),

		PlatformFeePercentage:    platformFee,
		DefaultAdvancePercentage: defaultAdvance,
		MinInvoiceAmount:         minInvoice,
		MaxInvoiceAmount:         maxInvoice,

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", ""),

		FrontendURL: getEnv("FRONTEND_URL", ""),

		// SMTP Settings for OTP
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMTP_USERNAME", getEnv("SMTP_USER", "")),
		SMTPPassword: strings.ReplaceAll(getEnv("SMTP_PASSWORD", getEnv("SMTP_PASS", "")), " ", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		// OTP Settings
		OTPExpiryMinutes: otpExpiry,
		OTPMaxAttempts:   otpMaxAttempts,

		// Currency Settings
		DefaultBufferRate: bufferRate,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
