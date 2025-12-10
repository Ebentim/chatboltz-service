package repository

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(workspace *entity.Workspace) error
	GetByID(id string) (*entity.Workspace, error)
	GetByUserID(userID string) ([]entity.Workspace, error)
	GetByAgentID(agentID string) (*entity.Workspace, error)
	Update(workspace *entity.Workspace) error
	Delete(id string) error
	AddMember(member *entity.WorkspaceMember) error
	RemoveMember(workspaceID, userID string) error
}

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) Create(workspace *entity.Workspace) error {
	return r.db.Create(workspace).Error
}

func (r *workspaceRepository) GetByID(id string) (*entity.Workspace, error) {
	var workspace entity.Workspace
	if err := r.db.Preload("Members").Preload("Agents").First(&workspace, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) GetByAgentID(agentID string) (*entity.Workspace, error) {
	var workspace entity.Workspace
	if err := r.db.Joins("JOIN agents ON agents.workspace_id = workspaces.id").
		Where("agents.id = ?", agentID).
		Preload("Members").
		First(&workspace).Error; err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) GetByUserID(userID string) ([]entity.Workspace, error) {
	var workspaces []entity.Workspace
	// Find workspaces where the user is a member or owner
	// Assuming logic: User is owner OR User is in members list
	// For simplicity, we can query via WorkspaceMember if we enforce owner is also a member,
	// or query both. Let's assume owner is added as a member upon creation.
	err := r.db.Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Find(&workspaces).Error
	return workspaces, err
}

func (r *workspaceRepository) Update(workspace *entity.Workspace) error {
	return r.db.Save(workspace).Error
}

func (r *workspaceRepository) Delete(id string) error {
	return r.db.Delete(&entity.Workspace{}, "id = ?", id).Error
}

func (r *workspaceRepository) AddMember(member *entity.WorkspaceMember) error {
	return r.db.Create(member).Error
}

func (r *workspaceRepository) RemoveMember(workspaceID, userID string) error {
	return r.db.Delete(&entity.WorkspaceMember{}, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error
}
