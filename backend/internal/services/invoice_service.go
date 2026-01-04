package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/config"
	"github.com/receiv3/backend/internal/models"
	"github.com/receiv3/backend/internal/repository"
)

type InvoiceService struct {
	invoiceRepo repository.InvoiceRepositoryInterface
	buyerRepo   repository.BuyerRepositoryInterface
	fundingRepo repository.FundingRepositoryInterface
	pinata      PinataServiceInterface
	cfg         *config.Config
}

func NewInvoiceService(
	invoiceRepo repository.InvoiceRepositoryInterface,
	buyerRepo repository.BuyerRepositoryInterface,
	fundingRepo repository.FundingRepositoryInterface,
	pinata PinataServiceInterface,
	cfg *config.Config,
) *InvoiceService {
	return &InvoiceService{
		invoiceRepo: invoiceRepo,
		buyerRepo:   buyerRepo,
		fundingRepo: fundingRepo,
		pinata:      pinata,
		cfg:         cfg,
	}
}

func (s *InvoiceService) Create(exporterID uuid.UUID, req *models.CreateInvoiceRequest) (*models.Invoice, error) {
	// Verify buyer exists and belongs to exporter
	buyer, err := s.buyerRepo.FindByID(req.BuyerID)
	if err != nil {
		return nil, err
	}
	if buyer == nil {
		return nil, errors.New("buyer not found")
	}
	if buyer.CreatedBy != exporterID {
		return nil, errors.New("buyer does not belong to you")
	}

	// Validate amount
	if req.Amount < s.cfg.MinInvoiceAmount || req.Amount > s.cfg.MaxInvoiceAmount {
		return nil, errors.New("invoice amount out of allowed range")
	}

	// Parse dates
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, errors.New("invalid issue date format (use YYYY-MM-DD)")
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New("invalid due date format (use YYYY-MM-DD)")
	}

	if dueDate.Before(issueDate) {
		return nil, errors.New("due date must be after issue date")
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	invoice := &models.Invoice{
		ExporterID:        exporterID,
		BuyerID:           req.BuyerID,
		InvoiceNumber:     req.InvoiceNumber,
		Currency:          currency,
		Amount:            req.Amount,
		IssueDate:         issueDate,
		DueDate:           dueDate,
		Description:       req.Description,
		Status:            models.StatusDraft,
		AdvancePercentage: s.cfg.DefaultAdvancePercentage,
	}

	if err := s.invoiceRepo.Create(invoice); err != nil {
		return nil, err
	}

	invoice.Buyer = buyer
	return invoice, nil
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
