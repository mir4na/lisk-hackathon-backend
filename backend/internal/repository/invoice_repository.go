package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(invoice *models.Invoice) error {
	query := `
		INSERT INTO invoices (exporter_id, buyer_id, invoice_number, currency, amount, issue_date, due_date, description, status, advance_percentage)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		invoice.ExporterID,
		invoice.BuyerID,
		invoice.InvoiceNumber,
		invoice.Currency,
		invoice.Amount,
		invoice.IssueDate,
		invoice.DueDate,
		invoice.Description,
		invoice.Status,
		invoice.AdvancePercentage,
	).Scan(&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt)
}

func (r *InvoiceRepository) FindByID(id uuid.UUID) (*models.Invoice, error) {
	invoice := &models.Invoice{}
	query := `
		SELECT i.id, i.exporter_id, i.buyer_id, i.invoice_number, i.currency, i.amount, i.issue_date, i.due_date,
		       i.description, i.status, i.interest_rate, i.advance_percentage, i.advance_amount, i.document_hash,
		       i.created_at, i.updated_at,
		       b.id, b.company_name, b.country, b.credit_score
		FROM invoices i
		LEFT JOIN buyers b ON i.buyer_id = b.id
		WHERE i.id = $1
	`
	var buyer models.Buyer
	err := r.db.QueryRow(query, id).Scan(
		&invoice.ID,
		&invoice.ExporterID,
		&invoice.BuyerID,
		&invoice.InvoiceNumber,
		&invoice.Currency,
		&invoice.Amount,
		&invoice.IssueDate,
		&invoice.DueDate,
		&invoice.Description,
		&invoice.Status,
		&invoice.InterestRate,
		&invoice.AdvancePercentage,
		&invoice.AdvanceAmount,
		&invoice.DocumentHash,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
		&buyer.ID,
		&buyer.CompanyName,
		&buyer.Country,
		&buyer.CreditScore,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	invoice.Buyer = &buyer

	// Load documents
	docs, _ := r.FindDocumentsByInvoiceID(id)
	invoice.Documents = docs

	// Load NFT
	nft, _ := r.FindNFTByInvoiceID(id)
	invoice.NFT = nft

	return invoice, nil
}

func (r *InvoiceRepository) FindByExporter(exporterID uuid.UUID, filter *models.InvoiceFilter) ([]models.Invoice, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM invoices WHERE exporter_id = $1`
	args := []interface{}{exporterID}
	argCount := 1

	if filter.Status != nil {
		argCount++
		countQuery += ` AND status = $` + string(rune('0'+argCount))
		args = append(args, *filter.Status)
	}

	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	query := `
		SELECT i.id, i.exporter_id, i.buyer_id, i.invoice_number, i.currency, i.amount, i.issue_date, i.due_date,
		       i.description, i.status, i.interest_rate, i.advance_percentage, i.advance_amount, i.document_hash,
		       i.created_at, i.updated_at,
		       b.id, b.company_name, b.country, b.credit_score
		FROM invoices i
		LEFT JOIN buyers b ON i.buyer_id = b.id
		WHERE i.exporter_id = $1
	`
	if filter.Status != nil {
		query += ` AND i.status = $2 ORDER BY i.created_at DESC LIMIT $3 OFFSET $4`
		args = append(args, filter.PerPage, offset)
	} else {
		query += ` ORDER BY i.created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{exporterID, filter.PerPage, offset}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		var buyer models.Buyer
		if err := rows.Scan(
			&invoice.ID,
			&invoice.ExporterID,
			&invoice.BuyerID,
			&invoice.InvoiceNumber,
			&invoice.Currency,
			&invoice.Amount,
			&invoice.IssueDate,
			&invoice.DueDate,
			&invoice.Description,
			&invoice.Status,
			&invoice.InterestRate,
			&invoice.AdvancePercentage,
			&invoice.AdvanceAmount,
			&invoice.DocumentHash,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&buyer.ID,
			&buyer.CompanyName,
			&buyer.Country,
			&buyer.CreditScore,
		); err != nil {
			return nil, 0, err
		}
		invoice.Buyer = &buyer
		invoices = append(invoices, invoice)
	}
	return invoices, total, nil
}

func (r *InvoiceRepository) FindFundable(page, perPage int) ([]models.Invoice, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM invoices WHERE status = 'funding'`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT i.id, i.exporter_id, i.buyer_id, i.invoice_number, i.currency, i.amount, i.issue_date, i.due_date,
		       i.description, i.status, i.interest_rate, i.advance_percentage, i.advance_amount, i.document_hash,
		       i.created_at, i.updated_at,
		       b.id, b.company_name, b.country, b.credit_score
		FROM invoices i
		LEFT JOIN buyers b ON i.buyer_id = b.id
		WHERE i.status = 'funding'
		ORDER BY i.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		var buyer models.Buyer
		if err := rows.Scan(
			&invoice.ID,
			&invoice.ExporterID,
			&invoice.BuyerID,
			&invoice.InvoiceNumber,
			&invoice.Currency,
			&invoice.Amount,
			&invoice.IssueDate,
			&invoice.DueDate,
			&invoice.Description,
			&invoice.Status,
			&invoice.InterestRate,
			&invoice.AdvancePercentage,
			&invoice.AdvanceAmount,
			&invoice.DocumentHash,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&buyer.ID,
			&buyer.CompanyName,
			&buyer.Country,
			&buyer.CreditScore,
		); err != nil {
			return nil, 0, err
		}
		invoice.Buyer = &buyer
		invoices = append(invoices, invoice)
	}
	return invoices, total, nil
}

func (r *InvoiceRepository) Update(invoice *models.Invoice) error {
	query := `
		UPDATE invoices
		SET invoice_number = $1, currency = $2, amount = $3, issue_date = $4, due_date = $5,
		    description = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(
		query,
		invoice.InvoiceNumber,
		invoice.Currency,
		invoice.Amount,
		invoice.IssueDate,
		invoice.DueDate,
		invoice.Description,
		time.Now(),
		invoice.ID,
	)
	return err
}

func (r *InvoiceRepository) UpdateStatus(id uuid.UUID, status models.InvoiceStatus) error {
	query := `UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *InvoiceRepository) SetInterestRate(id uuid.UUID, rate float64) error {
	query := `UPDATE invoices SET interest_rate = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, rate, time.Now(), id)
	return err
}

func (r *InvoiceRepository) SetAdvanceAmount(id uuid.UUID, amount float64) error {
	query := `UPDATE invoices SET advance_amount = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, amount, time.Now(), id)
	return err
}

func (r *InvoiceRepository) SetDocumentHash(id uuid.UUID, hash string) error {
	query := `UPDATE invoices SET document_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, hash, time.Now(), id)
	return err
}

func (r *InvoiceRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM invoices WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Document methods
func (r *InvoiceRepository) CreateDocument(doc *models.InvoiceDocument) error {
	query := `
		INSERT INTO invoice_documents (invoice_id, document_type, file_name, file_url, file_hash, file_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uploaded_at
	`
	return r.db.QueryRow(
		query,
		doc.InvoiceID,
		doc.DocumentType,
		doc.FileName,
		doc.FileURL,
		doc.FileHash,
		doc.FileSize,
	).Scan(&doc.ID, &doc.UploadedAt)
}

func (r *InvoiceRepository) FindDocumentsByInvoiceID(invoiceID uuid.UUID) ([]models.InvoiceDocument, error) {
	query := `
		SELECT id, invoice_id, document_type, file_name, file_url, file_hash, file_size, uploaded_at
		FROM invoice_documents
		WHERE invoice_id = $1
	`
	rows, err := r.db.Query(query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.InvoiceDocument
	for rows.Next() {
		var doc models.InvoiceDocument
		if err := rows.Scan(
			&doc.ID,
			&doc.InvoiceID,
			&doc.DocumentType,
			&doc.FileName,
			&doc.FileURL,
			&doc.FileHash,
			&doc.FileSize,
			&doc.UploadedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *InvoiceRepository) DeleteDocument(id uuid.UUID) error {
	query := `DELETE FROM invoice_documents WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// NFT methods
func (r *InvoiceRepository) CreateNFT(nft *models.InvoiceNFT) error {
	query := `
		INSERT INTO invoice_nfts (invoice_id, token_id, contract_address, chain_id, owner_address, mint_tx_hash, metadata_uri, minted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		nft.InvoiceID,
		nft.TokenID,
		nft.ContractAddress,
		nft.ChainID,
		nft.OwnerAddress,
		nft.MintTxHash,
		nft.MetadataURI,
		nft.MintedAt,
	).Scan(&nft.ID, &nft.CreatedAt, &nft.UpdatedAt)
}

func (r *InvoiceRepository) FindNFTByInvoiceID(invoiceID uuid.UUID) (*models.InvoiceNFT, error) {
	nft := &models.InvoiceNFT{}
	query := `
		SELECT id, invoice_id, token_id, contract_address, chain_id, owner_address, mint_tx_hash, metadata_uri,
		       minted_at, burned_at, burn_tx_hash, created_at, updated_at
		FROM invoice_nfts
		WHERE invoice_id = $1
	`
	err := r.db.QueryRow(query, invoiceID).Scan(
		&nft.ID,
		&nft.InvoiceID,
		&nft.TokenID,
		&nft.ContractAddress,
		&nft.ChainID,
		&nft.OwnerAddress,
		&nft.MintTxHash,
		&nft.MetadataURI,
		&nft.MintedAt,
		&nft.BurnedAt,
		&nft.BurnTxHash,
		&nft.CreatedAt,
		&nft.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return nft, nil
}

func (r *InvoiceRepository) UpdateNFTOwner(id uuid.UUID, owner string) error {
	query := `UPDATE invoice_nfts SET owner_address = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, owner, time.Now(), id)
	return err
}

func (r *InvoiceRepository) BurnNFT(id uuid.UUID, txHash string) error {
	now := time.Now()
	query := `UPDATE invoice_nfts SET burned_at = $1, burn_tx_hash = $2, updated_at = $1 WHERE id = $3`
	_, err := r.db.Exec(query, now, txHash, id)
	return err
}

func (r *InvoiceRepository) ApproveWithTransaction(id uuid.UUID, interestRate, advanceAmount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	now := time.Now()

	_, err = tx.Exec(`UPDATE invoices SET interest_rate = $1, updated_at = $2 WHERE id = $3`, interestRate, now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`UPDATE invoices SET advance_amount = $1, updated_at = $2 WHERE id = $3`, advanceAmount, now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3`, "approved", now, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// FindAll finds invoices with optional filters
func (r *InvoiceRepository) FindAll(filter *models.InvoiceFilter) ([]models.Invoice, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM invoices WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if filter.Status != nil {
		argCount++
		countQuery += ` AND status = $` + string(rune('0'+argCount))
		args = append(args, *filter.Status)
	}
	if filter.ExporterID != nil {
		argCount++
		countQuery += ` AND exporter_id = $` + string(rune('0'+argCount))
		args = append(args, *filter.ExporterID)
	}

	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	query := `
		SELECT i.id, i.exporter_id, i.buyer_id, i.invoice_number, i.currency, i.amount, i.issue_date, i.due_date,
		       i.description, i.status, i.interest_rate, i.advance_percentage, i.advance_amount, i.document_hash,
		       i.created_at, i.updated_at,
		       b.id, b.company_name, b.country, b.credit_score
		FROM invoices i
		LEFT JOIN buyers b ON i.buyer_id = b.id
		WHERE 1=1
	`

	queryArgs := []interface{}{}
	queryArgCount := 0

	if filter.Status != nil {
		queryArgCount++
		query += ` AND i.status = $` + string(rune('0'+queryArgCount))
		queryArgs = append(queryArgs, *filter.Status)
	}
	if filter.ExporterID != nil {
		queryArgCount++
		query += ` AND i.exporter_id = $` + string(rune('0'+queryArgCount))
		queryArgs = append(queryArgs, *filter.ExporterID)
	}

	queryArgCount++
	query += ` ORDER BY i.created_at DESC LIMIT $` + string(rune('0'+queryArgCount))
	queryArgs = append(queryArgs, filter.PerPage)

	queryArgCount++
	query += ` OFFSET $` + string(rune('0'+queryArgCount))
	queryArgs = append(queryArgs, offset)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var invoice models.Invoice
		var buyer models.Buyer
		if err := rows.Scan(
			&invoice.ID,
			&invoice.ExporterID,
			&invoice.BuyerID,
			&invoice.InvoiceNumber,
			&invoice.Currency,
			&invoice.Amount,
			&invoice.IssueDate,
			&invoice.DueDate,
			&invoice.Description,
			&invoice.Status,
			&invoice.InterestRate,
			&invoice.AdvancePercentage,
			&invoice.AdvanceAmount,
			&invoice.DocumentHash,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&buyer.ID,
			&buyer.CompanyName,
			&buyer.Country,
			&buyer.CreditScore,
		); err != nil {
			return nil, 0, err
		}
		invoice.Buyer = &buyer
		invoices = append(invoices, invoice)
	}

	return invoices, total, nil
}

// CountByExporter counts invoices by exporter
func (r *InvoiceRepository) CountByExporter(exporterID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM invoices WHERE exporter_id = $1`
	err := r.db.QueryRow(query, exporterID).Scan(&count)
	return count, err
}

// CountByBuyerID counts invoices by buyer
func (r *InvoiceRepository) CountByBuyerID(buyerID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM invoices WHERE buyer_id = $1`
	err := r.db.QueryRow(query, buyerID).Scan(&count)
	return count, err
}
