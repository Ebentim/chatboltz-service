package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	Admin      UserRole = "admin"
	SuperAdmin UserRole = "superadmin"
	Staff      UserRole = "staff"
	Usr        UserRole = "user"
	Owner      UserRole = "owner"
	Member     UserRole = "member"
)

// NOTE: super admin can only create admin and staff, can create agent, train agent
// NOTE: admin can only create staff, create and train agent
// NOTE: staff can only handle escalations,

type Users struct {
	ID                string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FirebaseUID       string     `json:"firebase_uid" gorm:"uniqueIndex;type:varchar(128);not null"`
	Name              string     `json:"name" gorm:"type:varchar(255);not null"`
	Agents            []Agent    `json:"agents" gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE,-:save,-:update"`
	Email             string     `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Role              string     `json:"role" gorm:"type:varchar(50);not null;default:admin"`
	Avatar            *string    `json:"avatar" gorm:"type:text"`
	OTPSecret         *string    `json:"-" gorm:"type:varchar(64);index"`
	OTPEnabled        bool       `json:"otp_enabled" gorm:"type:boolean;not null;default:false"`
	OTPLastVerifiedAt *time.Time `json:"otp_last_verified_at" gorm:"type:timestamp"`
	CreatedAt         string     `json:"created_at" gorm:"not null"`
	UpdatedAt         string     `json:"updated_at" gorm:"not null"`
}
type Token struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	Token     string    `json:"token" gorm:"type:text;not null"`
	Purpose   string    `json:"purpose" gorm:"type:varchar(20);not null"`
	Used      bool      `json:"used" gorm:"type:boolean;not null;default:false"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt string    `json:"created_at" gorm:"not null"`
	UpdatedAt string    `json:"updated_at" gorm:"not null"`
	gorm.DeletedAt
}

type AuthRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Token    string `json:"token" binding:"required"`
}

func (Users) TableName() string {
	return "users"
}

func (Token) TableName() string {
	return "tokens"
}

func (t *Token) BeforeCreate(tx *gorm.DB) (err error) {
	var existing Token
	err = tx.Unscoped().Where("email =? AND purpose =?", t.Email, t.Purpose).First(&existing).Error
	if err == nil {
		existing.Token = t.Token
		existing.ExpiresAt = t.ExpiresAt
		existing.Used = false
		existing.DeletedAt = gorm.DeletedAt{}
		existing.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		return tx.Save(&existing).Error
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func (t *Token) AfterFind(tx *gorm.DB) (err error) {
	now := time.Now()

	if now.After(t.ExpiresAt) {
		return tx.Delete(&Token{}, "id = ?", t.ID).Error
	}

	if !t.Used {
		t.Used = true
		t.UpdatedAt = now.UTC().Format(time.RFC3339)

		if err = tx.Save(&t).Error; err != nil {
			return err
		}
		return tx.Delete(&Token{}, "id = ?", t.ID).Error
	}
	return gorm.ErrRecordNotFound
}
