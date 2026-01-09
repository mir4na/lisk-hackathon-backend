package models

import (
	"time"

	"github.com/google/uuid"
)

// PoolDetailResponse is the response for pool detail page (Flow 6 - Project Detail)
type PoolDetailResponse struct {
	PoolID        uuid.UUID `json:"pool_id"`
	InvoiceID     uuid.UUID `json:"invoice_id"`
	ProjectTitle  string    `json:"project_title"`
	InvoiceNumber string    `json:"invoice_number"`

	// Grade & Risk
	Grade      string `json:"grade"`
	GradeScore int    `json:"grade_score"`
	IsInsured  bool   `json:"is_insured"`

	// Amounts
	TargetAmount    float64 `json:"target_amount"`
	FundedAmount    float64 `json:"funded_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	FundingProgress float64 `json:"funding_progress"`

	// Tenor
	TenorDays      int        `json:"tenor_days"`
	TenorDisplay   string     `json:"tenor_display"`
	DueDate        *time.Time `json:"due_date"`
	Deadline       *time.Time `json:"deadline"`
	RemainingTime  string     `json:"remaining_time"`
	RemainingHours int        `json:"remaining_hours"`

	// Status
	Status   string `json:"status"`
	Currency string `json:"currency"`

	// Detail Info
	BuyerInfo    BuyerDetailInfo    `json:"buyer_info"`
	ExporterInfo ExporterDetailInfo `json:"exporter_info"`
	Documents    []DocumentInfo     `json:"documents"`

	// Tranche Info
	PriorityTranche TrancheInfo `json:"priority_tranche"`
	CatalystTranche TrancheInfo `json:"catalyst_tranche"`
}

// BuyerDetailInfo contains buyer info for detail page
type BuyerDetailInfo struct {
	CompanyName  string `json:"company_name"`
	Country      string `json:"country"`
	CountryFlag  string `json:"country_flag"`
	CountryRisk  string `json:"country_risk"`
	Industry     string `json:"industry"`
	IsRepeat     bool   `json:"is_repeat"`
	TotalHistory int    `json:"total_history"`
}

// ExporterDetailInfo contains exporter info for detail page
type ExporterDetailInfo struct {
	CompanyName     string  `json:"company_name"`
	IsVerified      bool    `json:"is_verified"`
	CreditLimit     float64 `json:"credit_limit"`
	AvailableCredit float64 `json:"available_credit"`
	TotalInvoices   int     `json:"total_invoices"`
	SuccessRate     float64 `json:"success_rate"`
}

// DocumentInfo contains document info for preview/download
type DocumentInfo struct {
	Type       string `json:"type"`
	TypeLabel  string `json:"type_label"`
	IsVerified bool   `json:"is_verified"`
	IPFSHash   string `json:"ipfs_hash"`
}

// TrancheInfo contains tranche-specific info for UI
type TrancheInfo struct {
	Type                string  `json:"type"`
	TypeDisplay         string  `json:"type_display"`
	Description         string  `json:"description"`
	TargetAmount        float64 `json:"target_amount"`
	FundedAmount        float64 `json:"funded_amount"`
	RemainingAmount     float64 `json:"remaining_amount"`
	ProgressPercent     float64 `json:"progress_percent"`
	InterestRate        float64 `json:"interest_rate"`
	InterestRateDisplay string  `json:"interest_rate_display"`
	RiskLevel           string  `json:"risk_level"`
	RiskLevelDisplay    string  `json:"risk_level_display"`
	InfoBox             string  `json:"info_box"`
}

// InvestmentCalculatorRequest is the request for calculating investment returns
type InvestmentCalculatorRequest struct {
	PoolID  uuid.UUID `json:"pool_id" binding:"required"`
	Amount  float64   `json:"amount" binding:"required,gt=0"`
	Tranche string    `json:"tranche" binding:"required,oneof=priority catalyst"`
}

// InvestmentCalculatorResponse is the response for investment calculator
type InvestmentCalculatorResponse struct {
	PoolID         uuid.UUID `json:"pool_id"`
	Tranche        string    `json:"tranche"`
	TrancheDisplay string    `json:"tranche_display"`
	Principal      float64   `json:"principal"`
	InterestRate   float64   `json:"interest_rate"`
	TenorDays      int       `json:"tenor_days"`
	GrossInterest  float64   `json:"gross_interest"`
	PlatformFee    float64   `json:"platform_fee"`
	NetInterest    float64   `json:"net_interest"`
	TotalReturn    float64   `json:"total_return"`
	NetTotalReturn float64   `json:"net_total_return"`
	EffectiveRate  float64   `json:"effective_rate"`
	MaxInvestable  float64   `json:"max_investable"`
	CanInvest      bool      `json:"can_invest"`
	Message        string    `json:"message,omitempty"`
}

// InvestConfirmationRequest is the request for confirming investment (Flow 6)
type InvestConfirmationRequest struct {
	PoolID  uuid.UUID `json:"pool_id" binding:"required"`
	Amount  float64   `json:"amount" binding:"required,gt=0"`
	Tranche string    `json:"tranche" binding:"required,oneof=priority catalyst"`

	// For Priority Tranche
	TermsAccepted bool `json:"terms_accepted"`

	// For Catalyst Tranche - 2 mandatory checkboxes
	CatalystWarning1 bool `json:"catalyst_warning_1,omitempty"` // "Saya sadar dana ini menjadi jaminan pertama jika gagal bayar"
	CatalystWarning2 bool `json:"catalyst_warning_2,omitempty"` // "Saya siap menanggung risiko kehilangan modal"
}

// InvestmentConfirmationData contains data for confirmation UI
type InvestmentConfirmationData struct {
	PoolID          uuid.UUID `json:"pool_id"`
	ProjectTitle    string    `json:"project_title"`
	Amount          float64   `json:"amount"`
	Tranche         string    `json:"tranche"`
	TrancheDisplay  string    `json:"tranche_display"`
	InterestRate    float64   `json:"interest_rate"`
	EstimatedReturn float64   `json:"estimated_return"`
	TotalExpected   float64   `json:"total_expected"`
	DueDate         time.Time `json:"due_date"`

	// UI guidance
	RequiredCheckboxes  []string `json:"required_checkboxes"`
	WarningLevel        string   `json:"warning_level"` // normal, high (for catalyst)
	ConfirmationTitle   string   `json:"confirmation_title"`
	ConfirmationMessage string   `json:"confirmation_message"`
}

// PaymentMethodOption represents a payment method for midtrans
type PaymentMethodOption struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"` // bank_transfer, qris, etc.
	Icon        string `json:"icon"`
	IsAvailable bool   `json:"is_available"`
}

// CreatePaymentRequest is the request for creating payment
type CreatePaymentRequest struct {
	InvestmentID  uuid.UUID `json:"investment_id" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
}

// PaymentResponse is the response after creating payment
type PaymentResponse struct {
	InvestmentID   uuid.UUID `json:"investment_id"`
	PaymentID      string    `json:"payment_id"`
	PaymentURL     string    `json:"payment_url"` // Midtrans redirect URL
	VirtualAccount string    `json:"virtual_account,omitempty"`
	QRCodeURL      string    `json:"qr_code_url,omitempty"`
	Amount         float64   `json:"amount"`
	ExpiresAt      time.Time `json:"expires_at"`
	Status         string    `json:"status"`
}

// BalanceInfo represents user balance information (Flow 3)
type BalanceInfo struct {
	UserID     uuid.UUID `json:"user_id"`
	Role       string    `json:"role"` // investor or mitra
	BalanceIDR float64   `json:"balance_idr"`
	Currency   string    `json:"currency"`

	// For Investor: funds in active funding
	ActiveFunding float64 `json:"active_funding,omitempty"`

	// For Mitra: amount owed to investors
	TotalOwed     float64 `json:"total_owed,omitempty"`
	TotalInterest float64 `json:"total_interest,omitempty"`

	Description string `json:"description"`
}
