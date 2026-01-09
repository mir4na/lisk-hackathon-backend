package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/utils"
)

type UserHandler struct {
	userRepo *repository.UserRepository
	kycRepo  *repository.KYCRepository
}

func NewUserHandler(userRepo *repository.UserRepository, kycRepo *repository.KYCRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		kycRepo:  kycRepo,
	}
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get profile")
		return
	}
	if user == nil {
		utils.NotFoundError(c, "User not found")
		return
	}

	utils.SuccessResponse(c, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "Profile update"
// @Success 200 {object} models.User
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	if err := h.userRepo.UpdateProfile(userID, &req); err != nil {
		utils.InternalServerError(c, "Failed to update profile")
		return
	}

	user, _ := h.userRepo.FindByID(userID)
	utils.SuccessResponse(c, user)
}

// SubmitKYC godoc
// @Summary Submit KYC verification
// @Description Submit KYC/KYB verification request
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.SubmitKYCRequest true "KYC details"
// @Success 201 {object} models.KYCVerification
// @Router /user/kyc [post]
func (h *UserHandler) SubmitKYC(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.SubmitKYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Check for existing pending KYC
	existing, _ := h.kycRepo.FindByUserID(userID)
	if existing != nil && existing.Status == models.KYCStatusPending {
		utils.ConflictError(c, "You already have a pending KYC verification")
		return
	}

	kyc := &models.KYCVerification{
		UserID:           userID,
		VerificationType: req.VerificationType,
		Status:           models.KYCStatusPending,
		IDType:           &req.IDType,
		IDNumber:         &req.IDNumber,
	}

	if err := h.kycRepo.Create(kyc); err != nil {
		utils.InternalServerError(c, "Failed to submit KYC")
		return
	}

	utils.CreatedResponse(c, kyc)
}

// GetKYCStatus godoc
// @Summary Get KYC status
// @Description Get the current KYC verification status
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.KYCVerification
// @Router /user/kyc [get]
func (h *UserHandler) GetKYCStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	kyc, err := h.kycRepo.FindByUserID(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get KYC status")
		return
	}
	if kyc == nil {
		utils.NotFoundError(c, "No KYC verification found")
		return
	}

	utils.SuccessResponse(c, kyc)
}

// Admin handlers

// GetPendingKYC godoc
// @Summary Get pending KYC requests (Admin)
// @Description Get all pending KYC verification requests
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {array} models.KYCVerification
// @Router /admin/kyc/pending [get]
func (h *UserHandler) GetPendingKYC(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	kycs, total, err := h.kycRepo.FindPending(params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get pending KYC")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"kyc_requests": kycs,
		"total":        total,
		"page":         params.Page,
		"per_page":     params.PerPage,
		"total_pages":  models.CalculateTotalPages(total, params.PerPage),
	})
}

// ApproveKYC godoc
// @Summary Approve KYC request (Admin)
// @Description Approve a pending KYC verification request
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "KYC ID"
// @Success 200 {object} models.KYCVerification
// @Router /admin/kyc/{id}/approve [post]
func (h *UserHandler) ApproveKYC(c *gin.Context) {
	adminID := c.MustGet("user_id").(uuid.UUID)
	kycID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid KYC ID")
		return
	}

	kyc, err := h.kycRepo.FindByID(kycID)
	if err != nil || kyc == nil {
		utils.NotFoundError(c, "KYC request not found")
		return
	}

	if err := h.kycRepo.Approve(kycID, adminID); err != nil {
		utils.InternalServerError(c, "Failed to approve KYC")
		return
	}

	// Update user verified status
	h.userRepo.SetVerified(kyc.UserID, true)

	kyc, _ = h.kycRepo.FindByID(kycID)
	utils.SuccessResponse(c, kyc)
}

// RejectKYC godoc
// @Summary Reject KYC request (Admin)
// @Description Reject a pending KYC verification request
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "KYC ID"
// @Param request body map[string]string true "Rejection reason"
// @Success 200 {object} models.KYCVerification
// @Router /admin/kyc/{id}/reject [post]
func (h *UserHandler) RejectKYC(c *gin.Context) {
	adminID := c.MustGet("user_id").(uuid.UUID)
	kycID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid KYC ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "Reason is required")
		return
	}

	if err := h.kycRepo.Reject(kycID, adminID, req.Reason); err != nil {
		utils.InternalServerError(c, "Failed to reject KYC")
		return
	}

	kyc, _ := h.kycRepo.FindByID(kycID)
	utils.SuccessResponse(c, kyc)
}

// ==================== Profile Management (Flow 2) ====================

// GetPersonalData godoc
// @Summary Get personal data (read-only from KTP)
// @Description Get user's personal data from KYC - read only
// @Tags User Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ProfileDataResponse
// @Router /user/profile/data [get]
func (h *UserHandler) GetPersonalData(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		utils.NotFoundError(c, "User not found")
		return
	}

	identity, _ := h.userRepo.FindIdentityByUserID(userID)

	response := models.ProfileDataResponse{
		Email:        user.Email,
		MemberStatus: string(user.MemberStatus),
		Role:         string(user.Role),
		IsVerified:   user.IsVerified,
		JoinedAt:     user.CreatedAt.Format("02 January 2006"),
	}

	if user.Username != nil {
		response.Username = *user.Username
	}

	if identity != nil {
		response.FullName = identity.FullName
		response.NIKMasked = identity.MaskNIK()
	} else if user.Profile != nil {
		response.FullName = user.Profile.FullName
		response.NIKMasked = "****"
	}

	utils.SuccessResponse(c, response)
}

// GetBankAccount godoc
// @Summary Get bank account info
// @Description Get user's primary bank account for disbursement
// @Tags User Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.BankAccountResponse
// @Router /user/profile/bank-account [get]
func (h *UserHandler) GetBankAccount(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	account, err := h.userRepo.FindPrimaryBankAccount(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get bank account")
		return
	}
	if account == nil {
		utils.NotFoundError(c, "No bank account found")
		return
	}

	response := models.BankAccountResponse{
		BankCode:      account.BankCode,
		BankName:      account.BankName,
		AccountNumber: models.MaskAccountNumber(account.AccountNumber),
		AccountName:   account.AccountName,
		IsPrimary:     account.IsPrimary,
		IsVerified:    account.IsVerified,
		Microcopy:     models.BankAccountMicrocopy,
	}

	utils.SuccessResponse(c, response)
}

// ChangeBankAccount godoc
// @Summary Change bank account (requires OTP)
// @Description Change user's primary bank account - requires OTP verification for security
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ChangeBankAccountRequest true "New bank account"
// @Success 200 {object} models.BankAccountResponse
// @Router /user/profile/bank-account [put]
func (h *UserHandler) ChangeBankAccount(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.ChangeBankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Verify OTP token (abstracted for MVP - always pass)
	// In production: validate OTP via OTPService
	if req.OTPToken == "" {
		utils.BadRequestError(c, "OTP verification required for security")
		return
	}

	// Validate bank code
	var bankName string
	for _, bank := range models.GetSupportedBanks() {
		if bank.Code == req.BankCode {
			bankName = bank.Name
			break
		}
	}
	if bankName == "" {
		utils.BadRequestError(c, "Bank tidak didukung")
		return
	}

	// Create new bank account as primary
	newAccount := &models.BankAccount{
		UserID:        userID,
		BankCode:      req.BankCode,
		BankName:      bankName,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		IsVerified:    true, // Auto-verified for MVP
		IsPrimary:     true,
	}

	// Unset previous primary and create new
	if err := h.userRepo.CreateBankAccount(newAccount); err != nil {
		utils.InternalServerError(c, "Failed to update bank account")
		return
	}

	// Set as primary (this also unsets previous primary)
	if err := h.userRepo.SetPrimaryBankAccount(userID, newAccount.ID); err != nil {
		utils.InternalServerError(c, "Failed to set primary bank account")
		return
	}

	response := models.BankAccountResponse{
		BankCode:      newAccount.BankCode,
		BankName:      newAccount.BankName,
		AccountNumber: models.MaskAccountNumber(newAccount.AccountNumber),
		AccountName:   newAccount.AccountName,
		IsPrimary:     true,
		IsVerified:    true,
		Microcopy:     "Rekening berhasil diubah. Rekening ini akan digunakan untuk pencairan dana.",
	}

	utils.SuccessResponse(c, response)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user's password - requires current password verification
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]string
// @Router /user/profile/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Get current user
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		utils.NotFoundError(c, "User not found")
		return
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		utils.BadRequestError(c, "Password saat ini salah")
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalServerError(c, "Failed to process password")
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		utils.InternalServerError(c, "Failed to update password")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Password berhasil diubah",
	})
}

// GetSupportedBanks godoc
// @Summary Get list of supported banks
// @Description Get list of banks supported for disbursement
// @Tags User Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.SupportedBank
// @Router /user/profile/banks [get]
func (h *UserHandler) GetSupportedBanks(c *gin.Context) {
	banks := models.GetSupportedBanks()
	utils.SuccessResponse(c, banks)
}
