package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (invoice_id, user_id, type, amount, currency, tx_hash, status, from_address, to_address, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		tx.InvoiceID,
		tx.UserID,
		tx.Type,
		tx.Amount,
		tx.Currency,
		tx.TxHash,
		tx.Status,
		tx.FromAddress,
		tx.ToAddress,
		tx.Notes,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
}

func (r *TransactionRepository) FindByID(id uuid.UUID) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `
		SELECT id, invoice_id, user_id, type, amount, currency, tx_hash, status,
		       from_address, to_address, block_number, gas_used, notes, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID,
		&tx.InvoiceID,
		&tx.UserID,
		&tx.Type,
		&tx.Amount,
		&tx.Currency,
		&tx.TxHash,
		&tx.Status,
		&tx.FromAddress,
		&tx.ToAddress,
		&tx.BlockNumber,
		&tx.GasUsed,
		&tx.Notes,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepository) FindByTxHash(txHash string) (*models.Transaction, error) {
	tx := &models.Transaction{}
	query := `
		SELECT id, invoice_id, user_id, type, amount, currency, tx_hash, status,
		       from_address, to_address, block_number, gas_used, notes, created_at, updated_at
		FROM transactions
		WHERE tx_hash = $1
	`
	err := r.db.QueryRow(query, txHash).Scan(
		&tx.ID,
		&tx.InvoiceID,
		&tx.UserID,
		&tx.Type,
		&tx.Amount,
		&tx.Currency,
		&tx.TxHash,
		&tx.Status,
		&tx.FromAddress,
		&tx.ToAddress,
		&tx.BlockNumber,
		&tx.GasUsed,
		&tx.Notes,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepository) FindByUser(userID uuid.UUID, page, perPage int) ([]models.Transaction, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	if err := r.db.QueryRow(countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, invoice_id, user_id, type, amount, currency, tx_hash, status,
		       from_address, to_address, block_number, gas_used, notes, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID,
			&tx.InvoiceID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.Currency,
			&tx.TxHash,
			&tx.Status,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.BlockNumber,
			&tx.GasUsed,
			&tx.Notes,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		txs = append(txs, tx)
	}
	return txs, total, nil
}

func (r *TransactionRepository) FindByInvoice(invoiceID uuid.UUID) ([]models.Transaction, error) {
	query := `
		SELECT id, invoice_id, user_id, type, amount, currency, tx_hash, status,
		       from_address, to_address, block_number, gas_used, notes, created_at, updated_at
		FROM transactions
		WHERE invoice_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID,
			&tx.InvoiceID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.Currency,
			&tx.TxHash,
			&tx.Status,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.BlockNumber,
			&tx.GasUsed,
			&tx.Notes,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (r *TransactionRepository) UpdateStatus(id uuid.UUID, status models.TransactionStatus) error {
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *TransactionRepository) UpdateBlockInfo(id uuid.UUID, blockNumber, gasUsed int64) error {
	query := `UPDATE transactions SET block_number = $1, gas_used = $2, status = 'confirmed', updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, blockNumber, gasUsed, time.Now(), id)
	return err
}
