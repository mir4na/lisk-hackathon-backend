package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTManager(secret string, accessExpiry, refreshExpiry int) *JWTManager {
	return &JWTManager{
		secretKey:          secret,
		accessTokenExpiry:  time.Duration(accessExpiry) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshExpiry) * time.Hour,
	}
}

func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vessel",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vessel-refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
