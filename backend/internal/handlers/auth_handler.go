package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/models"
	"github.com/vessel/backend/internal/services"
	"github.com/vessel/backend/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
	otpService  *services.OTPService
}

func NewAuthHandler(authService *services.AuthService, otpService *services.OTPService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		otpService:  otpService,
	}
}

// SendOTP godoc
// @Summary Send OTP to email
// @Description Send a 6-digit OTP code to the specified email for verification
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.SendOTPRequest true "Email and purpose"
// @Success 200 {object} models.SendOTPResponse
// @Failure 400 {object} models.APIError
// @Failure 429 {object} models.APIError "Rate limit exceeded"
// @Router /auth/send-otp [post]
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req models.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.otpService.GenerateAndSendOTP(req.Email, req.Purpose)
	if err != nil {
		if err == services.ErrOTPRateLimit {
			utils.TooManyRequestsError(c, "Too many OTP requests. Please wait before trying again.")
			return
		}
		utils.BadRequestError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// VerifyOTP godoc
// @Summary Verify OTP code
// @Description Verify the OTP code sent to email. Returns a token for registration.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Email and OTP code"
// @Success 200 {object} models.VerifyOTPResponse
// @Failure 400 {object} models.APIError
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.otpService.VerifyOTP(req.Email, req.Code)
	if err != nil {
		switch err {
		case services.ErrOTPNotFound:
			utils.BadRequestError(c, "OTP code not found or expired")
		case services.ErrOTPExpired:
			utils.BadRequestError(c, "OTP code has expired")
		case services.ErrOTPMaxAttempts:
			utils.BadRequestError(c, "Maximum attempts reached. Please request a new OTP.")
		default:
			utils.BadRequestError(c, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new exporter or investor account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} models.APIError
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.APIError
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.APIError
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err.Error())
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.UnauthorizedError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}
