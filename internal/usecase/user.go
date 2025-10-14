package usecase

import (
	"errors"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
)

type UserUsecase struct {
	userRepo        repository.UserRepositoryInterface
	firebaseService *FirebaseService
}

func NewUserUsecase(user_repo repository.UserRepositoryInterface, firebase_service *FirebaseService) *UserUsecase {
	return &UserUsecase{
		userRepo:        user_repo,
		firebaseService: firebase_service,
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
