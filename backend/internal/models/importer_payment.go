package models

import (
	"time"

	"github.com/google/uuid"
)

// ImporterPaymentStatus represents the payment status from importer
type ImporterPaymentStatus string

const (
	ImporterPaymentStatusPending  ImporterPaymentStatus = "pending"
	ImporterPaymentStatusPaid     ImporterPaymentStatus = "paid"
	ImporterPaymentStatusOverdue  ImporterPaymentStatus = "overdue"
	ImporterPaymentStatusCanceled ImporterPaymentStatus = "canceled"
)

// ImporterPayment represents a payment that importir (buyer) needs to make
// This is for non-user importers who pay via payment ID sent to their email
type ImporterPayment struct {
	ID            uuid.UUID             `json:"id"`
	InvoiceID     uuid.UUID             `json:"invoice_id"`
	PoolID        uuid.UUID             `json:"pool_id"`
	BuyerEmail    string                `json:"buyer_email"`
	BuyerName     string                `json:"buyer_name"`
	AmountDue     float64               `json:"amount_due"` // Total (target + interest)
	AmountPaid    float64               `json:"amount_paid"`
	Currency      string                `json:"currency"`
	PaymentStatus ImporterPaymentStatus `json:"payment_status"`
	DueDate       time.Time             `json:"due_date"`
	PaidAt        *time.Time            `json:"paid_at,omitempty"`
	TxHash        *string               `json:"tx_hash,omitempty"` // Blockchain tx hash after payment
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`

	// Relations
	Invoice *Invoice     `json:"invoice,omitempty"`
	Pool    *FundingPool `json:"pool,omitempty"`
}

// ImporterPaymentRequest is request body for importer to pay
type ImporterPaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// ImporterPaymentResponse is the response after payment
type ImporterPaymentResponse struct {
	PaymentID  uuid.UUID  `json:"payment_id"`
	Status     string     `json:"status"`
	AmountPaid float64    `json:"amount_paid"`
	TxHash     *string    `json:"tx_hash,omitempty"`
	Message    string     `json:"message"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
}

// PaymentNotificationData is data for email notification to importer
type PaymentNotificationData struct {
	PaymentID     string
	InvoiceNumber string
	BuyerName     string
	ExporterName  string
	AmountDue     float64
	Currency      string
	DueDate       time.Time
	PaymentLink   string
}

// ExporterPaymentNotificationData is data for email notification to exporter
// when funding pool ends - contains invoice details for importer to pay
type ExporterPaymentNotificationData struct {
	InvoiceID       string
	InvoiceNumber   string
	ExporterName    string
	BuyerName       string
	BuyerEmail      string
	PrincipalAmount float64    // Original invoice amount funded
	TotalInterest   float64    // Total interest to be paid by importer
	TotalAmountDue  float64    // Principal + Total Interest
	Currency        string
	DueDate         time.Time  // Invoice due date
	InvestorDetails []InvestorPaymentDetail
	PaymentID       string
	PaymentLink     string
}

// InvestorPaymentDetail contains details for each investor's share
type InvestorPaymentDetail struct {
	InvestorID     string
	Amount         float64
	InterestRate   float64
	ExpectedReturn float64
	Tranche        string
}
