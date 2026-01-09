package handlers

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type MitraHandler struct {
	mitraService *services.MitraService
}

func NewMitraHandler(mitraService *services.MitraService) *MitraHandler {
	return &MitraHandler{mitraService: mitraService}
}

// Apply godoc
// @Summary Apply for MITRA status
// @Description Submit an application to become a MITRA (funding recipient)
// @Tags MITRA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.SubmitMitraApplicationRequest true "Application details"
// @Success 201 {object} models.MitraApplication
// @Failure 400 {object} models.APIError
// @Router /user/mitra/apply [post]
func (h *MitraHandler) Apply(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not authenticated")
		return
	}

	var req models.SubmitMitraApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	app, err := h.mitraService.Apply(userID.(uuid.UUID), &req)
	if err != nil {
		switch err {
		case services.ErrAlreadyApplied:
			utils.ConflictError(c, err.Error())
		case services.ErrAlreadyMitra:
			utils.ConflictError(c, err.Error())
		case services.ErrInvalidNPWP:
			utils.BadRequestError(c, err.Error())
		default:
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.CreatedResponse(c, app)
}

// GetStatus godoc
// @Summary Get MITRA application status
// @Description Get the current MITRA application status for the authenticated user
// @Tags MITRA
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.MitraApplicationResponse
// @Failure 404 {object} models.APIError
// @Router /user/mitra/status [get]
func (h *MitraHandler) GetStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not authenticated")
		return
	}

	response, err := h.mitraService.GetApplicationStatus(userID.(uuid.UUID))
	if err != nil {
		if err == services.ErrApplicationNotFound {
			utils.NotFoundError(c, "no MITRA application found")
		} else {
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, response)
}

// UploadDocument godoc
// @Summary Upload MITRA document
// @Description Upload a document (NIB, Akta Pendirian, or KTP Direktur) for MITRA application
// @Tags MITRA
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param document_type formData string true "Document type (nib, akta_pendirian, ktp_direktur)"
// @Param file formData file true "Document file (PDF/ZIP)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.APIError
// @Router /user/mitra/documents [post]
func (h *MitraHandler) UploadDocument(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not authenticated")
		return
	}

	// Get document type
	docType := c.PostForm("document_type")
	if docType != "nib" && docType != "akta_pendirian" && docType != "ktp_direktur" {
		utils.BadRequestError(c, "invalid document type")
		return
	}

	// Get file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequestError(c, "file not found")
		return
	}
	defer file.Close()

	// Read file data
	fileData, err := io.ReadAll(file)
	if err != nil {
		utils.BadRequestError(c, "failed to read file")
		return
	}

	// Validate file size (max 10MB)
	if len(fileData) > 10*1024*1024 {
		utils.BadRequestError(c, "file size exceeds 10MB")
		return
	}

	// Upload document
	if err := h.mitraService.UploadDocument(userID.(uuid.UUID), docType, fileData, header.Filename); err != nil {
		switch err {
		case services.ErrApplicationNotFound:
			utils.NotFoundError(c, "application not found, please apply first")
		case services.ErrApplicationNotPending:
			utils.BadRequestError(c, "application already processed, cannot upload documents")
		default:
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":       "document uploaded successfully",
		"document_type": docType,
	})
}

// GetPendingApplications godoc
// @Summary List pending MITRA applications
// @Description Get all pending MITRA applications (admin only)
// @Tags Admin MITRA
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} models.APIError
// @Router /admin/mitra/pending [get]
func (h *MitraHandler) GetPendingApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	applications, total, err := h.mitraService.GetPendingApplications(page, perPage)
	if err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"applications": applications,
		"total":        total,
		"page":         page,
		"per_page":     perPage,
	})
}

// GetApplication godoc
// @Summary Get MITRA application details
// @Description Get details of a specific MITRA application (admin only)
// @Tags Admin MITRA
// @Security BearerAuth
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.MitraApplication
// @Failure 404 {object} models.APIError
// @Router /admin/mitra/{id} [get]
func (h *MitraHandler) GetApplication(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "invalid ID")
		return
	}

	app, err := h.mitraService.GetApplicationByID(id)
	if err != nil {
		if err == services.ErrApplicationNotFound {
			utils.NotFoundError(c, "application not found")
		} else {
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, app)
}

// Approve godoc
// @Summary Approve MITRA application
// @Description Approve a pending MITRA application (admin only)
// @Tags Admin MITRA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.APIError
// @Router /admin/mitra/{id}/approve [post]
func (h *MitraHandler) Approve(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "invalid ID")
		return
	}

	if err := h.mitraService.Approve(id, adminID.(uuid.UUID)); err != nil {
		switch err {
		case services.ErrApplicationNotFound:
			utils.NotFoundError(c, "application not found")
		case services.ErrApplicationNotPending:
			utils.BadRequestError(c, "application is not in pending status")
		case services.ErrIncompleteDocuments:
			utils.BadRequestError(c, "documents are incomplete")
		default:
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "MITRA application approved successfully",
	})
}

// Reject godoc
// @Summary Reject MITRA application
// @Description Reject a pending MITRA application with reason (admin only)
// @Tags Admin MITRA
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param request body models.RejectMitraRequest true "Rejection reason"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.APIError
// @Router /admin/mitra/{id}/reject [post]
func (h *MitraHandler) Reject(c *gin.Context) {
	adminID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedError(c, "user not authenticated")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestError(c, "invalid ID")
		return
	}

	var req models.RejectMitraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, "rejection reason is required")
		return
	}

	if err := h.mitraService.Reject(id, adminID.(uuid.UUID), req.Reason); err != nil {
		switch err {
		case services.ErrApplicationNotFound:
			utils.NotFoundError(c, "application not found")
		case services.ErrApplicationNotPending:
			utils.BadRequestError(c, "application is not in pending status")
		default:
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "MITRA application rejected",
	})
}
