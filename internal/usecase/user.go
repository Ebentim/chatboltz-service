package usecase

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/provider/smtp"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type UserUsecase struct {
	userRepo        repository.UserRepositoryInterface
	firebaseService *FirebaseService
	smtpClient      *smtp.Client
}

func NewUserUsecase(user_repo repository.UserRepositoryInterface, firebase_service *FirebaseService, smtpClient *smtp.Client) *UserUsecase {
	return &UserUsecase{
		userRepo:        user_repo,
		firebaseService: firebase_service,
		smtpClient:      smtpClient,
	}
}

func (u *UserUsecase) SignupWithEmail(req entity.SignupRequest) (*entity.Users, error) {
	// Validate input
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := utils.ValidateRequired(req.Name, "Name"); err != nil {
		return nil, err
	}

	// Create user in Firebase
	firebaseUser, err := u.firebaseService.CreateUser(req.Email, req.Password)
	if err != nil {
		return nil, appErrors.WrapExternalAPIError(err, "Firebase")
	}

	// Create user profile in PostgreSQL
	user, err := u.userRepo.CreateUser(firebaseUser.UID, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	// Generate per-user TOTP secret (disabled by default until explicitly enabled)
	key, err := totp.Generate(totp.GenerateOpts{Issuer: "ChatBoltz", AccountName: req.Email})
	if err == nil {
		secret := key.Secret()
		user.OTPSecret = &secret
		// Do not enable yet; OTPEnabled remains false.
		_ = u.userRepo.UpdateUser(user)
	}

	// Send welcome email (ignore error)
	if u.smtpClient != nil {
		_ = u.smtpClient.Send(req.Email, "Welcome to ChatBoltz", "Thanks for signing up! Secure your account by enabling OTP in settings.")
	}

	return user, nil
}

func (u *UserUsecase) LoginWithEmail(req entity.LoginRequest) (*entity.Users, error) {
	// Validate input
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidateRequired(req.Password, "Password"); err != nil {
		return nil, err
	}

	// Get Firebase user by email to verify existence
	_, err := u.firebaseService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, appErrors.NewAuthenticationError("Invalid credentials")
	}

	// Get user profile from PostgreSQL
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// EnableOTP sets otp_enabled=true and (re)generates secret if missing, sending notification email.
func (u *UserUsecase) EnableOTP(email string) (*entity.Users, error) {
	if err := utils.ValidateEmail(email); err != nil { return nil, err }
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil { return nil, err }
	if user.OTPSecret == nil || *user.OTPSecret == "" {
		key, errGen := totp.Generate(totp.GenerateOpts{Issuer: "ChatBoltz", AccountName: email})
		if errGen == nil { secret := key.Secret(); user.OTPSecret = &secret }
	}
	user.OTPEnabled = true
	now := time.Now()
	user.OTPLastVerifiedAt = &now // mark when enabling (optional)
	if err := u.userRepo.UpdateUser(user); err != nil { return nil, err }
	if u.smtpClient != nil { _ = u.smtpClient.Send(email, "OTP Enabled", "You have enabled OTP. If this wasn't you, disable it immediately.") }
	return user, nil
}

// DisableOTP sets otp_enabled=false (keeps secret for possible re-enable) and notifies user.
func (u *UserUsecase) DisableOTP(email string) (*entity.Users, error) {
	if err := utils.ValidateEmail(email); err != nil { return nil, err }
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil { return nil, err }
	user.OTPEnabled = false
	if err := u.userRepo.UpdateUser(user); err != nil { return nil, err }
	if u.smtpClient != nil { _ = u.smtpClient.Send(email, "OTP Disabled", "You have disabled OTP. Your account is less protected.") }
	return user, nil
}

// ChangePasswordNotification is a hook to send password change email
func (u *UserUsecase) SendPasswordChangedEmail(email string) {
	if u.smtpClient != nil {
		_ = u.smtpClient.Send(email, "Password Changed", "Your password was changed. If this wasn't you, reset immediately.")
	}
}

func (u *UserUsecase) AuthenticateWithToken(id_token string) (*entity.Users, error) {
	// Validate input
	if err := utils.ValidateRequired(id_token, "ID Token"); err != nil {
		return nil, err
	}

	// Verify Firebase ID token
	token, err := u.firebaseService.VerifyIDToken(id_token)
	if err != nil {
		return nil, appErrors.NewAuthenticationError("Invalid token")
	}

	// Get user profile from PostgreSQL
	user, err := u.userRepo.GetUserByFirebaseUID(token.UID)
	if err != nil {
		var appErr *appErrors.AppError
		if errors.As(err, &appErr) && appErr.Type == appErrors.NotFoundError {
			// Create user profile if doesn't exist (for social login)
			displayName := ""
			email := ""
			if name, ok := token.Claims["displayName"].(string); ok {
				displayName = name
			}
			if emailClaim, ok := token.Claims["email"].(string); ok {
				email = emailClaim
			}
			return u.userRepo.CreateUser(token.UID, displayName, email)
		}
		return nil, err
	}

	return user, nil
}
