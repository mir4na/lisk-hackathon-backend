package services

import (
	"errors"
	"fmt"

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
