package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type CurrencyHandler struct {
	currencyService *services.CurrencyService
}

func NewCurrencyHandler(currencyService *services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{currencyService: currencyService}
}

// GetLockedExchangeRate godoc
// @Summary Get locked exchange rate with buffer (BE-4)
// @Description Get locked exchange rate for currency conversion. Applies buffer to protect disbursement value.
// @Tags Currency
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CurrencyConversionRequest true "Currency conversion request"
// @Success 200 {object} models.CurrencyConversionResponse
// @Failure 400 {object} models.APIError
// @Router /currency/convert [post]
func (h *CurrencyHandler) GetLockedExchangeRate(c *gin.Context) {
	var req models.CurrencyConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.currencyService.GetLockedExchangeRate(&req)
	if err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// GetSupportedCurrencies godoc
// @Summary Get list of supported currencies
// @Description Get all supported currencies with their current exchange rates
// @Tags Currency
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.SupportedCurrency
// @Router /currency/supported [get]
func (h *CurrencyHandler) GetSupportedCurrencies(c *gin.Context) {
	currencies := h.currencyService.GetSupportedCurrencies()
	utils.SuccessResponse(c, currencies)
}

// CalculateEstimatedDisbursement godoc
// @Summary Calculate estimated net disbursement
// @Description Calculate net disbursement after platform fee deduction
// @Tags Currency
// @Security BearerAuth
// @Produce json
// @Param amount query float64 true "Amount in IDRX"
// @Success 200 {object} models.EstimatedDisbursement
// @Router /currency/disbursement-estimate [get]
func (h *CurrencyHandler) CalculateEstimatedDisbursement(c *gin.Context) {
	var params struct {
		Amount float64 `form:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.BadRequestError(c, "Amount is required and must be positive")
		return
	}

	result := h.currencyService.CalculateEstimatedDisbursement(params.Amount)
	utils.SuccessResponse(c, result)
}
