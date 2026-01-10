package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/utils"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedError(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedError(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			utils.UnauthorizedError(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.UnauthorizedError(c, "User role not found")
			c.Abort()
			return
		}

		roleStr := role.(string)
		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		utils.ForbiddenError(c, "You don't have permission to access this resource")
		c.Abort()
	}
}

func ExporterOnly() gin.HandlerFunc {
	return RoleMiddleware("exporter", "mitra", "admin")
}

func InvestorOnly() gin.HandlerFunc {
	return RoleMiddleware("investor", "admin")
}

func MitraOnly() gin.HandlerFunc {
	return RoleMiddleware("mitra", "admin")
}

func GuestOnly() gin.HandlerFunc {
	// Guest can only view marketplace and pool details
	return RoleMiddleware("guest", "investor", "mitra", "admin")
}

func GuestRestricted() gin.HandlerFunc {
	// Restrict guests from performing transactions
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.UnauthorizedError(c, "User role not found")
			c.Abort()
			return
		}

		if role.(string) == "guest" {
			utils.ForbiddenError(c, "Guest accounts cannot perform transactions. Please upgrade your account.")
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware("admin")
}
