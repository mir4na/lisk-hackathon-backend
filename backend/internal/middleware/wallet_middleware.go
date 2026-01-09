package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/repository"
)

// WalletMiddleware - DEPRECATED
// Wallet connection is no longer required. Users register bank accounts instead.
// This middleware is kept for backwards compatibility but does nothing.
type WalletMiddleware struct {
	userRepo repository.UserRepositoryInterface
}

func NewWalletMiddleware(userRepo repository.UserRepositoryInterface) *WalletMiddleware {
	return &WalletMiddleware{userRepo: userRepo}
}

// RequireWallet - DEPRECATED
// Wallet is no longer required. This middleware now just passes through.
// Kept for backwards compatibility in case any code still references it.
func (m *WalletMiddleware) RequireWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No-op: wallet no longer required
		// Users now have bank accounts registered at signup for receiving funds
		c.Next()
	}
}
