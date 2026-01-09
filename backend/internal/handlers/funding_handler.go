package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type FundingHandler struct {
	fundingService *services.FundingService
}

func NewFundingHandler(fundingService *services.FundingService) *FundingHandler {
	return &FundingHandler{fundingService: fundingService}
}

// GetPools godoc
// @Summary Get open funding pools
// @Description Get all open funding pools available for investment
// @Tags Funding
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.PoolListResponse
// @Router /pools [get]
func (h *FundingHandler) ListPools(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	response, err := h.fundingService.GetOpenPools(params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get pools")
		return
	}

	utils.SuccessResponse(c, response)
}

// GetPool godoc
// @Summary Get funding pool by ID
// @Description Get a specific funding pool by ID
// @Tags Funding
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pool ID"
// @Success 200 {object} models.FundingPoolResponse
// @Router /pools/{id} [get]
func (h *FundingHandler) GetPool(c *gin.Context) {
	poolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid pool ID")
		return
	}

	response, err := h.fundingService.GetPool(poolID)
	if err != nil {
		utils.NotFoundError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// CreatePool godoc
// @Summary Create funding pool for invoice
// @Description Create a funding pool for a tokenized invoice
// @Tags Funding
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 201 {object} models.FundingPool
// @Router /invoices/{id}/pool [post]
func (h *FundingHandler) CreatePool(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	pool, err := h.fundingService.CreatePool(invoiceID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.CreatedResponse(c, pool)
}

// Invest godoc
// @Summary Invest in a funding pool
// @Description Make an investment in a funding pool
// @Tags Funding
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.InvestRequest true "Investment details"
// @Success 201 {object} models.Investment
// @Router /investments [post]
func (h *FundingHandler) Invest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.InvestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	investment, err := h.fundingService.Invest(userID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.CreatedResponse(c, investment)
}

// GetMyInvestments godoc
// @Summary Get my investments
// @Description Get all investments made by the current investor
// @Tags Funding
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.InvestmentListResponse
// @Router /investments [get]
func (h *FundingHandler) GetMyInvestments(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	response, err := h.fundingService.GetInvestmentsByInvestor(userID, params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get investments")
		return
	}

	utils.SuccessResponse(c, response)
}

// DisburseToExporter godoc
// @Summary Disburse funds to exporter (Admin)
// @Description Release funds from a filled pool to the exporter
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pool ID"
// @Success 200 {object} map[string]string
// @Router /admin/pools/{id}/disburse [post]
func (h *FundingHandler) Disburse(c *gin.Context) {
	poolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid pool ID")
		return
	}

	if err := h.fundingService.DisburseToExporter(poolID); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Funds disbursed to exporter"})
}

// ProcessRepayment godoc
// @Summary Process invoice repayment (Admin)
// @Description Process buyer repayment and distribute to investors
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body map[string]float64 true "Repayment amount"
// @Success 200 {object} map[string]string
// @Router /admin/invoices/{id}/repay [post]
func (h *FundingHandler) ProcessRepayment(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "Amount is required")
		return
	}

	if err := h.fundingService.ProcessRepayment(invoiceID, req.Amount); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Repayment processed successfully"})
}

// GetPortfolio godoc
// @Summary Get investor portfolio summary
// @Description Get a summary of investor's portfolio including total funding, gains, and allocation
// @Tags Funding
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.InvestorPortfolio
// @Router /investments/portfolio [get]
func (h *FundingHandler) GetPortfolio(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	portfolio, err := h.fundingService.GetInvestorPortfolio(userID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get portfolio")
		return
	}

	utils.SuccessResponse(c, portfolio)
}

// GetMitraDashboard godoc
// @Summary Get mitra dashboard data
// @Description Get dashboard data for mitra including active financing and invoice status
// @Tags Mitra
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.MitraDashboard
// @Router /mitra/dashboard [get]
func (h *FundingHandler) GetMitraDashboard(c *gin.Context) {
	// Return placeholder data for now
	// In a full implementation, this would query actual mitra financing data
	dashboard := models.MitraDashboard{
		TotalActiveFinancing:  0,
		AverageRemainingTenor: 0,
		ActiveInvoices:        []models.InvoiceDashboard{},
		TimelineStatus:        "Belum ada pembiayaan aktif",
	}

	utils.SuccessResponse(c, dashboard)
}

// ClosePoolAndNotify godoc
// @Summary Close funding pool and notify exporter (Admin)
// @Description Close a funding pool when deadline ends and send payment notification to exporter
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pool ID"
// @Success 200 {object} models.ExporterPaymentNotificationData
// @Router /admin/pools/{id}/close [post]
func (h *FundingHandler) ClosePoolAndNotify(c *gin.Context) {
	poolID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid pool ID")
		return
	}

	notificationData, err := h.fundingService.ClosePoolAndNotifyExporter(poolID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":           "Pool closed and payment notification sent to exporter",
		"notification_data": notificationData,
	})
}

// GetMarketplace godoc
// @Summary Get marketplace pools with filters
// @Description Get open funding pools for marketplace with grade and insured filters
// @Tags Marketplace
// @Security BearerAuth
// @Produce json
// @Param grade query string false "Filter by grade (A, B, C)"
// @Param is_insured query bool false "Filter by insured status"
// @Param min_amount query number false "Minimum pool amount"
// @Param max_amount query number false "Maximum pool amount"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.MarketplaceListResponse
// @Router /marketplace [get]
func (h *FundingHandler) GetMarketplace(c *gin.Context) {
	var filter models.MarketplaceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PerPage == 0 {
		filter.PerPage = 10
	}

	response, err := h.fundingService.GetMarketplacePools(&filter)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// ExporterDisbursement godoc
// @Summary Exporter disburse funds to investors
// @Description Exporter pays back investors via escrow. Priority tranche paid first, catalyst gets remainder.
// @Tags Exporter
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body services.ExporterDisbursementRequest true "Disbursement request"
// @Success 200 {object} services.ExporterDisbursementResponse
// @Router /exporter/disbursement [post]
func (h *FundingHandler) ExporterDisbursement(c *gin.Context) {
	exporterID := c.MustGet("user_id").(uuid.UUID)

	var req services.ExporterDisbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.fundingService.ExporterDisbursementToInvestors(exporterID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}
