package models

import (
	"time"

	"github.com/google/uuid"
)

type PoolStatus string

const (
	PoolStatusOpen      PoolStatus = "open"
	PoolStatusFilled    PoolStatus = "filled"
	PoolStatusDisbursed PoolStatus = "disbursed"
	PoolStatusClosed    PoolStatus = "closed"
)

// TrancheType represents the type of investment tranche
type TrancheType string

const (
	TranchePriority TrancheType = "priority" // Senior tranche - paid first, lower risk/return
	TrancheCatalyst TrancheType = "catalyst" // Junior tranche - paid last, higher risk/return
)

type FundingPool struct {
	ID            uuid.UUID  `json:"id"`
	InvoiceID     uuid.UUID  `json:"invoice_id"`
	TargetAmount  float64    `json:"target_amount"`
	FundedAmount  float64    `json:"funded_amount"`
	InvestorCount int        `json:"investor_count"`
	Status        PoolStatus `json:"status"`
	OpenedAt      *time.Time `json:"opened_at,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"` // Pool funding deadline
	FilledAt      *time.Time `json:"filled_at,omitempty"`
	DisbursedAt   *time.Time `json:"disbursed_at,omitempty"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Tranche fields
	PriorityTarget       float64 `json:"priority_target"`
	PriorityFunded       float64 `json:"priority_funded"`
	CatalystTarget       float64 `json:"catalyst_target"`
	CatalystFunded       float64 `json:"catalyst_funded"`
	PriorityInterestRate float64 `json:"priority_interest_rate"`
	CatalystInterestRate float64 `json:"catalyst_interest_rate"`
	PoolCurrency         string  `json:"pool_currency"`

	// Relations
	Invoice     *Invoice     `json:"invoice,omitempty"`
	Investments []Investment `json:"investments,omitempty"`
}

type InvestmentStatus string

const (
	InvestmentStatusActive    InvestmentStatus = "active"
	InvestmentStatusRepaid    InvestmentStatus = "repaid"
	InvestmentStatusDefaulted InvestmentStatus = "defaulted"
)

type Investment struct {
	ID             uuid.UUID        `json:"id"`
	PoolID         uuid.UUID        `json:"pool_id"`
	InvestorID     uuid.UUID        `json:"investor_id"`
	Amount         float64          `json:"amount"`
	ExpectedReturn float64          `json:"expected_return"`
	ActualReturn   *float64         `json:"actual_return,omitempty"`
	Status         InvestmentStatus `json:"status"`
	Tranche        TrancheType      `json:"tranche"`
	TxHash         *string          `json:"tx_hash,omitempty"`
	InvestedAt     time.Time        `json:"invested_at"`
	RepaidAt       *time.Time       `json:"repaid_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`

	// Relations
	Pool     *FundingPool `json:"pool,omitempty"`
	Investor *User        `json:"investor,omitempty"`
}

type InvestRequest struct {
	PoolID  uuid.UUID   `json:"pool_id" binding:"required"`
	Amount  float64     `json:"amount" binding:"required,gt=0"`
	Tranche TrancheType `json:"tranche" binding:"required,oneof=priority catalyst"`

	// Consent fields - inline per investment
	// For Priority: Only tnc_accepted required
	// For Catalyst: All 3 catalyst_consents + tnc_accepted required
	TncAccepted bool `json:"tnc_accepted" binding:"required"` // "Saya menyetujui Syarat & Ketentuan"

	// Catalyst-specific consents (required only when tranche=catalyst)
	// consent_1: "Saya sadar dana ini menjadi jaminan pertama jika gagal bayar."
	// consent_2: "Saya siap menanggung risiko kehilangan modal."
	// consent_3: "Saya paham ini bukan produk bank."
	CatalystConsents *CatalystConsents `json:"catalyst_consents,omitempty"`
}

// CatalystConsents contains the 3 required consents for Catalyst tranche
type CatalystConsents struct {
	FirstLossConsent bool `json:"first_loss_consent"` // "Saya sadar dana ini menjadi jaminan pertama jika gagal bayar."
	RiskLossConsent  bool `json:"risk_loss_consent"`  // "Saya siap menanggung risiko kehilangan modal."
	NotBankConsent   bool `json:"not_bank_consent"`   // "Saya paham ini bukan produk bank."
}

// AllAccepted checks if all catalyst consents are accepted
func (c *CatalystConsents) AllAccepted() bool {
	if c == nil {
		return false
	}
	return c.FirstLossConsent && c.RiskLossConsent && c.NotBankConsent
}

type FundingPoolResponse struct {
	Pool                     FundingPool `json:"pool"`
	RemainingAmount          float64     `json:"remaining_amount"`
	PercentageFunded         float64     `json:"percentage_funded"`
	PriorityRemaining        float64     `json:"priority_remaining"`
	CatalystRemaining        float64     `json:"catalyst_remaining"`
	PriorityPercentageFunded float64     `json:"priority_percentage_funded"`
	CatalystPercentageFunded float64     `json:"catalyst_percentage_funded"`
	Invoice                  *Invoice    `json:"invoice,omitempty"`
}

type PoolListResponse struct {
	Pools      []FundingPoolResponse `json:"pools"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PerPage    int                   `json:"per_page"`
	TotalPages int                   `json:"total_pages"`
}

type InvestmentListResponse struct {
	Investments []Investment `json:"investments"`
	Total       int          `json:"total"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
	TotalPages  int          `json:"total_pages"`
}

type PoolFilter struct {
	Status  *PoolStatus `json:"status,omitempty"`
	Grade   *string     `json:"grade,omitempty"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// MarketplaceFilter for filtering funding pools in marketplace (Flow 6)
type MarketplaceFilter struct {
	Grade     *string  `json:"grade,omitempty" form:"grade"`           // A, B, or C
	IsInsured *bool    `json:"is_insured,omitempty" form:"is_insured"` // Filter by insured status
	MinAmount *float64 `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount *float64 `json:"max_amount,omitempty" form:"max_amount"`
	SortBy    string   `json:"sort_by,omitempty" form:"sort_by"` // yield_desc, tenor_asc, newest
	Page      int      `json:"page" form:"page"`
	PerPage   int      `json:"per_page" form:"per_page"`
}

// MarketplacePoolResponse enhanced pool response for marketplace (Flow 6)
type MarketplacePoolResponse struct {
	FundingPoolResponse
	// Card Info
	ProjectTitle string `json:"project_title"` // e.g., "Kopi Arabika Gayo Batch #12"

	// Invoice grading info
	Grade            string `json:"grade"` // A, B, C with different colors
	GradeScore       int    `json:"grade_score"`
	IsInsured        bool   `json:"is_insured"`
	BuyerCountry     string `json:"buyer_country"`
	BuyerCountryFlag string `json:"buyer_country_flag"` // Flag emoji
	BuyerCompanyName string `json:"buyer_company_name"`
	BuyerCountryRisk string `json:"buyer_country_risk"`

	// Yield range display
	YieldRange string  `json:"yield_range"` // "10% - 15% p.a"
	MinYield   float64 `json:"min_yield"`   // Priority rate
	MaxYield   float64 `json:"max_yield"`   // Catalyst rate

	// Tenor
	TenorDays    int    `json:"tenor_days"`    // Days until due
	TenorDisplay string `json:"tenor_display"` // "60 Hari"

	// Progress info
	FundingProgress float64 `json:"funding_progress"` // Percentage funded
	RemainingAmount float64 `json:"remaining_amount"`
	RemainingTime   string  `json:"remaining_time"` // Human readable
	RemainingHours  int     `json:"remaining_hours"`
	IsFullyFunded   bool    `json:"is_fully_funded"` // For overlay display

	// Tranche info
	PriorityProgress float64 `json:"priority_progress"` // % of priority filled
	CatalystProgress float64 `json:"catalyst_progress"` // % of catalyst filled
}

// MarketplaceListResponse for marketplace listing
type MarketplaceListResponse struct {
	Pools      []MarketplacePoolResponse `json:"pools"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	PerPage    int                       `json:"per_page"`
	TotalPages int                       `json:"total_pages"`
}

// InvestorPortfolio represents the portfolio summary for an investor (Flow 9)
type InvestorPortfolio struct {
	// Summary Cards
	TotalFunding      float64 `json:"total_funding"`       // Total Pembiayaan (Saldo Dana yang Sedang Disalurkan)
	TotalExpectedGain float64 `json:"total_expected_gain"` // Total imbal hasil yang diharapkan
	TotalRealizedGain float64 `json:"total_realized_gain"` // Akumulasi profit yang sudah terealisasi

	// Donut Chart: Sebaran Dana
	PriorityAllocation float64 `json:"priority_allocation"` // Biru: Prioritas (Senior)
	CatalystAllocation float64 `json:"catalyst_allocation"` // Oranye: Katalis (Junior)

	// Counts
	ActiveInvestments int `json:"active_investments"`
	CompletedDeals    int `json:"completed_deals"`

	// Balance info
	AvailableBalance float64 `json:"available_balance"` // Saldo tersedia untuk funding
}

// InvestorActiveInvestment represents a single active investment for listing (Flow 10)
type InvestorActiveInvestment struct {
	InvestmentID    uuid.UUID `json:"investment_id"`
	ProjectName     string    `json:"project_name"` // e.g., "Kopi Gayo #12"
	InvoiceNumber   string    `json:"invoice_number"`
	BuyerName       string    `json:"buyer_name"`
	BuyerCountry    string    `json:"buyer_country"`
	BuyerFlag       string    `json:"buyer_flag"`       // Flag emoji
	Tranche         string    `json:"tranche"`          // priority / catalyst
	TrancheDisplay  string    `json:"tranche_display"`  // Prioritas / Katalis
	Principal       float64   `json:"principal"`        // Modal Disalurkan (Rp)
	InterestRate    float64   `json:"interest_rate"`    // Interest rate
	EstimatedReturn float64   `json:"estimated_return"` // Estimasi Hasil (Rp)
	TotalExpected   float64   `json:"total_expected"`   // Principal + Return
	DueDate         time.Time `json:"due_date"`
	DaysRemaining   int       `json:"days_remaining"`
	Status          string    `json:"status"`         // lancar, perhatian, gagal_bayar
	StatusDisplay   string    `json:"status_display"` // Lancar, Perhatian, Gagal Bayar
	StatusColor     string    `json:"status_color"`   // green, yellow, red
	InvestedAt      time.Time `json:"invested_at"`
}

// InvestorActiveInvestmentList is the response for investor active investments (Flow 10)
type InvestorActiveInvestmentList struct {
	Investments []InvestorActiveInvestment `json:"investments"`
	Total       int                        `json:"total"`
	Page        int                        `json:"page"`
	PerPage     int                        `json:"per_page"`
	TotalPages  int                        `json:"total_pages"`
	Summary     InvestorPortfolio          `json:"summary"`
}

// MitraDashboard represents the dashboard data for a mitra (Flow 8)
type MitraDashboard struct {
	TotalActiveFinancing  float64            `json:"total_active_financing"`  // Total pembiayaan aktif (Rp)
	TotalOwedToInvestors  float64            `json:"total_owed_to_investors"` // Total hutang ke investor (termasuk bunga)
	AverageRemainingTenor int                `json:"average_remaining_tenor"` // Sisa tenor rata-rata (hari)
	ActiveInvoices        []InvoiceDashboard `json:"active_invoices"`
	TimelineStatus        TimelineStatus     `json:"timeline_status"`
}

// TimelineStatus represents the timeline tracker status
type TimelineStatus struct {
	FundraisingComplete  bool   `json:"fundraising_complete"`  // Penggalangan Dana
	DisbursementComplete bool   `json:"disbursement_complete"` // Dana Cair
	RepaymentComplete    bool   `json:"repayment_complete"`    // Disbursement ke funder
	CurrentStep          string `json:"current_step"`          // Current step description
}

type InvoiceDashboard struct {
	InvoiceID     uuid.UUID `json:"invoice_id"`
	InvoiceNumber string    `json:"invoice_number"`
	BuyerName     string    `json:"buyer_name"`
	BuyerCountry  string    `json:"buyer_country"`
	DueDate       time.Time `json:"due_date"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`       // Aktif, Dalam Pengawasan
	StatusColor   string    `json:"status_color"` // green, yellow, red
	DaysRemaining int       `json:"days_remaining"`
	FundedAmount  float64   `json:"funded_amount"`
	TotalOwed     float64   `json:"total_owed"` // Termasuk bunga
}

// ActiveInvestmentListResponse is the paginated response for active investments
type ActiveInvestmentListResponse struct {
	Investments []InvestorActiveInvestment `json:"investments"`
	Total       int                        `json:"total"`
	Page        int                        `json:"page"`
	PerPage     int                        `json:"per_page"`
	TotalPages  int                        `json:"total_pages"`
}

// MitraInvoiceListResponse is the paginated response for mitra's invoices
type MitraInvoiceListResponse struct {
	Invoices   []InvoiceDashboard `json:"invoices"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}
