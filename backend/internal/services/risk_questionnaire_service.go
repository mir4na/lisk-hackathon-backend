package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
)

type RiskQuestionnaireService struct {
	rqRepo repository.RiskQuestionnaireRepositoryInterface
}

func NewRiskQuestionnaireService(rqRepo repository.RiskQuestionnaireRepositoryInterface) *RiskQuestionnaireService {
	return &RiskQuestionnaireService{rqRepo: rqRepo}
}

// SubmitQuestionnaire submits or updates risk questionnaire answers
func (s *RiskQuestionnaireService) SubmitQuestionnaire(userID uuid.UUID, req *models.RiskQuestionnaireRequest) (*models.RiskQuestionnaireResponse, error) {
	// Validate answers
	if req.Q1Answer < 1 || req.Q1Answer > 3 {
		return nil, errors.New("invalid answer for question 1")
	}
	if req.Q2Answer < 1 || req.Q2Answer > 2 {
		return nil, errors.New("invalid answer for question 2")
	}
	if req.Q3Answer < 1 || req.Q3Answer > 2 {
		return nil, errors.New("invalid answer for question 3")
	}

	// Check if questionnaire already exists
	existing, err := s.rqRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	rq := &models.RiskQuestionnaire{
		UserID:      userID,
		Q1Answer:    req.Q1Answer,
		Q2Answer:    req.Q2Answer,
		Q3Answer:    req.Q3Answer,
		CompletedAt: &now,
	}

	// Check if catalyst is unlocked based on answers
	// Logic: Must answer "Spekulasi (3)" for Q1 AND "Siap Rugi (2)" for Q2 AND "Mengerti (2)" for Q3
	rq.CatalystUnlocked = rq.IsCatalystEligible()

	if existing != nil {
		// Update existing
		rq.ID = existing.ID
		if err := s.rqRepo.Update(rq); err != nil {
			return nil, err
		}
	} else {
		// Create new
		if err := s.rqRepo.Create(rq); err != nil {
			return nil, err
		}
	}

	message := "Questionnaire saved successfully."
	if rq.CatalystUnlocked {
		message = "Congratulations! You are eligible to invest in the Catalyst Tranche (Junior)."
	} else {
		message = "You can invest in the Priority Tranche (Senior). Catalyst Tranche is not available for your risk profile."
	}

	return &models.RiskQuestionnaireResponse{
		Completed:        true,
		CatalystUnlocked: rq.CatalystUnlocked,
		CompletedAt:      rq.CompletedAt,
		Message:          message,
	}, nil
}

// GetStatus gets the current questionnaire status for a user
func (s *RiskQuestionnaireService) GetStatus(userID uuid.UUID) (*models.RiskQuestionnaireResponse, error) {
	rq, err := s.rqRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if rq == nil {
		return &models.RiskQuestionnaireResponse{
			Completed:        false,
			CatalystUnlocked: false,
			Message:          "You have not completed the risk questionnaire. Please fill it out to unlock investment features.",
		}, nil
	}

	message := "Priority Tranche is available."
	if rq.CatalystUnlocked {
		message = "Priority and Catalyst Tranche are available for you."
	}

	return &models.RiskQuestionnaireResponse{
		Completed:        true,
		CatalystUnlocked: rq.CatalystUnlocked,
		CompletedAt:      rq.CompletedAt,
		Message:          message,
	}, nil
}

// IsCatalystUnlocked checks if user has catalyst tranche unlocked
func (s *RiskQuestionnaireService) IsCatalystUnlocked(userID uuid.UUID) (bool, error) {
	return s.rqRepo.IsCatalystUnlocked(userID)
}

// GetQuestions returns the list of questions and options
func (s *RiskQuestionnaireService) GetQuestions() map[string]interface{} {
	return map[string]interface{}{
		"questions": []map[string]interface{}{
			{
				"id":       1,
				"question": models.RiskQuestions[1],
				"options": []map[string]interface{}{
					{"value": 1, "label": models.RiskAnswerOptions[1][1]},
					{"value": 2, "label": models.RiskAnswerOptions[1][2]},
					{"value": 3, "label": models.RiskAnswerOptions[1][3]},
				},
			},
			{
				"id":       2,
				"question": models.RiskQuestions[2],
				"options": []map[string]interface{}{
					{"value": 1, "label": models.RiskAnswerOptions[2][1]},
					{"value": 2, "label": models.RiskAnswerOptions[2][2]},
				},
				"required_for_catalyst": true,
				"required_answer":       2,
			},
			{
				"id":       3,
				"question": models.RiskQuestions[3],
				"options": []map[string]interface{}{
					{"value": 1, "label": models.RiskAnswerOptions[3][1]},
					{"value": 2, "label": models.RiskAnswerOptions[3][2]},
				},
				"required_for_catalyst": true,
				"required_answer":       2,
			},
		},
		"catalyst_unlock_rules": map[string]interface{}{
			"description": "To unlock Catalyst Tranche, you must answer: Q1=Speculation (3), Q2=Yes ready (2), Q3=Understand (2)",
			"required_q1": 3,
			"required_q2": 2,
			"required_q3": 2,
		},
	}
}
