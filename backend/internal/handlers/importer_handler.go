package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type ImporterHandler struct {
	paymentRepo    *repository.ImporterPaymentRepository
	fundingService *services.FundingService
	fundingRepo    repository.FundingRepositoryInterface
	invoiceRepo    repository.InvoiceRepositoryInterface
}

func NewImporterHandler(
	paymentRepo *repository.ImporterPaymentRepository,
	fundingService *services.FundingService,
	fundingRepo repository.FundingRepositoryInterface,
	invoiceRepo repository.InvoiceRepositoryInterface,
) *ImporterHandler {
	return &ImporterHandler{
		paymentRepo:    paymentRepo,
		fundingService: fundingService,
		fundingRepo:    fundingRepo,
		invoiceRepo:    invoiceRepo,
	}
}

// GetPaymentInfo godoc
// @Summary Get payment info (PUBLIC)
// @Description Get payment information for importer. No authentication required.
// @Tags Public
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} models.ImporterPayment
// @Router /public/payments/{payment_id} [get]
func (h *ImporterHandler) GetPaymentInfo(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("payment_id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid payment ID")
		return
	}

	payment, err := h.paymentRepo.FindByID(paymentID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get payment info")
		return
	}
	if payment == nil {
		utils.NotFoundError(c, "Payment not found")
		return
	}

	// Get associated invoice info
	invoice, _ := h.invoiceRepo.FindByID(payment.InvoiceID)
	if invoice != nil {
		payment.Invoice = invoice
	}

	// Get pool info
	pool, _ := h.fundingRepo.FindPoolByID(payment.PoolID)
	if pool != nil {
		payment.Pool = pool
	}

	utils.SuccessResponse(c, payment)
}

// Pay godoc
// @Summary Process payment from importer (PUBLIC)
// @Description Process payment from importer. No authentication required. This is for non-user importers.
// @Tags Public
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Param request body models.ImporterPaymentRequest true "Payment amount"
// @Success 200 {object} models.ImporterPaymentResponse
// @Router /public/payments/{payment_id}/pay [post]
func (h *ImporterHandler) Pay(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("payment_id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid payment ID")
		return
	}

	var req models.ImporterPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	// Get payment record
	payment, err := h.paymentRepo.FindByID(paymentID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get payment info")
		return
	}
	if payment == nil {
		utils.NotFoundError(c, "Payment not found")
		return
	}

	// Check if already paid
	if payment.PaymentStatus == models.ImporterPaymentStatusPaid {
		utils.BadRequestError(c, "This invoice has already been paid")
		return
	}

	// Validate amount
	if req.Amount < payment.AmountDue {
		utils.BadRequestError(c, "Payment amount is less than amount due")
		return
	}

	// Process repayment through funding service
	// This distributes funds to investors (priority-first)
	if err := h.fundingService.ProcessRepayment(payment.InvoiceID, req.Amount); err != nil {
		utils.HandleAppError(c, err)
		return
	}

	// Generate simulated tx hash (in production, this comes from blockchain)
	txHash := generateTxHash()

	// Update payment record
	if err := h.paymentRepo.UpdatePayment(paymentID, req.Amount, txHash); err != nil {
		utils.InternalServerError(c, "Failed to update payment record")
		return
	}

	// Get updated payment
	updatedPayment, _ := h.paymentRepo.FindByID(paymentID)

	utils.SuccessResponse(c, models.ImporterPaymentResponse{
		PaymentID:  paymentID,
		Status:     "paid",
		AmountPaid: req.Amount,
		TxHash:     &txHash,
		Message:    "Payment processed successfully. Funds have been distributed to investors.",
		PaidAt:     updatedPayment.PaidAt,
	})
}

// generateTxHash generates a simulated transaction hash for prototype
func generateTxHash() string {
	id := uuid.New()
	return "0x" + id.String()[:40] + id.String()[24:]
}
