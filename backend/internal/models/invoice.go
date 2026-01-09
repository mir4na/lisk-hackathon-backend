package models

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	StatusDraft         InvoiceStatus = "draft"
	StatusPendingReview InvoiceStatus = "pending_review"
	StatusApproved      InvoiceStatus = "approved"
	StatusRejected      InvoiceStatus = "rejected"
	StatusTokenized     InvoiceStatus = "tokenized"
	StatusFunding       InvoiceStatus = "funding"
	StatusFunded        InvoiceStatus = "funded"
	StatusMatured       InvoiceStatus = "matured"
	StatusRepaid        InvoiceStatus = "repaid"
	StatusDefaulted     InvoiceStatus = "defaulted"
)

type Invoice struct {
	ID         uuid.UUID `json:"id"`
	ExporterID uuid.UUID `json:"exporter_id"`
	// Buyer Details (Flattened)
	BuyerName    string `json:"buyer_name"`
	BuyerCountry string `json:"buyer_country"`

	InvoiceNumber     string        `json:"invoice_number"`
	Currency          string        `json:"currency"`
	Amount            float64       `json:"amount"`
	IssueDate         time.Time     `json:"issue_date"`
	DueDate           time.Time     `json:"due_date"`
	Description       *string       `json:"description,omitempty"`
	Status            InvoiceStatus `json:"status"`
	InterestRate      *float64      `json:"interest_rate,omitempty"`
	AdvancePercentage float64       `json:"advance_percentage"`
	AdvanceAmount     *float64      `json:"advance_amount,omitempty"`
	DocumentHash      *string       `json:"document_hash,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`

	// VESSEL Grading Fields
	Grade                  *string `json:"grade,omitempty"`              // A, B, or C
	GradeScore             *int    `json:"grade_score,omitempty"`        // Numeric score 0-100
	IsRepeatBuyer          bool    `json:"is_repeat_buyer"`              // True if buyer has previous transactions
	FundingLimitPercentage float64 `json:"funding_limit_percentage"`     // Max funding % (60% for new, 100% for repeat)
	IsInsured              bool    `json:"is_insured"`                   // True if invoice has insurance
	DocumentCompleteScore  int     `json:"document_complete_score"`      // 0-100 based on document completeness
	BuyerCountryRisk       *string `json:"buyer_country_risk,omitempty"` // low, medium, high

	// VESSEL Tranche/Interest Fields
	PriorityRatio        float64  `json:"priority_ratio"`                   // Default 80%
	CatalystRatio        float64  `json:"catalyst_ratio"`                   // Default 20%
	PriorityInterestRate *float64 `json:"priority_interest_rate,omitempty"` // e.g., 10%
	CatalystInterestRate *float64 `json:"catalyst_interest_rate,omitempty"` // e.g., 15%

	// VESSEL Currency Conversion Fields
	OriginalCurrency *string  `json:"original_currency,omitempty"`
	OriginalAmount   *float64 `json:"original_amount,omitempty"`
	IDRAmount        *float64 `json:"idr_amount,omitempty"`
	ExchangeRate     *float64 `json:"exchange_rate,omitempty"`
	BufferRate       float64  `json:"buffer_rate"` // Default 1.5%

	// VESSEL Additional Fields
	FundingDurationDays int     `json:"funding_duration_days"` // Default 14 days
	PaymentLink         *string `json:"payment_link,omitempty"`

	// Relations

	Exporter  *User             `json:"exporter,omitempty"`
	Documents []InvoiceDocument `json:"documents,omitempty"`
	NFT       *InvoiceNFT       `json:"nft,omitempty"`
}

type DocumentType string

const (
	DocTypeInvoicePDF          DocumentType = "invoice_pdf"
	DocTypeBillOfLading        DocumentType = "bill_of_lading"
	DocTypePackingList         DocumentType = "packing_list"
	DocTypeCertificateOfOrigin DocumentType = "certificate_of_origin"
	DocTypeInsurance           DocumentType = "insurance"
	DocTypeCustoms             DocumentType = "customs"
	DocTypeOther               DocumentType = "other"
	DocTypePurchaseOrder       DocumentType = "purchase_order"
	DocTypeCommercialInvoice   DocumentType = "commercial_invoice"
)

type InvoiceDocument struct {
	ID           uuid.UUID    `json:"id"`
	InvoiceID    uuid.UUID    `json:"invoice_id"`
	DocumentType DocumentType `json:"document_type"`
	FileName     string       `json:"file_name"`
	FileURL      string       `json:"file_url"`
	FileHash     string       `json:"file_hash"`
	FileSize     int          `json:"file_size"`
	UploadedAt   time.Time    `json:"uploaded_at"`
}

type InvoiceNFT struct {
	ID              uuid.UUID  `json:"id"`
	InvoiceID       uuid.UUID  `json:"invoice_id"`
	TokenID         *int64     `json:"token_id,omitempty"`
	ContractAddress *string    `json:"contract_address,omitempty"`
	ChainID         int        `json:"chain_id"`
	OwnerAddress    *string    `json:"owner_address,omitempty"`
	MintTxHash      *string    `json:"mint_tx_hash,omitempty"`
	MetadataURI     *string    `json:"metadata_uri,omitempty"`
	MintedAt        *time.Time `json:"minted_at,omitempty"`
	BurnedAt        *time.Time `json:"burned_at,omitempty"`
	BurnTxHash      *string    `json:"burn_tx_hash,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateInvoiceFundingRequest is the request for mitra to create invoice funding request (Flow 4)
// UI Group 1: Data Pembayar (Buyer)
// UI Group 2: Nilai Pengajuan
type CreateInvoiceFundingRequest struct {
	// UI Group 1: Data Pembayar (Buyer)
	BuyerCompanyName string `json:"buyer_company_name" binding:"required"` // Nama Perusahaan Importir
	BuyerCountry     string `json:"buyer_country" binding:"required"`      // Negara Tujuan
	BuyerEmail       string `json:"buyer_email" binding:"required,email"`  // Email utama importir

	// UI Group 2: Nilai Pengajuan
	InvoiceNumber       string  `json:"invoice_number" binding:"required"`
	OriginalCurrency    string  `json:"original_currency" binding:"required"`    // Mata Uang Invoice (USD, EUR, etc.)
	OriginalAmount      float64 `json:"original_amount" binding:"required,gt=0"` // Nominal Invoice
	LockedExchangeRate  float64 `json:"locked_exchange_rate" binding:"required"` // Kurs konversi yang dikunci
	IDRAmount           float64 `json:"idr_amount" binding:"required,gt=0"`      // Nominal dalam IDR
	DueDate             string  `json:"due_date" binding:"required"`             // Tanggal Jatuh Tempo (Tenor)
	FundingDurationDays int     `json:"funding_duration_days"`                   // Funding Duration (Default 14 Hari)

	// Tranche Configuration
	PriorityRatio        float64 `json:"priority_ratio"`                                         // Default 80%
	CatalystRatio        float64 `json:"catalyst_ratio"`                                         // Default 20%
	PriorityInterestRate float64 `json:"priority_interest_rate" binding:"required,gt=0,lte=100"` // e.g., 10%
	CatalystInterestRate float64 `json:"catalyst_interest_rate" binding:"required,gt=0,lte=100"` // e.g., 15%

	// Repeat buyer info
	IsRepeatBuyer    bool   `json:"is_repeat_buyer"`              // Auto-checked by system or manual
	RepeatBuyerProof string `json:"repeat_buyer_proof,omitempty"` // URL jika manual upload bukti

	// Confirmation
	DataConfirmation bool `json:"data_confirmation" binding:"required"` // Checkbox: "Data yang saya berikan adalah benar dan asli."

	Description *string `json:"description,omitempty"`
}

// CreateInvoiceRequest - legacy request, kept for backward compatibility
type CreateInvoiceRequest struct {
	BuyerID       uuid.UUID `json:"buyer_id" binding:"required"`
	InvoiceNumber string    `json:"invoice_number" binding:"required"`
	Currency      string    `json:"currency"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	IssueDate     string    `json:"issue_date" binding:"required"`
	DueDate       string    `json:"due_date" binding:"required"`
	Description   *string   `json:"description,omitempty"`
}

type UpdateInvoiceRequest struct {
	InvoiceNumber string  `json:"invoice_number"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	IssueDate     string  `json:"issue_date"`
	DueDate       string  `json:"due_date"`
	Description   *string `json:"description,omitempty"`
}

type SubmitInvoiceRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" binding:"required"`
}

// AdminApproveInvoiceRequest is the request for admin to approve invoice with grade (Flow 5)
type AdminApproveInvoiceRequest struct {
	Grade                string  `json:"grade" binding:"required,oneof=A B C"` // Grade konfirmasi
	PriorityInterestRate float64 `json:"priority_interest_rate,omitempty"`     // Override if needed
	CatalystInterestRate float64 `json:"catalyst_interest_rate,omitempty"`     // Override if needed
	Notes                string  `json:"notes,omitempty"`                      // Catatan admin
}

// AdminGradeSuggestionResponse is the response for grade suggestion (BE-ADM-1)
type AdminGradeSuggestionResponse struct {
	InvoiceID         string  `json:"invoice_id"`
	SuggestedGrade    string  `json:"suggested_grade"` // A, B, or C
	GradeScore        int     `json:"grade_score"`     // 0-100
	CountryRisk       string  `json:"country_risk"`    // low, medium, high
	CountryScore      int     `json:"country_score"`   // Score contribution from country
	HistoryScore      int     `json:"history_score"`   // Score contribution from repeat buyer
	DocumentScore     int     `json:"document_score"`  // Score from document completeness
	IsRepeatBuyer     bool    `json:"is_repeat_buyer"`
	DocumentsComplete bool    `json:"documents_complete"`
	FundingLimit      float64 `json:"funding_limit"` // 60% for new, 100% for repeat
}

// ApproveInvoiceRequest - legacy
type ApproveInvoiceRequest struct {
	InvoiceID    uuid.UUID `json:"invoice_id" binding:"required"`
	InterestRate float64   `json:"interest_rate" binding:"required,gt=0,lte=100"`
}

type RejectInvoiceRequest struct {
	InvoiceID uuid.UUID `json:"invoice_id" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}

// DocumentValidationStatus represents the validation status of each document
type DocumentValidationStatus struct {
	DocumentID    string `json:"document_id"`
	DocumentType  string `json:"document_type"`
	FileName      string `json:"file_name"`
	FileURL       string `json:"file_url"`
	IsValid       bool   `json:"is_valid"`       // Admin marked as valid
	NeedsRevision bool   `json:"needs_revision"` // Admin marked as needs revision
	RevisionNote  string `json:"revision_note,omitempty"`
}

// InvoiceReviewData contains all data for admin review (Flow 5 - Split Screen)
type InvoiceReviewData struct {
	Invoice         Invoice                      `json:"invoice"`
	Exporter        *UserProfile                 `json:"exporter"`
	Documents       []DocumentValidationStatus   `json:"documents"`
	GradeSuggestion AdminGradeSuggestionResponse `json:"grade_suggestion"`
}

// ValidateDocumentRequest is the request to validate/revise a document
type ValidateDocumentRequest struct {
	DocumentID    string `json:"document_id" binding:"required"`
	IsValid       bool   `json:"is_valid"`
	NeedsRevision bool   `json:"needs_revision"`
	RevisionNote  string `json:"revision_note,omitempty"`
}

type InvoiceListResponse struct {
	Invoices   []Invoice `json:"invoices"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}

type InvoiceFilter struct {
	Status     *InvoiceStatus `json:"status,omitempty"`
	BuyerID    *uuid.UUID     `json:"buyer_id,omitempty"`
	ExporterID *uuid.UUID     `json:"exporter_id,omitempty"`
	MinAmount  *float64       `json:"min_amount,omitempty"`
	MaxAmount  *float64       `json:"max_amount,omitempty"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
}

// RepeatBuyerCheckRequest is for checking if buyer is repeat buyer
type RepeatBuyerCheckRequest struct {
	BuyerCompanyName string `json:"buyer_company_name" binding:"required"`
}

// RepeatBuyerCheckResponse contains the result of repeat buyer check
type RepeatBuyerCheckResponse struct {
	IsRepeatBuyer        bool    `json:"is_repeat_buyer"`
	Message              string  `json:"message"`
	PreviousTransactions int     `json:"previous_transactions,omitempty"`
	FundingLimit         float64 `json:"funding_limit"` // 60% for new, 100% for repeat
}

// EstimatedDisbursement contains the net disbursement calculation
type EstimatedDisbursement struct {
	GrossAmount     float64 `json:"gross_amount"`     // Total IDR amount
	PlatformFee     float64 `json:"platform_fee"`     // 2%
	NetDisbursement float64 `json:"net_disbursement"` // After fee
	Currency        string  `json:"currency"`         // IDR
}
