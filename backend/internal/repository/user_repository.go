package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vessel/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User, profile *models.UserProfile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (email, username, phone_number, password_hash, role, is_verified, is_active, cooperative_agreement, member_status, balance_idr, email_verified, profile_completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(
		query,
		user.Email,
		user.Username,
		user.PhoneNumber,
		user.PasswordHash,
		user.Role,
		user.IsVerified,
		user.IsActive,
		user.CooperativeAgreement,
		user.MemberStatus,
		user.BalanceIDR,
		user.EmailVerified,
		user.ProfileCompleted,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	// Handle nil profile (initialize default)
	if profile == nil {
		profile = &models.UserProfile{
			FullName: "",
			// Other fields are pointers (*string), so they default to nil
			// which is correct for database NULL
		}
	}

	profileQuery := `
		INSERT INTO user_profiles (user_id, full_name, phone, country, company_name, company_type, business_sector)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	profile.UserID = user.ID
	err = tx.QueryRow(
		profileQuery,
		profile.UserID,
		profile.FullName,
		profile.Phone,
		profile.Country,
		profile.CompanyName,
		profile.CompanyType,
		profile.BusinessSector,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return err
	}

	// Create initial credit score for exporters
	if user.Role == models.RoleMitra {
		creditQuery := `
			INSERT INTO credit_scores (user_id, score)
			VALUES ($1, 50)
		`
		_, err = tx.Exec(creditQuery, user.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *UserRepository) CompleteUserRegistration(userID uuid.UUID, profile *models.UserProfile, identity *models.UserIdentity, bankAccount *models.BankAccount) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update Profile
	profileQuery := `
		UPDATE user_profiles
		SET full_name = $1, phone = $2, country = $3, company_name = $4, updated_at = $5
		WHERE user_id = $6
		RETURNING id
	`
	err = tx.QueryRow(
		profileQuery,
		profile.FullName,
		profile.Phone,
		profile.Country,
		profile.CompanyName,
		time.Now(),
		userID,
	).Scan(&profile.ID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// 2. Create Identity (KYC)
	identity.ID = uuid.New()
	now := time.Now()
	identity.CreatedAt = now
	identity.UpdatedAt = now
	if identity.IsVerified {
		identity.VerifiedAt = &now
	}

	identityQuery := `
		INSERT INTO user_identities (id, user_id, nik, full_name, ktp_photo_url, selfie_url, is_verified, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id) DO UPDATE SET
			nik = EXCLUDED.nik,
			full_name = EXCLUDED.full_name,
			ktp_photo_url = EXCLUDED.ktp_photo_url,
			selfie_url = EXCLUDED.selfie_url,
			updated_at = EXCLUDED.updated_at,
			is_verified = EXCLUDED.is_verified,
			verified_at = EXCLUDED.verified_at
	`
	_, err = tx.Exec(identityQuery,
		identity.ID, userID, identity.NIK, identity.FullName,
		identity.KTPPhotoURL, identity.SelfieURL, identity.IsVerified,
		identity.VerifiedAt, identity.CreatedAt, identity.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create identity: %w", err)
	}

	// 3. Create Bank Account
	bankAccount.ID = uuid.New()
	bankAccount.CreatedAt = now
	bankAccount.UpdatedAt = now
	if bankAccount.IsVerified {
		bankAccount.VerifiedAt = &now
	}
	// Ensure it's primary
	bankAccount.IsPrimary = true

	bankQuery := `
		INSERT INTO bank_accounts (id, user_id, bank_code, bank_name, account_number, account_name, is_verified, is_primary, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(bankQuery,
		bankAccount.ID, userID, bankAccount.BankCode, bankAccount.BankName,
		bankAccount.AccountNumber, bankAccount.AccountName, bankAccount.IsVerified,
		bankAccount.IsPrimary, bankAccount.VerifiedAt, bankAccount.CreatedAt, bankAccount.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	// 4. Mark user as verified (if auto-verify logic applies, or just rely on KYC status)
	// For MVP: Set IsVerified = true if KYC is verified
	if identity.IsVerified {
		userUpdateQuery := `UPDATE users SET is_verified = true, profile_completed = true, updated_at = $1 WHERE id = $2`
		_, err = tx.Exec(userUpdateQuery, time.Now(), userID)
		if err != nil {
			return fmt.Errorf("failed to update user verification: %w", err)
		}
	}

	return tx.Commit()
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), COALESCE(profile_completed, false), created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`
	var username, phoneNumber string
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&username,
		&phoneNumber,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
		&user.ProfileCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if username != "" {
		user.Username = &username
	}
	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), COALESCE(profile_completed, false), created_at, updated_at
		FROM users
		WHERE username = $1 AND is_active = true
	`
	var uname, phoneNumber string
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Email,
		&uname,
		&phoneNumber,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
		&user.ProfileCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if uname != "" {
		user.Username = &uname
	}
	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}
	return user, nil
}

func (r *UserRepository) FindByEmailOrUsername(identifier string) (*models.User, error) {
	// Try email first
	user, err := r.FindByEmail(identifier)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	// Try username
	return r.FindByUsername(identifier)
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), COALESCE(profile_completed, false), created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var username, phoneNumber string
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&username,
		&phoneNumber,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
		&user.ProfileCompleted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if username != "" {
		user.Username = &username
	}
	if phoneNumber != "" {
		user.PhoneNumber = &phoneNumber
	}

	profile, _ := r.FindProfileByUserID(id)
	user.Profile = profile

	return user, nil
}

func (r *UserRepository) FindProfileByUserID(userID uuid.UUID) (*models.UserProfile, error) {
	profile := &models.UserProfile{}
	query := `
		SELECT id, user_id, full_name, phone, country, company_name, company_type, business_sector, avatar_url, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.Phone,
		&profile.Country,
		&profile.CompanyName,
		&profile.CompanyType,
		&profile.BusinessSector,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return profile, nil
}

func (r *UserRepository) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error {
	query := `
		UPDATE user_profiles
		SET full_name = $1, phone = $2, country = $3, company_name = $4, company_type = $5, business_sector = $6, updated_at = $7
		WHERE user_id = $8
	`
	_, err := r.db.Exec(
		query,
		req.FullName,
		req.Phone,
		req.Country,
		req.CompanyName,
		req.CompanyType,
		req.BusinessSector,
		time.Now(),
		userID,
	)
	return err
}

func (r *UserRepository) SetVerified(userID uuid.UUID, verified bool) error {
	query := `UPDATE users SET is_verified = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, verified, time.Now(), userID)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateBalance(userID uuid.UUID, amount float64) error {
	query := `UPDATE users SET balance_idr = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, amount, time.Now(), userID)
	return err
}

func (r *UserRepository) SetEmailVerified(userID uuid.UUID, verified bool) error {
	query := `UPDATE users SET email_verified = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, verified, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateMemberStatus(userID uuid.UUID, status models.MemberStatus) error {
	query := `UPDATE users SET member_status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateRole(userID uuid.UUID, role models.UserRole) error {
	query := `UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, role, time.Now(), userID)
	return err
}

// ==================== Identity/KYC Methods ====================

func (r *UserRepository) CreateIdentity(identity *models.UserIdentity) error {
	identity.ID = uuid.New()
	now := time.Now()
	identity.CreatedAt = now
	identity.UpdatedAt = now
	if identity.IsVerified {
		identity.VerifiedAt = &now
	}

	query := `
		INSERT INTO user_identities (id, user_id, nik, full_name, ktp_photo_url, selfie_url, is_verified, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query,
		identity.ID, identity.UserID, identity.NIK, identity.FullName,
		identity.KTPPhotoURL, identity.SelfieURL, identity.IsVerified,
		identity.VerifiedAt, identity.CreatedAt, identity.UpdatedAt,
	)
	return err
}

func (r *UserRepository) FindIdentityByUserID(userID uuid.UUID) (*models.UserIdentity, error) {
	identity := &models.UserIdentity{}
	query := `
		SELECT id, user_id, nik, full_name, ktp_photo_url, selfie_url, is_verified, verified_at, created_at, updated_at
		FROM user_identities
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&identity.ID, &identity.UserID, &identity.NIK, &identity.FullName,
		&identity.KTPPhotoURL, &identity.SelfieURL, &identity.IsVerified,
		&identity.VerifiedAt, &identity.CreatedAt, &identity.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return identity, nil
}

func (r *UserRepository) UpdateIdentity(identity *models.UserIdentity) error {
	identity.UpdatedAt = time.Now()
	query := `
		UPDATE user_identities
		SET nik = $1, full_name = $2, ktp_photo_url = $3, selfie_url = $4, is_verified = $5, verified_at = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(query,
		identity.NIK, identity.FullName, identity.KTPPhotoURL, identity.SelfieURL,
		identity.IsVerified, identity.VerifiedAt, identity.UpdatedAt, identity.ID,
	)
	return err
}

// ==================== Bank Account Methods ====================

func (r *UserRepository) CreateBankAccount(account *models.BankAccount) error {
	account.ID = uuid.New()
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now
	if account.IsVerified {
		account.VerifiedAt = &now
	}

	query := `
		INSERT INTO bank_accounts (id, user_id, bank_code, bank_name, account_number, account_name, is_verified, is_primary, verified_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		account.ID, account.UserID, account.BankCode, account.BankName,
		account.AccountNumber, account.AccountName, account.IsVerified,
		account.IsPrimary, account.VerifiedAt, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

func (r *UserRepository) FindBankAccountsByUserID(userID uuid.UUID) ([]models.BankAccount, error) {
	query := `
		SELECT id, user_id, bank_code, bank_name, account_number, account_name, is_verified, is_primary, verified_at, created_at, updated_at
		FROM bank_accounts
		WHERE user_id = $1
		ORDER BY is_primary DESC, created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.BankAccount
	for rows.Next() {
		var account models.BankAccount
		if err := rows.Scan(
			&account.ID, &account.UserID, &account.BankCode, &account.BankName,
			&account.AccountNumber, &account.AccountName, &account.IsVerified,
			&account.IsPrimary, &account.VerifiedAt, &account.CreatedAt, &account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *UserRepository) FindPrimaryBankAccount(userID uuid.UUID) (*models.BankAccount, error) {
	account := &models.BankAccount{}
	query := `
		SELECT id, user_id, bank_code, bank_name, account_number, account_name, is_verified, is_primary, verified_at, created_at, updated_at
		FROM bank_accounts
		WHERE user_id = $1 AND is_primary = true
		LIMIT 1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&account.ID, &account.UserID, &account.BankCode, &account.BankName,
		&account.AccountNumber, &account.AccountName, &account.IsVerified,
		&account.IsPrimary, &account.VerifiedAt, &account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}

func (r *UserRepository) UpdateBankAccount(account *models.BankAccount) error {
	account.UpdatedAt = time.Now()
	query := `
		UPDATE bank_accounts
		SET bank_code = $1, bank_name = $2, account_number = $3, account_name = $4, is_verified = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(query,
		account.BankCode, account.BankName, account.AccountNumber, account.AccountName,
		account.IsVerified, account.UpdatedAt, account.ID,
	)
	return err
}

func (r *UserRepository) SetPrimaryBankAccount(userID, accountID uuid.UUID) error {
	// First, unset all primary flags
	_, err := r.db.Exec(`UPDATE bank_accounts SET is_primary = false, updated_at = $1 WHERE user_id = $2`, time.Now(), userID)
	if err != nil {
		return err
	}
	// Set the new primary
	_, err = r.db.Exec(`UPDATE bank_accounts SET is_primary = true, updated_at = $1 WHERE id = $2`, time.Now(), accountID)
	return err
}

// ==================== Password Methods ====================

// ==================== Password Methods ====================

func (r *UserRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

// ==================== Wallet Methods ====================

func (r *UserRepository) UpdateWalletAddress(userID uuid.UUID, walletAddress string) error {
	query := `UPDATE users SET wallet_address = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, walletAddress, time.Now(), userID)
	return err
}
