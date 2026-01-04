package models

import (
	"time"

	"github.com/google/uuid"
)

type KYCStatus string
type KYCType string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusApproved KYCStatus = "approved"
	KYCStatusRejected KYCStatus = "rejected"

	KYCTypeKYC KYCType = "kyc"
	KYCTypeKYB KYCType = "kyb"
)

type KYCVerification struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	VerificationType KYCType    `json:"verification_type"`
	Status           KYCStatus  `json:"status"`
	IDType           *string    `json:"id_type,omitempty"`
	IDNumber         *string    `json:"id_number,omitempty"`
	IDDocumentURL    *string    `json:"id_document_url,omitempty"`
	SelfieURL        *string    `json:"selfie_url,omitempty"`
	RejectionReason  *string    `json:"rejection_reason,omitempty"`
	VerifiedBy       *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type SubmitKYCRequest struct {
	VerificationType KYCType `json:"verification_type" binding:"required,oneof=kyc kyb"`
	IDType           string  `json:"id_type" binding:"required"`
	IDNumber         string  `json:"id_number" binding:"required"`
}

type ApproveKYCRequest struct {
	KYCID uuid.UUID `json:"kyc_id" binding:"required"`
}

type RejectKYCRequest struct {
	KYCID  uuid.UUID `json:"kyc_id" binding:"required"`
	Reason string    `json:"reason" binding:"required"`
}
