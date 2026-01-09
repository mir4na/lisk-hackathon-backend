package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
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
		INSERT INTO funding_pools (
			invoice_id, target_amount, status, opened_at, deadline,
			priority_target, priority_funded, catalyst_target, catalyst_funded,
			priority_interest_rate, catalyst_interest_rate, pool_currency
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, funded_amount, investor_count, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		pool.InvoiceID,
		pool.TargetAmount,
		pool.Status,
		now,
		pool.Deadline,
		pool.PriorityTarget,
		pool.PriorityFunded,
		pool.CatalystTarget,
		pool.CatalystFunded,
		pool.PriorityInterestRate,
		pool.CatalystInterestRate,
		pool.PoolCurrency,
	).Scan(&pool.ID, &pool.FundedAmount, &pool.InvestorCount, &pool.CreatedAt, &pool.UpdatedAt)
}

func (r *FundingRepository) FindPoolByID(id uuid.UUID) (*models.FundingPool, error) {
	pool := &models.FundingPool{}
	query := `
		SELECT id, invoice_id, target_amount, funded_amount, investor_count, status,
		       opened_at, deadline, filled_at, disbursed_at, closed_at, created_at, updated_at,
		       COALESCE(priority_target, 0), COALESCE(priority_funded, 0),
		       COALESCE(catalyst_target, 0), COALESCE(catalyst_funded, 0),
		       COALESCE(priority_interest_rate, 0), COALESCE(catalyst_interest_rate, 0),
		       COALESCE(pool_currency, 'IDRX')
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
		&pool.Deadline,
		&pool.FilledAt,
		&pool.DisbursedAt,
		&pool.ClosedAt,
		&pool.CreatedAt,
		&pool.UpdatedAt,
		&pool.PriorityTarget,
		&pool.PriorityFunded,
		&pool.CatalystTarget,
		&pool.CatalystFunded,
		&pool.PriorityInterestRate,
		&pool.CatalystInterestRate,
		&pool.PoolCurrency,
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
		       opened_at, deadline, filled_at, disbursed_at, closed_at, created_at, updated_at,
		       COALESCE(priority_target, 0), COALESCE(priority_funded, 0),
		       COALESCE(catalyst_target, 0), COALESCE(catalyst_funded, 0),
		       COALESCE(priority_interest_rate, 0), COALESCE(catalyst_interest_rate, 0),
		       COALESCE(pool_currency, 'IDRX')
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
		&pool.Deadline,
		&pool.FilledAt,
		&pool.DisbursedAt,
		&pool.ClosedAt,
		&pool.CreatedAt,
		&pool.UpdatedAt,
		&pool.PriorityTarget,
		&pool.PriorityFunded,
		&pool.CatalystTarget,
		&pool.CatalystFunded,
		&pool.PriorityInterestRate,
		&pool.CatalystInterestRate,
		&pool.PoolCurrency,
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
		       fp.opened_at, fp.deadline, fp.filled_at, fp.disbursed_at, fp.closed_at, fp.created_at, fp.updated_at,
		       COALESCE(fp.priority_target, 0), COALESCE(fp.priority_funded, 0),
		       COALESCE(fp.catalyst_target, 0), COALESCE(fp.catalyst_funded, 0),
		       COALESCE(fp.priority_interest_rate, 0), COALESCE(fp.catalyst_interest_rate, 0),
		       COALESCE(fp.pool_currency, 'IDRX')
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
			&pool.Deadline,
			&pool.FilledAt,
			&pool.DisbursedAt,
			&pool.ClosedAt,
			&pool.CreatedAt,
			&pool.UpdatedAt,
			&pool.PriorityTarget,
			&pool.PriorityFunded,
			&pool.CatalystTarget,
			&pool.CatalystFunded,
			&pool.PriorityInterestRate,
			&pool.CatalystInterestRate,
			&pool.PoolCurrency,
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

// UpdatePoolTrancheFunding updates funding for a specific tranche
func (r *FundingRepository) UpdatePoolTrancheFunding(id uuid.UUID, amount float64, tranche models.TrancheType) error {
	now := time.Now()
	var query string
	if tranche == models.TranchePriority {
		query = `
			UPDATE funding_pools
			SET funded_amount = funded_amount + $1,
			    priority_funded = COALESCE(priority_funded, 0) + $1,
			    investor_count = investor_count + 1,
			    updated_at = $2
			WHERE id = $3
		`
	} else {
		query = `
			UPDATE funding_pools
			SET funded_amount = funded_amount + $1,
			    catalyst_funded = COALESCE(catalyst_funded, 0) + $1,
			    investor_count = investor_count + 1,
			    updated_at = $2
			WHERE id = $3
		`
	}
	_, err := r.db.Exec(query, amount, now, id)
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
		INSERT INTO investments (pool_id, investor_id, amount, expected_return, status, tranche, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, invested_at, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		inv.PoolID,
		inv.InvestorID,
		inv.Amount,
		inv.ExpectedReturn,
		inv.Status,
		inv.Tranche,
		inv.TxHash,
	).Scan(&inv.ID, &inv.InvestedAt, &inv.CreatedAt, &inv.UpdatedAt)
}

func (r *FundingRepository) FindInvestmentByID(id uuid.UUID) (*models.Investment, error) {
	inv := &models.Investment{}
	query := `
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status,
		       COALESCE(tranche, 'priority'), tx_hash, invested_at, repaid_at, created_at, updated_at
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
		&inv.Tranche,
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
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status,
		       COALESCE(tranche, 'priority'), tx_hash, invested_at, repaid_at, created_at, updated_at
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
			&inv.Tranche,
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
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status,
		       COALESCE(tranche, 'priority'), tx_hash, invested_at, repaid_at, created_at, updated_at
		FROM investments
		WHERE pool_id = $1
		ORDER BY tranche ASC, invested_at ASC
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
			&inv.Tranche,
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

// FindInvestmentsByPoolAndTranche returns investments for a pool filtered by tranche
func (r *FundingRepository) FindInvestmentsByPoolAndTranche(poolID uuid.UUID, tranche models.TrancheType) ([]models.Investment, error) {
	query := `
		SELECT id, pool_id, investor_id, amount, expected_return, actual_return, status,
		       COALESCE(tranche, 'priority'), tx_hash, invested_at, repaid_at, created_at, updated_at
		FROM investments
		WHERE pool_id = $1 AND tranche = $2
		ORDER BY invested_at ASC
	`
	rows, err := r.db.Query(query, poolID, tranche)
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
			&inv.Tranche,
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

// GetInvestorPortfolio calculates portfolio summary for an investor
func (r *FundingRepository) GetInvestorPortfolio(investorID uuid.UUID) (*models.InvestorPortfolio, error) {
	portfolio := &models.InvestorPortfolio{}

	query := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_funding,
			COALESCE(SUM(expected_return - amount), 0) as total_expected_gain,
			COALESCE(SUM(CASE WHEN status = 'repaid' THEN COALESCE(actual_return, 0) - amount ELSE 0 END), 0) as total_realized_gain,
			COALESCE(SUM(CASE WHEN tranche = 'priority' THEN amount ELSE 0 END), 0) as priority_allocation,
			COALESCE(SUM(CASE WHEN tranche = 'catalyst' THEN amount ELSE 0 END), 0) as catalyst_allocation,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_investments,
			COUNT(CASE WHEN status = 'repaid' THEN 1 END) as completed_deals
		FROM investments
		WHERE investor_id = $1
	`

	err := r.db.QueryRow(query, investorID).Scan(
		&portfolio.TotalFunding,
		&portfolio.TotalExpectedGain,
		&portfolio.TotalRealizedGain,
		&portfolio.PriorityAllocation,
		&portfolio.CatalystAllocation,
		&portfolio.ActiveInvestments,
		&portfolio.CompletedDeals,
	)
	if err != nil {
		return nil, err
	}

	return portfolio, nil
}
