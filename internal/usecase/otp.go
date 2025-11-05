package usecase

import (
	"errors"
	"time"

	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
	"github.com/pquerna/otp/totp"
)

// OTPUsecase handles generation and verification of OTP codes using the token repository.
// It reuses the existing tokens table. Because the tokens.purpose column is an enum
// (login|register|password_reset|2fa), we MUST store the OTP under one of those base purposes
// and cannot introduce new values (like otp_login) without a schema migration.
// OTP codes are encrypted before storage and are one-time because of the Token.AfterFind hook.
type OTPUsecase struct {
	tokenRepo repository.TokenRepositoryInterface
	secret    string        // TOTP secret per project (could be per-user in future)
	otpTTL    time.Duration // validity window for manual expiration (tokens table) separate from TOTP step period
	period    uint          // TOTP period seconds
	digits    int           // number of digits
}

func NewOTPUsecase(tokenRepo repository.TokenRepositoryInterface, secret string, ttl time.Duration) *OTPUsecase {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &OTPUsecase{tokenRepo: tokenRepo, secret: secret, otpTTL: ttl, period: 30, digits: 6}
}

// Generate issues a new OTP for given email & basePurpose (e.g., "login"). It overwrites previous.
func (o *OTPUsecase) Generate(email, basePurpose string, length int) (string, error) {
	if err := utils.ValidateEmail(email); err != nil {
		return "", err
	}
	if err := utils.ValidateTokenPurpose(basePurpose); err != nil {
		return "", err
	}
	if length < 4 || length > 10 {
		length = int(o.digits)
	}
	// Use pquerna totp GenerateCodeCustom with custom digits/period anchored to current time
	code, err := totp.GenerateCodeCustom(o.secret, time.Now(), totp.ValidateOpts{Period: o.period, Digits: totp.Digits(length), Skew: 1})
	if err != nil {
		return "", err
	}
	// Hash the code before storing
	hash, err := utils.CreateHash([]byte(code))
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(o.otpTTL)
	if err := o.tokenRepo.CreateToken(email, basePurpose, hash, expiresAt); err != nil {
		return "", err
	}
	return code, nil
}

// Verify checks provided otp. Returns nil if valid.
func (o *OTPUsecase) Verify(email, basePurpose, provided string) error {
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}
	if err := utils.ValidateTokenPurpose(basePurpose); err != nil {
		return err
	}
	if provided == "" {
		return appErrors.NewValidationError("otp is required")
	}
	t, err := o.tokenRepo.GetToken(email, basePurpose)
	if err != nil {
		return err
	}
	if time.Now().After(t.ExpiresAt) {
		return appErrors.NewValidationError("otp expired")
	}
	if err := utils.ValidateHash(provided, t.Token); err != nil {
		return appErrors.NewValidationError("invalid otp")
	}
	// Optional: verify TOTP window (defense in depth)
	valid := totp.Validate(provided, o.secret)
	if !valid {
		return appErrors.NewValidationError("invalid otp")
	}
	return nil
}

// Consume convenience alias for Verify for semantics
func (o *OTPUsecase) Consume(email, basePurpose, provided string) error {
	return o.Verify(email, basePurpose, provided)
}

// Helper to map base purposes if needed in handlers
var ErrOTPInvalid = errors.New("invalid otp")
