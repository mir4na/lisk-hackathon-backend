package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/models"
	"github.com/receiv3/backend/internal/repository"
	"github.com/receiv3/backend/internal/utils"
)

type BuyerHandler struct {
	buyerRepo *repository.BuyerRepository
}

func NewBuyerHandler(buyerRepo *repository.BuyerRepository) *BuyerHandler {
	return &BuyerHandler{buyerRepo: buyerRepo}
}

// CreateBuyer godoc
// @Summary Create a new buyer
// @Description Create a new overseas buyer record
// @Tags Buyers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateBuyerRequest true "Buyer details"
// @Success 201 {object} models.Buyer
// @Router /buyers [post]
func (h *BuyerHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateBuyerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	buyer := &models.Buyer{
		CreatedBy:    userID,
		CompanyName:  req.CompanyName,
		Country:      req.Country,
		Address:      req.Address,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Website:      req.Website,
	}

	if err := h.buyerRepo.Create(buyer); err != nil {
		utils.InternalServerError(c, "Failed to create buyer")
		return
	}

	utils.CreatedResponse(c, buyer)
}

// GetBuyers godoc
// @Summary Get all buyers
// @Description Get all buyers created by the current exporter
// @Tags Buyers
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} models.BuyerListResponse
// @Router /buyers [get]
func (h *BuyerHandler) List(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = models.PaginationParams{Page: 1, PerPage: 10}
	}
	params.Normalize()

	buyers, total, err := h.buyerRepo.FindByExporter(userID, params.Page, params.PerPage)
	if err != nil {
		utils.InternalServerError(c, "Failed to get buyers")
		return
	}

	utils.SuccessResponse(c, models.BuyerListResponse{
		Buyers:     buyers,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: models.CalculateTotalPages(total, params.PerPage),
	})
}

// GetBuyer godoc
// @Summary Get buyer by ID
// @Description Get a specific buyer by ID
// @Tags Buyers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Buyer ID"
// @Success 200 {object} models.Buyer
// @Router /buyers/{id} [get]
func (h *BuyerHandler) Get(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	buyerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid buyer ID")
		return
	}

	buyer, err := h.buyerRepo.FindByID(buyerID)
	if err != nil {
		utils.InternalServerError(c, "Failed to get buyer")
		return
	}
	if buyer == nil {
		utils.NotFoundError(c, "Buyer not found")
		return
	}
	if buyer.CreatedBy != userID {
		utils.ForbiddenError(c, "Not authorized to view this buyer")
		return
	}

	utils.SuccessResponse(c, buyer)
}

// UpdateBuyer godoc
// @Summary Update buyer
// @Description Update an existing buyer
// @Tags Buyers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Buyer ID"
// @Param request body models.UpdateBuyerRequest true "Buyer update"
// @Success 200 {object} models.Buyer
// @Router /buyers/{id} [put]
func (h *BuyerHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	buyerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid buyer ID")
		return
	}

	buyer, err := h.buyerRepo.FindByID(buyerID)
	if err != nil || buyer == nil {
		utils.NotFoundError(c, "Buyer not found")
		return
	}
	if buyer.CreatedBy != userID {
		utils.ForbiddenError(c, "Not authorized to update this buyer")
		return
	}

	var req models.UpdateBuyerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	if req.CompanyName != "" {
		buyer.CompanyName = req.CompanyName
	}
	if req.Country != "" {
		buyer.Country = req.Country
	}
	buyer.Address = req.Address
	buyer.ContactEmail = req.ContactEmail
	buyer.ContactPhone = req.ContactPhone
	buyer.Website = req.Website

	if err := h.buyerRepo.Update(buyer); err != nil {
		utils.InternalServerError(c, "Failed to update buyer")
		return
	}

	utils.SuccessResponse(c, buyer)
}

// DeleteBuyer godoc
// @Summary Delete buyer
// @Description Delete an existing buyer
// @Tags Buyers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Buyer ID"
// @Success 200 {object} map[string]string
// @Router /buyers/{id} [delete]
func (h *BuyerHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	buyerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "Invalid buyer ID")
		return
	}

	buyer, err := h.buyerRepo.FindByID(buyerID)
	if err != nil || buyer == nil {
		utils.NotFoundError(c, "Buyer not found")
		return
	}
	if buyer.CreatedBy != userID {
		utils.ForbiddenError(c, "Not authorized to delete this buyer")
		return
	}

	if err := h.buyerRepo.Delete(buyerID); err != nil {
		utils.InternalServerError(c, "Failed to delete buyer")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Buyer deleted successfully"})
}
