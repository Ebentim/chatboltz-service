package handler

import (
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

// EmailServiceInterface defines the interface for email services
type EmailServiceInterface interface {
	SendOTP(to, purpose, code string) error
	SendOTPWithClient(client interface{}, to, purpose, code string) error
	CreateClientFromConfig(config interface{}) interface{}
}

// OTPHandler handles HTTP requests for OTP operations
type OTPHandler struct {
	otpUsecase   *usecase.OTPUsecase
	emailService EmailServiceInterface
}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler(otpUsecase *usecase.OTPUsecase, emailService EmailServiceInterface) *OTPHandler {
	return &OTPHandler{
		otpUsecase:   otpUsecase,
		emailService: emailService,
	}
}

// RequestOTP handles OTP generation requests
func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req entity.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "RequestOTP")
		return
	}

	code, response, err := h.otpUsecase.GenerateOTPWithCode(&req)
	if err != nil {
		appErrors.HandleError(c, err, "RequestOTP")
		return
	}

	// Send email with actual OTP code using templates
	if h.emailService != nil {
		if err := h.emailService.SendOTP(req.Email, req.Purpose.String(), code); err != nil {
			// Log error but don't fail the request
			// In production, you might want to handle this differently
		}
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP handles OTP verification requests
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req entity.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "VerifyOTP")
		return
	}

	response, err := h.otpUsecase.VerifyOTP(&req)
	if err != nil {
		appErrors.HandleError(c, err, "VerifyOTP")
		return
	}

	c.JSON(http.StatusOK, response)
}

// Enable2FA enables two-factor authentication
func (h *OTPHandler) Enable2FA(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("User ID required"), "Enable2FA")
		return
	}

	if err := h.otpUsecase.Enable2FA(userID); err != nil {
		appErrors.HandleError(c, err, "Enable2FA")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA disables two-factor authentication
func (h *OTPHandler) Disable2FA(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("User ID required"), "Disable2FA")
		return
	}

	if err := h.otpUsecase.Disable2FA(userID); err != nil {
		appErrors.HandleError(c, err, "Disable2FA")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// CompletePasswordReset handles password reset completion
func (h *OTPHandler) CompletePasswordReset(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CompletePasswordReset")
		return
	}

	if err := h.otpUsecase.CompletePasswordReset(req.Email, req.Code, req.NewPassword); err != nil {
		appErrors.HandleError(c, err, "CompletePasswordReset")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// CompleteOTPLogin handles OTP login completion
func (h *OTPHandler) CompleteOTPLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "CompleteOTPLogin")
		return
	}

	user, err := h.otpUsecase.CompleteOTPLogin(req.Email, req.Code)
	if err != nil {
		appErrors.HandleError(c, err, "CompleteOTPLogin")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "message": "Login successful"})
}
