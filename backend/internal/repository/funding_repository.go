package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/receiv3/backend/internal/models"
)

type FundingRepository struct {
	db *sql.DB
}

func NewFundingRepository(db *sql.DB) *FundingRepository {
	return &FundingRepository{db: db}
}

// Pool methods
func (r *FundingRepository) CreatePool(pool *models.FundingPool) error {
	now := time.Now()
	query := `
		INSERT INTO funding_pools (invoice_id, target_amount, status, opened_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, funded_amount, investor_count, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		pool.InvoiceID,
		pool.TargetAmount,
		pool.Status,
		now,
	).Scan(&pool.ID, &pool.FundedAmount, &pool.InvestorCount, &pool.CreatedAt, &pool.UpdatedAt)
}

func (r *FundingRepository) FindPoolByID(id uuid.UUID) (*models.FundingPool, error) {
	pool := &models.FundingPool{}
	query := `
		SELECT id, invoice_id, target_amount, funded_amount, investor_count, status,
		       opened_at, filled_at, disbursed_at, closed_at, created_at, updated_at
		FROM funding_pools
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&pool.ID,
		&pool.InvoiceID,
		&pool.TargetAmount,
		&pool.FundedAmount,
		&pool.InvestorCount,
		&pool.Status,
		&pool.OpenedAt,
		&pool.FilledAt,
		&pool.DisbursedAt,
		&pool.ClosedAt,
		&pool.CreatedAt,
		&pool.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (r *FundingRepository) FindPoolByInvoiceID(invoiceID uuid.UUID) (*models.FundingPool, error) {
	pool := &models.FundingPool{}
	query := `
		SELECT id, invoice_id, target_amount, funded_amount, investor_count, status,
		       opened_at, filled_at, disbursed_at, closed_at, created_at, updated_at
		FROM funding_pools
		WHERE invoice_id = $1
	`
	err := r.db.QueryRow(query, invoiceID).Scan(
		&pool.ID,
		&pool.InvoiceID,
		&pool.TargetAmount,
		&pool.FundedAmount,
		&pool.InvestorCount,
		&pool.Status,
		&pool.OpenedAt,
		&pool.FilledAt,
		&pool.DisbursedAt,
		&pool.ClosedAt,
		&pool.CreatedAt,
		&pool.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (r *FundingRepository) FindOpenPools(page, perPage int) ([]models.FundingPool, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM funding_pools WHERE status = 'open'`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT fp.id, fp.invoice_id, fp.target_amount, fp.funded_amount, fp.investor_count, fp.status,
		       fp.opened_at, fp.filled_at, fp.disbursed_at, fp.closed_at, fp.created_at, fp.updated_at
		FROM funding_pools fp
		WHERE fp.status = 'open'
		ORDER BY fp.opened_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pools []models.FundingPool
	for rows.Next() {
		var pool models.FundingPool
		if err := rows.Scan(
			&pool.ID,
			&pool.InvoiceID,
			&pool.TargetAmount,
			&pool.FundedAmount,
			&pool.InvestorCount,
			&pool.Status,
			&pool.OpenedAt,
			&pool.FilledAt,
			&pool.DisbursedAt,
			&pool.ClosedAt,
			&pool.CreatedAt,
			&pool.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		pools = append(pools, pool)
	}
	return pools, total, nil
}

func (r *FundingRepository) UpdatePoolFunding(id uuid.UUID, amount float64) error {
	query := `
		UPDATE funding_pools
		SET funded_amount = funded_amount + $1, investor_count = investor_count + 1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *FundingRepository) UpdatePoolStatus(id uuid.UUID, status models.PoolStatus) error {
	now := time.Now()
	query := `UPDATE funding_pools SET status = $1, updated_at = $2`

	switch status {
	case models.PoolStatusFilled:
		query += `, filled_at = $2`
	case models.PoolStatusDisbursed:
		query += `, disbursed_at = $2`
	case models.PoolStatusClosed:
		query += `, closed_at = $2`
	}

	query += ` WHERE id = $3`
	_, err := r.db.Exec(query, status, now, id)
	return err
}

// Investment methods
func (r *FundingRepository) CreateInvestment(inv *models.Investment) error {
	query := `
		INSERT INTO investments (pool_id, investor_id, amount, expected_return, status, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, invested_at, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		inv.PoolID,
		inv.InvestorID,
		inv.Amount,
		inv.ExpectedReturn,
		inv.Status,
		inv.TxHash,
	).Scan(&inv.ID, &inv.InvestedAt, &inv.CreatedAt, &inv.UpdatedAt)
}

func (r *FundingRepository) FindInvestmentByID(id uuid.UUID) (*models.Investment, error) {
	inv := &models.Investment{}
	query := `
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status, tx_hash,
		       invested_at, repaid_at, created_at, updated_at
		FROM investments
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&inv.ID,
		&inv.PoolID,
		&inv.InvestorID,
		&inv.Amount,
		&inv.ExpectedReturn,
		&inv.ActualReturn,
		&inv.Status,
		&inv.TxHash,
		&inv.InvestedAt,
		&inv.RepaidAt,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return inv, nil
}

func (r *FundingRepository) FindInvestmentsByInvestor(investorID uuid.UUID, page, perPage int) ([]models.Investment, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM investments WHERE investor_id = $1`
	if err := r.db.QueryRow(countQuery, investorID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status, tx_hash,
		       invested_at, repaid_at, created_at, updated_at
		FROM investments
		WHERE investor_id = $1
		ORDER BY invested_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, investorID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var investments []models.Investment
	for rows.Next() {
		var inv models.Investment
		if err := rows.Scan(
			&inv.ID,
			&inv.PoolID,
			&inv.InvestorID,
			&inv.Amount,
			&inv.ExpectedReturn,
			&inv.ActualReturn,
			&inv.Status,
			&inv.TxHash,
			&inv.InvestedAt,
			&inv.RepaidAt,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		investments = append(investments, inv)
	}
	return investments, total, nil
}

func (r *FundingRepository) FindInvestmentsByPool(poolID uuid.UUID) ([]models.Investment, error) {
	query := `
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status, tx_hash,
		       invested_at, repaid_at, created_at, updated_at
		FROM investments
		WHERE pool_id = $1
	`
	rows, err := r.db.Query(query, poolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []models.Investment
	for rows.Next() {
		var inv models.Investment
		if err := rows.Scan(
			&inv.ID,
			&inv.PoolID,
			&inv.InvestorID,
			&inv.Amount,
			&inv.ExpectedReturn,
			&inv.ActualReturn,
			&inv.Status,
			&inv.TxHash,
			&inv.InvestedAt,
			&inv.RepaidAt,
			&inv.CreatedAt,
			&inv.UpdatedAt,
		); err != nil {
			return nil, err
		}
		investments = append(investments, inv)
	}
	return investments, nil
}

func (r *FundingRepository) UpdateInvestmentStatus(id uuid.UUID, status models.InvestmentStatus, actualReturn *float64) error {
	now := time.Now()
	query := `UPDATE investments SET status = $1, actual_return = $2, repaid_at = $3, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, actualReturn, now, id)
	return err
}
