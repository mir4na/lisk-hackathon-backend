package utils

import (
	"errors"
	"fmt"
)

// Custom error types for better error handling

// AppError represents an application-level error
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeBlockchain     = "BLOCKCHAIN_ERROR"
	ErrCodeExternalAPI    = "EXTERNAL_API_ERROR"
)

// Common errors
var (
	ErrUserNotFound      = &AppError{Code: ErrCodeNotFound, Message: "User not found"}
	ErrInvoiceNotFound   = &AppError{Code: ErrCodeNotFound, Message: "Invoice not found"}
	ErrBuyerNotFound     = &AppError{Code: ErrCodeNotFound, Message: "Buyer not found"}
	ErrPoolNotFound      = &AppError{Code: ErrCodeNotFound, Message: "Funding pool not found"}
	ErrKYCNotFound       = &AppError{Code: ErrCodeNotFound, Message: "KYC verification not found"}

	ErrEmailExists       = &AppError{Code: ErrCodeConflict, Message: "Email already registered"}
	ErrWalletExists      = &AppError{Code: ErrCodeConflict, Message: "Wallet already linked to another account"}
	ErrInvoiceExists     = &AppError{Code: ErrCodeConflict, Message: "Invoice already exists"}
	ErrPoolExists        = &AppError{Code: ErrCodeConflict, Message: "Funding pool already exists"}
	ErrKYCPending        = &AppError{Code: ErrCodeConflict, Message: "KYC verification already pending"}

	ErrInvalidCredentials = &AppError{Code: ErrCodeUnauthorized, Message: "Invalid email or password"}
	ErrInvalidToken       = &AppError{Code: ErrCodeUnauthorized, Message: "Invalid or expired token"}
	ErrAccountDisabled    = &AppError{Code: ErrCodeUnauthorized, Message: "Account is deactivated"}

	ErrNotAuthorized     = &AppError{Code: ErrCodeForbidden, Message: "Not authorized to perform this action"}
	ErrInsufficientRole  = &AppError{Code: ErrCodeForbidden, Message: "Insufficient permissions"}

	ErrInvalidStatus     = &AppError{Code: ErrCodeBadRequest, Message: "Invalid status for this operation"}
	ErrInvalidAmount     = &AppError{Code: ErrCodeBadRequest, Message: "Invalid amount"}
	ErrInvalidDate       = &AppError{Code: ErrCodeBadRequest, Message: "Invalid date format"}
	ErrMissingDocument   = &AppError{Code: ErrCodeBadRequest, Message: "Required document missing"}
	ErrExceedsCapacity   = &AppError{Code: ErrCodeBadRequest, Message: "Amount exceeds pool capacity"}

	ErrBlockchainFailed  = &AppError{Code: ErrCodeBlockchain, Message: "Blockchain operation failed"}
	ErrIPFSUploadFailed  = &AppError{Code: ErrCodeExternalAPI, Message: "IPFS upload failed"}
)

// NewAppError creates a new AppError with wrapped error
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
	}
}

// IsAppError checks if error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// WrapError wraps an error with context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
