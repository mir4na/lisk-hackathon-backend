package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/models"
)

type KYCRepository struct {
	db *sql.DB
}

func NewKYCRepository(db *sql.DB) *KYCRepository {
	return &KYCRepository{db: db}
}

func (r *KYCRepository) Create(kyc *models.KYCVerification) error {
	query := `
		INSERT INTO kyc_verifications (user_id, verification_type, status, id_type, id_number, id_document_url, selfie_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		kyc.UserID,
		kyc.VerificationType,
		kyc.Status,
		kyc.IDType,
		kyc.IDNumber,
		kyc.IDDocumentURL,
		kyc.SelfieURL,
	).Scan(&kyc.ID, &kyc.CreatedAt, &kyc.UpdatedAt)
}

func (r *KYCRepository) FindByID(id uuid.UUID) (*models.KYCVerification, error) {
	kyc := &models.KYCVerification{}
	query := `
		SELECT id, user_id, verification_type, status, id_type, id_number, id_document_url, selfie_url,
		       rejection_reason, verified_by, verified_at, created_at, updated_at
		FROM kyc_verifications
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&kyc.ID,
		&kyc.UserID,
		&kyc.VerificationType,
		&kyc.Status,
		&kyc.IDType,
		&kyc.IDNumber,
		&kyc.IDDocumentURL,
		&kyc.SelfieURL,
		&kyc.RejectionReason,
		&kyc.VerifiedBy,
		&kyc.VerifiedAt,
		&kyc.CreatedAt,
		&kyc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return kyc, nil
}

func (r *KYCRepository) FindByUserID(userID uuid.UUID) (*models.KYCVerification, error) {
	kyc := &models.KYCVerification{}
	query := `
		SELECT id, user_id, verification_type, status, id_type, id_number, id_document_url, selfie_url,
		       rejection_reason, verified_by, verified_at, created_at, updated_at
		FROM kyc_verifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&kyc.ID,
		&kyc.UserID,
		&kyc.VerificationType,
		&kyc.Status,
		&kyc.IDType,
		&kyc.IDNumber,
		&kyc.IDDocumentURL,
		&kyc.SelfieURL,
		&kyc.RejectionReason,
		&kyc.VerifiedBy,
		&kyc.VerifiedAt,
		&kyc.CreatedAt,
		&kyc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return kyc, nil
}

func (r *KYCRepository) FindPending(page, perPage int) ([]models.KYCVerification, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM kyc_verifications WHERE status = 'pending'`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, verification_type, status, id_type, id_number, id_document_url, selfie_url,
		       rejection_reason, verified_by, verified_at, created_at, updated_at
		FROM kyc_verifications
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var kycs []models.KYCVerification
	for rows.Next() {
		var kyc models.KYCVerification
		if err := rows.Scan(
			&kyc.ID,
			&kyc.UserID,
			&kyc.VerificationType,
			&kyc.Status,
			&kyc.IDType,
			&kyc.IDNumber,
			&kyc.IDDocumentURL,
			&kyc.SelfieURL,
			&kyc.RejectionReason,
			&kyc.VerifiedBy,
			&kyc.VerifiedAt,
			&kyc.CreatedAt,
			&kyc.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		kycs = append(kycs, kyc)
	}
	return kycs, total, nil
}

func (r *KYCRepository) Approve(kycID, adminID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE kyc_verifications
		SET status = 'approved', verified_by = $1, verified_at = $2, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, adminID, now, kycID)
	return err
}

func (r *KYCRepository) Reject(kycID, adminID uuid.UUID, reason string) error {
	now := time.Now()
	query := `
		UPDATE kyc_verifications
		SET status = 'rejected', verified_by = $1, verified_at = $2, rejection_reason = $3, updated_at = $2
		WHERE id = $4
	`
	_, err := r.db.Exec(query, adminID, now, reason, kycID)
	return err
}

func (r *KYCRepository) UpdateDocumentURL(kycID uuid.UUID, docURL string) error {
	query := `UPDATE kyc_verifications SET id_document_url = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, docURL, time.Now(), kycID)
	return err
}

func (r *KYCRepository) UpdateSelfieURL(kycID uuid.UUID, selfieURL string) error {
	query := `UPDATE kyc_verifications SET selfie_url = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, selfieURL, time.Now(), kycID)
	return err
}
