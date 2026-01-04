package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Enable UUID extension
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL CHECK (role IN ('exporter', 'investor', 'admin')),
			wallet_address VARCHAR(42) UNIQUE,
			is_verified BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// User profiles table
		`CREATE TABLE IF NOT EXISTS user_profiles (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			full_name VARCHAR(255) NOT NULL,
			phone VARCHAR(20),
			country VARCHAR(100),
			company_name VARCHAR(255),
			company_type VARCHAR(100),
			business_sector VARCHAR(100),
			avatar_url TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// KYC verifications table
		`CREATE TABLE IF NOT EXISTS kyc_verifications (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			verification_type VARCHAR(10) NOT NULL CHECK (verification_type IN ('kyc', 'kyb')),
			status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
			id_type VARCHAR(50),
			id_number VARCHAR(100),
			id_document_url TEXT,
			selfie_url TEXT,
			rejection_reason TEXT,
			verified_by UUID REFERENCES users(id),
			verified_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Buyers table
		`CREATE TABLE IF NOT EXISTS buyers (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			created_by UUID REFERENCES users(id) ON DELETE CASCADE,
			company_name VARCHAR(255) NOT NULL,
			country VARCHAR(100) NOT NULL,
			address TEXT,
			contact_email VARCHAR(255),
			contact_phone VARCHAR(50),
			website VARCHAR(255),
			credit_score INTEGER DEFAULT 50,
			total_invoices INTEGER DEFAULT 0,
			total_paid DECIMAL(20,2) DEFAULT 0,
			total_defaulted DECIMAL(20,2) DEFAULT 0,
			is_verified BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Invoices table
		`CREATE TABLE IF NOT EXISTS invoices (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			exporter_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
			buyer_id UUID REFERENCES buyers(id) NOT NULL,
			invoice_number VARCHAR(100) NOT NULL,
			currency VARCHAR(10) DEFAULT 'USD',
			amount DECIMAL(20,2) NOT NULL,
			issue_date DATE NOT NULL,
			due_date DATE NOT NULL,
			description TEXT,
			status VARCHAR(30) NOT NULL DEFAULT 'draft' CHECK (status IN (
				'draft', 'pending_review', 'approved', 'rejected',
				'tokenized', 'funding', 'funded', 'matured', 'repaid', 'defaulted'
			)),
			interest_rate DECIMAL(5,2),
			advance_percentage DECIMAL(5,2) DEFAULT 80.00,
			advance_amount DECIMAL(20,2),
			document_hash VARCHAR(66),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Invoice documents table
		`CREATE TABLE IF NOT EXISTS invoice_documents (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
			document_type VARCHAR(30) NOT NULL CHECK (document_type IN (
				'invoice_pdf', 'bill_of_lading', 'packing_list',
				'certificate_of_origin', 'insurance', 'customs', 'other'
			)),
			file_name VARCHAR(255),
			file_url TEXT,
			file_hash VARCHAR(66),
			file_size INTEGER,
			uploaded_at TIMESTAMP DEFAULT NOW()
		);`,

		// Invoice NFTs table
		`CREATE TABLE IF NOT EXISTS invoice_nfts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			invoice_id UUID UNIQUE REFERENCES invoices(id) ON DELETE CASCADE,
			token_id BIGINT,
			contract_address VARCHAR(42),
			chain_id INTEGER,
			owner_address VARCHAR(42),
			mint_tx_hash VARCHAR(66),
			metadata_uri TEXT,
			minted_at TIMESTAMP,
			burned_at TIMESTAMP,
			burn_tx_hash VARCHAR(66),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Funding pools table
		`CREATE TABLE IF NOT EXISTS funding_pools (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			invoice_id UUID UNIQUE REFERENCES invoices(id) ON DELETE CASCADE,
			target_amount DECIMAL(20,2) NOT NULL,
			funded_amount DECIMAL(20,2) DEFAULT 0,
			investor_count INTEGER DEFAULT 0,
			status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'filled', 'disbursed', 'closed')),
			opened_at TIMESTAMP,
			filled_at TIMESTAMP,
			disbursed_at TIMESTAMP,
			closed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Investments table
		`CREATE TABLE IF NOT EXISTS investments (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			pool_id UUID REFERENCES funding_pools(id) ON DELETE CASCADE,
			investor_id UUID REFERENCES users(id) ON DELETE CASCADE,
			amount DECIMAL(20,2) NOT NULL,
			expected_return DECIMAL(20,2),
			actual_return DECIMAL(20,2),
			status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'repaid', 'defaulted')),
			tx_hash VARCHAR(66),
			invested_at TIMESTAMP DEFAULT NOW(),
			repaid_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Transactions table
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			invoice_id UUID REFERENCES invoices(id),
			user_id UUID REFERENCES users(id),
			type VARCHAR(30) NOT NULL CHECK (type IN (
				'investment', 'advance_payment', 'buyer_repayment',
				'investor_return', 'platform_fee', 'refund'
			)),
			amount DECIMAL(20,2) NOT NULL,
			currency VARCHAR(10) DEFAULT 'USDC',
			tx_hash VARCHAR(66),
			status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'failed')),
			from_address VARCHAR(42),
			to_address VARCHAR(42),
			block_number BIGINT,
			gas_used BIGINT,
			notes TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Credit scores table
		`CREATE TABLE IF NOT EXISTS credit_scores (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			score INTEGER DEFAULT 50,
			total_invoices INTEGER DEFAULT 0,
			successful_invoices INTEGER DEFAULT 0,
			defaulted_invoices INTEGER DEFAULT 0,
			total_volume DECIMAL(20,2) DEFAULT 0,
			avg_payment_delay INTEGER DEFAULT 0,
			last_updated TIMESTAMP DEFAULT NOW(),
			created_at TIMESTAMP DEFAULT NOW()
		);`,

		// Credit score history table
		`CREATE TABLE IF NOT EXISTS credit_score_history (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			previous_score INTEGER,
			new_score INTEGER,
			reason VARCHAR(255),
			invoice_id UUID REFERENCES invoices(id),
			created_at TIMESTAMP DEFAULT NOW()
		);`,

		// Notifications table
		`CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			type VARCHAR(50),
			title VARCHAR(255) NOT NULL,
			message TEXT,
			data JSONB,
			is_read BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT NOW()
		);`,

		// Audit logs table
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID REFERENCES users(id),
			action VARCHAR(100),
			entity_type VARCHAR(50),
			entity_id UUID,
			old_data JSONB,
			new_data JSONB,
			ip_address INET,
			user_agent TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);`,

		// Shipment verifications table (Oracle)
		`CREATE TABLE IF NOT EXISTS shipment_verifications (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
			container_id VARCHAR(50),
			bill_of_lading_number VARCHAR(100),
			carrier VARCHAR(100),
			origin_port VARCHAR(100),
			destination_port VARCHAR(100),
			status VARCHAR(50),
			verified_at TIMESTAMP,
			api_response JSONB,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_wallet ON users(wallet_address);`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`,
		`CREATE INDEX IF NOT EXISTS idx_invoices_exporter ON invoices(exporter_id);`,
		`CREATE INDEX IF NOT EXISTS idx_invoices_buyer ON invoices(buyer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);`,
		`CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);`,
		`CREATE INDEX IF NOT EXISTS idx_investments_pool ON investments(pool_id);`,
		`CREATE INDEX IF NOT EXISTS idx_investments_investor ON investments(investor_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_invoice ON transactions(invoice_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash ON transactions(tx_hash);`,
		`CREATE INDEX IF NOT EXISTS idx_nfts_token_id ON invoice_nfts(token_id);`,
		`CREATE INDEX IF NOT EXISTS idx_nfts_owner ON invoice_nfts(owner_address);`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration %d failed: %v", i+1, err)
			return err
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}
