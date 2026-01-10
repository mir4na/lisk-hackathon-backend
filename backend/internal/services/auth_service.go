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
	if s.otpService != nil && !s.otpService.ValidateOTPToken(req.OTPToken, req.Email) {
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

	// Validate role selection (only investor or mitra)
	role := req.Role
	if role != models.RoleInvestor && role != models.RoleMitra {
		return nil, errors.New("invalid role, must be investor or mitra")
	}

	// Determine member status based on role
	memberStatus := models.MemberStatusCalonAnggotaPendana
	if role == models.RoleMitra {
		memberStatus = models.MemberStatusCalonAnggotaMitra
	}

	// Create user without profile (profile will be completed later)
	username := req.Username
	user := &models.User{
		Email:                req.Email,
		Username:             &username,
		PhoneNumber:          nil,
		PasswordHash:         hashedPassword,
		Role:                 role,
		IsVerified:           false, // Not verified until profile is completed
		IsActive:             true,
		CooperativeAgreement: req.CooperativeAgreement,
		MemberStatus:         memberStatus,
		BalanceIDR:           0,
		EmailVerified:        true, // Email verified via OTP
	}

	// Create user without profile - profile will be created during profile completion
	if err := s.userRepo.Create(user, nil); err != nil {
		return nil, err
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

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600 * 24, // 24 hours in seconds
	}, nil
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
