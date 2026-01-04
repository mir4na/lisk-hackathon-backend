package config

import (
	"os"
	"strconv"

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
	JWTSecret            string
	JWTExpiryHours       int
	JWTRefreshExpiryHours int

	// Blockchain
	PrivateKey               string
	BlockchainRPCURL         string
	ChainID                  int64
	InvoiceNFTContractAddr   string
	InvoicePoolContractAddr  string
	USDCContractAddr         string

	// Pinata (IPFS)
	PinataAPIKey     string
	PinataSecretKey  string
	PinataJWT        string
	PinataGatewayURL string

	// File Upload
	MaxFileSizeMB    int
	AllowedFileTypes string

	// Logistics API
	JSONCargoAPIKey string
	JSONCargoAPIURL string

	// Platform Settings
	PlatformFeePercentage     float64
	DefaultAdvancePercentage  float64
	MinInvoiceAmount          float64
	MaxInvoiceAmount          float64

	// CORS
	CORSAllowedOrigins string
}

func Load() (*Config, error) {
	godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))
	chainID, _ := strconv.ParseInt(getEnv("CHAIN_ID", "4202"), 10, 64)
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	maxFileSize, _ := strconv.Atoi(getEnv("MAX_FILE_SIZE_MB", "10"))
	platformFee, _ := strconv.ParseFloat(getEnv("PLATFORM_FEE_PERCENTAGE", "2.0"), 64)
	defaultAdvance, _ := strconv.ParseFloat(getEnv("DEFAULT_ADVANCE_PERCENTAGE", "80.0"), 64)
	minInvoice, _ := strconv.ParseFloat(getEnv("MIN_INVOICE_AMOUNT", "1000"), 64)
	maxInvoice, _ := strconv.ParseFloat(getEnv("MAX_INVOICE_AMOUNT", "1000000"), 64)

	return &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		DatabaseURL:  getEnv("DATABASE_URL", ""),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "receiv3"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "receiv3"),
		PostgresDB:   getEnv("POSTGRES_DB", "receiv3"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTSecret:             getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiryHours:        jwtExpiry,
		JWTRefreshExpiryHours: jwtRefreshExpiry,

		PrivateKey:               getEnv("PRIVATE_KEY", ""),
		BlockchainRPCURL:         getEnv("BLOCKCHAIN_RPC_URL", "https://rpc.sepolia-api.lisk.com"),
		ChainID:                  chainID,
		InvoiceNFTContractAddr:   getEnv("INVOICE_NFT_CONTRACT_ADDRESS", ""),
		InvoicePoolContractAddr:  getEnv("INVOICE_POOL_CONTRACT_ADDRESS", ""),
		USDCContractAddr:         getEnv("USDC_CONTRACT_ADDRESS", ""),

		PinataAPIKey:     getEnv("PINATA_API_KEY", ""),
		PinataSecretKey:  getEnv("PINATA_SECRET_KEY", ""),
		PinataJWT:        getEnv("PINATA_JWT", ""),
		PinataGatewayURL: getEnv("PINATA_GATEWAY_URL", "https://gateway.pinata.cloud/ipfs/"),

		MaxFileSizeMB:    maxFileSize,
		AllowedFileTypes: getEnv("ALLOWED_FILE_TYPES", "pdf,png,jpg,jpeg"),

		JSONCargoAPIKey: getEnv("JSONCARGO_API_KEY", ""),
		JSONCargoAPIURL: getEnv("JSONCARGO_API_URL", "https://api.jsoncargo.com/v1"),

		PlatformFeePercentage:    platformFee,
		DefaultAdvancePercentage: defaultAdvance,
		MinInvoiceAmount:         minInvoice,
		MaxInvoiceAmount:         maxInvoice,

		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
