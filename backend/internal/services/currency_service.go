package services

import (
	"errors"
	"fmt"

	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
)

// CurrencyService handles currency conversion with buffer rate
type CurrencyService struct {
	cfg *config.Config
}

func NewCurrencyService(cfg *config.Config) *CurrencyService {
	return &CurrencyService{cfg: cfg}
}

// ExchangeRates - For MVP, we use static rates. In production, fetch from API
var ExchangeRates = map[string]float64{
	"USD": 15500.0, // 1 USD = Rp 15,500
	"EUR": 16800.0, // 1 EUR = Rp 16,800
	"GBP": 19500.0, // 1 GBP = Rp 19,500
	"SGD": 11500.0, // 1 SGD = Rp 11,500
	"JPY": 105.0,   // 1 JPY = Rp 105
	"CNY": 2150.0,  // 1 CNY = Rp 2,150
	"AUD": 10200.0, // 1 AUD = Rp 10,200
	"MYR": 3450.0,  // 1 MYR = Rp 3,450
}

// SupportedCurrencies list
var SupportedCurrencies = []models.SupportedCurrency{
	{Code: "USD", Name: "US Dollar", FlagEmoji: "🇺🇸"},
	{Code: "EUR", Name: "Euro", FlagEmoji: "🇪🇺"},
	{Code: "GBP", Name: "British Pound", FlagEmoji: "🇬🇧"},
	{Code: "SGD", Name: "Singapore Dollar", FlagEmoji: "🇸🇬"},
	{Code: "JPY", Name: "Japanese Yen", FlagEmoji: "🇯🇵"},
	{Code: "CNY", Name: "Chinese Yuan", FlagEmoji: "🇨🇳"},
	{Code: "AUD", Name: "Australian Dollar", FlagEmoji: "🇦🇺"},
	{Code: "MYR", Name: "Malaysian Ringgit", FlagEmoji: "🇲🇾"},
}

// GetLockedExchangeRate implements BE-4 logic:
// 1. Fetch RealTimeRate
// 2. Add BufferConfig (e.g., 1.5%)
// 3. LockedRate = rate * (1 - buffer)
// This protects the disbursement value
func (s *CurrencyService) GetLockedExchangeRate(req *models.CurrencyConversionRequest) (*models.CurrencyConversionResponse, error) {
	if req.OriginalCurrency == "" {
		return nil, errors.New("original currency is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Get real-time rate (static for MVP)
	realTimeRate, exists := ExchangeRates[req.OriginalCurrency]
	if !exists {
		return nil, fmt.Errorf("unsupported currency: %s", req.OriginalCurrency)
	}

	// Apply buffer rate (default 1.5%)
	bufferPercentage := s.cfg.DefaultBufferRate
	if bufferPercentage == 0 {
		bufferPercentage = 0.015 // 1.5%
	}

	// LockedRate = realTimeRate * (1 - buffer)
	// This gives exporter slightly less to protect against rate fluctuation
	lockedRate := realTimeRate * (1 - bufferPercentage)

	// Convert amount
	convertedAmount := req.Amount * lockedRate

	return &models.CurrencyConversionResponse{
		OriginalCurrency: req.OriginalCurrency,
		OriginalAmount:   req.Amount,
		TargetCurrency:   "IDR",
		RealTimeRate:     realTimeRate,
		BufferPercentage: bufferPercentage * 100, // Display as percentage
		LockedRate:       lockedRate,
		ConvertedAmount:  convertedAmount,
		Microcopy:        "Kurs dikunci final untuk melindungi nilai pencairan.",
	}, nil
}

// GetSupportedCurrencies returns list of supported currencies with current rates
func (s *CurrencyService) GetSupportedCurrencies() []models.SupportedCurrency {
	currencies := make([]models.SupportedCurrency, len(SupportedCurrencies))
	for i, c := range SupportedCurrencies {
		currencies[i] = c
		if rate, exists := ExchangeRates[c.Code]; exists {
			currencies[i].RateToIDR = rate
		}
	}
	return currencies
}

// CalculateEstimatedDisbursement calculates net disbursement after platform fee
func (s *CurrencyService) CalculateEstimatedDisbursement(idrAmount float64) *models.EstimatedDisbursement {
	platformFeePercentage := s.cfg.PlatformFeePercentage
	if platformFeePercentage == 0 {
		platformFeePercentage = 2.0 // Default 2%
	}

	platformFee := idrAmount * (platformFeePercentage / 100)
	netDisbursement := idrAmount - platformFee

	return &models.EstimatedDisbursement{
		GrossAmount:     idrAmount,
		PlatformFee:     platformFee,
		NetDisbursement: netDisbursement,
		Currency:        "IDR",
	}
}
