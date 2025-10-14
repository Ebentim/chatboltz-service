package repository

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(firebaseUID, name, email string) (*entity.Users, error) {
	user := &entity.Users{
		ID:          uuid.New().String(),
		FirebaseUID: firebaseUID,
		Name:        name,
		Email:       email,
		Role:        string(entity.Admin),
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "create user")
	}

	return user, nil
}

func (r *UserRepository) GetUserByFirebaseUID(firebaseUID string) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("User not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get user by firebase UID")
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("User not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get user by email")
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id string) (*entity.Users, error) {
	var user entity.Users
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError("User not found")
		}
		return nil, appErrors.WrapDatabaseError(err, "get user by ID")
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *entity.Users) error {
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := r.db.Save(user).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "update user")
	}
	return nil
}

func (r *UserRepository) DeleteUser(id string) error {
	if err := r.db.Delete(&entity.Users{}, "id = ?", id).Error; err != nil {
		return appErrors.WrapDatabaseError(err, "delete user")
	}
	return nil
}

func (r *UserRepository) ListUsers() ([]*entity.Users, error) {
	var users []*entity.Users
	if err := r.db.Find(&users).Error; err != nil {
		return nil, appErrors.WrapDatabaseError(err, "list users")
	}
	return users, nil
}
