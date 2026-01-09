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

// UpdateWallet godoc
// @Summary Update wallet address
// @Description Link a wallet address to the user account
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UpdateWalletRequest true "Wallet address"
// @Success 200 {object} models.User
// @Router /user/wallet [put]
func (h *UserHandler) UpdateWallet(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.UpdateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Check if wallet already exists
	exists, _ := h.userRepo.WalletExists(req.WalletAddress)
	if exists {
		utils.ConflictError(c, "Wallet address already linked to another account")
		return
	}

	if err := h.userRepo.UpdateWallet(userID, req.WalletAddress); err != nil {
		utils.InternalServerError(c, "Failed to update wallet")
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
