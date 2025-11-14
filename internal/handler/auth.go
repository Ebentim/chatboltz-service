package handler

import (
	"log"
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/middleware"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUsecase *usecase.UserUsecase
	jwtSecret   []byte
}

func NewAuthHandler(userUsecase *usecase.UserUsecase, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{userUsecase: userUsecase, jwtSecret: jwtSecret}
}

func (h *AuthHandler) SignupWithEmail(c *gin.Context) {
	var req entity.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "SignupWithEmail - JSON binding")
		return
	}

	user, err := h.userUsecase.SignupWithEmail(req)
	if err != nil {
		appErrors.HandleError(c, err, "SignupWithEmail")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *AuthHandler) LoginWithEmail(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(&req)
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "LoginWithEmail - JSON binding")
		return
	}

	user, err := h.userUsecase.LoginWithEmail(req)
	if err != nil {
		appErrors.HandleError(c, err, "LoginWithEmail")
		return
	}

	token, err := middleware.GenerateToken(*user, h.jwtSecret)
	if err != nil {
		appErrors.HandleError(c, err, "LoginWithEmail - Generate Token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *AuthHandler) AuthenticateWithToken(c *gin.Context) {
	var req entity.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "AuthenticateWithToken - JSON binding")
		return
	}

	user, err := h.userUsecase.AuthenticateWithToken(req.IDToken)
	if err != nil {
		appErrors.HandleError(c, err, "AuthenticateWithToken")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// EnableOTP endpoint
func (h *AuthHandler) EnableOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request"), "EnableOTP - bind")
		return
	}
	user, err := h.userUsecase.EnableOTP(req.Email)
	if err != nil {
		appErrors.HandleError(c, err, "EnableOTP")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "otp enabled", "user": user.ID})
}

// DisableOTP endpoint
func (h *AuthHandler) DisableOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request"), "DisableOTP - bind")
		return
	}
	user, err := h.userUsecase.DisableOTP(req.Email)
	if err != nil {
		appErrors.HandleError(c, err, "DisableOTP")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "otp disabled", "user": user.ID})
}
