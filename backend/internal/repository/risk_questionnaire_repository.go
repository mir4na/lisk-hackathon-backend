package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

type RiskQuestionnaireRepository struct {
	db *sql.DB
}

func NewRiskQuestionnaireRepository(db *sql.DB) *RiskQuestionnaireRepository {
	return &RiskQuestionnaireRepository{db: db}
}

func (r *RiskQuestionnaireRepository) Create(rq *models.RiskQuestionnaire) error {
	now := time.Now()
	query := `
		INSERT INTO risk_questionnaires (user_id, q1_answer, q2_answer, q3_answer, catalyst_unlocked, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		rq.UserID,
		rq.Q1Answer,
		rq.Q2Answer,
		rq.Q3Answer,
		rq.CatalystUnlocked,
		now,
	).Scan(&rq.ID, &rq.CreatedAt)
}

func (r *RiskQuestionnaireRepository) FindByUserID(userID uuid.UUID) (*models.RiskQuestionnaire, error) {
	rq := &models.RiskQuestionnaire{}
	query := `
		SELECT id, user_id, q1_answer, q2_answer, q3_answer, catalyst_unlocked, completed_at, created_at
		FROM risk_questionnaires
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&rq.ID,
		&rq.UserID,
		&rq.Q1Answer,
		&rq.Q2Answer,
		&rq.Q3Answer,
		&rq.CatalystUnlocked,
		&rq.CompletedAt,
		&rq.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return rq, nil
}

func (r *RiskQuestionnaireRepository) Update(rq *models.RiskQuestionnaire) error {
	now := time.Now()
	query := `
		UPDATE risk_questionnaires
		SET q1_answer = $1, q2_answer = $2, q3_answer = $3, catalyst_unlocked = $4, completed_at = $5
		WHERE user_id = $6
	`
	_, err := r.db.Exec(
		query,
		rq.Q1Answer,
		rq.Q2Answer,
		rq.Q3Answer,
		rq.CatalystUnlocked,
		now,
		rq.UserID,
	)
	return err
}

func (r *RiskQuestionnaireRepository) IsCatalystUnlocked(userID uuid.UUID) (bool, error) {
	var unlocked bool
	query := `SELECT catalyst_unlocked FROM risk_questionnaires WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&unlocked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return unlocked, nil
}
