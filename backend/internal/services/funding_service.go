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
		PoolCurrency:         "IDRX",
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

	// Check catalyst unlock for catalyst tranche
	if req.Tranche == models.TrancheCatalyst {
		if s.rqRepo != nil {
			unlocked, err := s.rqRepo.IsCatalystUnlocked(investorID)
			if err != nil {
				return nil, errors.New("failed to check catalyst eligibility")
			}
			if !unlocked {
				return nil, errors.New("catalyst tranche not unlocked. Please complete risk questionnaire first")
			}
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

	// Get invoice for tenor calculation
	invoice, err := s.invoiceRepo.FindByID(pool.InvoiceID)
	if err != nil {
		return nil, err
	}

	// Calculate days until maturity
	daysToMaturity := time.Until(invoice.DueDate).Hours() / 24
	if daysToMaturity < 0 {
		daysToMaturity = 0
	}

	// Pro-rata interest calculation: Principal * (1 + rate/100 * days/360)
	expectedReturn := req.Amount * (1 + (interestRate/100)*(daysToMaturity/360))

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

	// Create transaction record
	tx := &models.Transaction{
		InvoiceID: &pool.InvoiceID,
		UserID:    &investorID,
		Type:      models.TxTypeInvestment,
		Amount:    req.Amount,
		Currency:  "IDRX",
		Status:    models.TxStatusPending,
	}
	s.txRepo.Create(tx)

	// Check if pool is now filled (both tranches)
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
	PoolID               uuid.UUID                  `json:"pool_id"`
	InvoiceID            uuid.UUID                  `json:"invoice_id"`
	TotalRequired        float64                    `json:"total_required"`         // Total amount required (principal + interest)
	TotalPaid            float64                    `json:"total_paid"`             // Amount exporter paid
	TotalDisbursed       float64                    `json:"total_disbursed"`        // Amount disbursed to investors
	PriorityDisbursed    float64                    `json:"priority_disbursed"`     // Amount to priority investors
	CatalystDisbursed    float64                    `json:"catalyst_disbursed"`     // Amount to catalyst investors
	PriorityFullyPaid    bool                       `json:"priority_fully_paid"`
	CatalystFullyPaid    bool                       `json:"catalyst_fully_paid"`
	InvestorDisbursements []InvestorDisbursement    `json:"investor_disbursements"`
	Status               string                     `json:"status"`
	Message              string                     `json:"message"`
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
		walletAddress := ""
		if investor != nil && investor.WalletAddress != nil {
			walletAddress = *investor.WalletAddress
		}

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
				WalletAddress: walletAddress,
				Amount:        actualReturn,
				Principal:     inv.Amount,
				ReturnAmount:  actualReturn - inv.Amount,
				Tranche:       string(models.TranchePriority),
				Currency:      "IDRX",
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
			WalletAddress:  walletAddress,
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
		walletAddress := ""
		if investor != nil && investor.WalletAddress != nil {
			walletAddress = *investor.WalletAddress
		}

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
				WalletAddress: walletAddress,
				Amount:        actualReturn,
				Principal:     inv.Amount,
				ReturnAmount:  actualReturn - inv.Amount,
				Tranche:       string(models.TrancheCatalyst),
				Currency:      "IDRX",
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
			WalletAddress:  walletAddress,
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
