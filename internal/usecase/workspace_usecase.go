package usecase

import (
	"errors"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/google/uuid"
)

type WorkspaceUsecase interface {
	CreateWorkspace(name, description, ownerID string) (*entity.Workspace, error)
	GetWorkspace(id string) (*entity.Workspace, error)
	GetUserWorkspaces(userID string) ([]entity.Workspace, error)
	AddMember(workspaceID, userID, role string) error
}

type workspaceUsecase struct {
	repo repository.WorkspaceRepository
}

func NewWorkspaceUsecase(repo repository.WorkspaceRepository) WorkspaceUsecase {
	return &workspaceUsecase{repo: repo}
}

func (u *workspaceUsecase) CreateWorkspace(name, description, ownerID string) (*entity.Workspace, error) {
	if name == "" {
		return nil, errors.New("workspace name is required")
	}

	workspaceID := uuid.New().String()
	workspace := &entity.Workspace{
		ID:          workspaceID,
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.repo.Create(workspace); err != nil {
		return nil, err
	}

	// Add owner as member
	member := &entity.WorkspaceMember{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		UserID:      ownerID,
		Role:        "owner",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.repo.AddMember(member); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (u *workspaceUsecase) GetWorkspace(id string) (*entity.Workspace, error) {
	return u.repo.GetByID(id)
}

func (u *workspaceUsecase) GetUserWorkspaces(userID string) ([]entity.Workspace, error) {
	return u.repo.GetByUserID(userID)
}

func (u *workspaceUsecase) AddMember(workspaceID, userID, role string) error {
	member := &entity.WorkspaceMember{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return u.repo.AddMember(member)
}
