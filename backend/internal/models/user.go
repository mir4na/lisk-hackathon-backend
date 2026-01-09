package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleInvestor UserRole = "investor"
	RoleAdmin    UserRole = "admin"
	RoleMitra    UserRole = "mitra"
	// Note: Guest is an unregistered user viewing the app, not a role for registered users
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
	IsVerified           bool         `json:"is_verified"`
	IsActive             bool         `json:"is_active"`
	CooperativeAgreement bool         `json:"cooperative_agreement"`
	MemberStatus         MemberStatus `json:"member_status"`
	BalanceIDR           float64      `json:"balance_idr"`
	EmailVerified        bool         `json:"email_verified"`
	ProfileCompleted     bool         `json:"profile_completed"`
	WalletAddress        *string      `json:"wallet_address,omitempty"`
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
	// Account Credentials - Only basic fields for registration
	Email                string   `json:"email" binding:"required,email"`
	Username             string   `json:"username" binding:"required,min=3,max=50"`
	Password             string   `json:"password" binding:"required,min=8"`
	ConfirmPassword      string   `json:"confirm_password" binding:"required,eqfield=Password"`
	Role                 UserRole `json:"role" binding:"required,oneof=investor mitra"`
	CooperativeAgreement bool     `json:"cooperative_agreement" binding:"required"`
	OTPToken             string   `json:"otp_token" binding:"required"`
}

// CompleteProfileRequest is used to complete user profile after registration
// User must complete profile before accessing most features
type CompleteProfileRequest struct {
	// Required for profile completion
	FullName string  `json:"full_name" binding:"required,min=3"`
	Phone    *string `json:"phone,omitempty"`

	// Optional KYC - can be added later
	NIK         *string `json:"nik,omitempty"`
	KTPPhotoURL *string `json:"ktp_photo_url,omitempty"`
	SelfieURL   *string `json:"selfie_url,omitempty"`

	// Optional Bank Account - can be added later
	BankCode      *string `json:"bank_code,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
	AccountName   *string `json:"account_name,omitempty"`

	// Mitra specific
	CompanyName *string `json:"company_name,omitempty"`
	Country     *string `json:"country,omitempty"`
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

// ==================== Bank Account Models ====================

// SupportedBank represents a bank that users can use for disbursement
type SupportedBank struct {
	Code string `json:"code"` // e.g., "bca", "mandiri", "bni"
	Name string `json:"name"` // e.g., "Bank Central Asia"
}

// GetSupportedBanks returns list of supported banks for disbursement
func GetSupportedBanks() []SupportedBank {
	return []SupportedBank{
		{Code: "bca", Name: "Bank Central Asia (BCA)"},
		{Code: "mandiri", Name: "Bank Mandiri"},
		{Code: "bni", Name: "Bank Negara Indonesia (BNI)"},
		{Code: "bri", Name: "Bank Rakyat Indonesia (BRI)"},
		{Code: "cimb", Name: "CIMB Niaga"},
		{Code: "danamon", Name: "Bank Danamon"},
		{Code: "permata", Name: "Bank Permata"},
		{Code: "bsi", Name: "Bank Syariah Indonesia (BSI)"},
		{Code: "btn", Name: "Bank Tabungan Negara (BTN)"},
		{Code: "ocbc", Name: "OCBC NISP"},
	}
}

// BankAccount represents user's bank account for disbursement
type BankAccount struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	BankCode      string     `json:"bank_code"`
	BankName      string     `json:"bank_name"`
	AccountNumber string     `json:"account_number"`
	AccountName   string     `json:"account_name"`
	IsVerified    bool       `json:"is_verified"`
	IsPrimary     bool       `json:"is_primary"` // Primary account for disbursement
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
}

// Microcopy for bank account
const BankAccountMicrocopy = "Rekening ini akan menjadi satu-satunya tujuan pencairan dana demi keamanan. Kamu bisa mengubahnya nanti di bagian profile."

// ChangeBankAccountRequest is used to change bank account (requires OTP verification)
type ChangeBankAccountRequest struct {
	OTPToken      string `json:"otp_token" binding:"required"`      // OTP verification for security
	BankCode      string `json:"bank_code" binding:"required"`      // New bank code
	AccountNumber string `json:"account_number" binding:"required"` // New account number
	AccountName   string `json:"account_name" binding:"required"`   // New account holder name
}

// ChangePasswordRequest for security section
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ==================== KYC/Identity Models ====================

// UserIdentity stores KYC identity data (from registration)
type UserIdentity struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	NIK         string     `json:"nik"`                     // 16-digit NIK (masked in response)
	FullName    string     `json:"full_name"`               // Name as per KTP
	KTPPhotoURL string     `json:"ktp_photo_url,omitempty"` // URL of KTP photo (hidden in profile view)
	SelfieURL   string     `json:"selfie_url,omitempty"`    // URL of selfie (hidden in profile view)
	IsVerified  bool       `json:"is_verified"`             // KYC verification status
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// MaskNIK masks NIK for display (show first 6 and last 4 digits)
func (u *UserIdentity) MaskNIK() string {
	if len(u.NIK) != 16 {
		return "****"
	}
	return u.NIK[:6] + "******" + u.NIK[12:]
}

// ProfileDataResponse for "Data Diri" section (read-only)
type ProfileDataResponse struct {
	FullName     string `json:"full_name"`     // From KTP
	NIKMasked    string `json:"nik_masked"`    // Masked NIK (e.g., 320101******1234)
	Email        string `json:"email"`         // Read-only
	Phone        string `json:"phone"`         // Added Phone
	Username     string `json:"username"`      // Read-only
	MemberStatus string `json:"member_status"` // Member status
	Role         string `json:"role"`          // investor/mitra
	IsVerified   bool   `json:"is_verified"`   // KYC verified
	JoinedAt     string `json:"joined_at"`     // Registration date
}

// BankAccountResponse for "Rekening Bank" section
type BankAccountResponse struct {
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"` // Partially masked
	AccountName   string `json:"account_name"`
	IsPrimary     bool   `json:"is_primary"`
	IsVerified    bool   `json:"is_verified"`
	Microcopy     string `json:"microcopy"`
}

// MaskAccountNumber masks bank account number for display
func MaskAccountNumber(number string) string {
	if len(number) <= 4 {
		return number
	}
	return "****" + number[len(number)-4:]
}
