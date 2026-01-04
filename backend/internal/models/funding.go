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

type FundingPool struct {
	ID            uuid.UUID   `json:"id"`
	InvoiceID     uuid.UUID   `json:"invoice_id"`
	TargetAmount  float64     `json:"target_amount"`
	FundedAmount  float64     `json:"funded_amount"`
	InvestorCount int         `json:"investor_count"`
	Status        PoolStatus  `json:"status"`
	OpenedAt      *time.Time  `json:"opened_at,omitempty"`
	FilledAt      *time.Time  `json:"filled_at,omitempty"`
	DisbursedAt   *time.Time  `json:"disbursed_at,omitempty"`
	ClosedAt      *time.Time  `json:"closed_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`

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
	PoolID uuid.UUID `json:"pool_id" binding:"required"`
	Amount float64   `json:"amount" binding:"required,gt=0"`
}

type FundingPoolResponse struct {
	Pool             FundingPool `json:"pool"`
	RemainingAmount  float64     `json:"remaining_amount"`
	PercentageFunded float64     `json:"percentage_funded"`
	Invoice          *Invoice    `json:"invoice,omitempty"`
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
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}
