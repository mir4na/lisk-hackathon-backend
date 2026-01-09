package models

// CurrencyConversionRequest is the request for getting locked exchange rate
type CurrencyConversionRequest struct {
	OriginalCurrency string  `json:"original_currency" binding:"required"` // USD, EUR, etc.
	Amount           float64 `json:"amount" binding:"required,gt=0"`
}

// CurrencyConversionResponse contains the locked exchange rate with buffer
type CurrencyConversionResponse struct {
	OriginalCurrency string  `json:"original_currency"`
	OriginalAmount   float64 `json:"original_amount"`
	TargetCurrency   string  `json:"target_currency"` // IDR
	RealTimeRate     float64 `json:"realtime_rate"`   // Current rate before buffer
	BufferPercentage float64 `json:"buffer_percentage"`
	LockedRate       float64 `json:"locked_rate"`      // Rate after buffer: rate * (1 - buffer)
	ConvertedAmount  float64 `json:"converted_amount"` // Amount in IDR
	Microcopy        string  `json:"microcopy"`
}

// SupportedCurrency represents a supported currency
type SupportedCurrency struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	RateToIDR float64 `json:"rate_to_idr"`
	FlagEmoji string  `json:"flag_emoji"`
}
