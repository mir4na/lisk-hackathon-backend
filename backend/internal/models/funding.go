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

// MarketplaceFilter for filtering funding pools in marketplace
type MarketplaceFilter struct {
	Grade     *string  `json:"grade,omitempty" form:"grade"`           // A, B, or C
	IsInsured *bool    `json:"is_insured,omitempty" form:"is_insured"` // Filter by insured status
	MinAmount *float64 `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount *float64 `json:"max_amount,omitempty" form:"max_amount"`
	Page      int      `json:"page" form:"page"`
	PerPage   int      `json:"per_page" form:"per_page"`
}

// MarketplacePoolResponse enhanced pool response for marketplace
type MarketplacePoolResponse struct {
	FundingPoolResponse
	// Invoice grading info
	Grade            string `json:"grade"`
	GradeScore       int    `json:"grade_score"`
	IsInsured        bool   `json:"is_insured"`
	BuyerCountryRisk string `json:"buyer_country_risk"`
	// Progress info
	FundingProgress float64 `json:"funding_progress"` // Percentage funded
	RemainingAmount float64 `json:"remaining_amount"`
	RemainingTime   string  `json:"remaining_time"` // Human readable
	RemainingHours  int     `json:"remaining_hours"`
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

// InvestorPortfolio represents the portfolio summary for an investor
type InvestorPortfolio struct {
	TotalFunding       float64 `json:"total_funding"`
	TotalExpectedGain  float64 `json:"total_expected_gain"`
	TotalRealizedGain  float64 `json:"total_realized_gain"`
	PriorityAllocation float64 `json:"priority_allocation"`
	CatalystAllocation float64 `json:"catalyst_allocation"`
	ActiveInvestments  int     `json:"active_investments"`
	CompletedDeals     int     `json:"completed_deals"`
}

// MitraDashboard represents the dashboard data for a mitra
type MitraDashboard struct {
	TotalActiveFinancing  float64            `json:"total_active_financing"`
	AverageRemainingTenor int                `json:"average_remaining_tenor"`
	ActiveInvoices        []InvoiceDashboard `json:"active_invoices"`
	TimelineStatus        string             `json:"timeline_status"`
}

type InvoiceDashboard struct {
	InvoiceID     uuid.UUID `json:"invoice_id"`
	BuyerName     string    `json:"buyer_name"`
	DueDate       time.Time `json:"due_date"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	DaysRemaining int       `json:"days_remaining"`
}
