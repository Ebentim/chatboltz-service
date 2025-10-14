package usecase

import (
	"errors"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"gorm.io/gorm"
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
	// Create user in Firebase
	firebaseUser, err := u.firebaseService.CreateUser(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Create user profile in PostgreSQL
	user, err := u.userRepo.CreateUser(firebaseUser.UID, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) LoginWithEmail(req entity.LoginRequest) (*entity.Users, error) {
	// Get Firebase user by email to verify existence
	_, err := u.firebaseService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user profile from PostgreSQL
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) AuthenticateWithToken(id_token string) (*entity.Users, error) {
	// Verify Firebase ID token
	token, err := u.firebaseService.VerifyIDToken(id_token)
	if err != nil {
		return nil, err
	}

	// Get user profile from PostgreSQL
	user, err := u.userRepo.GetUserByFirebaseUID(token.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create user profile if doesn't exist (for social login)
			return u.userRepo.CreateUser(token.UID, token.Claims["name"].(string), token.Claims["email"].(string))
		}
		return nil, err
	}

	return user, nil
}
