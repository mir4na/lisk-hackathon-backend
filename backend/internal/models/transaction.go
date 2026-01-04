package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string
type TransactionStatus string

const (
	TxTypeInvestment     TransactionType = "investment"
	TxTypeAdvancePayment TransactionType = "advance_payment"
	TxTypeBuyerRepayment TransactionType = "buyer_repayment"
	TxTypeInvestorReturn TransactionType = "investor_return"
	TxTypePlatformFee    TransactionType = "platform_fee"
	TxTypeRefund         TransactionType = "refund"

	TxStatusPending   TransactionStatus = "pending"
	TxStatusConfirmed TransactionStatus = "confirmed"
	TxStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID          uuid.UUID         `json:"id"`
	InvoiceID   *uuid.UUID        `json:"invoice_id,omitempty"`
	UserID      *uuid.UUID        `json:"user_id,omitempty"`
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	TxHash      *string           `json:"tx_hash,omitempty"`
	Status      TransactionStatus `json:"status"`
	FromAddress *string           `json:"from_address,omitempty"`
	ToAddress   *string           `json:"to_address,omitempty"`
	BlockNumber *int64            `json:"block_number,omitempty"`
	GasUsed     *int64            `json:"gas_used,omitempty"`
	Notes       *string           `json:"notes,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	PerPage      int           `json:"per_page"`
	TotalPages   int           `json:"total_pages"`
}

type TransactionFilter struct {
	Type      *TransactionType   `json:"type,omitempty"`
	Status    *TransactionStatus `json:"status,omitempty"`
	InvoiceID *uuid.UUID         `json:"invoice_id,omitempty"`
	UserID    *uuid.UUID         `json:"user_id,omitempty"`
	Page      int                `json:"page"`
	PerPage   int                `json:"per_page"`
}
