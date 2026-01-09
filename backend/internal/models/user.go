package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleExporter UserRole = "exporter"
	RoleInvestor UserRole = "investor"
	RoleAdmin    UserRole = "admin"
	RoleMitra    UserRole = "mitra"
)

type MemberStatus string

const (
	MemberStatusCalonAnggotaPendana MemberStatus = "calon_anggota_pendana"
	MemberStatusMemberMitra         MemberStatus = "member_mitra"
	MemberStatusAdmin               MemberStatus = "admin"
)

type User struct {
	ID                   uuid.UUID    `json:"id"`
	Email                string       `json:"email"`
	Username             *string      `json:"username,omitempty"`
	PhoneNumber          *string      `json:"phone_number,omitempty"`
	PasswordHash         string       `json:"-"`
	Role                 UserRole     `json:"role"`
	WalletAddress        *string      `json:"wallet_address,omitempty"`
	IsVerified           bool         `json:"is_verified"`
	IsActive             bool         `json:"is_active"`
	CooperativeAgreement bool         `json:"cooperative_agreement"`
	MemberStatus         MemberStatus `json:"member_status"`
	BalanceIDR           float64      `json:"balance_idr"`
	EmailVerified        bool         `json:"email_verified"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	Profile              *UserProfile `json:"profile,omitempty"`
}

type UserProfile struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	FullName       string    `json:"full_name"`
	Phone          *string   `json:"phone,omitempty"`
	Country        *string   `json:"country,omitempty"`
	CompanyName    *string   `json:"company_name,omitempty"`
	CompanyType    *string   `json:"company_type,omitempty"`
	BusinessSector *string   `json:"business_sector,omitempty"`
	AvatarURL      *string   `json:"avatar_url,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Email                string   `json:"email" binding:"required,email"`
	Username             string   `json:"username" binding:"required,min=3,max=50"`
	Password             string   `json:"password" binding:"required,min=8"`
	ConfirmPassword      string   `json:"confirm_password" binding:"required,eqfield=Password"`
	Role                 UserRole `json:"role" binding:"required,oneof=exporter investor"`
	FullName             string   `json:"full_name" binding:"required"`
	PhoneNumber          string   `json:"phone_number" binding:"required"`
	CooperativeAgreement bool     `json:"cooperative_agreement" binding:"required"`
	OTPToken             string   `json:"otp_token" binding:"required"` // Token from OTP verification
	CompanyName          *string  `json:"company_name,omitempty"`
	Country              *string  `json:"country,omitempty"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"` // Can be email or username
	Password        string `json:"password" binding:"required"`
}

// UserBalanceResponse represents user balance info
type UserBalanceResponse struct {
	BalanceIDR float64 `json:"balance_idr"`
	Currency   string  `json:"currency"`
}

type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UpdateProfileRequest struct {
	FullName       string  `json:"full_name"`
	Phone          *string `json:"phone,omitempty"`
	Country        *string `json:"country,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	CompanyType    *string `json:"company_type,omitempty"`
	BusinessSector *string `json:"business_sector,omitempty"`
}

type UpdateWalletRequest struct {
	WalletAddress string `json:"wallet_address" binding:"required,len=42"`
}
