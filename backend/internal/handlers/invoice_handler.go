package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type InvoiceHandler struct {
	invoiceService    *services.InvoiceService
	blockchainService *services.BlockchainService
}

func NewInvoiceHandler(invoiceService *services.InvoiceService, blockchainService *services.BlockchainService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService:    invoiceService,
		blockchainService: blockchainService,
	}
}

// CheckRepeatBuyer godoc
// @Summary Check if buyer is repeat buyer (Flow 4 Pre-condition)
// @Description Check transaction history to determine if buyer is a repeat buyer
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.RepeatBuyerCheckRequest true "Buyer company name"
// @Success 200 {object} models.RepeatBuyerCheckResponse
// @Router /invoices/check-repeat-buyer [post]
func (h *InvoiceHandler) CheckRepeatBuyer(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.RepeatBuyerCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.invoiceService.CheckRepeatBuyer(userID, req.BuyerCompanyName)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// CreateFundingRequest godoc
// @Summary Create invoice funding request (Flow 4)
// @Description Mitra creates a funding request for their invoice
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateInvoiceFundingRequest true "Funding request details"
// @Success 201 {object} models.Invoice
// @Router /invoices/funding-request [post]
func (h *InvoiceHandler) CreateFundingRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateInvoiceFundingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Validate data confirmation checkbox
	if !req.DataConfirmation {
		utils.BadRequestError(c, "Anda harus menyetujui bahwa data yang diberikan adalah benar dan asli")
		return
	}

	invoice, err := h.invoiceService.CreateFundingRequest(userID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.CreatedResponse(c, invoice)
}

// CreateInvoice godoc
// @Summary Create a new invoice
// @Description Create a new invoice for export
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateInvoiceRequest true "Invoice details"
// @Success 201 {object} models.Invoice
// @Router /invoices [post]
func (h *InvoiceHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	invoice, err := h.invoiceService.Create(userID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.CreatedResponse(c, invoice)
}

// GetInvoices godoc
// @Summary Get all invoices
// @Description Get all invoices for the current exporter
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Param status query string false "Filter by status"
// @Success 200 {object} models.InvoiceListResponse
// @Router /invoices [get]
func (h *InvoiceHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var filter models.InvoiceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		filter = models.InvoiceFilter{Page: 1, PerPage: 10}
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}

	response, err := h.invoiceService.GetByExporter(userID, &filter)
	if err != nil {
		utils.InternalServerError(c, "Failed to get invoices")
		return
	}

	utils.SuccessResponse(c, response)
}

// GetFundableInvoices godoc
// @Summary Get fundable invoices
// @Description Get all invoices available for investment
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.InvoiceListResponse
// @Router /invoices/fundable [get]
func (h *InvoiceHandler) ListFundable(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	response, err := h.invoiceService.GetFundable(params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get fundable invoices")
		return
	}

	utils.SuccessResponse(c, response)
}

// GetInvoice godoc
// @Summary Get invoice by ID
// @Description Get a specific invoice by ID
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.Invoice
// @Router /invoices/{id} [get]
func (h *InvoiceHandler) Get(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	invoice, err := h.invoiceService.GetByID(invoiceID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get invoice")
		return
	}
	if invoice == nil {
		utils.NotFoundError(c, "Invoice not found")
		return
	}

	utils.SuccessResponse(c, invoice)
}

// UpdateInvoice godoc
// @Summary Update invoice
// @Description Update an existing invoice (draft only)
// @Tags Invoices
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body models.UpdateInvoiceRequest true "Invoice update"
// @Success 200 {object} models.Invoice
// @Router /invoices/{id} [put]
func (h *InvoiceHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	var req models.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	invoice, err := h.invoiceService.Update(invoiceID, userID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, invoice)
}

// DeleteInvoice godoc
// @Summary Delete invoice
// @Description Delete an invoice (draft only)
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]string
// @Router /invoices/{id} [delete]
func (h *InvoiceHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	if err := h.invoiceService.Delete(invoiceID, userID); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Invoice deleted successfully"})
}

// SubmitInvoice godoc
// @Summary Submit invoice for review
// @Description Submit a draft invoice for review
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} map[string]string
// @Router /invoices/{id}/submit [post]
func (h *InvoiceHandler) Submit(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	if err := h.invoiceService.Submit(invoiceID, userID); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Invoice submitted for review"})
}

// UploadDocument godoc
// @Summary Upload invoice document
// @Description Upload a document for an invoice (PDF, Bill of Lading, etc.)
// @Tags Invoices
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Invoice ID"
// @Param document_type formData string true "Document type"
// @Param file formData file true "Document file"
// @Success 201 {object} models.InvoiceDocument
// @Router /invoices/{id}/documents [post]
func (h *InvoiceHandler) UploadDocument(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	docType := c.PostForm("document_type")
	if docType == "" {
		utils.BadRequestError(c, "document_type is required")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequestError(c, "File is required")
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		utils.InternalServerError(c, "Failed to read file")
		return
	}

	doc, err := h.invoiceService.UploadDocument(invoiceID, userID, models.DocumentType(docType), fileData, header.Filename)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.CreatedResponse(c, doc)
}

// GetDocuments godoc
// @Summary Get invoice documents
// @Description Get all documents for an invoice
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {array} models.InvoiceDocument
// @Router /invoices/{id}/documents [get]
func (h *InvoiceHandler) GetDocuments(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	docs, err := h.invoiceService.GetDocuments(invoiceID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get documents")
		return
	}

	utils.SuccessResponse(c, docs)
}

// TokenizeInvoice godoc
// @Summary Tokenize invoice as NFT
// @Description Mint an NFT representing the approved invoice
// @Tags Invoices
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.InvoiceNFT
// @Router /invoices/{id}/tokenize [post]
func (h *InvoiceHandler) Tokenize(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	// Get wallet address from query or use system default
	ownerAddress := c.Query("owner_address")
	if ownerAddress == "" {
		ownerAddress = "0x0000000000000000000000000000000000000000"
	}

	nft, err := h.blockchainService.TokenizeInvoice(invoiceID, ownerAddress)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, nft)
}

// Admin handlers

// ApproveInvoice godoc
// @Summary Approve invoice (Admin) - Flow 5
// @Description Approve an invoice, mint NFT, and create funding pool
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body models.AdminApproveInvoiceRequest true "Approval details with grade"
// @Success 200 {object} map[string]interface{}
// @Router /admin/invoices/{id}/approve [post]
func (h *InvoiceHandler) Approve(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	var req models.AdminApproveInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Approve with grade
	if err := h.invoiceService.ApproveWithGrade(invoiceID, &req); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	// Auto-tokenize after approval
	invoice, _ := h.invoiceService.GetByID(invoiceID)
	var nft *models.InvoiceNFT
	if h.blockchainService != nil && invoice != nil {
		// Use system/platform address for NFT ownership since users don't have wallets
		// NFT serves as on-chain proof of invoice, owned by platform
		ownerAddress := "0x0000000000000000000000000000000000000000"
		nft, _ = h.blockchainService.TokenizeInvoice(invoiceID, ownerAddress)
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Invoice approved, tokenized, and ready for pool creation",
		"invoice_id": invoiceID,
		"grade":      req.Grade,
		"nft":        nft,
	})
}

// GetGradeSuggestion godoc
// @Summary Get grade suggestion for invoice (Admin) - BE-ADM-1
// @Description Get AI-suggested grade based on risk matrix
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.AdminGradeSuggestionResponse
// @Router /admin/invoices/{id}/grade-suggestion [get]
func (h *InvoiceHandler) GetGradeSuggestion(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	suggestion, err := h.invoiceService.GetGradeSuggestion(invoiceID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, suggestion)
}

// GetInvoiceReviewData godoc
// @Summary Get invoice data for admin review (Flow 5 - Split Screen)
// @Description Get all data needed for admin to review an invoice
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.InvoiceReviewData
// @Router /admin/invoices/{id}/review [get]
func (h *InvoiceHandler) GetInvoiceReviewData(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	reviewData, err := h.invoiceService.GetInvoiceReviewData(invoiceID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, reviewData)
}

// GetPendingInvoices godoc
// @Summary Get pending invoices for admin review
// @Description Get all invoices pending admin review
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.InvoiceListResponse
// @Router /admin/invoices/pending [get]
func (h *InvoiceHandler) GetPendingInvoices(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	response, err := h.invoiceService.GetPendingInvoices(params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get pending invoices")
		return
	}

	utils.SuccessResponse(c, response)
}

// @Summary Reject invoice (Admin)
// @Description Reject an invoice pending review
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body map[string]string true "Rejection reason"
// @Success 200 {object} map[string]string
// @Router /admin/invoices/{id}/reject [post]
func (h *InvoiceHandler) Reject(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid invoice ID")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "Reason is required")
		return
	}

	if err := h.invoiceService.Reject(invoiceID, req.Reason); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Invoice rejected"})
}
