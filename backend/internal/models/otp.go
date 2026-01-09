package models

import (
	"time"

	"github.com/google/uuid"
)

type OTPPurpose string

const (
	OTPPurposeRegistration  OTPPurpose = "registration"
	OTPPurposeLogin         OTPPurpose = "login"
	OTPPurposePasswordReset OTPPurpose = "password_reset"
)

type OTPCode struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Code      string     `json:"-"` // Never expose the code in JSON
	Purpose   OTPPurpose `json:"purpose"`
	ExpiresAt time.Time  `json:"expires_at"`
	Verified  bool       `json:"verified"`
	Attempts  int        `json:"attempts"`
	CreatedAt time.Time  `json:"created_at"`
}

// SendOTPRequest is the request to send an OTP
type SendOTPRequest struct {
	Email   string     `json:"email" binding:"required,email"`
	Purpose OTPPurpose `json:"purpose" binding:"required,oneof=registration login password_reset"`
}

// VerifyOTPRequest is the request to verify an OTP
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// SendOTPResponse is the response after sending an OTP
type SendOTPResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VerifyOTPResponse is the response after verifying an OTP
type VerifyOTPResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"` // Temporary token for registration flow
}

// IsExpired checks if the OTP has expired
func (o *OTPCode) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// CanRetry checks if the OTP can be retried (max 5 attempts)
func (o *OTPCode) CanRetry() bool {
	return o.Attempts < 5
}
