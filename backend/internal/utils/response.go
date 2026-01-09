package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vessel/backend/internal/models"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Created successfully",
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    code,
			Message: message,
		},
	})
}

func BadRequestError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func UnauthorizedError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func ForbiddenError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFoundError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message)
}

func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

func ConflictError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message)
}

func TooManyRequestsError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

func ValidationError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message)
}

func HandleAppError(c *gin.Context, err error) {
	if appErr := GetAppError(err); appErr != nil {
		switch appErr.Code {
		case ErrCodeNotFound:
			NotFoundError(c, appErr.Message)
		case ErrCodeUnauthorized:
			UnauthorizedError(c, appErr.Message)
		case ErrCodeForbidden:
			ForbiddenError(c, appErr.Message)
		case ErrCodeConflict:
			ConflictError(c, appErr.Message)
		case ErrCodeValidation:
			ValidationError(c, appErr.Message)
		case ErrCodeBadRequest:
			BadRequestError(c, appErr.Message)
		default:
			InternalServerError(c, appErr.Message)
		}
		return
	}
	BadRequestError(c, err.Error())
}
