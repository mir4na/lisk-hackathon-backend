package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/utils"
)

var (
	ErrOTPRateLimit     = errors.New("too many OTP requests, please wait before trying again")
	ErrOTPNotFound      = errors.New("OTP not found or expired")
	ErrOTPInvalid       = errors.New("invalid OTP code")
	ErrOTPMaxAttempts   = errors.New("maximum verification attempts exceeded")
	ErrOTPExpired       = errors.New("OTP has expired")
	ErrEmailNotVerified = errors.New("email not verified")
)

type OTPService struct {
	otpRepo      *repository.OTPRepository
	emailService *EmailService
	config       *config.Config
	jwtManager   *utils.JWTManager
}

func NewOTPService(otpRepo *repository.OTPRepository, emailService *EmailService, cfg *config.Config, jwtManager *utils.JWTManager) *OTPService {
	return &OTPService{
		otpRepo:      otpRepo,
		emailService: emailService,
		config:       cfg,
		jwtManager:   jwtManager,
	}
}

// GenerateAndSendOTP generates a new OTP and sends it via email
func (s *OTPService) GenerateAndSendOTP(email string, purpose models.OTPPurpose) (*models.SendOTPResponse, error) {
	// Check rate limit (max 3 OTPs per 2 minutes)
	count, err := s.otpRepo.CountRecentOTPs(email, purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}
	if count >= 3 {
		return nil, ErrOTPRateLimit
	}

	// Invalidate any existing OTPs for this email
	if err := s.otpRepo.InvalidateAllForEmail(email, purpose); err != nil {
		return nil, fmt.Errorf("failed to invalidate existing OTPs: %w", err)
	}

	// Generate 6-digit OTP
	code, err := generateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	// [DEBUG] Log OTP for hackathon/testing
	fmt.Printf("[OTP] Generated Code for %s: %s\n", email, code)

	// Create OTP record
	expiresAt := time.Now().Add(time.Duration(s.config.OTPExpiryMinutes) * time.Minute)
	otp := &models.OTPCode{
		Email:     email,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		Verified:  false,
		Attempts:  0,
	}

	if err := s.otpRepo.Create(otp); err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	// Send OTP via email
	if err := s.emailService.SendOTPEmail(email, code, string(purpose)); err != nil {
		return nil, fmt.Errorf("failed to send OTP email: %w", err)
	}

	return &models.SendOTPResponse{
		Message:   "OTP code has been sent to your email",
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyOTP verifies the OTP code and returns a verification token
func (s *OTPService) VerifyOTP(email, code string) (*models.VerifyOTPResponse, error) {
	// Find the latest unverified OTP for this email (for registration purpose)
	otp, err := s.otpRepo.FindLatestByEmail(email, models.OTPPurposeRegistration)
	if err != nil {
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}
	if otp == nil {
		return nil, ErrOTPNotFound
	}

	// Check if expired
	if otp.IsExpired() {
		return nil, ErrOTPExpired
	}

	// Check max attempts
	if !otp.CanRetry() {
		return nil, ErrOTPMaxAttempts
	}

	// Increment attempt count
	if err := s.otpRepo.IncrementAttempts(otp.ID); err != nil {
		return nil, fmt.Errorf("failed to increment attempts: %w", err)
	}

	// Verify code
	if otp.Code != code {
		return &models.VerifyOTPResponse{
			Verified: false,
			Message:  fmt.Sprintf("Invalid OTP code. Remaining attempts: %d", s.config.OTPMaxAttempts-otp.Attempts-1),
		}, nil
	}

	// Mark as verified
	if err := s.otpRepo.MarkVerified(otp.ID); err != nil {
		return nil, fmt.Errorf("failed to mark OTP as verified: %w", err)
	}

	// Generate verification token (signed JWT)
	token, err := s.jwtManager.GenerateVerificationToken(email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	return &models.VerifyOTPResponse{
		Verified: true,
		Message:  "Email successfully verified",
		Token:    token,
	}, nil
}

// ValidateOTPToken validates the OTP verification token
// It checks if the token is valid and belongs to the specified email
func (s *OTPService) ValidateOTPToken(token string, email string) bool {
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return false
	}

	// Ensure the token was issued for this specific email
	if claims.Email != email {
		return false
	}

	// Ensure it is a verification token
	if claims.Issuer != "vessel-verify" {
		return false
	}

	return true
}

// generateOTPCode generates a 6-digit random OTP code
func generateOTPCode() (string, error) {
	const digits = "0123456789"
	result := make([]byte, 6)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[n.Int64()]
	}

	return string(result), nil
}

// generateVerificationToken generates a secure random token
func generateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
