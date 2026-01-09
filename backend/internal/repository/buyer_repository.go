package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

type BuyerRepository struct {
	db *sql.DB
}

func NewBuyerRepository(db *sql.DB) *BuyerRepository {
	return &BuyerRepository{db: db}
}

func (r *BuyerRepository) Create(buyer *models.Buyer) error {
	query := `
		INSERT INTO buyers (created_by, company_name, country, address, contact_email, contact_phone, website)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, credit_score, total_invoices, total_paid, total_defaulted, is_verified, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		buyer.CreatedBy,
		buyer.CompanyName,
		buyer.Country,
		buyer.Address,
		buyer.ContactEmail,
		buyer.ContactPhone,
		buyer.Website,
	).Scan(
		&buyer.ID,
		&buyer.CreditScore,
		&buyer.TotalInvoices,
		&buyer.TotalPaid,
		&buyer.TotalDefaulted,
		&buyer.IsVerified,
		&buyer.CreatedAt,
		&buyer.UpdatedAt,
	)
}

func (r *BuyerRepository) FindByID(id uuid.UUID) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	query := `
		SELECT id, created_by, company_name, country, address, contact_email, contact_phone, website,
		       credit_score, total_invoices, total_paid, total_defaulted, is_verified, created_at, updated_at
		FROM buyers
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&buyer.ID,
		&buyer.CreatedBy,
		&buyer.CompanyName,
		&buyer.Country,
		&buyer.Address,
		&buyer.ContactEmail,
		&buyer.ContactPhone,
		&buyer.Website,
		&buyer.CreditScore,
		&buyer.TotalInvoices,
		&buyer.TotalPaid,
		&buyer.TotalDefaulted,
		&buyer.IsVerified,
		&buyer.CreatedAt,
		&buyer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return buyer, nil
}

func (r *BuyerRepository) FindByExporter(exporterID uuid.UUID, page, perPage int) ([]models.Buyer, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM buyers WHERE created_by = $1`
	if err := r.db.QueryRow(countQuery, exporterID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, created_by, company_name, country, address, contact_email, contact_phone, website,
		       credit_score, total_invoices, total_paid, total_defaulted, is_verified, created_at, updated_at
		FROM buyers
		WHERE created_by = $1
		ORDER BY company_name ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, exporterID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var buyers []models.Buyer
	for rows.Next() {
		var buyer models.Buyer
		if err := rows.Scan(
			&buyer.ID,
			&buyer.CreatedBy,
			&buyer.CompanyName,
			&buyer.Country,
			&buyer.Address,
			&buyer.ContactEmail,
			&buyer.ContactPhone,
			&buyer.Website,
			&buyer.CreditScore,
			&buyer.TotalInvoices,
			&buyer.TotalPaid,
			&buyer.TotalDefaulted,
			&buyer.IsVerified,
			&buyer.CreatedAt,
			&buyer.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		buyers = append(buyers, buyer)
	}
	return buyers, total, nil
}

func (r *BuyerRepository) Update(buyer *models.Buyer) error {
	query := `
		UPDATE buyers
		SET company_name = $1, country = $2, address = $3, contact_email = $4, contact_phone = $5, website = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(
		query,
		buyer.CompanyName,
		buyer.Country,
		buyer.Address,
		buyer.ContactEmail,
		buyer.ContactPhone,
		buyer.Website,
		time.Now(),
		buyer.ID,
	)
	return err
}

func (r *BuyerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM buyers WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *BuyerRepository) UpdateCreditScore(id uuid.UUID, score int) error {
	query := `UPDATE buyers SET credit_score = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, score, time.Now(), id)
	return err
}

func (r *BuyerRepository) IncrementInvoiceStats(id uuid.UUID, amount float64, isPaid bool) error {
	var query string
	if isPaid {
		query = `UPDATE buyers SET total_invoices = total_invoices + 1, total_paid = total_paid + $1, updated_at = $2 WHERE id = $3`
	} else {
		query = `UPDATE buyers SET total_invoices = total_invoices + 1, total_defaulted = total_defaulted + $1, updated_at = $2 WHERE id = $3`
	}
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *BuyerRepository) SetVerified(id uuid.UUID, verified bool) error {
	query := `UPDATE buyers SET is_verified = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, verified, time.Now(), id)
	return err
}
