package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/vessel/backend/internal/models"
)

type MitraRepository struct {
	db *sql.DB
}

func NewMitraRepository(db *sql.DB) *MitraRepository {
	return &MitraRepository{db: db}
}

// Create creates a new MITRA application
func (r *MitraRepository) Create(app *models.MitraApplication) error {
	app.ID = uuid.New()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.Status = models.MitraStatusPending

	query := `
		INSERT INTO mitra_applications (id, user_id, company_name, company_type, npwp, annual_revenue, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		app.ID,
		app.UserID,
		app.CompanyName,
		app.CompanyType,
		app.NPWP,
		app.AnnualRevenue,
		app.Status,
		app.CreatedAt,
		app.UpdatedAt,
	)

	return err
}

// FindByID finds a MITRA application by ID
func (r *MitraRepository) FindByID(id uuid.UUID) (*models.MitraApplication, error) {
	app := &models.MitraApplication{}
	query := `
		SELECT id, user_id, company_name, company_type, npwp, annual_revenue,
		       nib_document_url, akta_pendirian_url, ktp_direktur_url,
		       status, rejection_reason, reviewed_by, reviewed_at, created_at, updated_at
		FROM mitra_applications
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&app.ID,
		&app.UserID,
		&app.CompanyName,
		&app.CompanyType,
		&app.NPWP,
		&app.AnnualRevenue,
		&app.NIBDocumentURL,
		&app.AktaPendirianURL,
		&app.KTPDirekturURL,
		&app.Status,
		&app.RejectionReason,
		&app.ReviewedBy,
		&app.ReviewedAt,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return app, nil
}

// FindByUserID finds a MITRA application by user ID
func (r *MitraRepository) FindByUserID(userID uuid.UUID) (*models.MitraApplication, error) {
	app := &models.MitraApplication{}
	query := `
		SELECT id, user_id, company_name, company_type, npwp, annual_revenue,
		       nib_document_url, akta_pendirian_url, ktp_direktur_url,
		       status, rejection_reason, reviewed_by, reviewed_at, created_at, updated_at
		FROM mitra_applications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.QueryRow(query, userID).Scan(
		&app.ID,
		&app.UserID,
		&app.CompanyName,
		&app.CompanyType,
		&app.NPWP,
		&app.AnnualRevenue,
		&app.NIBDocumentURL,
		&app.AktaPendirianURL,
		&app.KTPDirekturURL,
		&app.Status,
		&app.RejectionReason,
		&app.ReviewedBy,
		&app.ReviewedAt,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return app, nil
}

// FindPending finds all pending MITRA applications
func (r *MitraRepository) FindPending(page, perPage int) ([]models.MitraApplication, int, error) {
	offset := (page - 1) * perPage

	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM mitra_applications WHERE status = 'pending'`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get applications
	query := `
		SELECT id, user_id, company_name, company_type, npwp, annual_revenue,
		       nib_document_url, akta_pendirian_url, ktp_direktur_url,
		       status, rejection_reason, reviewed_by, reviewed_at, created_at, updated_at
		FROM mitra_applications
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var applications []models.MitraApplication
	for rows.Next() {
		var app models.MitraApplication
		err := rows.Scan(
			&app.ID,
			&app.UserID,
			&app.CompanyName,
			&app.CompanyType,
			&app.NPWP,
			&app.AnnualRevenue,
			&app.NIBDocumentURL,
			&app.AktaPendirianURL,
			&app.KTPDirekturURL,
			&app.Status,
			&app.RejectionReason,
			&app.ReviewedBy,
			&app.ReviewedAt,
			&app.CreatedAt,
			&app.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		applications = append(applications, app)
	}

	return applications, total, nil
}

// UpdateDocumentURL updates the document URL for a specific document type
func (r *MitraRepository) UpdateDocumentURL(id uuid.UUID, docType, url string) error {
	var column string
	switch docType {
	case "nib":
		column = "nib_document_url"
	case "akta_pendirian":
		column = "akta_pendirian_url"
	case "ktp_direktur":
		column = "ktp_direktur_url"
	default:
		return errors.New("invalid document type")
	}

	query := `UPDATE mitra_applications SET ` + column + ` = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, url, time.Now(), id)
	return err
}

// Approve approves a MITRA application
func (r *MitraRepository) Approve(id, adminID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE mitra_applications
		SET status = 'approved', reviewed_by = $1, reviewed_at = $2, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, adminID, now, id)
	return err
}

// Reject rejects a MITRA application with a reason
func (r *MitraRepository) Reject(id, adminID uuid.UUID, reason string) error {
	now := time.Now()
	query := `
		UPDATE mitra_applications
		SET status = 'rejected', rejection_reason = $1, reviewed_by = $2, reviewed_at = $3, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, reason, adminID, now, id)
	return err
}

// HasPendingApplication checks if user has a pending MITRA application
func (r *MitraRepository) HasPendingApplication(userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM mitra_applications WHERE user_id = $1 AND status = 'pending')`
	err := r.db.QueryRow(query, userID).Scan(&exists)
	return exists, err
}

// DeleteRejected deletes a rejected application (so user can reapply)
func (r *MitraRepository) DeleteRejected(userID uuid.UUID) error {
	query := `DELETE FROM mitra_applications WHERE user_id = $1 AND status = 'rejected'`
	_, err := r.db.Exec(query, userID)
	return err
}
