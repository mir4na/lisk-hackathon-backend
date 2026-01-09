package repository

import (
	"database/sql"
	"errors"
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
		INSERT INTO users (email, username, phone_number, password_hash, role, wallet_address, is_verified, is_active, cooperative_agreement, member_status, balance_idr, email_verified)
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
		user.WalletAddress,
		user.IsVerified,
		user.IsActive,
		user.CooperativeAgreement,
		user.MemberStatus,
		user.BalanceIDR,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
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
	if user.Role == models.RoleExporter {
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

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role, wallet_address,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), created_at, updated_at
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
		&user.WalletAddress,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
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
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role, wallet_address,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), created_at, updated_at
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
		&user.WalletAddress,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
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
		SELECT id, email, COALESCE(username, ''), COALESCE(phone_number, ''), password_hash, role, wallet_address,
		       is_verified, is_active, COALESCE(cooperative_agreement, false),
		       COALESCE(member_status, 'calon_anggota_pendana'), COALESCE(balance_idr, 0),
		       COALESCE(email_verified, false), created_at, updated_at
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
		&user.WalletAddress,
		&user.IsVerified,
		&user.IsActive,
		&user.CooperativeAgreement,
		&user.MemberStatus,
		&user.BalanceIDR,
		&user.EmailVerified,
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

func (r *UserRepository) UpdateWallet(userID uuid.UUID, walletAddress string) error {
	query := `UPDATE users SET wallet_address = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, walletAddress, time.Now(), userID)
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

func (r *UserRepository) WalletExists(wallet string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE wallet_address = $1)`
	err := r.db.QueryRow(query, wallet).Scan(&exists)
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
