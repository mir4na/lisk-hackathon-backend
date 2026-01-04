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
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Role          UserRole   `json:"role"`
	WalletAddress *string    `json:"wallet_address,omitempty"`
	IsVerified    bool       `json:"is_verified"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Profile       *UserProfile `json:"profile,omitempty"`
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
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=8"`
	Role        UserRole `json:"role" binding:"required,oneof=exporter investor"`
	FullName    string   `json:"full_name" binding:"required"`
	CompanyName *string  `json:"company_name,omitempty"`
	Country     *string  `json:"country,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
