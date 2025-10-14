package entity

type UserRole string

const (
	Admin      UserRole = "admin"
	SuperAdmin UserRole = "superadmin" 
	Staff      UserRole = "staff"
	Usr        UserRole = "user"
)

// NOTE: super admin can only create admin and staff, can create agent, train agent
// NOTE: admin can only create staff, create and train agent
// NOTE: staff can only handle escalations,

type Users struct {
	ID          string  `json:"id" gorm:"primaryKey;type:varchar(36)"`
	FirebaseUID string  `json:"firebase_uid" gorm:"uniqueIndex;type:varchar(128);not null"`
	Name        string  `json:"name" gorm:"type:varchar(255);not null"`
	Email       string  `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Role        string  `json:"role" gorm:"type:varchar(50);not null;default:admin"`
	Avatar      *string `json:"avatar" gorm:"type:text"`
	CreatedAt   string  `json:"created_at" gorm:"not null"`
	UpdatedAt   string  `json:"updated_at" gorm:"not null"`
}

func (Users) TableName() string {
	return "users"
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
