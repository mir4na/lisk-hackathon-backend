package models

import (
	"time"

	"github.com/google/uuid"
)

type Buyer struct {
	ID            uuid.UUID `json:"id"`
	CreatedBy     uuid.UUID `json:"created_by"`
	CompanyName   string    `json:"company_name"`
	Country       string    `json:"country"`
	Address       *string   `json:"address,omitempty"`
	ContactEmail  *string   `json:"contact_email,omitempty"`
	ContactPhone  *string   `json:"contact_phone,omitempty"`
	Website       *string   `json:"website,omitempty"`
	CreditScore   int       `json:"credit_score"`
	TotalInvoices int       `json:"total_invoices"`
	TotalPaid     float64   `json:"total_paid"`
	TotalDefaulted float64  `json:"total_defaulted"`
	IsVerified    bool      `json:"is_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateBuyerRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	Country      string  `json:"country" binding:"required"`
	Address      *string `json:"address,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Website      *string `json:"website,omitempty"`
}

type UpdateBuyerRequest struct {
	CompanyName  string  `json:"company_name"`
	Country      string  `json:"country"`
	Address      *string `json:"address,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Website      *string `json:"website,omitempty"`
}

type BuyerListResponse struct {
	Buyers     []Buyer `json:"buyers"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	PerPage    int     `json:"per_page"`
	TotalPages int     `json:"total_pages"`
}
