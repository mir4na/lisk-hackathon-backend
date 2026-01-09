package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

var (
	ErrAlreadyApplied        = errors.New("you already have a pending MITRA application")
	ErrAlreadyMitra          = errors.New("you are already a MITRA member")
	ErrApplicationNotFound   = errors.New("application not found")
	ErrInvalidNPWP           = errors.New("invalid NPWP format (must be 15 or 16 digits)")
	ErrIncompleteDocuments   = errors.New("documents are incomplete")
	ErrApplicationNotPending = errors.New("application is not in pending status")
)

type MitraService struct {
	mitraRepo     *repository.MitraRepository
	userRepo      repository.UserRepositoryInterface
	emailService  *EmailService
	pinataService *PinataService
}

func NewMitraService(
	mitraRepo *repository.MitraRepository,
	userRepo repository.UserRepositoryInterface,
	emailService *EmailService,
	pinataService *PinataService,
) *MitraService {
	return &MitraService{
		mitraRepo:     mitraRepo,
		userRepo:      userRepo,
		emailService:  emailService,
		pinataService: pinataService,
	}
}

// Apply submits a new MITRA application
func (s *MitraService) Apply(userID uuid.UUID, req *models.SubmitMitraApplicationRequest) (*models.MitraApplication, error) {
	// Check user exists and is not already MITRA
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if already a MITRA
	if user.MemberStatus == models.MemberStatusMemberMitra {
		return nil, ErrAlreadyMitra
	}

	// Check if has pending application
	hasPending, err := s.mitraRepo.HasPendingApplication(userID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, ErrAlreadyApplied
	}

	// Delete any rejected applications to allow reapply
	if err := s.mitraRepo.DeleteRejected(userID); err != nil {
		return nil, err
	}

	// Validate NPWP
	if !models.ValidateNPWP(req.NPWP) {
		return nil, ErrInvalidNPWP
	}

	// Create application
	app := &models.MitraApplication{
		UserID:        userID,
		CompanyName:   req.CompanyName,
		CompanyType:   req.CompanyType,
		NPWP:          req.NPWP,
		AnnualRevenue: req.AnnualRevenue,
		Status:        models.MitraStatusPending,
	}

	if err := s.mitraRepo.Create(app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

// GetApplicationStatus gets the current MITRA application status for a user
func (s *MitraService) GetApplicationStatus(userID uuid.UUID) (*models.MitraApplicationResponse, error) {
	app, err := s.mitraRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrApplicationNotFound
	}

	response := &models.MitraApplicationResponse{
		Application: app,
	}

	// Check document status
	response.DocumentsStatus.NIB = app.NIBDocumentURL != nil && *app.NIBDocumentURL != ""
	response.DocumentsStatus.AktaPendirian = app.AktaPendirianURL != nil && *app.AktaPendirianURL != ""
	response.DocumentsStatus.KTPDirektur = app.KTPDirekturURL != nil && *app.KTPDirekturURL != ""

	// Check if all documents are complete
	response.IsComplete = response.DocumentsStatus.NIB &&
		response.DocumentsStatus.AktaPendirian &&
		response.DocumentsStatus.KTPDirektur

	return response, nil
}

// UploadDocument uploads a document for MITRA application
func (s *MitraService) UploadDocument(userID uuid.UUID, docType string, fileData []byte, fileName string) error {
	// Get application
	app, err := s.mitraRepo.FindByUserID(userID)
	if err != nil {
		return err
	}
	if app == nil {
		return ErrApplicationNotFound
	}

	// Check if application is still pending
	if app.Status != models.MitraStatusPending {
		return ErrApplicationNotPending
	}

	// Upload to IPFS via Pinata
	metadata := map[string]string{
		"user_id":  userID.String(),
		"doc_type": docType,
	}
	pinataResp, _, err := s.pinataService.UploadFile(fileData, fileName, metadata)
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}

	url := s.pinataService.GetIPFSURL(pinataResp.IpfsHash)

	// Update document URL
	if err := s.mitraRepo.UpdateDocumentURL(app.ID, docType, url); err != nil {
		return fmt.Errorf("failed to save document: %w", err)
	}

	return nil
}

// GetPendingApplications gets all pending MITRA applications (admin)
func (s *MitraService) GetPendingApplications(page, perPage int) ([]models.MitraApplication, int, error) {
	return s.mitraRepo.FindPending(page, perPage)
}

// GetApplicationByID gets a MITRA application by ID (admin)
func (s *MitraService) GetApplicationByID(id uuid.UUID) (*models.MitraApplication, error) {
	app, err := s.mitraRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, ErrApplicationNotFound
	}

	// Load user info
	user, err := s.userRepo.FindByID(app.UserID)
	if err == nil && user != nil {
		app.User = user
	}

	return app, nil
}

// Approve approves a MITRA application
func (s *MitraService) Approve(applicationID, adminID uuid.UUID) error {
	// Get application
	app, err := s.mitraRepo.FindByID(applicationID)
	if err != nil {
		return err
	}
	if app == nil {
		return ErrApplicationNotFound
	}

	// Check if pending
	if app.Status != models.MitraStatusPending {
		return ErrApplicationNotPending
	}

	// Check documents are complete
	if app.NIBDocumentURL == nil || app.AktaPendirianURL == nil || app.KTPDirekturURL == nil {
		return ErrIncompleteDocuments
	}

	// Approve application
	if err := s.mitraRepo.Approve(applicationID, adminID); err != nil {
		return fmt.Errorf("failed to approve application: %w", err)
	}

	// Update user status to MEMBER_MITRA
	if err := s.userRepo.UpdateMemberStatus(app.UserID, models.MemberStatusMemberMitra); err != nil {
		return fmt.Errorf("failed to update member status: %w", err)
	}

	// Update user role to mitra
	if err := s.userRepo.UpdateRole(app.UserID, models.RoleMitra); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Send approval email
	user, _ := s.userRepo.FindByID(app.UserID)
	if user != nil && s.emailService != nil {
		go s.emailService.SendMitraApprovalEmail(user.Email, app.CompanyName)
	}

	return nil
}

// Reject rejects a MITRA application
func (s *MitraService) Reject(applicationID, adminID uuid.UUID, reason string) error {
	// Get application
	app, err := s.mitraRepo.FindByID(applicationID)
	if err != nil {
		return err
	}
	if app == nil {
		return ErrApplicationNotFound
	}

	// Check if pending
	if app.Status != models.MitraStatusPending {
		return ErrApplicationNotPending
	}

	// Reject application
	if err := s.mitraRepo.Reject(applicationID, adminID, reason); err != nil {
		return fmt.Errorf("failed to reject application: %w", err)
	}

	// Send rejection email
	user, _ := s.userRepo.FindByID(app.UserID)
	if user != nil && s.emailService != nil {
		go s.emailService.SendMitraRejectionEmail(user.Email, app.CompanyName, reason)
	}

	return nil
}

// ==================== Mitra Repayment (VA Payment) ====================

// GetActiveInvoicesForRepayment gets all active invoices that need repayment
func (s *MitraService) GetActiveInvoicesForRepayment(userID uuid.UUID) ([]models.MitraActiveInvoice, error) {
	// This would query invoices where pool is filled/closed and needs repayment
	// For MVP, we'll return data from the funding pools
	// In production: s.invoiceRepo.FindActiveByExporter(userID)

	// Placeholder - return empty for now, actual implementation depends on invoice repository
	return []models.MitraActiveInvoice{}, nil
}

// GetRepaymentBreakdown calculates the breakdown of repayment by tranche
func (s *MitraService) GetRepaymentBreakdown(userID, poolID uuid.UUID) (*models.MitraRepaymentBreakdown, error) {
	// Get pool and verify ownership
	// For MVP: abstracted to return sample breakdown

	breakdown := &models.MitraRepaymentBreakdown{
		PoolID:               poolID,
		PriorityInterestRate: 10.0,
		CatalystInterestRate: 15.0,
		Currency:             "IDR",
	}

	// Calculate sample values (in production, get from actual investments)
	// Priority: principal 40M, interest 4M (10%), total 44M
	// Catalyst: principal 10M, interest 1.5M (15%), total 11.5M
	breakdown.PriorityPrincipal = 40000000
	breakdown.CatalystPrincipal = 10000000
	breakdown.TotalPrincipal = breakdown.PriorityPrincipal + breakdown.CatalystPrincipal

	// Flat interest: principal × rate / 100
	breakdown.PriorityInterest = breakdown.PriorityPrincipal * breakdown.PriorityInterestRate / 100
	breakdown.CatalystInterest = breakdown.CatalystPrincipal * breakdown.CatalystInterestRate / 100
	breakdown.TotalInterest = breakdown.PriorityInterest + breakdown.CatalystInterest

	breakdown.PriorityTotal = breakdown.PriorityPrincipal + breakdown.PriorityInterest
	breakdown.CatalystTotal = breakdown.CatalystPrincipal + breakdown.CatalystInterest
	breakdown.GrandTotal = breakdown.TotalPrincipal + breakdown.TotalInterest

	return breakdown, nil
}

// CreateVAPayment creates a Virtual Account for mitra to pay
func (s *MitraService) CreateVAPayment(userID uuid.UUID, req *models.CreateVARequest) (*models.VAPaymentResponse, error) {
	// Validate bank code
	var bankName string
	for _, method := range models.GetVAPaymentMethods() {
		if method.BankCode == req.BankCode {
			bankName = method.BankName
			break
		}
	}
	if bankName == "" {
		return nil, errors.New("metode pembayaran tidak didukung")
	}

	// Get repayment breakdown
	breakdown, err := s.GetRepaymentBreakdown(userID, req.PoolID)
	if err != nil {
		return nil, err
	}

	// Generate VA number (abstracted for MVP)
	vaNumber := fmt.Sprintf("8%s%d", req.BankCode[:3], rand.Intn(9999999999))

	// Create VA with 24-hour expiry
	expiresAt := time.Now().Add(24 * time.Hour)
	va := &models.VirtualAccount{
		ID:        uuid.New(),
		PoolID:    req.PoolID,
		UserID:    userID,
		VANumber:  vaNumber,
		BankCode:  req.BankCode,
		BankName:  bankName,
		Amount:    breakdown.GrandTotal,
		Status:    models.VAStatusPending,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// In production: store VA to database
	// s.vaRepo.Create(va)

	// Calculate remaining time
	remainingDuration := time.Until(expiresAt)
	hours := int(remainingDuration.Hours())
	minutes := int(remainingDuration.Minutes()) % 60
	seconds := int(remainingDuration.Seconds()) % 60
	remainingTime := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return &models.VAPaymentResponse{
		VA:             *va,
		Breakdown:      *breakdown,
		RemainingTime:  remainingTime,
		RemainingHours: hours,
		Microcopy:      "Selesaikan pembayaran dalam waktu 24 jam. VA akan otomatis kadaluarsa setelah batas waktu.",
	}, nil
}

// GetVAPaymentStatus gets VA payment details for payment page
func (s *MitraService) GetVAPaymentStatus(userID, vaID uuid.UUID) (*models.VAPaymentPageResponse, error) {
	// In production: fetch from database
	// va, err := s.vaRepo.FindByID(vaID)

	// For MVP: return sample data
	expiresAt := time.Now().Add(23 * time.Hour)
	remainingDuration := time.Until(expiresAt)
	hours := int(remainingDuration.Hours())
	minutes := int(remainingDuration.Minutes()) % 60
	seconds := int(remainingDuration.Seconds()) % 60

	breakdown, _ := s.GetRepaymentBreakdown(userID, uuid.New())

	return &models.VAPaymentPageResponse{
		VANumber:        "88001234567890",
		BankCode:        "bca",
		BankName:        "Bank Central Asia (BCA)",
		Amount:          breakdown.GrandTotal,
		AmountFormatted: fmt.Sprintf("Rp %.0f", breakdown.GrandTotal),
		Status:          models.VAStatusPending,
		ExpiresAt:       expiresAt,
		RemainingTime:   fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds),
		Breakdown:       *breakdown,
		Microcopy:       "Nominal pembayaran bersifat tetap dan tidak dapat diubah.",
	}, nil
}

// SimulateVAPayment simulates VA payment for MVP testing
func (s *MitraService) SimulateVAPayment(userID, vaID uuid.UUID) (map[string]interface{}, error) {
	// In production: this would be a webhook from payment gateway
	// For MVP: directly trigger the disbursement flow

	// Get VA details and update status to paid
	// Then trigger ExporterDisbursementToInvestors

	return map[string]interface{}{
		"status":  "paid",
		"message": "Pembayaran berhasil! Dana sedang didistribusikan ke investor.",
		"va_id":   vaID,
		"paid_at": time.Now(),
	}, nil
}
