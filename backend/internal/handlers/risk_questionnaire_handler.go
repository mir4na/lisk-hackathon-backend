package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type RiskQuestionnaireHandler struct {
	rqService *services.RiskQuestionnaireService
}

func NewRiskQuestionnaireHandler(rqService *services.RiskQuestionnaireService) *RiskQuestionnaireHandler {
	return &RiskQuestionnaireHandler{rqService: rqService}
}

// GetQuestions godoc
// @Summary Get risk questionnaire questions
// @Description Get the list of risk assessment questions and options
// @Tags Risk Questionnaire
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /risk-questionnaire/questions [get]
func (h *RiskQuestionnaireHandler) GetQuestions(c *gin.Context) {
	questions := h.rqService.GetQuestions()
	utils.SuccessResponse(c, questions)
}

// Submit godoc
// @Summary Submit risk questionnaire
// @Description Submit answers to unlock catalyst tranche investment
// @Tags Risk Questionnaire
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.RiskQuestionnaireRequest true "Questionnaire answers"
// @Success 200 {object} models.RiskQuestionnaireResponse
// @Router /risk-questionnaire [post]
func (h *RiskQuestionnaireHandler) Submit(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.RiskQuestionnaireRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.rqService.SubmitQuestionnaire(userID, &req)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// GetStatus godoc
// @Summary Get risk questionnaire status
// @Description Get current user's risk questionnaire status and catalyst unlock status
// @Tags Risk Questionnaire
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.RiskQuestionnaireResponse
// @Router /risk-questionnaire/status [get]
func (h *RiskQuestionnaireHandler) GetStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	response, err := h.rqService.GetStatus(userID)
	if err != nil {
		utils.HandleAppError(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}
