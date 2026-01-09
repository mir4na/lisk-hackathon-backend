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
	userRepo    repository.UserRepositoryInterface
	txRepo      repository.TransactionRepositoryInterface
	fundingRepo repository.FundingRepositoryInterface
	invoiceRepo repository.InvoiceRepositoryInterface
}

func NewPaymentService(
	userRepo repository.UserRepositoryInterface,
	txRepo repository.TransactionRepositoryInterface,
	fundingRepo repository.FundingRepositoryInterface,
	invoiceRepo repository.InvoiceRepositoryInterface,
) *PaymentService {
	return &PaymentService{
		userRepo:    userRepo,
		txRepo:      txRepo,
		fundingRepo: fundingRepo,
		invoiceRepo: invoiceRepo,
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

// BalanceResponse represents user balance info with role-specific data (Flow 3)
type BalanceResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Role         string    `json:"role"`
	MemberStatus string    `json:"member_status"`
	BalanceIDR   float64   `json:"balance_idr"`
	Currency     string    `json:"currency"`

	// For Investor: funds currently in active investments
	ActiveFunding  float64 `json:"active_funding,omitempty"`
	ExpectedReturn float64 `json:"expected_return,omitempty"`

	// For Mitra: amount owed to investors
	TotalOwed     float64 `json:"total_owed,omitempty"`
	TotalInterest float64 `json:"total_interest,omitempty"`

	// Description based on role
	Description string `json:"description"`
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

// GetBalance returns user's current balance with role-specific info (Flow 3)
func (s *PaymentService) GetBalance(userID uuid.UUID) (*BalanceResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	response := &BalanceResponse{
		UserID:       userID,
		Role:         string(user.Role),
		MemberStatus: string(user.MemberStatus),
		BalanceIDR:   user.BalanceIDR,
		Currency:     "IDR",
	}

	// Role-specific balance info
	if user.Role == models.RoleInvestor {
		// For Investor: show funds in active investments
		portfolio, err := s.fundingRepo.GetInvestorPortfolio(userID)
		if err == nil && portfolio != nil {
			response.ActiveFunding = portfolio.TotalFunding
			response.ExpectedReturn = portfolio.TotalExpectedGain
		}
		response.Description = "Saldo dana yang tersedia untuk pendanaan"
	} else if user.Role == models.RoleExporter || user.Role == models.RoleMitra {
		// For Mitra: show total owed to investors
		filter := &models.InvoiceFilter{
			Page:    1,
			PerPage: 1000,
		}
		invoices, _, _ := s.invoiceRepo.FindByExporter(userID, filter)

		var totalOwed, totalInterest float64
		for _, invoice := range invoices {
			// Check for active invoices (funded or funding status)
			if invoice.Status == models.StatusFunding || invoice.Status == models.StatusFunded || invoice.Status == models.StatusMatured {
				pool, _ := s.fundingRepo.FindPoolByInvoiceID(invoice.ID)
				if pool != nil {
					// Calculate owed amount (principal + interest)
					priorityOwed := pool.PriorityFunded * (1 + pool.PriorityInterestRate/100)
					catalystOwed := pool.CatalystFunded * (1 + pool.CatalystInterestRate/100)
					totalOwed += priorityOwed + catalystOwed
					totalInterest += (pool.PriorityFunded * pool.PriorityInterestRate / 100) +
						(pool.CatalystFunded * pool.CatalystInterestRate / 100)
				}
			}
		}
		response.TotalOwed = totalOwed
		response.TotalInterest = totalInterest
		response.Description = "Total kewajiban kepada pendana (pokok + bunga)"
	}

	return response, nil
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
