package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Validation helper functions

// Email validation
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// Ethereum address validation (0x + 40 hex chars)
var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func IsValidEthAddress(address string) bool {
	return ethAddressRegex.MatchString(address)
}

// Transaction hash validation (0x + 64 hex chars)
var txHashRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{64}$`)

func IsValidTxHash(hash string) bool {
	return txHashRegex.MatchString(hash)
}

// Password validation
func IsValidPassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters"
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasDigit {
		return false, "Password must contain at least one digit"
	}

	return true, ""
}

// Date validation
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func IsValidDateFormat(dateStr string) bool {
	_, err := ParseDate(dateStr)
	return err == nil
}

func IsFutureDate(dateStr string) bool {
	date, err := ParseDate(dateStr)
	if err != nil {
		return false
	}
	return date.After(time.Now())
}

func IsDateAfter(date1, date2 string) bool {
	d1, err1 := ParseDate(date1)
	d2, err2 := ParseDate(date2)
	if err1 != nil || err2 != nil {
		return false
	}
	return d1.After(d2)
}

// String validation
func IsNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

func IsWithinLength(s string, min, max int) bool {
	length := len(strings.TrimSpace(s))
	return length >= min && length <= max
}

// Numeric validation
func IsPositiveFloat(f float64) bool {
	return f > 0
}

func IsWithinRange(f, min, max float64) bool {
	return f >= min && f <= max
}

func IsPositiveInt(i int) bool {
	return i > 0
}

// Invoice number validation (alphanumeric with dashes)
var invoiceNumberRegex = regexp.MustCompile(`^[A-Za-z0-9\-]+$`)

func IsValidInvoiceNumber(number string) bool {
	return invoiceNumberRegex.MatchString(number) && len(number) >= 3 && len(number) <= 50
}

// Phone number validation (basic)
var phoneRegex = regexp.MustCompile(`^\+?[0-9]{8,15}$`)

func IsValidPhone(phone string) bool {
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return phoneRegex.MatchString(cleaned)
}

// Country code validation (ISO 3166-1 alpha-2 style)
func IsValidCountryCode(code string) bool {
	return len(code) >= 2 && len(code) <= 100
}

// Currency validation
var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"IDR": true,
	"SGD": true,
	"JPY": true,
	"CNY": true,
	"AUD": true,
}

func IsValidCurrency(currency string) bool {
	return validCurrencies[strings.ToUpper(currency)]
}

// File validation
var allowedMimeTypes = map[string]bool{
	"application/pdf":  true,
	"image/jpeg":       true,
	"image/jpg":        true,
	"image/png":        true,
}

var allowedExtensions = map[string]bool{
	".pdf":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func IsAllowedFileType(mimeType string) bool {
	return allowedMimeTypes[strings.ToLower(mimeType)]
}

func IsAllowedExtension(filename string) bool {
	lower := strings.ToLower(filename)
	for ext := range allowedExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func IsValidFileSize(size int64, maxMB int) bool {
	maxBytes := int64(maxMB) * 1024 * 1024
	return size > 0 && size <= maxBytes
}

// ValidationResult holds multiple validation errors
type ValidationResult struct {
	Valid  bool
	Errors map[string]string
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: make(map[string]string),
	}
}

func (v *ValidationResult) AddError(field, message string) {
	v.Valid = false
	v.Errors[field] = message
}

func (v *ValidationResult) HasErrors() bool {
	return !v.Valid
}

func (v *ValidationResult) FirstError() string {
	for _, msg := range v.Errors {
		return msg
	}
	return ""
}
