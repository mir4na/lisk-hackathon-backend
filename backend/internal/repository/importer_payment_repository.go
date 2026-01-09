package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

type ImporterPaymentRepository struct {
	db *sql.DB
}

func NewImporterPaymentRepository(db *sql.DB) *ImporterPaymentRepository {
	return &ImporterPaymentRepository{db: db}
}

// Create creates a new importer payment record
func (r *ImporterPaymentRepository) Create(payment *models.ImporterPayment) error {
	query := `
		INSERT INTO importer_payments (
			invoice_id, pool_id, buyer_email, buyer_name, amount_due, currency, 
			payment_status, due_date
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, amount_paid, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		payment.InvoiceID,
		payment.PoolID,
		payment.BuyerEmail,
		payment.BuyerName,
		payment.AmountDue,
		payment.Currency,
		payment.PaymentStatus,
		payment.DueDate,
	).Scan(&payment.ID, &payment.AmountPaid, &payment.CreatedAt, &payment.UpdatedAt)
}

// FindByID finds an importer payment by ID
func (r *ImporterPaymentRepository) FindByID(id uuid.UUID) (*models.ImporterPayment, error) {
	payment := &models.ImporterPayment{}
	query := `
		SELECT id, invoice_id, pool_id, buyer_email, buyer_name, amount_due, amount_paid,
		       currency, payment_status, due_date, paid_at, tx_hash, created_at, updated_at
		FROM importer_payments
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.InvoiceID,
		&payment.PoolID,
		&payment.BuyerEmail,
		&payment.BuyerName,
		&payment.AmountDue,
		&payment.AmountPaid,
		&payment.Currency,
		&payment.PaymentStatus,
		&payment.DueDate,
		&payment.PaidAt,
		&payment.TxHash,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return payment, nil
}

// FindByInvoiceID finds payment by invoice ID
func (r *ImporterPaymentRepository) FindByInvoiceID(invoiceID uuid.UUID) (*models.ImporterPayment, error) {
	payment := &models.ImporterPayment{}
	query := `
		SELECT id, invoice_id, pool_id, buyer_email, buyer_name, amount_due, amount_paid,
		       currency, payment_status, due_date, paid_at, tx_hash, created_at, updated_at
		FROM importer_payments
		WHERE invoice_id = $1
	`
	err := r.db.QueryRow(query, invoiceID).Scan(
		&payment.ID,
		&payment.InvoiceID,
		&payment.PoolID,
		&payment.BuyerEmail,
		&payment.BuyerName,
		&payment.AmountDue,
		&payment.AmountPaid,
		&payment.Currency,
		&payment.PaymentStatus,
		&payment.DueDate,
		&payment.PaidAt,
		&payment.TxHash,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return payment, nil
}

// UpdatePayment updates payment after importer pays
func (r *ImporterPaymentRepository) UpdatePayment(id uuid.UUID, amountPaid float64, txHash string) error {
	now := time.Now()
	query := `
		UPDATE importer_payments
		SET amount_paid = $1, payment_status = $2, paid_at = $3, tx_hash = $4, updated_at = $3
		WHERE id = $5
	`
	_, err := r.db.Exec(query, amountPaid, models.ImporterPaymentStatusPaid, now, txHash, id)
	return err
}

// UpdateStatus updates payment status
func (r *ImporterPaymentRepository) UpdateStatus(id uuid.UUID, status models.ImporterPaymentStatus) error {
	query := `UPDATE importer_payments SET payment_status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// FindPendingByDueDate finds pending payments that are overdue
func (r *ImporterPaymentRepository) FindPendingByDueDate(before time.Time) ([]models.ImporterPayment, error) {
	query := `
		SELECT id, invoice_id, pool_id, buyer_email, buyer_name, amount_due, amount_paid,
		       currency, payment_status, due_date, paid_at, tx_hash, created_at, updated_at
		FROM importer_payments
		WHERE payment_status = 'pending' AND due_date < $1
	`
	rows, err := r.db.Query(query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.ImporterPayment
	for rows.Next() {
		var p models.ImporterPayment
		if err := rows.Scan(
			&p.ID, &p.InvoiceID, &p.PoolID, &p.BuyerEmail, &p.BuyerName,
			&p.AmountDue, &p.AmountPaid, &p.Currency, &p.PaymentStatus,
			&p.DueDate, &p.PaidAt, &p.TxHash, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
