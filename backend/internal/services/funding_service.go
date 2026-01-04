package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/config"
	"github.com/receiv3/backend/internal/models"
	"github.com/receiv3/backend/internal/repository"
)

type FundingService struct {
	fundingRepo repository.FundingRepositoryInterface
	invoiceRepo repository.InvoiceRepositoryInterface
	txRepo      repository.TransactionRepositoryInterface
	cfg         *config.Config
}

func NewFundingService(
	fundingRepo repository.FundingRepositoryInterface,
	invoiceRepo repository.InvoiceRepositoryInterface,
	txRepo repository.TransactionRepositoryInterface,
	cfg *config.Config,
) *FundingService {
	return &FundingService{
		fundingRepo: fundingRepo,
		invoiceRepo: invoiceRepo,
		txRepo:      txRepo,
		cfg:         cfg,
	}
}

func (s *FundingService) CreatePool(invoiceID uuid.UUID) (*models.FundingPool, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.Status != models.StatusTokenized {
		return nil, errors.New("invoice must be tokenized before creating funding pool")
	}
	if invoice.AdvanceAmount == nil {
		return nil, errors.New("advance amount not set")
	}

	// Check if pool already exists
	existing, _ := s.fundingRepo.FindPoolByInvoiceID(invoiceID)
	if existing != nil {
		return nil, errors.New("funding pool already exists for this invoice")
	}

	pool := &models.FundingPool{
		InvoiceID:    invoiceID,
		TargetAmount: *invoice.AdvanceAmount,
		Status:       models.PoolStatusOpen,
	}

	if err := s.fundingRepo.CreatePool(pool); err != nil {
		return nil, err
	}

	// Update invoice status
	if err := s.invoiceRepo.UpdateStatus(invoiceID, models.StatusFunding); err != nil {
		return nil, err
	}

	return pool, nil
}

func (s *FundingService) GetPool(poolID uuid.UUID) (*models.FundingPoolResponse, error) {
	pool, err := s.fundingRepo.FindPoolByID(poolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}

	invoice, _ := s.invoiceRepo.FindByID(pool.InvoiceID)

	remaining := pool.TargetAmount - pool.FundedAmount
	percentage := (pool.FundedAmount / pool.TargetAmount) * 100

	return &models.FundingPoolResponse{
		Pool:             *pool,
		RemainingAmount:  remaining,
		PercentageFunded: percentage,
		Invoice:          invoice,
	}, nil
}

func (s *FundingService) GetOpenPools(page, perPage int) (*models.PoolListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	pools, total, err := s.fundingRepo.FindOpenPools(page, perPage)
	if err != nil {
		return nil, err
	}

	var responses []models.FundingPoolResponse
	for _, pool := range pools {
		invoice, _ := s.invoiceRepo.FindByID(pool.InvoiceID)
		remaining := pool.TargetAmount - pool.FundedAmount
		percentage := (pool.FundedAmount / pool.TargetAmount) * 100

		responses = append(responses, models.FundingPoolResponse{
			Pool:             pool,
			RemainingAmount:  remaining,
			PercentageFunded: percentage,
			Invoice:          invoice,
		})
	}

	return &models.PoolListResponse{
		Pools:      responses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: models.CalculateTotalPages(total, perPage),
	}, nil
}

func (s *FundingService) Invest(investorID uuid.UUID, req *models.InvestRequest) (*models.Investment, error) {
	pool, err := s.fundingRepo.FindPoolByID(req.PoolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}
	if pool.Status != models.PoolStatusOpen {
		return nil, errors.New("pool is not open for investment")
	}

	remaining := pool.TargetAmount - pool.FundedAmount
	if req.Amount > remaining {
		return nil, errors.New("investment amount exceeds remaining pool capacity")
	}

	// Get invoice for interest rate calculation
	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}

	// Calculate expected return
	interestRate := 10.0 // Default
	if invoice.InterestRate != nil {
		interestRate = *invoice.InterestRate
	}

	// Calculate days until maturity
	daysToMaturity := time.Until(invoice.DueDate).Hours() / 24
	if daysToMaturity < 0 {
		daysToMaturity = 0
	}

	// Pro-rata interest calculation
	expectedReturn := req.Amount * (1 + (interestRate/100)*(daysToMaturity/365))

	investment := &models.Investment{
		PoolID:         req.PoolID,
		InvestorID:     investorID,
		Amount:         req.Amount,
		ExpectedReturn: expectedReturn,
		Status:         models.InvestmentStatusActive,
	}

	if err := s.fundingRepo.CreateInvestment(investment); err != nil {
		return nil, err
	}

	// Update pool funding
	if err := s.fundingRepo.UpdatePoolFunding(req.PoolID, req.Amount); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &investorID,
		Type:      models.TxTypeInvestment,
		Amount:    req.Amount,
		Currency:  "USDC",
		Status:    models.TxStatusPending,
	}
	s.txRepo.Create(tx)

	// Check if pool is now filled
	updatedPool, _ := s.fundingRepo.FindPoolByID(req.PoolID)
	if updatedPool.FundedAmount >= updatedPool.TargetAmount {
		s.fundingRepo.UpdatePoolStatus(req.PoolID, models.PoolStatusFilled)
	}

	return investment, nil
}

func (s *FundingService) GetInvestmentsByInvestor(investorID uuid.UUID, page, perPage int) (*models.InvestmentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	investments, total, err := s.fundingRepo.FindInvestmentsByInvestor(investorID, page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.InvestmentListResponse{
		Investments: investments,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  models.CalculateTotalPages(total, perPage),
	}, nil
}

func (s *FundingService) DisburseToExporter(poolID uuid.UUID) error {
	pool, err := s.fundingRepo.FindPoolByID(poolID)
	if err != nil {
		return err
	}
	if pool == nil {
		return errors.New("pool not found")
	}
	if pool.Status != models.PoolStatusFilled {
		return errors.New("pool must be filled before disbursement")
	}

	// Update pool status
	if err := s.fundingRepo.UpdatePoolStatus(poolID, models.PoolStatusDisbursed); err != nil {
		return err
	}

	// Update invoice status
	if err := s.invoiceRepo.UpdateStatus(pool.InvoiceID, models.StatusFunded); err != nil {
		return err
	}

	// Get invoice for exporter info
	invoice, _ := s.invoiceRepo.FindByID(pool.InvoiceID)

	// Create advance payment transaction
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &invoice.ExporterID,
		Type:      models.TxTypeAdvancePayment,
		Amount:    pool.FundedAmount,
		Currency:  "USDC",
		Status:    models.TxStatusPending,
	}
	s.txRepo.Create(tx)

	return nil
}

func (s *FundingService) ProcessRepayment(invoiceID uuid.UUID, amount float64) error {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}

	pool, err := s.fundingRepo.FindPoolByInvoiceID(invoiceID)
	if err != nil {
		return err
	}
	if pool == nil {
		return errors.New("pool not found")
	}

	// Get all investments for this pool
	investments, err := s.fundingRepo.FindInvestmentsByPool(pool.ID)
	if err != nil {
		return err
	}

	// Calculate platform fee
	platformFee := amount * (s.cfg.PlatformFeePercentage / 100)
	remainingAmount := amount - platformFee

	// Distribute to investors proportionally
	for _, inv := range investments {
		proportion := inv.Amount / pool.FundedAmount
		actualReturn := remainingAmount * proportion

		s.fundingRepo.UpdateInvestmentStatus(inv.ID, models.InvestmentStatusRepaid, &actualReturn)

		// Create return transaction
		tx := &models.Transaction{
			InvoiceID: &invoiceID,
			UserID:    &inv.InvestorID,
			Type:      models.TxTypeInvestorReturn,
			Amount:    actualReturn,
			Currency:  "USDC",
			Status:    models.TxStatusPending,
		}
		s.txRepo.Create(tx)
	}

	// Close pool and update invoice
	s.fundingRepo.UpdatePoolStatus(pool.ID, models.PoolStatusClosed)
	s.invoiceRepo.UpdateStatus(invoiceID, models.StatusRepaid)

	return nil
}
