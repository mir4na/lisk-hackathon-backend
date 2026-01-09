package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/vessel/backend/internal/models"
)

type OTPRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// Create creates a new OTP record
func (r *OTPRepository) Create(otp *models.OTPCode) error {
	otp.ID = uuid.New()
	otp.CreatedAt = time.Now()

	query := `
		INSERT INTO otp_codes (id, email, code, purpose, expires_at, verified, attempts, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		otp.ID,
		otp.Email,
		otp.Code,
		otp.Purpose,
		otp.ExpiresAt,
		otp.Verified,
		otp.Attempts,
		otp.CreatedAt,
	)

	return err
}

// FindLatestByEmail finds the most recent OTP for an email and purpose
func (r *OTPRepository) FindLatestByEmail(email string, purpose models.OTPPurpose) (*models.OTPCode, error) {
	query := `
		SELECT id, email, code, purpose, expires_at, verified, attempts, created_at
		FROM otp_codes
		WHERE email = $1 AND purpose = $2 AND verified = false
		ORDER BY created_at DESC
		LIMIT 1
	`

	otp := &models.OTPCode{}
	err := r.db.QueryRow(query, email, purpose).Scan(
		&otp.ID,
		&otp.Email,
		&otp.Code,
		&otp.Purpose,
		&otp.ExpiresAt,
		&otp.Verified,
		&otp.Attempts,
		&otp.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return otp, nil
}

// IncrementAttempts increments the attempt count for an OTP
func (r *OTPRepository) IncrementAttempts(id uuid.UUID) error {
	query := `UPDATE otp_codes SET attempts = attempts + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// MarkVerified marks an OTP as verified
func (r *OTPRepository) MarkVerified(id uuid.UUID) error {
	query := `UPDATE otp_codes SET verified = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// CountRecentOTPs counts OTPs sent to an email in the last hour
func (r *OTPRepository) CountRecentOTPs(email string, purpose models.OTPPurpose) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM otp_codes
		WHERE email = $1 AND purpose = $2 AND created_at > NOW() - INTERVAL '1 hour'
	`

	var count int
	err := r.db.QueryRow(query, email, purpose).Scan(&count)
	return count, err
}

// DeleteExpired deletes expired OTP records (cleanup)
func (r *OTPRepository) DeleteExpired() error {
	query := `DELETE FROM otp_codes WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}

// InvalidateAllForEmail invalidates all unverified OTPs for an email
func (r *OTPRepository) InvalidateAllForEmail(email string, purpose models.OTPPurpose) error {
	query := `UPDATE otp_codes SET verified = true WHERE email = $1 AND purpose = $2 AND verified = false`
	_, err := r.db.Exec(query, email, purpose)
	return err
}
