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
	PrincipalAmount float64 // Original invoice amount funded
	TotalInterest   float64 // Total interest to be paid by importer
	PlatformFee     float64 // Platform fee for the application (2%)
	TotalAmountDue  float64 // Principal + Total Interest + Platform Fee
	Currency        string
	DueDate         time.Time // Invoice due date
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

// ==================== Virtual Account Payment (Mitra Repayment) ====================

// VAStatus represents Virtual Account payment status
type VAStatus string

const (
	VAStatusPending   VAStatus = "pending"
	VAStatusPaid      VAStatus = "paid"
	VAStatusExpired   VAStatus = "expired"
	VAStatusCancelled VAStatus = "cancelled"
)

// VirtualAccount represents a VA for Mitra to pay back investors
type VirtualAccount struct {
	ID        uuid.UUID  `json:"id"`
	PoolID    uuid.UUID  `json:"pool_id"`
	UserID    uuid.UUID  `json:"user_id"` // Mitra user ID
	VANumber  string     `json:"va_number"`
	BankCode  string     `json:"bank_code"`
	BankName  string     `json:"bank_name"`
	Amount    float64    `json:"amount"` // Total amount due (principal + interest)
	Status    VAStatus   `json:"status"`
	ExpiresAt time.Time  `json:"expires_at"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CreateVARequest is request to create VA for mitra repayment
type CreateVARequest struct {
	PoolID   uuid.UUID `json:"pool_id" binding:"required"`
	BankCode string    `json:"bank_code" binding:"required"` // bca, mandiri, bni
}

// VAPaymentMethodOption represents available VA payment method
type VAPaymentMethodOption struct {
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
	LogoURL  string `json:"logo_url,omitempty"`
}

// GetVAPaymentMethods returns available VA payment methods
func GetVAPaymentMethods() []VAPaymentMethodOption {
	return []VAPaymentMethodOption{
		{BankCode: "bca", BankName: "Bank Central Asia (BCA)"},
		{BankCode: "mandiri", BankName: "Bank Mandiri"},
		{BankCode: "bni", BankName: "Bank Negara Indonesia (BNI)"},
	}
}

// MitraRepaymentBreakdown contains the breakdown of repayment by tranche
type MitraRepaymentBreakdown struct {
	PoolID    uuid.UUID `json:"pool_id"`
	InvoiceID uuid.UUID `json:"invoice_id"`
	InvoiceNo string    `json:"invoice_number"`
	BuyerName string    `json:"buyer_name"`
	DueDate   time.Time `json:"due_date"`

	// Principal amounts
	PriorityPrincipal float64 `json:"priority_principal"`
	CatalystPrincipal float64 `json:"catalyst_principal"`
	TotalPrincipal    float64 `json:"total_principal"`

	// Interest amounts (flat rate: principal × rate/100)
	PriorityInterestRate float64 `json:"priority_interest_rate"` // e.g., 10%
	CatalystInterestRate float64 `json:"catalyst_interest_rate"` // e.g., 15%
	PriorityInterest     float64 `json:"priority_interest"`
	CatalystInterest     float64 `json:"catalyst_interest"`
	TotalInterest        float64 `json:"total_interest"`

	// Total amounts
	PriorityTotal float64 `json:"priority_total"` // Priority Principal + Interest
	CatalystTotal float64 `json:"catalyst_total"` // Catalyst Principal + Interest
	PlatformFee   float64 `json:"platform_fee"`   // Platform fee for the application (2%)
	GrandTotal    float64 `json:"grand_total"`    // Total to pay (including platform fee)

	Currency string `json:"currency"`
}

// VAPaymentResponse is the response after creating VA
type VAPaymentResponse struct {
	VA             VirtualAccount          `json:"virtual_account"`
	Breakdown      MitraRepaymentBreakdown `json:"breakdown"`
	RemainingTime  string                  `json:"remaining_time"`  // e.g., "23:59:59"
	RemainingHours int                     `json:"remaining_hours"` // e.g., 24
	Microcopy      string                  `json:"microcopy"`
}

// VAPaymentPageResponse is for the VA payment page UI
type VAPaymentPageResponse struct {
	VANumber        string                  `json:"va_number"`
	BankCode        string                  `json:"bank_code"`
	BankName        string                  `json:"bank_name"`
	Amount          float64                 `json:"amount"`           // Closed amount, user cannot change
	AmountFormatted string                  `json:"amount_formatted"` // e.g., "Rp 55.000.000"
	Status          VAStatus                `json:"status"`
	ExpiresAt       time.Time               `json:"expires_at"`
	RemainingTime   string                  `json:"remaining_time"` // Timer: "23:59:59"
	Breakdown       MitraRepaymentBreakdown `json:"breakdown"`
	Microcopy       string                  `json:"microcopy"`
}

// MitraActiveInvoice is for displaying active invoices with pay button
type MitraActiveInvoice struct {
	InvoiceID     uuid.UUID `json:"invoice_id"`
	PoolID        uuid.UUID `json:"pool_id"`
	InvoiceNumber string    `json:"invoice_number"`
	BuyerName     string    `json:"buyer_name"`
	Amount        float64   `json:"amount"`    // Original invoice amount
	TotalDue      float64   `json:"total_due"` // Principal + Interest
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	DaysUntilDue  int       `json:"days_until_due"` // Negative if overdue
	CanPay        bool      `json:"can_pay"`        // True if pool is filled and can pay
}
