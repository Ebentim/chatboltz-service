package handler

import (
	"net/http"
	"time"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/provider/smtp"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	otpUsecase *usecase.OTPUsecase
	smtpClient *smtp.Client
	defaultTTL time.Duration
}

func NewOTPHandler(uc *usecase.OTPUsecase, smtpClient *smtp.Client) *OTPHandler {
	return &OTPHandler{otpUsecase: uc, smtpClient: smtpClient, defaultTTL: 10 * time.Minute}
}

// RequestOTP generates and emails an OTP.
func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
		Length  int    `json:"length"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "RequestOTP - JSON binding")
		return
	}
	code, err := h.otpUsecase.Generate(req.Email, req.Purpose, req.Length)
	if err != nil {
		appErrors.HandleError(c, err, "RequestOTP - generate")
		return
	}
	_ = h.smtpClient.SendOTP(req.Email, req.Purpose, code) // ignore delivery errors for now
	c.JSON(http.StatusOK, gin.H{"message": "otp sent", "ttl_minutes": h.defaultTTL.Minutes()})
}

// VerifyOTP verifies provided OTP.
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
		Code    string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "VerifyOTP - JSON binding")
		return
	}
	if err := h.otpUsecase.Verify(req.Email, req.Purpose, req.Code); err != nil {
		appErrors.HandleError(c, err, "VerifyOTP - verify")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "otp verified"})
}
