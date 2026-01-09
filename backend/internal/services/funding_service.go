package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

type FundingService struct {
	fundingRepo   repository.FundingRepositoryInterface
	invoiceRepo   repository.InvoiceRepositoryInterface
	txRepo        repository.TransactionRepositoryInterface
	userRepo      repository.UserRepositoryInterface
	buyerRepo     repository.BuyerRepositoryInterface
	rqRepo        repository.RiskQuestionnaireRepositoryInterface
	emailService  *EmailService
	escrowService *EscrowService
	cfg           *config.Config
}

func NewFundingService(
	fundingRepo repository.FundingRepositoryInterface,
	invoiceRepo repository.InvoiceRepositoryInterface,
	txRepo repository.TransactionRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	buyerRepo repository.BuyerRepositoryInterface,
	rqRepo repository.RiskQuestionnaireRepositoryInterface,
	emailService *EmailService,
	escrowService *EscrowService,
	cfg *config.Config,
) *FundingService {
	return &FundingService{
		fundingRepo:   fundingRepo,
		invoiceRepo:   invoiceRepo,
		txRepo:        txRepo,
		userRepo:      userRepo,
		buyerRepo:     buyerRepo,
		rqRepo:        rqRepo,
		emailService:  emailService,
		escrowService: escrowService,
		cfg:           cfg,
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

	// Calculate tranche targets based on invoice config
	totalTarget := *invoice.AdvanceAmount
	priorityRatio := invoice.PriorityRatio
	catalystRatio := invoice.CatalystRatio

	// Use defaults if not set
	if priorityRatio == 0 {
		priorityRatio = 80.0
	}
	if catalystRatio == 0 {
		catalystRatio = 20.0
	}

	priorityTarget := totalTarget * (priorityRatio / 100)
	catalystTarget := totalTarget * (catalystRatio / 100)

	// Get interest rates from invoice or use defaults
	priorityInterestRate := 10.0
	catalystInterestRate := 15.0
	if invoice.PriorityInterestRate != nil {
		priorityInterestRate = *invoice.PriorityInterestRate
	}
	if invoice.CatalystInterestRate != nil {
		catalystInterestRate = *invoice.CatalystInterestRate
	}

	// Calculate deadline based on funding duration days
	fundingDurationDays := invoice.FundingDurationDays
	if fundingDurationDays == 0 {
		fundingDurationDays = 14 // Default 14 days
	}
	deadline := time.Now().AddDate(0, 0, fundingDurationDays)

	pool := &models.FundingPool{
		InvoiceID:            invoiceID,
		TargetAmount:         totalTarget,
		Status:               models.PoolStatusOpen,
		Deadline:             &deadline,
		PriorityTarget:       priorityTarget,
		PriorityFunded:       0,
		CatalystTarget:       catalystTarget,
		CatalystFunded:       0,
		PriorityInterestRate: priorityInterestRate,
		CatalystInterestRate: catalystInterestRate,
		PoolCurrency:         "IDR",
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
	percentage := 0.0
	if pool.TargetAmount > 0 {
		percentage = (pool.FundedAmount / pool.TargetAmount) * 100
	}

	priorityRemaining := pool.PriorityTarget - pool.PriorityFunded
	catalystRemaining := pool.CatalystTarget - pool.CatalystFunded
	priorityPercentage := 0.0
	catalystPercentage := 0.0
	if pool.PriorityTarget > 0 {
		priorityPercentage = (pool.PriorityFunded / pool.PriorityTarget) * 100
	}
	if pool.CatalystTarget > 0 {
		catalystPercentage = (pool.CatalystFunded / pool.CatalystTarget) * 100
	}

	return &models.FundingPoolResponse{
		Pool:                     *pool,
		RemainingAmount:          remaining,
		PercentageFunded:         percentage,
		PriorityRemaining:        priorityRemaining,
		CatalystRemaining:        catalystRemaining,
		PriorityPercentageFunded: priorityPercentage,
		CatalystPercentageFunded: catalystPercentage,
		Invoice:                  invoice,
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
		percentage := 0.0
		if pool.TargetAmount > 0 {
			percentage = (pool.FundedAmount / pool.TargetAmount) * 100
		}

		priorityRemaining := pool.PriorityTarget - pool.PriorityFunded
		catalystRemaining := pool.CatalystTarget - pool.CatalystFunded
		priorityPercentage := 0.0
		catalystPercentage := 0.0
		if pool.PriorityTarget > 0 {
			priorityPercentage = (pool.PriorityFunded / pool.PriorityTarget) * 100
		}
		if pool.CatalystTarget > 0 {
			catalystPercentage = (pool.CatalystFunded / pool.CatalystTarget) * 100
		}

		responses = append(responses, models.FundingPoolResponse{
			Pool:                     pool,
			RemainingAmount:          remaining,
			PercentageFunded:         percentage,
			PriorityRemaining:        priorityRemaining,
			CatalystRemaining:        catalystRemaining,
			PriorityPercentageFunded: priorityPercentage,
			CatalystPercentageFunded: catalystPercentage,
			Invoice:                  invoice,
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

// GetMarketplacePools gets open pools with grading info for marketplace
func (s *FundingService) GetMarketplacePools(filter *models.MarketplaceFilter) (*models.MarketplaceListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}

	// Get all open pools
	pools, total, err := s.fundingRepo.FindOpenPools(filter.Page, filter.PerPage)
	if err != nil {
		return nil, err
	}

	var responses []models.MarketplacePoolResponse
	for _, pool := range pools {
		invoice, _ := s.invoiceRepo.FindByID(pool.InvoiceID)
		if invoice == nil {
			continue
		}

		// Apply filters
		if filter.Grade != nil && invoice.Grade != nil && *invoice.Grade != *filter.Grade {
			total--
			continue
		}
		if filter.IsInsured != nil && invoice.IsInsured != *filter.IsInsured {
			total--
			continue
		}
		if filter.MinAmount != nil && pool.TargetAmount < *filter.MinAmount {
			total--
			continue
		}
		if filter.MaxAmount != nil && pool.TargetAmount > *filter.MaxAmount {
			total--
			continue
		}

		// Calculate progress
		remaining := pool.TargetAmount - pool.FundedAmount
		percentage := 0.0
		if pool.TargetAmount > 0 {
			percentage = (pool.FundedAmount / pool.TargetAmount) * 100
		}

		priorityRemaining := pool.PriorityTarget - pool.PriorityFunded
		catalystRemaining := pool.CatalystTarget - pool.CatalystFunded
		priorityPercentage := 0.0
		catalystPercentage := 0.0
		if pool.PriorityTarget > 0 {
			priorityPercentage = (pool.PriorityFunded / pool.PriorityTarget) * 100
		}
		if pool.CatalystTarget > 0 {
			catalystPercentage = (pool.CatalystFunded / pool.CatalystTarget) * 100
		}

		// Calculate remaining time
		remainingTime := "N/A"
		remainingHours := 0
		if pool.Deadline != nil {
			duration := time.Until(*pool.Deadline)
			if duration > 0 {
				remainingHours = int(duration.Hours())
				days := remainingHours / 24
				hours := remainingHours % 24
				if days > 0 {
					remainingTime = fmt.Sprintf("%d hari %d jam", days, hours)
				} else {
					remainingTime = fmt.Sprintf("%d jam", hours)
				}
			} else {
				remainingTime = "Berakhir"
			}
		}

		// Get grade info
		grade := "C"
		gradeScore := 0
		countryRisk := "medium"
		if invoice.Grade != nil {
			grade = *invoice.Grade
		}
		if invoice.GradeScore != nil {
			gradeScore = *invoice.GradeScore
		}
		if invoice.BuyerCountryRisk != nil {
			countryRisk = *invoice.BuyerCountryRisk
		}

		responses = append(responses, models.MarketplacePoolResponse{
			FundingPoolResponse: models.FundingPoolResponse{
				Pool:                     pool,
				RemainingAmount:          remaining,
				PercentageFunded:         percentage,
				PriorityRemaining:        priorityRemaining,
				CatalystRemaining:        catalystRemaining,
				PriorityPercentageFunded: priorityPercentage,
				CatalystPercentageFunded: catalystPercentage,
				Invoice:                  invoice,
			},
			Grade:            grade,
			GradeScore:       gradeScore,
			IsInsured:        invoice.IsInsured,
			BuyerCountryRisk: countryRisk,
			FundingProgress:  percentage,
			RemainingAmount:  remaining,
			RemainingTime:    remainingTime,
			RemainingHours:   remainingHours,
			PriorityProgress: priorityPercentage,
			CatalystProgress: catalystPercentage,
		})
	}

	return &models.MarketplaceListResponse{
		Pools:      responses,
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalPages: models.CalculateTotalPages(total, filter.PerPage),
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

	// Validate consent based on tranche type
	// All investments require T&C acceptance
	if !req.TncAccepted {
		return nil, errors.New("you must accept Terms & Conditions to invest")
	}

	// Catalyst tranche requires additional consents (3 checkboxes)
	if req.Tranche == models.TrancheCatalyst {
		if req.CatalystConsents == nil || !req.CatalystConsents.AllAccepted() {
			return nil, errors.New("untuk Pendanaan Katalis, Anda harus menyetujui semua 3 pernyataan risiko")
		}
	}

	// Validate tranche-specific quota
	var trancheRemaining float64
	var interestRate float64

	if req.Tranche == models.TranchePriority {
		trancheRemaining = pool.PriorityTarget - pool.PriorityFunded
		interestRate = pool.PriorityInterestRate
	} else if req.Tranche == models.TrancheCatalyst {
		trancheRemaining = pool.CatalystTarget - pool.CatalystFunded
		interestRate = pool.CatalystInterestRate
	} else {
		return nil, errors.New("invalid tranche type, must be 'priority' or 'catalyst'")
	}

	if req.Amount > trancheRemaining {
		return nil, errors.New("investment amount exceeds remaining tranche capacity")
	}

	// Flat interest calculation: Principal + (Principal * Rate/100)
	// No time factor - interest is calculated flat from total principal
	expectedReturn := req.Amount + (req.Amount * interestRate / 100)

	investment := &models.Investment{
		PoolID:         req.PoolID,
		InvestorID:     investorID,
		Amount:         req.Amount,
		ExpectedReturn: expectedReturn,
		Status:         models.InvestmentStatusActive,
		Tranche:        req.Tranche,
	}

	if err := s.fundingRepo.CreateInvestment(investment); err != nil {
		return nil, err
	}

	// Update pool funding for specific tranche
	if err := s.fundingRepo.UpdatePoolTrancheFunding(req.PoolID, req.Amount, req.Tranche); err != nil {
		return nil, err
	}

	// Create transaction record (using IDR - abstracted escrow for MVP)
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &investorID,
		Type:      models.TxTypeInvestment,
		Amount:    req.Amount,
		Currency:  "IDR",
		Status:    models.TxStatusPending,
	}
	s.txRepo.Create(tx)

	// Check if pool is now filled (both tranches) - Flow 7: Auto-Disbursement
	updatedPool, _ := s.fundingRepo.FindPoolByID(req.PoolID)
	if updatedPool.FundedAmount >= updatedPool.TargetAmount {
		s.fundingRepo.UpdatePoolStatus(req.PoolID, models.PoolStatusFilled)

		// Trigger automatic disbursement to mitra (Flow 7)
		go func() {
			if err := s.DisburseToExporter(req.PoolID); err != nil {
				// Log error but don't fail investment
				// In production, this would alert admins
				_ = err
			}
		}()
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

func (s *FundingService) GetInvestorPortfolio(investorID uuid.UUID) (*models.InvestorPortfolio, error) {
	return s.fundingRepo.GetInvestorPortfolio(investorID)
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

	// Calculate platform fee (dummy escrow)
	platformFee := pool.FundedAmount * (s.cfg.PlatformFeePercentage / 100)
	disbursementAmount := pool.FundedAmount - platformFee

	// Create advance payment transaction (dummy escrow simulation)
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &invoice.ExporterID,
		Type:      models.TxTypeAdvancePayment,
		Amount:    disbursementAmount,
		Currency:  "IDRX",
		Status:    models.TxStatusPending,
		Notes:     stringPtr("Disbursement to exporter via dummy escrow"),
	}
	s.txRepo.Create(tx)

	return nil
}

// ProcessRepayment processes repayment from exporter and distributes to investors
// Following priority-first rule: Priority investors are paid first, then Catalyst
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

	// Calculate platform fee
	platformFee := amount * (s.cfg.PlatformFeePercentage / 100)
	remainingAmount := amount - platformFee

	// Get investments grouped by tranche
	priorityInvestments, err := s.fundingRepo.FindInvestmentsByPoolAndTranche(pool.ID, models.TranchePriority)
	if err != nil {
		return err
	}
	catalystInvestments, err := s.fundingRepo.FindInvestmentsByPoolAndTranche(pool.ID, models.TrancheCatalyst)
	if err != nil {
		return err
	}

	// Calculate total expected returns for each tranche
	totalPriorityExpected := 0.0
	for _, inv := range priorityInvestments {
		totalPriorityExpected += inv.ExpectedReturn
	}
	totalCatalystExpected := 0.0
	for _, inv := range catalystInvestments {
		totalCatalystExpected += inv.ExpectedReturn
	}

	// Priority-first distribution
	// Step 1: Pay priority investors first
	priorityPaid := 0.0
	for _, inv := range priorityInvestments {
		// Calculate pro-rata share
		var actualReturn float64
		if remainingAmount >= totalPriorityExpected {
			// Full repayment for priority
			actualReturn = inv.ExpectedReturn
		} else {
			// Partial repayment - distribute proportionally
			proportion := inv.ExpectedReturn / totalPriorityExpected
			actualReturn = remainingAmount * proportion
		}
		priorityPaid += actualReturn

		s.fundingRepo.UpdateInvestmentStatus(inv.ID, models.InvestmentStatusRepaid, &actualReturn)

		// Create return transaction
		tx := &models.Transaction{
			InvoiceID: &invoiceID,
			UserID:    &inv.InvestorID,
			Type:      models.TxTypeInvestorReturn,
			Amount:    actualReturn,
			Currency:  "IDRX",
			Status:    models.TxStatusPending,
			Notes:     stringPtr("Priority tranche repayment"),
		}
		s.txRepo.Create(tx)
	}

	// Step 2: Pay catalyst investors with remaining amount (if any)
	remainingAfterPriority := remainingAmount - priorityPaid
	if remainingAfterPriority > 0 && len(catalystInvestments) > 0 {
		for _, inv := range catalystInvestments {
			var actualReturn float64
			if remainingAfterPriority >= totalCatalystExpected {
				// Full repayment for catalyst
				actualReturn = inv.ExpectedReturn
			} else {
				// Partial repayment - distribute proportionally
				proportion := inv.ExpectedReturn / totalCatalystExpected
				actualReturn = remainingAfterPriority * proportion
			}

			s.fundingRepo.UpdateInvestmentStatus(inv.ID, models.InvestmentStatusRepaid, &actualReturn)

			// Create return transaction
			tx := &models.Transaction{
				InvoiceID: &invoiceID,
				UserID:    &inv.InvestorID,
				Type:      models.TxTypeInvestorReturn,
				Amount:    actualReturn,
				Currency:  "IDRX",
				Status:    models.TxStatusPending,
				Notes:     stringPtr("Catalyst tranche repayment"),
			}
			s.txRepo.Create(tx)
		}
	} else if len(catalystInvestments) > 0 {
		// Mark catalyst as defaulted if no remaining funds
		for _, inv := range catalystInvestments {
			zero := 0.0
			s.fundingRepo.UpdateInvestmentStatus(inv.ID, models.InvestmentStatusDefaulted, &zero)
		}
	}

	// Close pool and update invoice
	s.fundingRepo.UpdatePoolStatus(pool.ID, models.PoolStatusClosed)
	s.invoiceRepo.UpdateStatus(invoiceID, models.StatusRepaid)

	return nil
}

func stringPtr(s string) *string {
	return &s
}

// ClosePoolAndNotifyExporter closes a funding pool when deadline ends
// It disburses funds to exporter and sends payment notification to exporter's email
// containing invoice details for the importer to pay
func (s *FundingService) ClosePoolAndNotifyExporter(poolID uuid.UUID) (*models.ExporterPaymentNotificationData, error) {
	pool, err := s.fundingRepo.FindPoolByID(poolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}
	if pool.Status != models.PoolStatusOpen && pool.Status != models.PoolStatusFilled {
		return nil, errors.New("pool is not open or filled")
	}

	// Check if pool has any funding
	if pool.FundedAmount == 0 {
		return nil, errors.New("pool has no funding")
	}

	// Get invoice details
	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Get exporter details
	exporter, err := s.userRepo.FindByID(invoice.ExporterID)
	if err != nil {
		return nil, err
	}
	if exporter == nil {
		return nil, errors.New("exporter not found")
	}

	// Get exporter profile for name
	exporterProfile, _ := s.userRepo.FindProfileByUserID(invoice.ExporterID)
	exporterName := "Exporter"
	if exporterProfile != nil && exporterProfile.FullName != "" {
		exporterName = exporterProfile.FullName
	} else if exporterProfile != nil && exporterProfile.CompanyName != nil {
		exporterName = *exporterProfile.CompanyName
	}

	// Get buyer details
	buyer, err := s.buyerRepo.FindByID(invoice.BuyerID)
	if err != nil {
		return nil, err
	}
	if buyer == nil {
		return nil, errors.New("buyer not found")
	}

	buyerEmail := ""
	if buyer.ContactEmail != nil {
		buyerEmail = *buyer.ContactEmail
	}

	// Get all investments for this pool
	investments, err := s.fundingRepo.FindInvestmentsByPool(pool.ID)
	if err != nil {
		return nil, err
	}

	// Calculate total interest and build investor details
	var investorDetails []models.InvestorPaymentDetail
	totalExpectedReturn := 0.0
	for _, inv := range investments {
		totalExpectedReturn += inv.ExpectedReturn

		var interestRate float64
		if inv.Tranche == models.TranchePriority {
			interestRate = pool.PriorityInterestRate
		} else {
			interestRate = pool.CatalystInterestRate
		}

		investorDetails = append(investorDetails, models.InvestorPaymentDetail{
			InvestorID:     inv.InvestorID.String(),
			Amount:         inv.Amount,
			InterestRate:   interestRate,
			ExpectedReturn: inv.ExpectedReturn,
			Tranche:        string(inv.Tranche),
		})
	}

	totalInterest := totalExpectedReturn - pool.FundedAmount
	totalAmountDue := totalExpectedReturn

	// Generate payment ID
	paymentID := uuid.New().String()
	paymentLink := s.cfg.FrontendURL + "/payment/" + paymentID

	// Update pool status to disbursed
	if err := s.fundingRepo.UpdatePoolStatus(poolID, models.PoolStatusDisbursed); err != nil {
		return nil, err
	}

	// Update invoice status to funded
	if err := s.invoiceRepo.UpdateStatus(pool.InvoiceID, models.StatusFunded); err != nil {
		return nil, err
	}

	// Calculate platform fee (dummy escrow)
	platformFee := pool.FundedAmount * (s.cfg.PlatformFeePercentage / 100)
	disbursementAmount := pool.FundedAmount - platformFee

	// Create advance payment transaction (funds to exporter)
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &invoice.ExporterID,
		Type:      models.TxTypeAdvancePayment,
		Amount:    disbursementAmount,
		Currency:  "IDRX",
		Status:    models.TxStatusPending,
		Notes:     stringPtr("Disbursement to exporter - pool deadline reached"),
	}
	s.txRepo.Create(tx)

	// Prepare notification data
	notificationData := &models.ExporterPaymentNotificationData{
		InvoiceID:       invoice.ID.String(),
		InvoiceNumber:   invoice.InvoiceNumber,
		ExporterName:    exporterName,
		BuyerName:       buyer.CompanyName,
		BuyerEmail:      buyerEmail,
		PrincipalAmount: pool.FundedAmount,
		TotalInterest:   totalInterest,
		TotalAmountDue:  totalAmountDue,
		Currency:        pool.PoolCurrency,
		DueDate:         invoice.DueDate,
		InvestorDetails: investorDetails,
		PaymentID:       paymentID,
		PaymentLink:     paymentLink,
	}

	// Send email to exporter
	if s.emailService != nil {
		if err := s.emailService.SendExporterPaymentNotification(exporter.Email, notificationData); err != nil {
			// Log error but don't fail the operation
			// Email sending failure shouldn't block the disbursement
		}
	}

	return notificationData, nil
}

// ExporterDisbursementRequest represents request for exporter to disburse to investors
type ExporterDisbursementRequest struct {
	PoolID uuid.UUID `json:"pool_id" binding:"required"`
	Amount float64   `json:"amount" binding:"required,gt=0"` // Amount exporter is paying
}

// ExporterDisbursementResponse represents the disbursement result
type ExporterDisbursementResponse struct {
	PoolID                uuid.UUID              `json:"pool_id"`
	InvoiceID             uuid.UUID              `json:"invoice_id"`
	TotalRequired         float64                `json:"total_required"`     // Total amount required (principal + interest)
	TotalPaid             float64                `json:"total_paid"`         // Amount exporter paid
	TotalDisbursed        float64                `json:"total_disbursed"`    // Amount disbursed to investors
	PriorityDisbursed     float64                `json:"priority_disbursed"` // Amount to priority investors
	CatalystDisbursed     float64                `json:"catalyst_disbursed"` // Amount to catalyst investors
	PriorityFullyPaid     bool                   `json:"priority_fully_paid"`
	CatalystFullyPaid     bool                   `json:"catalyst_fully_paid"`
	InvestorDisbursements []InvestorDisbursement `json:"investor_disbursements"`
	Status                string                 `json:"status"`
	Message               string                 `json:"message"`
}

// InvestorDisbursement represents disbursement to a single investor
type InvestorDisbursement struct {
	InvestorID     uuid.UUID `json:"investor_id"`
	WalletAddress  string    `json:"wallet_address"`
	Tranche        string    `json:"tranche"`
	Principal      float64   `json:"principal"`
	ExpectedReturn float64   `json:"expected_return"`
	ActualReturn   float64   `json:"actual_return"`
	Status         string    `json:"status"` // full, partial, none
}

// ExporterDisbursementToInvestors handles exporter paying back investors
// Flow: Exporter pays escrow -> System calculates shares -> Escrow disburses to investors
// Priority: Priority tranche gets paid first, Catalyst gets remainder
func (s *FundingService) ExporterDisbursementToInvestors(exporterID uuid.UUID, req *ExporterDisbursementRequest) (*ExporterDisbursementResponse, error) {
	// Get pool
	pool, err := s.fundingRepo.FindPoolByID(req.PoolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}

	// Pool must be disbursed (funds already given to exporter)
	if pool.Status != models.PoolStatusDisbursed {
		return nil, errors.New("pool must be in disbursed status for exporter to repay investors")
	}

	// Get invoice
	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Verify exporter owns this invoice
	if invoice.ExporterID != exporterID {
		return nil, errors.New("unauthorized: you do not own this invoice")
	}

	// Get all investments
	investments, err := s.fundingRepo.FindInvestmentsByPool(pool.ID)
	if err != nil {
		return nil, err
	}

	// Separate by tranche
	var priorityInvestments, catalystInvestments []models.Investment
	totalPriorityRequired := 0.0
	totalCatalystRequired := 0.0

	for _, inv := range investments {
		if inv.Tranche == models.TranchePriority {
			priorityInvestments = append(priorityInvestments, inv)
			totalPriorityRequired += inv.ExpectedReturn
		} else {
			catalystInvestments = append(catalystInvestments, inv)
			totalCatalystRequired += inv.ExpectedReturn
		}
	}

	totalRequired := totalPriorityRequired + totalCatalystRequired
	amountPaid := req.Amount
	remainingFunds := amountPaid

	// Verify with escrow (dummy)
	if s.escrowService != nil {
		verified, _, err := s.escrowService.VerifyExporterDeposit(invoice.ID, amountPaid)
		if err != nil || !verified {
			return nil, errors.New("failed to verify exporter deposit in escrow")
		}
	}

	var disbursementTargets []DisbursementTarget
	var investorDisbursements []InvestorDisbursement
	priorityDisbursed := 0.0
	catalystDisbursed := 0.0
	priorityFullyPaid := true
	catalystFullyPaid := true

	// STEP 1: Pay Priority Tranche First
	for _, inv := range priorityInvestments {
		investor, _ := s.userRepo.FindByID(inv.InvestorID)
		// Note: WalletAddress removed - funds go to registered bank account
		_ = investor // investor data available for bank account lookup

		var actualReturn float64
		var status string

		if remainingFunds >= inv.ExpectedReturn {
			// Full payment
			actualReturn = inv.ExpectedReturn
			remainingFunds -= inv.ExpectedReturn
			status = "full"
		} else if remainingFunds > 0 {
			// Partial payment
			actualReturn = remainingFunds
			remainingFunds = 0
			status = "partial"
			priorityFullyPaid = false
		} else {
			// No payment
			actualReturn = 0
			status = "none"
			priorityFullyPaid = false
		}

		if actualReturn > 0 {
			disbursementTargets = append(disbursementTargets, DisbursementTarget{
				InvestorID:    inv.InvestorID,
				WalletAddress: "", // Deprecated - funds go to bank account
				Amount:        actualReturn,
				Principal:     inv.Amount,
				ReturnAmount:  actualReturn - inv.Amount,
				Tranche:       string(models.TranchePriority),
				Currency:      "IDR",
			})
			priorityDisbursed += actualReturn

			// Update investment status
			invStatus := models.InvestmentStatusRepaid
			if status != "full" {
				invStatus = models.InvestmentStatusDefaulted
			}
			s.fundingRepo.UpdateInvestmentStatus(inv.ID, invStatus, &actualReturn)

			// Create transaction
			tx := &models.Transaction{
				InvoiceID: &invoice.ID,
				UserID:    &inv.InvestorID,
				Type:      models.TxTypeInvestorReturn,
				Amount:    actualReturn,
				Currency:  "IDRX",
				Status:    models.TxStatusPending,
				Notes:     stringPtr("Priority tranche return from exporter"),
			}
			s.txRepo.Create(tx)
		}

		investorDisbursements = append(investorDisbursements, InvestorDisbursement{
			InvestorID:     inv.InvestorID,
			WalletAddress:  "", // Deprecated - funds go to bank account
			Tranche:        string(models.TranchePriority),
			Principal:      inv.Amount,
			ExpectedReturn: inv.ExpectedReturn,
			ActualReturn:   actualReturn,
			Status:         status,
		})
	}

	// STEP 2: Pay Catalyst Tranche (only if there's remaining funds)
	for _, inv := range catalystInvestments {
		investor, _ := s.userRepo.FindByID(inv.InvestorID)
		// Note: WalletAddress removed - funds go to registered bank account
		_ = investor // investor data available for bank account lookup

		var actualReturn float64
		var status string

		if remainingFunds >= inv.ExpectedReturn {
			// Full payment
			actualReturn = inv.ExpectedReturn
			remainingFunds -= inv.ExpectedReturn
			status = "full"
		} else if remainingFunds > 0 {
			// Partial payment
			actualReturn = remainingFunds
			remainingFunds = 0
			status = "partial"
			catalystFullyPaid = false
		} else {
			// No payment (catalyst gets nothing if priority didn't get full)
			actualReturn = 0
			status = "none"
			catalystFullyPaid = false
		}

		if actualReturn > 0 {
			disbursementTargets = append(disbursementTargets, DisbursementTarget{
				InvestorID:    inv.InvestorID,
				WalletAddress: "", // Deprecated - funds go to bank account
				Amount:        actualReturn,
				Principal:     inv.Amount,
				ReturnAmount:  actualReturn - inv.Amount,
				Tranche:       string(models.TrancheCatalyst),
				Currency:      "IDR",
			})
			catalystDisbursed += actualReturn

			// Update investment status
			invStatus := models.InvestmentStatusRepaid
			if status != "full" {
				invStatus = models.InvestmentStatusDefaulted
			}
			s.fundingRepo.UpdateInvestmentStatus(inv.ID, invStatus, &actualReturn)

			// Create transaction
			tx := &models.Transaction{
				InvoiceID: &invoice.ID,
				UserID:    &inv.InvestorID,
				Type:      models.TxTypeInvestorReturn,
				Amount:    actualReturn,
				Currency:  "IDRX",
				Status:    models.TxStatusPending,
				Notes:     stringPtr("Catalyst tranche return from exporter"),
			}
			s.txRepo.Create(tx)
		} else {
			// Mark as defaulted
			zero := 0.0
			s.fundingRepo.UpdateInvestmentStatus(inv.ID, models.InvestmentStatusDefaulted, &zero)
		}

		investorDisbursements = append(investorDisbursements, InvestorDisbursement{
			InvestorID:     inv.InvestorID,
			WalletAddress:  "", // Deprecated - funds go to bank account
			Tranche:        string(models.TrancheCatalyst),
			Principal:      inv.Amount,
			ExpectedReturn: inv.ExpectedReturn,
			ActualReturn:   actualReturn,
			Status:         status,
		})
	}

	// Process disbursement via escrow (dummy)
	if s.escrowService != nil && len(disbursementTargets) > 0 {
		instruction, err := s.escrowService.CreateDisbursementInstruction(pool.ID, invoice.ID, disbursementTargets)
		if err != nil {
			return nil, err
		}
		_, err = s.escrowService.ProcessDisbursement(instruction)
		if err != nil {
			return nil, err
		}
	}

	// Update pool and invoice status
	s.fundingRepo.UpdatePoolStatus(pool.ID, models.PoolStatusClosed)
	if priorityFullyPaid && catalystFullyPaid {
		s.invoiceRepo.UpdateStatus(invoice.ID, models.StatusRepaid)
	} else {
		s.invoiceRepo.UpdateStatus(invoice.ID, models.StatusDefaulted)
	}

	// Build response
	totalDisbursed := priorityDisbursed + catalystDisbursed
	status := "completed"
	message := "Disbursement completed successfully"
	if amountPaid < totalRequired {
		status = "partial"
		message = "Partial disbursement completed. Priority tranche paid first."
	}

	return &ExporterDisbursementResponse{
		PoolID:                pool.ID,
		InvoiceID:             invoice.ID,
		TotalRequired:         totalRequired,
		TotalPaid:             amountPaid,
		TotalDisbursed:        totalDisbursed,
		PriorityDisbursed:     priorityDisbursed,
		CatalystDisbursed:     catalystDisbursed,
		PriorityFullyPaid:     priorityFullyPaid,
		CatalystFullyPaid:     catalystFullyPaid,
		InvestorDisbursements: investorDisbursements,
		Status:                status,
		Message:               message,
	}, nil
}

// GetPoolDetail returns comprehensive pool details for investor decision making (Flow 6)
func (s *FundingService) GetPoolDetail(poolID uuid.UUID) (*models.PoolDetailResponse, error) {
	pool, err := s.fundingRepo.FindPoolByID(poolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}

	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}

	var buyer *models.Buyer
	if invoice.BuyerID != uuid.Nil {
		buyer, _ = s.buyerRepo.FindByID(invoice.BuyerID)
	}

	exporter, _ := s.userRepo.FindByID(invoice.ExporterID)

	// Calculate remaining amounts
	priorityRemaining := pool.PriorityTarget - pool.PriorityFunded
	catalystRemaining := pool.CatalystTarget - pool.CatalystFunded
	totalRemaining := priorityRemaining + catalystRemaining

	// Calculate progress
	priorityProgress := 0.0
	if pool.PriorityTarget > 0 {
		priorityProgress = (pool.PriorityFunded / pool.PriorityTarget) * 100
	}
	catalystProgress := 0.0
	if pool.CatalystTarget > 0 {
		catalystProgress = (pool.CatalystFunded / pool.CatalystTarget) * 100
	}
	totalProgress := 0.0
	if pool.TargetAmount > 0 {
		totalProgress = (pool.FundedAmount / pool.TargetAmount) * 100
	}

	// Calculate remaining time
	remainingTime := ""
	remainingHours := 0
	if pool.Deadline != nil {
		duration := time.Until(*pool.Deadline)
		remainingHours = int(duration.Hours())
		if remainingHours > 24 {
			remainingTime = fmt.Sprintf("%d hari", remainingHours/24)
		} else if remainingHours > 0 {
			remainingTime = fmt.Sprintf("%d jam", remainingHours)
		} else {
			remainingTime = "Berakhir"
		}
	}

	// Calculate tenor
	tenorDays := int(time.Until(invoice.DueDate).Hours() / 24)
	if tenorDays < 0 {
		tenorDays = 0
	}

	// Build buyer info
	buyerInfo := models.BuyerDetailInfo{}
	if buyer != nil {
		buyerInfo = models.BuyerDetailInfo{
			CompanyName:  buyer.CompanyName,
			Country:      buyer.Country,
			CountryFlag:  getCountryFlag(buyer.Country),
			CountryRisk:  getCountryRisk(buyer.Country),
			Industry:     "", // Buyer doesn't have industry field
			IsRepeat:     s.isRepeatBuyer(buyer.ID),
			TotalHistory: s.getBuyerHistoryCount(buyer.ID),
		}
	}

	// Build exporter info
	exporterInfo := models.ExporterDetailInfo{}
	if exporter != nil {
		companyName := ""
		if exporter.Profile != nil {
			companyName = exporter.Profile.FullName
		}
		exporterInfo = models.ExporterDetailInfo{
			CompanyName:     companyName,
			IsVerified:      exporter.IsVerified,
			CreditLimit:     0, // Would need to add to User model
			AvailableCredit: 0, // Would need to add to User model
			TotalInvoices:   s.getExporterInvoiceCount(exporter.ID),
			SuccessRate:     s.getExporterSuccessRate(exporter.ID),
		}
	}

	// Build documents info
	documents := []models.DocumentInfo{}
	if invoice.DocumentHash != nil {
		documents = append(documents, models.DocumentInfo{
			Type:       "invoice",
			TypeLabel:  "Invoice Dokumen",
			IsVerified: true,
			IPFSHash:   *invoice.DocumentHash,
		})
	}

	// Build tranche info
	priorityTranche := models.TrancheInfo{
		Type:                "priority",
		TypeDisplay:         "Prioritas",
		Description:         "Pembayaran didahulukan saat eksportir melakukan pencairan.",
		TargetAmount:        pool.PriorityTarget,
		FundedAmount:        pool.PriorityFunded,
		RemainingAmount:     priorityRemaining,
		ProgressPercent:     priorityProgress,
		InterestRate:        pool.PriorityInterestRate,
		InterestRateDisplay: fmt.Sprintf("%.1f%% p.a", pool.PriorityInterestRate),
		RiskLevel:           "Low",
		RiskLevelDisplay:    "Risiko Rendah",
		InfoBox:             "Tranche Prioritas mendapat pembayaran terlebih dahulu dari hasil pelunasan. Cocok untuk pendana yang mengutamakan keamanan.",
	}

	catalystTranche := models.TrancheInfo{
		Type:                "catalyst",
		TypeDisplay:         "Katalis",
		Description:         "Pembayaran dilakukan setelah Prioritas, dengan imbal hasil lebih tinggi.",
		TargetAmount:        pool.CatalystTarget,
		FundedAmount:        pool.CatalystFunded,
		RemainingAmount:     catalystRemaining,
		ProgressPercent:     catalystProgress,
		InterestRate:        pool.CatalystInterestRate,
		InterestRateDisplay: fmt.Sprintf("%.1f%% p.a", pool.CatalystInterestRate),
		RiskLevel:           "Medium-High",
		RiskLevelDisplay:    "Risiko Menengah-Tinggi",
		InfoBox:             "Tranche Katalis dibayar setelah Prioritas. Imbal hasil lebih tinggi, namun dalam kondisi gagal bayar berisiko tidak menerima pembayaran penuh.",
	}

	// Get grade and score (handle nil pointers)
	grade := ""
	if invoice.Grade != nil {
		grade = *invoice.Grade
	}
	gradeScore := 0
	if invoice.GradeScore != nil {
		gradeScore = *invoice.GradeScore
	}

	// Handle Description pointer
	description := ""
	if invoice.Description != nil {
		description = *invoice.Description
	}

	return &models.PoolDetailResponse{
		PoolID:          pool.ID,
		InvoiceID:       invoice.ID,
		ProjectTitle:    fmt.Sprintf("%s #%s", description, invoice.InvoiceNumber),
		InvoiceNumber:   invoice.InvoiceNumber,
		Grade:           grade,
		GradeScore:      gradeScore,
		IsInsured:       invoice.IsInsured,
		TargetAmount:    pool.TargetAmount,
		FundedAmount:    pool.FundedAmount,
		RemainingAmount: totalRemaining,
		FundingProgress: totalProgress,
		TenorDays:       tenorDays,
		TenorDisplay:    fmt.Sprintf("%d Hari", tenorDays),
		DueDate:         &invoice.DueDate,
		Deadline:        pool.Deadline,
		RemainingTime:   remainingTime,
		RemainingHours:  remainingHours,
		Status:          string(pool.Status),
		Currency:        pool.PoolCurrency,
		BuyerInfo:       buyerInfo,
		ExporterInfo:    exporterInfo,
		Documents:       documents,
		PriorityTranche: priorityTranche,
		CatalystTranche: catalystTranche,
	}, nil
}

// Helper functions
func (s *FundingService) isRepeatBuyer(buyerID uuid.UUID) bool {
	count, _ := s.invoiceRepo.CountByBuyerID(buyerID)
	return count > 0
}

func (s *FundingService) getBuyerHistoryCount(buyerID uuid.UUID) int {
	count, _ := s.invoiceRepo.CountByBuyerID(buyerID)
	return count
}

func (s *FundingService) getExporterInvoiceCount(exporterID uuid.UUID) int {
	count, _ := s.invoiceRepo.CountByExporter(exporterID)
	return count
}

func (s *FundingService) getExporterSuccessRate(exporterID uuid.UUID) float64 {
	// Placeholder - would calculate from actual repayment data
	return 100.0
}

func getCountryFlag(country string) string {
	flags := map[string]string{
		"Indonesia":      "🇮🇩",
		"United States":  "🇺🇸",
		"Germany":        "🇩🇪",
		"Japan":          "🇯🇵",
		"South Korea":    "🇰🇷",
		"Singapore":      "🇸🇬",
		"Australia":      "🇦🇺",
		"Netherlands":    "🇳🇱",
		"United Kingdom": "🇬🇧",
		"France":         "🇫🇷",
		"China":          "🇨🇳",
		"India":          "🇮🇳",
		"Malaysia":       "🇲🇾",
		"Thailand":       "🇹🇭",
		"Vietnam":        "🇻🇳",
	}
	if flag, ok := flags[country]; ok {
		return flag
	}
	return "🏳️"
}

func getCountryRisk(country string) string {
	// Tier 1 - Low Risk
	tier1 := []string{"United States", "Germany", "Japan", "South Korea", "Singapore", "Australia", "Netherlands", "United Kingdom", "France"}
	for _, c := range tier1 {
		if c == country {
			return "Low"
		}
	}
	// Tier 2 - Medium Risk
	tier2 := []string{"China", "India", "Indonesia", "Malaysia", "Thailand", "Vietnam"}
	for _, c := range tier2 {
		if c == country {
			return "Medium"
		}
	}
	// Default - High Risk
	return "High"
}

// CalculateInvestmentReturns calculates projected returns for investment calculator (Flow 6)
func (s *FundingService) CalculateInvestmentReturns(req *models.InvestmentCalculatorRequest) (*models.InvestmentCalculatorResponse, error) {
	pool, err := s.fundingRepo.FindPoolByID(req.PoolID)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}

	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}

	// Get interest rate based on tranche
	var interestRate float64
	var maxInvestable float64
	var trancheDisplay string

	if req.Tranche == string(models.TranchePriority) {
		interestRate = pool.PriorityInterestRate
		maxInvestable = pool.PriorityTarget - pool.PriorityFunded
		trancheDisplay = "Prioritas"
	} else {
		interestRate = pool.CatalystInterestRate
		maxInvestable = pool.CatalystTarget - pool.CatalystFunded
		trancheDisplay = "Katalis"
	}

	// Calculate tenor in days (for display only)
	tenorDays := int(time.Until(invoice.DueDate).Hours() / 24)
	if tenorDays < 0 {
		tenorDays = 0
	}
	if tenorDays == 0 {
		tenorDays = 60 // Default
	}

	// Flat interest calculation: Interest = Principal × (Rate/100)
	// No time factor - bunga dihitung flat dari total nominal
	interestAmount := req.Amount * (interestRate / 100)
	totalReturn := req.Amount + interestAmount
	effectiveRate := interestRate // Same as nominal rate for flat calculation

	// Platform fee (2% of interest)
	platformFee := interestAmount * 0.02
	netInterest := interestAmount - platformFee
	netTotal := req.Amount + netInterest

	return &models.InvestmentCalculatorResponse{
		PoolID:         req.PoolID,
		Tranche:        req.Tranche,
		TrancheDisplay: trancheDisplay,
		Principal:      req.Amount,
		InterestRate:   interestRate,
		TenorDays:      tenorDays,
		GrossInterest:  interestAmount,
		PlatformFee:    platformFee,
		NetInterest:    netInterest,
		TotalReturn:    totalReturn,
		NetTotalReturn: netTotal,
		EffectiveRate:  effectiveRate,
		MaxInvestable:  maxInvestable,
		CanInvest:      req.Amount <= maxInvestable && req.Amount > 0,
		Message:        fmt.Sprintf("Estimasi imbal hasil untuk %s tranche", trancheDisplay),
	}, nil
}

// ConfirmInvestment processes investment after user confirms acknowledgements (Flow 6)
func (s *FundingService) ConfirmInvestment(userID uuid.UUID, req *models.InvestConfirmationRequest) (*models.Investment, error) {
	// Validate acknowledgements for Catalyst tranche
	if req.Tranche == string(models.TrancheCatalyst) {
		if !req.CatalystWarning1 || !req.CatalystWarning2 {
			return nil, errors.New("anda harus menyetujui semua risiko tranche Katalis sebelum melanjutkan")
		}
	}

	// Create investment request
	investReq := &models.InvestRequest{
		PoolID:  req.PoolID,
		Amount:  req.Amount,
		Tranche: models.TrancheType(req.Tranche),
	}

	return s.Invest(userID, investReq)
}

// GetActiveInvestments returns paginated list of investor's active investments (Flow 10)
func (s *FundingService) GetActiveInvestments(userID uuid.UUID, page, perPage int) (*models.ActiveInvestmentListResponse, error) {
	investments, total, err := s.fundingRepo.FindActiveInvestmentsByInvestor(userID, page, perPage)
	if err != nil {
		return nil, err
	}

	activeInvestments := make([]models.InvestorActiveInvestment, 0, len(investments))

	for _, inv := range investments {
		pool, _ := s.fundingRepo.FindPoolByID(inv.PoolID)
		var invoice *models.Invoice
		var buyer *models.Buyer

		if pool != nil {
			invoice, _ = s.invoiceRepo.FindByID(pool.InvoiceID)
			if invoice != nil && invoice.BuyerID != uuid.Nil {
				buyer, _ = s.buyerRepo.FindByID(invoice.BuyerID)
			}
		}

		projectName := "Unknown"
		invoiceNumber := ""
		buyerName := ""
		buyerCountry := ""
		buyerFlag := ""
		var dueDate time.Time

		if invoice != nil {
			if invoice.Description != nil {
				projectName = *invoice.Description
			}
			invoiceNumber = invoice.InvoiceNumber
			dueDate = invoice.DueDate
		}

		if buyer != nil {
			buyerName = buyer.CompanyName
			buyerCountry = buyer.Country
			buyerFlag = getCountryFlag(buyer.Country)
		}

		// Calculate days remaining
		daysRemaining := int(time.Until(dueDate).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Determine status
		status := "lancar"
		statusDisplay := "Lancar"
		statusColor := "green"
		if daysRemaining <= 7 && daysRemaining > 0 {
			status = "perhatian"
			statusDisplay = "Perhatian"
			statusColor = "yellow"
		} else if daysRemaining == 0 || inv.Status == models.InvestmentStatusDefaulted {
			status = "gagal_bayar"
			statusDisplay = "Gagal Bayar"
			statusColor = "red"
		}

		// Get interest rate
		interestRate := 0.0
		if pool != nil {
			if inv.Tranche == models.TranchePriority {
				interestRate = pool.PriorityInterestRate
			} else {
				interestRate = pool.CatalystInterestRate
			}
		}

		trancheDisplay := "Prioritas"
		if inv.Tranche == models.TrancheCatalyst {
			trancheDisplay = "Katalis"
		}

		activeInvestments = append(activeInvestments, models.InvestorActiveInvestment{
			InvestmentID:    inv.ID,
			ProjectName:     projectName,
			InvoiceNumber:   invoiceNumber,
			BuyerName:       buyerName,
			BuyerCountry:    buyerCountry,
			BuyerFlag:       buyerFlag,
			Tranche:         string(inv.Tranche),
			TrancheDisplay:  trancheDisplay,
			Principal:       inv.Amount,
			InterestRate:    interestRate,
			EstimatedReturn: inv.ExpectedReturn,
			TotalExpected:   inv.Amount + inv.ExpectedReturn,
			DueDate:         dueDate,
			DaysRemaining:   daysRemaining,
			Status:          status,
			StatusDisplay:   statusDisplay,
			StatusColor:     statusColor,
			InvestedAt:      inv.InvestedAt,
		})
	}

	totalPages := (total + perPage - 1) / perPage

	return &models.ActiveInvestmentListResponse{
		Investments: activeInvestments,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
	}, nil
}

// GetMitraActiveInvoices returns paginated list of mitra's invoices with funding status (Flow 8)
func (s *FundingService) GetMitraActiveInvoices(userID uuid.UUID, page, perPage int) (*models.MitraInvoiceListResponse, error) {
	// Get invoices by exporter
	filter := &models.InvoiceFilter{
		Page:    page,
		PerPage: perPage,
	}
	invoices, total, err := s.invoiceRepo.FindByExporter(userID, filter)
	if err != nil {
		return nil, err
	}

	invoiceDashboards := make([]models.InvoiceDashboard, 0, len(invoices))

	for _, invoice := range invoices {
		var buyer *models.Buyer
		if invoice.BuyerID != uuid.Nil {
			buyer, _ = s.buyerRepo.FindByID(invoice.BuyerID)
		}

		buyerName := ""
		buyerCountry := ""
		if buyer != nil {
			buyerName = buyer.CompanyName
			buyerCountry = buyer.Country
		}

		// Get pool for funded amount
		pool, _ := s.fundingRepo.FindPoolByInvoiceID(invoice.ID)
		fundedAmount := 0.0
		if pool != nil {
			fundedAmount = pool.FundedAmount
		}

		// Calculate days remaining
		dueDate := invoice.DueDate
		daysRemaining := int(time.Until(dueDate).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Determine status
		status := "Aktif"
		statusColor := "green"
		if daysRemaining <= 7 && daysRemaining > 0 {
			status = "Dalam Pengawasan"
			statusColor = "yellow"
		} else if daysRemaining < 0 {
			status = "Terlambat"
			statusColor = "red"
		}

		// Calculate total owed (with interest)
		amount := invoice.Amount
		totalOwed := amount // Would calculate with interest in real implementation

		invoiceDashboards = append(invoiceDashboards, models.InvoiceDashboard{
			InvoiceID:     invoice.ID,
			InvoiceNumber: invoice.InvoiceNumber,
			BuyerName:     buyerName,
			BuyerCountry:  buyerCountry,
			DueDate:       dueDate,
			Amount:        amount,
			Status:        status,
			StatusColor:   statusColor,
			DaysRemaining: daysRemaining,
			FundedAmount:  fundedAmount,
			TotalOwed:     totalOwed,
		})
	}

	totalPages := (total + perPage - 1) / perPage

	return &models.MitraInvoiceListResponse{
		Invoices:   invoiceDashboards,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetMitraDashboard returns comprehensive dashboard data for mitra (Flow 8)
func (s *FundingService) GetMitraDashboard(userID uuid.UUID) (*models.MitraDashboard, error) {
	// Get all active invoices for this mitra
	filter := &models.InvoiceFilter{
		Page:    1,
		PerPage: 100,
	}
	invoices, _, err := s.invoiceRepo.FindByExporter(userID, filter)
	if err != nil {
		return nil, err
	}

	var totalActiveFinancing float64
	var totalOwedToInvestors float64
	var totalDaysRemaining int
	var activeCount int

	activeInvoices := make([]models.InvoiceDashboard, 0)

	for _, invoice := range invoices {
		// Skip non-active invoices
		if invoice.Status == models.StatusDraft || invoice.Status == models.StatusRepaid || invoice.Status == models.StatusRejected {
			continue
		}

		// Get pool for this invoice
		pool, _ := s.fundingRepo.FindPoolByInvoiceID(invoice.ID)
		fundedAmount := 0.0
		if pool != nil && pool.Status != models.PoolStatusClosed {
			fundedAmount = pool.FundedAmount

			// Calculate total owed (principal + interest)
			priorityOwed := pool.PriorityFunded * (1 + pool.PriorityInterestRate/100)
			catalystOwed := pool.CatalystFunded * (1 + pool.CatalystInterestRate/100)
			totalOwedToInvestors += priorityOwed + catalystOwed
		}

		// Get buyer info
		var buyer *models.Buyer
		if invoice.BuyerID != uuid.Nil {
			buyer, _ = s.buyerRepo.FindByID(invoice.BuyerID)
		}

		buyerName := ""
		buyerCountry := ""
		if buyer != nil {
			buyerName = buyer.CompanyName
			buyerCountry = buyer.Country
		}

		// Calculate days remaining
		dueDate := invoice.DueDate
		daysRemaining := int(time.Until(dueDate).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Determine status
		status := "Aktif"
		statusColor := "green"
		if daysRemaining <= 7 && daysRemaining > 0 {
			status = "Dalam Pengawasan"
			statusColor = "yellow"
		} else if daysRemaining < 0 {
			status = "Terlambat"
			statusColor = "red"
		}

		if fundedAmount > 0 {
			totalActiveFinancing += fundedAmount
			totalDaysRemaining += daysRemaining
			activeCount++
		}

		// Calculate total owed
		amount := invoice.Amount
		totalOwed := amount

		activeInvoices = append(activeInvoices, models.InvoiceDashboard{
			InvoiceID:     invoice.ID,
			InvoiceNumber: invoice.InvoiceNumber,
			BuyerName:     buyerName,
			BuyerCountry:  buyerCountry,
			DueDate:       dueDate,
			Amount:        amount,
			Status:        status,
			StatusColor:   statusColor,
			DaysRemaining: daysRemaining,
			FundedAmount:  fundedAmount,
			TotalOwed:     totalOwed,
		})
	}

	// Calculate average remaining tenor
	averageTenor := 0
	if activeCount > 0 {
		averageTenor = totalDaysRemaining / activeCount
	}

	// Determine timeline status
	timelineStatus := models.TimelineStatus{
		FundraisingComplete:  totalActiveFinancing > 0,
		DisbursementComplete: false, // Would check actual disbursement status
		RepaymentComplete:    false,
		CurrentStep:          "Menunggu pembiayaan",
	}
	if totalActiveFinancing > 0 {
		timelineStatus.CurrentStep = "Pembiayaan aktif"
	}

	return &models.MitraDashboard{
		TotalActiveFinancing:  totalActiveFinancing,
		TotalOwedToInvestors:  totalOwedToInvestors,
		AverageRemainingTenor: averageTenor,
		ActiveInvoices:        activeInvoices,
		TimelineStatus:        timelineStatus,
	}, nil
}
