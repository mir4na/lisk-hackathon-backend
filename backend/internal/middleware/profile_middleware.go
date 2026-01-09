package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vessel/backend/internal/repository"
	"github.com/vessel/backend/internal/utils"
)

// ProfileMiddleware handles profile-related middleware
type ProfileMiddleware struct {
	userRepo repository.UserRepositoryInterface
}

// NewProfileMiddleware creates a new ProfileMiddleware instance
func NewProfileMiddleware(userRepo repository.UserRepositoryInterface) *ProfileMiddleware {
	return &ProfileMiddleware{userRepo: userRepo}
}

// RequireProfileComplete ensures the user has completed their profile before accessing protected routes
// Users must have a profile with full_name set before they can use most features
func (m *ProfileMiddleware) RequireProfileComplete() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uuid.UUID)

		// Check if profile exists and is complete
		profile, err := m.userRepo.FindProfileByUserID(userID)
		if err != nil {
			utils.InternalServerError(c, "Failed to verify profile status")
			c.Abort()
			return
		}

		// Profile must exist and have full_name
		if profile == nil || profile.FullName == "" {
			utils.ForbiddenError(c, "Please complete your profile first. Update your profile with full_name to continue.")
			c.Abort()
			return
		}

		c.Next()
	}
}
