package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

// PaymentService handles dummy payment gateway operations (prototype)
type PaymentService struct {
	userRepo repository.UserRepositoryInterface
	txRepo   repository.TransactionRepositoryInterface
}

func NewPaymentService(
	userRepo repository.UserRepositoryInterface,
	txRepo repository.TransactionRepositoryInterface,
) *PaymentService {
	return &PaymentService{
		userRepo: userRepo,
		txRepo:   txRepo,
	}
}

// DepositRequest represents a deposit request
type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// PaymentResponse represents a payment operation response
type PaymentResponse struct {
	Success       bool      `json:"success"`
	TransactionID uuid.UUID `json:"transaction_id,omitempty"`
	Message       string    `json:"message"`
	NewBalance    float64   `json:"new_balance"`
	Timestamp     time.Time `json:"timestamp"`
}

// BalanceResponse represents user balance info
type BalanceResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	BalanceIDR float64   `json:"balance_idr"`
	Currency   string    `json:"currency"`
}

// SimulateDeposit simulates depositing funds to user balance (PROTOTYPE)
// In production, this would integrate with Midtrans or other payment gateway
func (s *PaymentService) SimulateDeposit(userID uuid.UUID, amount float64) (*PaymentResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update user balance
	newBalance := user.BalanceIDR + amount
	if err := s.userRepo.UpdateBalance(userID, newBalance); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &models.Transaction{
		UserID:   &userID,
		Type:     models.TxTypeDeposit,
		Amount:   amount,
		Currency: "IDR",
		Status:   models.TxStatusConfirmed,
		Notes:    stringPtr("Simulated deposit via prototype payment gateway"),
	}
	s.txRepo.Create(tx)

	return &PaymentResponse{
		Success:       true,
		TransactionID: tx.ID,
		Message:       "Deposit successful (prototype)",
		NewBalance:    newBalance,
		Timestamp:     time.Now(),
	}, nil
}

// SimulateWithdraw simulates withdrawing funds from user balance (PROTOTYPE)
func (s *PaymentService) SimulateWithdraw(userID uuid.UUID, amount float64) (*PaymentResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check sufficient balance
	if user.BalanceIDR < amount {
		return nil, errors.New("insufficient balance")
	}

	// Update user balance
	newBalance := user.BalanceIDR - amount
	if err := s.userRepo.UpdateBalance(userID, newBalance); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &models.Transaction{
		UserID:   &userID,
		Type:     models.TxTypeWithdrawal,
		Amount:   amount,
		Currency: "IDR",
		Status:   models.TxStatusConfirmed,
		Notes:    stringPtr("Simulated withdrawal via prototype payment gateway"),
	}
	s.txRepo.Create(tx)

	return &PaymentResponse{
		Success:       true,
		TransactionID: tx.ID,
		Message:       "Withdrawal successful (prototype)",
		NewBalance:    newBalance,
		Timestamp:     time.Now(),
	}, nil
}

// GetBalance returns user's current balance
func (s *PaymentService) GetBalance(userID uuid.UUID) (*BalanceResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &BalanceResponse{
		UserID:     userID,
		BalanceIDR: user.BalanceIDR,
		Currency:   "IDR",
	}, nil
}

// SimulateInvestmentPayment simulates payment for investment (PROTOTYPE)
// This deducts from user balance for investment
func (s *PaymentService) SimulateInvestmentPayment(userID uuid.UUID, amount float64, poolID uuid.UUID) (*PaymentResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check sufficient balance
	if user.BalanceIDR < amount {
		return nil, errors.New("insufficient balance for investment")
	}

	// Update user balance
	newBalance := user.BalanceIDR - amount
	if err := s.userRepo.UpdateBalance(userID, newBalance); err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Success:    true,
		Message:    "Investment payment processed (prototype)",
		NewBalance: newBalance,
		Timestamp:  time.Now(),
	}, nil
}

// SimulateReturnPayment simulates crediting investment returns to user (PROTOTYPE)
func (s *PaymentService) SimulateReturnPayment(userID uuid.UUID, amount float64) (*PaymentResponse, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update user balance
	newBalance := user.BalanceIDR + amount
	if err := s.userRepo.UpdateBalance(userID, newBalance); err != nil {
		return nil, err
	}

	return &PaymentResponse{
		Success:    true,
		Message:    "Investment return credited (prototype)",
		NewBalance: newBalance,
		Timestamp:  time.Now(),
	}, nil
}

// AdminGrantBalanceRequest represents admin grant balance request
type AdminGrantBalanceRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

// AdminGrantBalance allows admin to grant balance to any user (MVP ONLY)
func (s *PaymentService) AdminGrantBalance(targetUserID uuid.UUID, amount float64) (*PaymentResponse, error) {
	// Get target user
	user, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update user balance (can be positive or negative)
	newBalance := user.BalanceIDR + amount
	if newBalance < 0 {
		return nil, errors.New("resulting balance cannot be negative")
	}

	if err := s.userRepo.UpdateBalance(targetUserID, newBalance); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &models.Transaction{
		UserID:   &targetUserID,
		Type:     models.TxTypeDeposit,
		Amount:   amount,
		Currency: "IDR",
		Status:   models.TxStatusConfirmed,
		Notes:    stringPtr("Admin granted balance (MVP)"),
	}
	s.txRepo.Create(tx)

	return &PaymentResponse{
		Success:       true,
		TransactionID: tx.ID,
		Message:       "Balance granted successfully by admin",
		NewBalance:    newBalance,
		Timestamp:     time.Now(),
	}, nil
}
