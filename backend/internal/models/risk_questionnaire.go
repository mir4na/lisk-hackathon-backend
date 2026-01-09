package models

import (
	"time"

	"github.com/google/uuid"
)

// RiskQuestionnaire stores investor's risk assessment answers for catalyst tranche unlock
type RiskQuestionnaire struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	Q1Answer         int        `json:"q1_answer"`         // 1=Emergency fund, 2=Savings, 3=Aggressive
	Q2Answer         int        `json:"q2_answer"`         // 1=Not ready, 2=Yes ready
	Q3Answer         int        `json:"q3_answer"`         // 1=Don't understand, 2=Understand
	CatalystUnlocked bool       `json:"catalyst_unlocked"` // True if user can invest in catalyst tranche
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// RiskQuestionnaireRequest is the request body for submitting risk questionnaire
type RiskQuestionnaireRequest struct {
	// Q1: "What is your investment objective?"
	// 1 = Emergency fund (Safe)
	// 2 = Long-term savings
	// 3 = Aggressive profit seeking (Speculation)
	Q1Answer int `json:"q1_answer" binding:"required,oneof=1 2 3"`

	// Q2: "If an Invoice defaults and insurance is rejected, are you prepared to lose 100% of your capital?"
	// 1 = No, I will sue the platform
	// 2 = Yes, I understand this is venture capital risk
	Q2Answer int `json:"q2_answer" binding:"required,oneof=1 2"`

	// Q3: "Do you understand that Junior Tranche funds are only paid after Senior Tranche is fully settled?"
	// 1 = I don't understand
	// 2 = I understand
	Q3Answer int `json:"q3_answer" binding:"required,oneof=1 2"`
}

// RiskQuestionnaireResponse represents the response for risk questionnaire status
type RiskQuestionnaireResponse struct {
	Completed        bool       `json:"completed"`
	CatalystUnlocked bool       `json:"catalyst_unlocked"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	Message          string     `json:"message,omitempty"`
}

// IsCatalystEligible checks if the user's answers qualify them for catalyst tranche
// Logic: Must answer "Spekulasi (3)" for Q1 AND "Siap Rugi (2)" for Q2 AND "Mengerti (2)" for Q3
func (r *RiskQuestionnaire) IsCatalystEligible() bool {
	return r.Q1Answer == 3 && r.Q2Answer == 2 && r.Q3Answer == 2
}

// Question texts for reference
var RiskQuestions = map[int]string{
	1: "What is your investment objective?",
	2: "If an Invoice defaults and insurance is rejected, are you prepared to lose 100% of your capital?",
	3: "Do you understand that Junior Tranche funds are only paid after Senior Tranche is fully settled?",
}

var RiskAnswerOptions = map[int]map[int]string{
	1: {
		1: "Emergency fund (Safe)",
		2: "Long-term savings",
		3: "Aggressive profit seeking (Speculation)",
	},
	2: {
		1: "No, I will sue the platform",
		2: "Yes, I understand this is venture capital risk",
	},
	3: {
		1: "I don't understand",
		2: "I understand",
	},
}
