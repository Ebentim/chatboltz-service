package usecase

import (
	"fmt"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
)

// OTPUsecase handles all OTP operations for the three use cases:
// 1. Two-Factor Authentication (2FA)
// 2. Forgot Password
// 3. Login
type OTPUsecase struct {
	tokenRepo repository.TokenRepositoryInterface
	userRepo  repository.UserRepositoryInterface
	config    *OTPConfig
}

// OTPConfig holds configuration for OTP service
type OTPConfig struct {
	DefaultLength int
	TTL           time.Duration
	MaxAttempts   int
}

// NewOTPUsecase creates a new OTP usecase
func NewOTPUsecase(tokenRepo repository.TokenRepositoryInterface, userRepo repository.UserRepositoryInterface, ttl time.Duration) *OTPUsecase {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &OTPUsecase{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
		config: &OTPConfig{
			DefaultLength: 6,
			TTL:           ttl,
			MaxAttempts:   3,
		},
	}
}

// GenerateOTP generates OTP for any of the three use cases
func (o *OTPUsecase) GenerateOTP(req *entity.OTPRequest) (*entity.OTPResponse, error) {
	_, response, err := o.GenerateOTPWithCode(req)
	return response, err
}

// GenerateOTPWithCode generates OTP and returns both the code and response
func (o *OTPUsecase) GenerateOTPWithCode(req *entity.OTPRequest) (string, *entity.OTPResponse, error) {
	if err := o.validateOTPRequest(req); err != nil {
		return "", nil, err
	}

	if err := o.validatePurposeRules(req); err != nil {
		return "", nil, err
	}

	length := req.Length
	if length == 0 {
		length = o.config.DefaultLength
	}

	code, err := utils.GenerateOTP(length)
	if err != nil {
		return "", nil, appErrors.NewInternalError("failed to generate OTP", err.Error())
	}

	hashedCode, err := utils.CreateHash([]byte(code))
	if err != nil {
		return "", nil, appErrors.NewInternalError("failed to hash OTP", err.Error())
	}

	expiresAt := time.Now().Add(o.config.TTL)
	if err := o.tokenRepo.CreateToken(req.Email, req.Purpose.String(), hashedCode, expiresAt); err != nil {
		return "", nil, appErrors.WrapDatabaseError(err, "create OTP token")
	}

	response := &entity.OTPResponse{
		Message:    fmt.Sprintf("OTP sent for %s", req.Purpose),
		TTLMinutes: int(o.config.TTL.Minutes()),
		Purpose:    req.Purpose.String(),
	}

	return code, response, nil
}

// VerifyOTP verifies OTP for any of the three use cases
func (o *OTPUsecase) VerifyOTP(req *entity.OTPVerifyRequest) (*entity.OTPVerifyResponse, error) {
	if err := o.validateVerifyRequest(req); err != nil {
		return nil, err
	}

	token, err := o.tokenRepo.GetToken(req.Email, req.Purpose.String())
	if err != nil {
		return nil, appErrors.NewValidationError("invalid or expired OTP")
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, appErrors.NewValidationError("OTP has expired")
	}

	if err := utils.ValidateHash(req.Code, token.Token); err != nil {
		return nil, appErrors.NewValidationError("invalid OTP code")
	}

	if err := o.handlePostVerification(req); err != nil {
		return nil, err
	}

	return &entity.OTPVerifyResponse{
		Message:   fmt.Sprintf("OTP verified for %s", req.Purpose),
		Purpose:   req.Purpose.String(),
		Verified:  true,
		ExpiresAt: &token.ExpiresAt,
	}, nil
}

// Enable2FA enables two-factor authentication for a user
func (o *OTPUsecase) Enable2FA(userID string) error {
	user, err := o.userRepo.GetUserByID(userID)
	if err != nil {
		return appErrors.WrapDatabaseError(err, "get user for 2FA enable")
	}
	user.OTPEnabled = true
	return o.userRepo.UpdateUser(user)
}

// Disable2FA disables two-factor authentication for a user
func (o *OTPUsecase) Disable2FA(userID string) error {
	user, err := o.userRepo.GetUserByID(userID)
	if err != nil {
		return appErrors.WrapDatabaseError(err, "get user for 2FA disable")
	}
	user.OTPEnabled = false
	user.OTPLastVerifiedAt = nil
	return o.userRepo.UpdateUser(user)
}

// CompletePasswordReset completes password reset with OTP verification
func (o *OTPUsecase) CompletePasswordReset(email, code, newPassword string) error {
	req := &entity.OTPVerifyRequest{
		Email:   email,
		Purpose: entity.OTPPurposeForgotPassword,
		Code:    code,
	}

	if _, err := o.VerifyOTP(req); err != nil {
		return err
	}

	if len(newPassword) < 6 {
		return appErrors.NewValidationError("password must be at least 6 characters long")
	}

	// Note: Password update logic would go here
	// This depends on your user authentication system
	return nil
}

// CompleteOTPLogin completes OTP-based login
func (o *OTPUsecase) CompleteOTPLogin(email, code string) (*entity.Users, error) {
	req := &entity.OTPVerifyRequest{
		Email:   email,
		Purpose: entity.OTPPurposeLogin,
		Code:    code,
	}

	if _, err := o.VerifyOTP(req); err != nil {
		return nil, err
	}

	return o.userRepo.GetUserByEmail(email)
}

// Legacy methods for backward compatibility
func (o *OTPUsecase) Generate(email, basePurpose string, length int) (string, error) {
	purpose := entity.OTPPurpose(basePurpose)
	if !purpose.IsValid() {
		return "", appErrors.NewValidationError("invalid purpose")
	}

	code, err := utils.GenerateOTP(length)
	if err != nil {
		return "", err
	}

	req := &entity.OTPRequest{
		Email:   email,
		Purpose: purpose,
		Length:  length,
	}

	if _, err := o.GenerateOTP(req); err != nil {
		return "", err
	}

	return code, nil
}

func (o *OTPUsecase) Verify(email, basePurpose, provided string) error {
	purpose := entity.OTPPurpose(basePurpose)
	if !purpose.IsValid() {
		return appErrors.NewValidationError("invalid purpose")
	}

	req := &entity.OTPVerifyRequest{
		Email:   email,
		Purpose: purpose,
		Code:    provided,
	}

	_, err := o.VerifyOTP(req)
	return err
}

func (o *OTPUsecase) Consume(email, basePurpose, provided string) error {
	return o.Verify(email, basePurpose, provided)
}

// Private helper methods
func (o *OTPUsecase) validateOTPRequest(req *entity.OTPRequest) error {
	if req == nil {
		return appErrors.NewValidationError("request cannot be nil")
	}
	if err := utils.ValidateEmail(req.Email); err != nil {
		return err
	}
	if !req.Purpose.IsValid() {
		return appErrors.NewValidationError("invalid OTP purpose")
	}
	if req.Length != 0 && (req.Length < 4 || req.Length > 10) {
		return appErrors.NewValidationError("OTP length must be between 4 and 10 digits")
	}
	return nil
}

func (o *OTPUsecase) validateVerifyRequest(req *entity.OTPVerifyRequest) error {
	if req == nil {
		return appErrors.NewValidationError("request cannot be nil")
	}
	return utils.ValidateOTPInput(req.Email, req.Purpose.String(), req.Code)
}

func (o *OTPUsecase) validatePurposeRules(req *entity.OTPRequest) error {
	switch req.Purpose {
	case entity.OTPPurpose2FA:
		user, err := o.userRepo.GetUserByEmail(req.Email)
		if err != nil {
			return appErrors.NewValidationError("user not found")
		}
		if !user.OTPEnabled {
			return appErrors.NewValidationError("2FA is not enabled for this user")
		}
	case entity.OTPPurposeForgotPassword, entity.OTPPurposeLogin:
		_, err := o.userRepo.GetUserByEmail(req.Email)
		if err != nil {
			return appErrors.NewValidationError("user not found")
		}
	}
	return nil
}

func (o *OTPUsecase) handlePostVerification(req *entity.OTPVerifyRequest) error {
	if req.Purpose == entity.OTPPurpose2FA {
		user, err := o.userRepo.GetUserByEmail(req.Email)
		if err != nil {
			return err
		}
		now := time.Now()
		user.OTPLastVerifiedAt = &now
		return o.userRepo.UpdateUser(user)
	}
	return nil
}
