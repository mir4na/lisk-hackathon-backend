package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/utils"
)

// WalletMiddleware creates middleware that requires wallet connection
type WalletMiddleware struct {
	userRepo repository.UserRepositoryInterface
}

func NewWalletMiddleware(userRepo repository.UserRepositoryInterface) *WalletMiddleware {
	return &WalletMiddleware{userRepo: userRepo}
}

// RequireWallet checks if user has connected wallet
// Must be used AFTER AuthMiddleware
func (m *WalletMiddleware) RequireWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.UnauthorizedError(c, "User not authenticated")
			c.Abort()
			return
		}

		user, err := m.userRepo.FindByID(userID.(uuid.UUID))
		if err != nil {
			utils.InternalServerError(c, "Failed to get user")
			c.Abort()
			return
		}

		if user == nil {
			utils.NotFoundError(c, "User not found")
			c.Abort()
			return
		}

		// Check if wallet is connected
		if user.WalletAddress == nil || *user.WalletAddress == "" {
			utils.ForbiddenError(c, "Wallet connection required. Please connect your wallet first via PUT /api/v1/user/wallet")
			c.Abort()
			return
		}

		// Add wallet address to context for later use
		c.Set("wallet_address", *user.WalletAddress)
		c.Next()
	}
}
