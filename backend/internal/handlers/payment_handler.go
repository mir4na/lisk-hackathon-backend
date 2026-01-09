package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	txRepo         *repository.TransactionRepository
}

func NewPaymentHandler(paymentService *services.PaymentService, txRepo *repository.TransactionRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		txRepo:         txRepo,
	}
}

// Deposit godoc
// @Summary Deposit funds to user balance (PROTOTYPE)
// @Description Simulate depositing funds. This is a prototype - no real payment gateway integration.
// @Tags Payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.DepositRequest true "Deposit amount"
// @Success 200 {object} services.PaymentResponse
// @Router /payments/deposit [post]
func (h *PaymentHandler) Deposit(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req services.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.paymentService.SimulateDeposit(userID, req.Amount)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// Withdraw godoc
// @Summary Withdraw funds from user balance (PROTOTYPE)
// @Description Simulate withdrawing funds. This is a prototype - no real payment gateway integration.
// @Tags Payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.WithdrawRequest true "Withdrawal amount"
// @Success 200 {object} services.PaymentResponse
// @Router /payments/withdraw [post]
func (h *PaymentHandler) Withdraw(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req services.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.paymentService.SimulateWithdraw(userID, req.Amount)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// GetBalance godoc
// @Summary Get user balance
// @Description Get current user's balance
// @Tags Payments
// @Security BearerAuth
// @Produce json
// @Success 200 {object} services.BalanceResponse
// @Router /payments/balance [get]
func (h *PaymentHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	response, err := h.paymentService.GetBalance(userID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// AdminGrantBalance godoc
// @Summary Grant balance to user (Admin Only - MVP)
// @Description Admin can grant balance to any user. This is for MVP testing purposes.
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.AdminGrantBalanceRequest true "User ID and amount to grant"
// @Success 200 {object} services.PaymentResponse
// @Router /admin/balance/grant [post]
func (h *PaymentHandler) AdminGrantBalance(c *gin.Context) {
	var req services.AdminGrantBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		utils.BadRequestError(c, "Invalid user ID format")
		return
	}

	response, err := h.paymentService.AdminGrantBalance(targetUserID, req.Amount)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// PlatformRevenueResponse represents the platform revenue data for admin dashboard
type PlatformRevenueResponse struct {
	TotalRevenue     float64              `json:"total_revenue"`
	Currency         string               `json:"currency"`
	FeePercentage    float64              `json:"fee_percentage"`
	TransactionCount int                  `json:"transaction_count"`
	Transactions     []models.Transaction `json:"transactions,omitempty"`
	Page             int                  `json:"page"`
	PerPage          int                  `json:"per_page"`
	TotalPages       int                  `json:"total_pages"`
}

// GetPlatformRevenue godoc
// @Summary Get total platform revenue (Admin Only)
// @Description Get total platform fees collected from mitra repayments
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} PlatformRevenueResponse
// @Router /admin/platform/revenue [get]
func (h *PaymentHandler) GetPlatformRevenue(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Get total platform fees
	totalRevenue, err := h.txRepo.GetTotalPlatformFees()
	if err != nil {
		utils.InternalServerError(c, "Failed to get platform revenue")
		return
	}

	// Get platform fee transactions
	transactions, total, err := h.txRepo.GetPlatformFeeTransactions(page, perPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get platform fee transactions")
		return
	}

	totalPages := (total + perPage - 1) / perPage

	utils.SuccessResponse(c, PlatformRevenueResponse{
		TotalRevenue:     totalRevenue,
		Currency:         "IDR",
		FeePercentage:    2.0, // Platform fee percentage
		TransactionCount: total,
		Transactions:     transactions,
		Page:             page,
		PerPage:          perPage,
		TotalPages:       totalPages,
	})
}
