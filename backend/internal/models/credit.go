package models

import (
	"time"

	"github.com/google/uuid"
)

type CreditScore struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	Score              int       `json:"score"`
	TotalInvoices      int       `json:"total_invoices"`
	SuccessfulInvoices int       `json:"successful_invoices"`
	DefaultedInvoices  int       `json:"defaulted_invoices"`
	TotalVolume        float64   `json:"total_volume"`
	AvgPaymentDelay    int       `json:"avg_payment_delay"`
	LastUpdated        time.Time `json:"last_updated"`
	CreatedAt          time.Time `json:"created_at"`
}

type CreditScoreHistory struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	PreviousScore int        `json:"previous_score"`
	NewScore      int        `json:"new_score"`
	Reason        string     `json:"reason"`
	InvoiceID     *uuid.UUID `json:"invoice_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreditScoreResponse struct {
	Score       CreditScore          `json:"score"`
	History     []CreditScoreHistory `json:"history"`
	Percentile  float64              `json:"percentile"`
	RiskLevel   string               `json:"risk_level"`
}

func CalculateRiskLevel(score int) string {
	switch {
	case score >= 80:
		return "low"
	case score >= 60:
		return "medium-low"
	case score >= 40:
		return "medium"
	case score >= 20:
		return "medium-high"
	default:
		return "high"
	}
}
