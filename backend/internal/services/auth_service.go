package services

import (
	"errors"

	"github.com/receiv3/backend/internal/models"
	"github.com/receiv3/backend/internal/repository"
	"github.com/receiv3/backend/internal/utils"
)

type AuthService struct {
	userRepo   repository.UserRepositoryInterface
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		IsVerified:   false,
		IsActive:     true,
	}

	profile := &models.UserProfile{
		FullName:    req.FullName,
		CompanyName: req.CompanyName,
		Country:     req.Country,
	}

	if err := s.userRepo.Create(user, profile); err != nil {
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

	user.Profile = profile

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600 * 24, // 24 hours in seconds
	}, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
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
