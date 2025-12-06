package entity

import (
	"time"

	"gorm.io/gorm"
)

type Workspace struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	OwnerID     string         `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Members []WorkspaceMember `gorm:"foreignKey:WorkspaceID" json:"members,omitempty"`
	Agents  []Agent           `gorm:"foreignKey:WorkspaceID" json:"agents,omitempty"`
}

type WorkspaceMember struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID string         `gorm:"type:uuid;not null;index" json:"workspace_id"`
	UserID      string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Role        string         `gorm:"type:varchar(50);not null;default:'member'" json:"role"` // owner, admin, member
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Workspace Workspace `gorm:"foreignKey:WorkspaceID" json:"-"`
	User      Users     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
