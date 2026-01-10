package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/config"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

type InvoiceService struct {
	invoiceRepo repository.InvoiceRepositoryInterface
	fundingRepo repository.FundingRepositoryInterface
	userRepo    repository.UserRepositoryInterface
	mitraRepo   *repository.MitraRepository
	pinata      PinataServiceInterface
	cfg         *config.Config
}

func NewInvoiceService(
	invoiceRepo repository.InvoiceRepositoryInterface,
	fundingRepo repository.FundingRepositoryInterface,
	pinata PinataServiceInterface,
	cfg *config.Config,
) *InvoiceService {
	return &InvoiceService{
		invoiceRepo: invoiceRepo,
		fundingRepo: fundingRepo,
		pinata:      pinata,
		cfg:         cfg,
	}
}

// SetUserRepo sets the user repository (for avoiding circular dependency)
func (s *InvoiceService) SetUserRepo(userRepo repository.UserRepositoryInterface) {
	s.userRepo = userRepo
}

// SetMitraRepo sets the mitra repository (for checking mitra approval)
func (s *InvoiceService) SetMitraRepo(mitraRepo *repository.MitraRepository) {
	s.mitraRepo = mitraRepo
}

// CheckRepeatBuyer checks if buyer is a repeat buyer based on transaction history (Flow 4 Pre-condition)
func (s *InvoiceService) CheckRepeatBuyer(mitraID uuid.UUID, buyerCompanyName string) (*models.RepeatBuyerCheckResponse, error) {
	// Simplified logic since Buyer table is removed.
	// In the future, we can query unique buyer names from invoices table with status=repaid

	// Default to treating as new buyer or relying on manual check for now
	return &models.RepeatBuyerCheckResponse{
		IsRepeatBuyer:        false,
		Message:              "⚠️ Untuk kemitraan baru, maksimal pembiayaan yang dapat dicairkan adalah 60% dari nilai tagihan.",
		PreviousTransactions: 0,
		FundingLimit:         60.0,
	}, nil
}

// CreateFundingRequest creates a new invoice funding request (Flow 4)
func (s *InvoiceService) CreateFundingRequest(mitraID uuid.UUID, req *models.CreateInvoiceFundingRequest) (*models.Invoice, error) {
	// Check if Mitra is approved before allowing invoice creation
	if s.mitraRepo != nil {
		mitraApp, err := s.mitraRepo.FindByUserID(mitraID)
		if err != nil {
			return nil, errors.New("failed to verify mitra status")
		}
		if mitraApp == nil || mitraApp.Status != models.MitraStatusApproved {
			return nil, errors.New("hanya mitra yang sudah disetujui yang dapat membuat invoice")
		}
	}

	// Validate funding duration
	fundingDurationDays := req.FundingDurationDays
	if fundingDurationDays <= 0 {
		fundingDurationDays = 14 // Default 14 days
	}

	// Validate tranche ratios
	priorityRatio := req.PriorityRatio
	catalystRatio := req.CatalystRatio
	if priorityRatio <= 0 {
		priorityRatio = 80.0
	}
	if catalystRatio <= 0 {
		catalystRatio = 20.0
	}
	// Use tolerance for floating point comparison
	sum := priorityRatio + catalystRatio
	if sum < 99.9 || sum > 100.1 {
		return nil, errors.New("priority and catalyst ratios must sum to 100%")
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New("invalid due date format (use YYYY-MM-DD)")
	}

	if dueDate.Before(time.Now()) {
		return nil, errors.New("due date must be in the future")
	}

	// Check repeat buyer and get funding limit
	repeatCheck, err := s.CheckRepeatBuyer(mitraID, req.BuyerCompanyName)
	if err != nil {
		return nil, err
	}

	fundingLimitPercentage := repeatCheck.FundingLimit
	if !req.IsRepeatBuyer && repeatCheck.IsRepeatBuyer {
		// System detected repeat buyer
		fundingLimitPercentage = 100.0
	} else if req.IsRepeatBuyer && !repeatCheck.IsRepeatBuyer && req.RepeatBuyerProof == "" {
		// User claims repeat but system doesn't detect - require proof
		return nil, errors.New("please upload proof of previous transactions for repeat buyer claim")
	}

	// Calculate advance amount based on funding limit
	advanceAmount := req.IDRAmount * (fundingLimitPercentage / 100)

	// Create invoice
	invoice := &models.Invoice{
		ExporterID:        mitraID,
		BuyerName:         req.BuyerCompanyName,
		BuyerCountry:      req.BuyerCountry,
		InvoiceNumber:     req.InvoiceNumber,
		Currency:          "IDR",
		Amount:            req.IDRAmount,
		IssueDate:         time.Now(),
		DueDate:           dueDate,
		Description:       req.Description,
		Status:            models.StatusDraft,
		AdvancePercentage: fundingLimitPercentage,
		AdvanceAmount:     &advanceAmount,

		// Currency conversion fields
		OriginalCurrency: &req.OriginalCurrency,
		OriginalAmount:   &req.OriginalAmount,
		IDRAmount:        &req.IDRAmount,
		ExchangeRate:     &req.LockedExchangeRate,
		BufferRate:       s.cfg.DefaultBufferRate,

		// Tranche configuration
		PriorityRatio:        priorityRatio,
		CatalystRatio:        catalystRatio,
		PriorityInterestRate: &req.PriorityInterestRate,
		CatalystInterestRate: &req.CatalystInterestRate,

		// Repeat buyer info
		IsRepeatBuyer:          repeatCheck.IsRepeatBuyer || req.IsRepeatBuyer,
		FundingLimitPercentage: fundingLimitPercentage,

		// Funding duration
		FundingDurationDays: fundingDurationDays,
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *InvoiceService) Create(exporterID uuid.UUID, req *models.CreateInvoiceRequest) (*models.Invoice, error) {
	return nil, errors.New("legacy create method deprecated, use CreateFundingRequest")
}

func (s *InvoiceService) GetByID(id uuid.UUID) (*models.Invoice, error) {
	return s.invoiceRepo.FindByID(id)
}

func (s *InvoiceService) GetByExporter(exporterID uuid.UUID, filter *models.InvoiceFilter) (*models.InvoiceListResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}

	invoices, total, err := s.invoiceRepo.FindByExporter(exporterID, filter)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceListResponse{
		Invoices:   invoices,
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalPages: models.CalculateTotalPages(total, filter.PerPage),
	}, nil
}

func (s *InvoiceService) GetFundable(page, perPage int) (*models.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	invoices, total, err := s.invoiceRepo.FindFundable(page, perPage)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceListResponse{
		Invoices:   invoices,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: models.CalculateTotalPages(total, perPage),
	}, nil
}

func (s *InvoiceService) Update(id, exporterID uuid.UUID, req *models.UpdateInvoiceRequest) (*models.Invoice, error) {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.ExporterID != exporterID {
		return nil, errors.New("not authorized to update this invoice")
	}
	if invoice.Status != models.StatusDraft {
		return nil, errors.New("can only update draft invoices")
	}

	if req.InvoiceNumber != "" {
		invoice.InvoiceNumber = req.InvoiceNumber
	}
	if req.Currency != "" {
		invoice.Currency = req.Currency
	}
	if req.Amount > 0 {
		invoice.Amount = req.Amount
	}
	if req.IssueDate != "" {
		issueDate, err := time.Parse("2006-01-02", req.IssueDate)
		if err == nil {
			invoice.IssueDate = issueDate
		}
	}
	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			invoice.DueDate = dueDate
		}
	}
	if req.Description != nil {
		invoice.Description = req.Description
	}

	if err := s.invoiceRepo.Update(invoice); err != nil {
		return nil, err
	}

	return s.invoiceRepo.FindByID(id)
}

func (s *InvoiceService) Delete(id, exporterID uuid.UUID) error {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}
	if invoice.ExporterID != exporterID {
		return errors.New("not authorized to delete this invoice")
	}
	if invoice.Status != models.StatusDraft {
		return errors.New("can only delete draft invoices")
	}

	return s.invoiceRepo.Delete(id)
}

func (s *InvoiceService) Submit(id, exporterID uuid.UUID) error {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}
	if invoice.ExporterID != exporterID {
		return errors.New("not authorized to submit this invoice")
	}
	if invoice.Status != models.StatusDraft {
		return errors.New("can only submit draft invoices")
	}

	// Check if documents are uploaded
	docs, _ := s.invoiceRepo.FindDocumentsByInvoiceID(id)
	if len(docs) == 0 {
		return errors.New("please upload at least one document before submitting")
	}

	return s.invoiceRepo.UpdateStatus(id, models.StatusPendingReview)
}

func (s *InvoiceService) Approve(id uuid.UUID, interestRate float64) error {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}
	if invoice.Status != models.StatusPendingReview {
		return errors.New("invoice is not pending review")
	}

	advanceAmount := invoice.Amount * (invoice.AdvancePercentage / 100)

	return s.invoiceRepo.ApproveWithTransaction(id, interestRate, advanceAmount)
}

func (s *InvoiceService) Reject(id uuid.UUID, reason string) error {
	invoice, err := s.invoiceRepo.FindByID(id)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}
	if invoice.Status != models.StatusPendingReview {
		return errors.New("invoice is not pending review")
	}

	return s.invoiceRepo.UpdateStatus(id, models.StatusRejected)
}

func (s *InvoiceService) UploadDocument(invoiceID, exporterID uuid.UUID, docType models.DocumentType, fileData []byte, fileName string) (*models.InvoiceDocument, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}
	if invoice.ExporterID != exporterID {
		return nil, errors.New("not authorized")
	}
	if invoice.Status != models.StatusDraft && invoice.Status != models.StatusPendingReview {
		return nil, errors.New("cannot upload documents at this stage")
	}

	// Upload to Pinata
	metadata := map[string]string{
		"invoice_id": invoiceID.String(),
		"doc_type":   string(docType),
	}

	pinataResp, fileHash, err := s.pinata.UploadFile(fileData, fileName, metadata)
	if err != nil {
		return nil, err
	}

	doc := &models.InvoiceDocument{
		InvoiceID:    invoiceID,
		DocumentType: docType,
		FileName:     fileName,
		FileURL:      s.pinata.GetIPFSURL(pinataResp.IpfsHash),
		FileHash:     fileHash,
		FileSize:     len(fileData),
	}

	if err := s.invoiceRepo.CreateDocument(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *InvoiceService) GetDocuments(invoiceID uuid.UUID) ([]models.InvoiceDocument, error) {
	return s.invoiceRepo.FindDocumentsByInvoiceID(invoiceID)
}

func (s *InvoiceService) DeleteDocument(docID, exporterID uuid.UUID) error {
	// Would need to verify ownership through invoice
	return s.invoiceRepo.DeleteDocument(docID)
}

// GetGradeSuggestion implements BE-ADM-1 logic for grade suggestion
func (s *InvoiceService) GetGradeSuggestion(invoiceID uuid.UUID) (*models.AdminGradeSuggestionResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Calculate scores using grading logic
	score := 0
	countryRisk := "medium"

	// Get documents for completeness check
	docs, _ := s.invoiceRepo.FindDocumentsByInvoiceID(invoiceID)
	documentCount := len(docs)

	// Get exporter invoice count for history
	exporterInvoiceCount, _ := s.invoiceRepo.CountByExporter(invoice.ExporterID)

	// 1. Country Risk Score (40 points max)
	countryScore := 25 // Default medium
	countryTier := 2
	if tier, ok := CountryRiskTier[invoice.BuyerCountry]; ok {
		countryTier = tier
	}

	switch countryTier {
	case 1:
		countryScore = 40
		countryRisk = "low"
	case 2:
		countryScore = 25
		countryRisk = "medium"
	case 3:
		countryScore = 10
		countryRisk = "high"
	}
	score += countryScore

	// 2. History Score (30 points max)
	historyScore := 10 // Default for first time
	isRepeatBuyer := invoice.IsRepeatBuyer
	if isRepeatBuyer {
		historyScore = 30
	} else if exporterInvoiceCount >= 1 {
		historyScore = 20
	}
	score += historyScore

	// 3. Document Completeness Score (30 points max)
	documentScore := 5
	documentsComplete := false
	if documentCount >= 3 {
		documentScore = 30
		documentsComplete = true
	} else if documentCount >= 1 {
		documentScore = 20
	}
	score += documentScore

	// Determine grade
	grade := "C"
	if score >= 80 {
		grade = "A"
	} else if score >= 50 {
		grade = "B"
	}

	// Funding limit based on repeat buyer status
	fundingLimit := 60.0
	if isRepeatBuyer {
		fundingLimit = 100.0
	}

	return &models.AdminGradeSuggestionResponse{
		InvoiceID:         invoiceID.String(),
		SuggestedGrade:    grade,
		GradeScore:        score,
		CountryRisk:       countryRisk,
		CountryScore:      countryScore,
		HistoryScore:      historyScore,
		DocumentScore:     documentScore,
		IsRepeatBuyer:     isRepeatBuyer,
		DocumentsComplete: documentsComplete,
		FundingLimit:      fundingLimit,
	}, nil
}

// GetInvoiceReviewData gets all data needed for admin review (Flow 5 - Split Screen)
func (s *InvoiceService) GetInvoiceReviewData(invoiceID uuid.UUID) (*models.InvoiceReviewData, error) {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	// Get buyer info (From invoice fields directly)

	// Get exporter profile
	var exporterProfile *models.UserProfile
	if s.userRepo != nil {
		exporterProfile, _ = s.userRepo.FindProfileByUserID(invoice.ExporterID)
	}

	// Get documents with validation status
	docs, _ := s.invoiceRepo.FindDocumentsByInvoiceID(invoiceID)
	var docStatuses []models.DocumentValidationStatus
	for _, doc := range docs {
		docStatuses = append(docStatuses, models.DocumentValidationStatus{
			DocumentID:    doc.ID.String(),
			DocumentType:  string(doc.DocumentType),
			FileName:      doc.FileName,
			FileURL:       doc.FileURL,
			IsValid:       false, // Default, can be updated by admin
			NeedsRevision: false,
		})
	}

	// Get grade suggestion
	gradeSuggestion, _ := s.GetGradeSuggestion(invoiceID)
	if gradeSuggestion == nil {
		gradeSuggestion = &models.AdminGradeSuggestionResponse{
			SuggestedGrade: "C",
			GradeScore:     0,
		}
	}

	return &models.InvoiceReviewData{
		Invoice:         *invoice,
		Exporter:        exporterProfile,
		Documents:       docStatuses,
		GradeSuggestion: *gradeSuggestion,
	}, nil
}

// ApproveWithGrade approves an invoice with the confirmed grade (Flow 5)
func (s *InvoiceService) ApproveWithGrade(invoiceID uuid.UUID, req *models.AdminApproveInvoiceRequest) error {
	invoice, err := s.invoiceRepo.FindByID(invoiceID)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}
	if invoice.Status != models.StatusPendingReview {
		return errors.New("invoice is not pending review")
	}

	// Update grade
	gradeScore := 0
	switch req.Grade {
	case "A":
		gradeScore = 90
	case "B":
		gradeScore = 65
	case "C":
		gradeScore = 35
	}

	// Update interest rates if provided
	priorityRate := 10.0
	catalystRate := 15.0
	if invoice.PriorityInterestRate != nil {
		priorityRate = *invoice.PriorityInterestRate
	}
	if invoice.CatalystInterestRate != nil {
		catalystRate = *invoice.CatalystInterestRate
	}
	if req.PriorityInterestRate > 0 {
		priorityRate = req.PriorityInterestRate
	}
	if req.CatalystInterestRate > 0 {
		catalystRate = req.CatalystInterestRate
	}

	// Calculate advance amount
	advanceAmount := invoice.Amount * (invoice.AdvancePercentage / 100)
	if invoice.AdvanceAmount != nil {
		advanceAmount = *invoice.AdvanceAmount
	}

	// Update invoice with grade and rates
	invoice.Grade = &req.Grade
	invoice.GradeScore = &gradeScore
	invoice.PriorityInterestRate = &priorityRate
	invoice.CatalystInterestRate = &catalystRate
	invoice.AdvanceAmount = &advanceAmount
	invoice.Status = models.StatusApproved

	return s.invoiceRepo.Update(invoice)
}

// GetPendingInvoices gets all invoices pending admin review
func (s *InvoiceService) GetPendingInvoices(page, perPage int) (*models.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	status := models.StatusPendingReview
	filter := &models.InvoiceFilter{
		Status:  &status,
		Page:    page,
		PerPage: perPage,
	}

	invoices, total, err := s.invoiceRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceListResponse{
		Invoices:   invoices,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: models.CalculateTotalPages(total, perPage),
	}, nil
}

func (s *InvoiceService) GetApprovedInvoices(page, perPage int) (*models.InvoiceListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 50
	}

	statuses := []models.InvoiceStatus{models.StatusApproved, models.StatusTokenized}
	filter := &models.InvoiceFilter{
		Statuses: statuses,
		Page:     page,
		PerPage:  perPage,
	}

	invoices, total, err := s.invoiceRepo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return &models.InvoiceListResponse{
		Invoices:   invoices,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: models.CalculateTotalPages(total, perPage),
	}, nil
}

// getCountryRiskTier returns the risk tier for a country (uses grading_service.CountryRiskTier)
func getCountryRiskTier(country string) int {
	if tier, ok := CountryRiskTier[country]; ok {
		return tier
	}
	return 3 // Default to high risk
}
