package models

import (
	"time"

	"github.com/google/uuid"
)

type MitraApplicationStatus string

const (
	MitraStatusPending  MitraApplicationStatus = "pending"
	MitraStatusApproved MitraApplicationStatus = "approved"
	MitraStatusRejected MitraApplicationStatus = "rejected"
)

type CompanyType string

const (
	CompanyTypePT CompanyType = "PT"
	CompanyTypeCV CompanyType = "CV"
	CompanyTypeUD CompanyType = "UD"
)

type AnnualRevenue string

const (
	RevenueLessThan1M   AnnualRevenue = "<1M"
	Revenue1Mto5M       AnnualRevenue = "1M-5M"
	Revenue5Mto25M      AnnualRevenue = "5M-25M"
	Revenue25Mto100M    AnnualRevenue = "25M-100M"
	RevenueMoreThan100M AnnualRevenue = ">100M"
)

type MitraApplication struct {
	ID               uuid.UUID              `json:"id"`
	UserID           uuid.UUID              `json:"user_id"`
	CompanyName      string                 `json:"company_name"`
	CompanyType      CompanyType            `json:"company_type"`
	NPWP             string                 `json:"npwp"`
	AnnualRevenue    AnnualRevenue          `json:"annual_revenue"`
	NIBDocumentURL   *string                `json:"nib_document_url,omitempty"`
	AktaPendirianURL *string                `json:"akta_pendirian_url,omitempty"`
	KTPDirekturURL   *string                `json:"ktp_direktur_url,omitempty"`
	Status           MitraApplicationStatus `json:"status"`
	RejectionReason  *string                `json:"rejection_reason,omitempty"`
	ReviewedBy       *uuid.UUID             `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time             `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`

	// Relations
	User     *User `json:"user,omitempty"`
	Reviewer *User `json:"reviewer,omitempty"`
}

// SubmitMitraApplicationRequest is the request to apply for MITRA status
type SubmitMitraApplicationRequest struct {
	CompanyName   string        `json:"company_name" binding:"required"`
	CompanyType   CompanyType   `json:"company_type" binding:"required,oneof=PT CV UD"`
	NPWP          string        `json:"npwp" binding:"required,min=15,max=16"`
	AnnualRevenue AnnualRevenue `json:"annual_revenue" binding:"required,oneof=<1M 1M-5M 5M-25M 25M-100M >100M"`
}

// UploadMitraDocumentRequest is the request for uploading MITRA documents
type UploadMitraDocumentRequest struct {
	DocumentType string `form:"document_type" binding:"required,oneof=nib akta_pendirian ktp_direktur"`
}

// ApproveMitraRequest is the request to approve MITRA application
type ApproveMitraRequest struct {
	// Admin can optionally add notes
	Notes string `json:"notes,omitempty"`
}

// RejectMitraRequest is the request to reject MITRA application
type RejectMitraRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// MitraApplicationResponse is the response for MITRA application status
type MitraApplicationResponse struct {
	Application     *MitraApplication `json:"application"`
	DocumentsStatus struct {
		NIB           bool `json:"nib"`
		AktaPendirian bool `json:"akta_pendirian"`
		KTPDirektur   bool `json:"ktp_direktur"`
	} `json:"documents_status"`
	IsComplete bool `json:"is_complete"`
}

// ValidateNPWP validates the NPWP format (15 or 16 digits)
func ValidateNPWP(npwp string) bool {
	if len(npwp) != 15 && len(npwp) != 16 {
		return false
	}
	for _, c := range npwp {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
