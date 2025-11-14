package usecase

import (
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/crypto"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
)

type TokenUsecase struct {
	tokenRepo repository.TokenRepositoryInterface
}

func NewTokenUsecase(tokenRepo repository.TokenRepositoryInterface) *TokenUsecase {
	if tokenRepo == nil {
		panic("Token repository is nil")
	}
	return &TokenUsecase{
		tokenRepo: tokenRepo,
	}
}

func (u *TokenUsecase) GetToken(email, purpose, providedToken string, EncryptionKey []byte) (string, error) {
	if err := utils.ValidateEmail(email); err != nil {
		return "", err
	}

	if err := utils.ValidateTokenPurpose(purpose); err != nil {
		return "", err
	}

	token, err := u.tokenRepo.GetToken(email, purpose)
	if err != nil {
		return "", err
	}

	init_encryption := crypto.NewEncryptionKey(EncryptionKey)

	decryptedToken, err := init_encryption.DecryptString(token.Token)

	if err != nil {
		return "", appErrors.NewValidationError("Invalid token")
	}

	if string(decryptedToken) == "" || string(decryptedToken) != providedToken {
		return "", appErrors.NewValidationError("Invalid token")
	}
	return string(decryptedToken), nil
}

func (u *TokenUsecase) CreateToken(email, purpose, token string, EncryptionKey []byte) error {
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	if err := utils.ValidateTokenPurpose(purpose); err != nil {
		return err
	}

	init_encryption := crypto.NewEncryptionKey(EncryptionKey)
	encryptedToken, err := init_encryption.EncryptString([]byte(token))
	expiresAt := time.Now().Add(10 * time.Minute)

	if err != nil {
		return err
	}

	return u.tokenRepo.CreateToken(email, purpose, string(encryptedToken), expiresAt)
}
