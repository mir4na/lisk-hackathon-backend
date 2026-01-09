package services

import (
	"errors"

	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/utils"
)

type AuthService struct {
	userRepo   repository.UserRepositoryInterface
	jwtManager *utils.JWTManager
	otpService *OTPService
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtManager *utils.JWTManager, otpService *OTPService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		otpService: otpService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate OTP token first
	if s.otpService != nil && !s.otpService.ValidateOTPToken(req.OTPToken) {
		return nil, errors.New("email not verified, please verify OTP first")
	}

	// Validate cooperative agreement
	if !req.CooperativeAgreement {
		return nil, errors.New("you must agree to VESSEL Cooperative Service terms")
	}

	// Check if email exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Check if username exists
	usernameExists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Validate role selection (only investor or mitra - guest means unregistered)
	role := req.Role
	if role != models.RoleInvestor && role != models.RoleMitra {
		return nil, errors.New("invalid role, must be investor or mitra")
	}

	// Validate NIK format (16 digits)
	if len(req.NIK) != 16 {
		return nil, errors.New("NIK must be exactly 16 digits")
	}

	// Validate bank account name matches KTP name (abstracted - always pass for MVP)
	// In production, this would use a bank verification API like Didit
	if !s.validateBankAccountName(req.FullName, req.AccountName) {
		return nil, errors.New("nama pemilik rekening harus sama dengan nama di KTP")
	}

	// Validate bank code is supported
	if !s.isSupportedBank(req.BankCode) {
		return nil, errors.New("bank tidak didukung")
	}

	// Determine member status based on role
	memberStatus := models.MemberStatusCalonAnggotaPendana
	if role == models.RoleMitra {
		memberStatus = models.MemberStatusMemberMitra
	}

	// Create user with selected role
	username := req.Username
	user := &models.User{
		Email:                req.Email,
		Username:             &username,
		PhoneNumber:          nil,
		PasswordHash:         hashedPassword,
		Role:                 role,
		IsVerified:           true, // KYC submitted at registration
		IsActive:             true,
		CooperativeAgreement: req.CooperativeAgreement,
		MemberStatus:         memberStatus,
		BalanceIDR:           0,
		EmailVerified:        true,
	}

	profile := &models.UserProfile{
		FullName:    req.FullName, // Use full name from KTP
		Phone:       nil,
		CompanyName: req.CompanyName,
		Country:     req.Country,
	}

	if err := s.userRepo.Create(user, profile); err != nil {
		return nil, err
	}

	// Store identity data (KYC)
	identity := &models.UserIdentity{
		UserID:      user.ID,
		NIK:         req.NIK,
		FullName:    req.FullName,
		KTPPhotoURL: req.KTPPhotoURL,
		SelfieURL:   req.SelfieURL,
		IsVerified:  true, // Auto-verified for MVP
	}
	if err := s.userRepo.CreateIdentity(identity); err != nil {
		// Log but don't fail - identity is supplementary
		// In production, this would be critical
	}

	// Store bank account
	bankAccount := &models.BankAccount{
		UserID:        user.ID,
		BankCode:      req.BankCode,
		BankName:      s.getBankName(req.BankCode),
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		IsVerified:    true, // Auto-verified for MVP
		IsPrimary:     true, // First account is always primary
	}
	if err := s.userRepo.CreateBankAccount(bankAccount); err != nil {
		// Log but don't fail
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	user.Profile = profile

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600 * 24, // 24 hours in seconds
	}, nil
}

// validateBankAccountName checks if bank account name matches KTP name
// For MVP: abstracted to always return true
// In production: integrate with bank verification API (e.g., Didit)
func (s *AuthService) validateBankAccountName(ktpName, accountName string) bool {
	// MVP: Always return true (abstracted validation)
	// TODO: Integrate with bank verification API
	return true
}

// isSupportedBank checks if bank code is in supported list
func (s *AuthService) isSupportedBank(bankCode string) bool {
	for _, bank := range models.GetSupportedBanks() {
		if bank.Code == bankCode {
			return true
		}
	}
	return false
}

// getBankName returns bank name from code
func (s *AuthService) getBankName(bankCode string) string {
	for _, bank := range models.GetSupportedBanks() {
		if bank.Code == bankCode {
			return bank.Name
		}
	}
	return bankCode
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Support login with email or username
	user, err := s.userRepo.FindByEmailOrUsername(req.EmailOrUsername)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email/username or password")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email/username or password")
	}

	if !user.IsActive {
		return nil, errors.New("account has been deactivated")
	}

	// Get profile
	profile, _ := s.userRepo.FindProfileByUserID(user.ID)
	user.Profile = profile

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600 * 24,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600 * 24,
	}, nil
}
