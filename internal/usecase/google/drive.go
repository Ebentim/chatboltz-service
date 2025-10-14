package usecase

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	gdrive "google.golang.org/api/drive/v3"
)

// DriveIntegration defines the interface for drive operations.
type DriveIntegration interface {
	ListFiles(pageSize int64, fields string) (*gdrive.FileList, error)
	CreateFolder(name string, parentIDs []string) (*gdrive.File, error)
	SearchFiles(query string, fields string) (*gdrive.FileList, error)
	UploadFile(name string, mimeType string, content io.Reader, parentIDs []string) (*gdrive.File, error)
	DownloadFile(fileID string) (*http.Response, error)
	GetFileMetadata(fileID string, fields string) (*gdrive.File, error)
	MoveFile(fileID string, newParentIDs []string, oldParentIDs []string) (*gdrive.File, error)
	TrashFile(fileID string) (*gdrive.File, error)
}

// DriveUseCase handles the business logic for drive operations.
type DriveUseCase struct {
	driveIntegration DriveIntegration
}

// NewDriveUseCase creates a new DriveUseCase.
func NewDriveUseCase(di DriveIntegration) *DriveUseCase {
	return &DriveUseCase{driveIntegration: di}
}

// ListFiles lists files in Google Drive.
func (uc *DriveUseCase) ListFiles(pageSize int64, fields string) (*gdrive.FileList, error) {
	return uc.driveIntegration.ListFiles(pageSize, fields)
}

// CreateFolder creates a new folder. If a folder with the same name and parents already exists,
// it returns the existing folder's information.
func (uc *DriveUseCase) CreateFolder(name string, parentIDs []string) (*gdrive.File, error) {
	query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and trashed = false", name)

	existing, err := uc.driveIntegration.SearchFiles(query, "files(id, name, parents)")
	if err != nil {
		return nil, err
	}

	for _, folder := range existing.Files {
		if parentsMatch(folder.Parents, parentIDs) {
			return folder, nil // Folder with same name and parents already exists
		}
	}

	return uc.driveIntegration.CreateFolder(name, parentIDs)
}

// SearchFiles searches for files and folders.
func (uc *DriveUseCase) SearchFiles(query string, fields string) (*gdrive.FileList, error) {
	return uc.driveIntegration.SearchFiles(query, fields)
}

// UploadFile uploads a file to Google Drive.
func (uc *DriveUseCase) UploadFile(name string, mimeType string, content io.Reader, parentIDs []string) (*gdrive.File, error) {
	// More complex business logic could be added here, e.g., checking for duplicates first.
	return uc.driveIntegration.UploadFile(name, mimeType, content, parentIDs)
}

// DownloadFile downloads a file from Google Drive.
func (uc *DriveUseCase) DownloadFile(fileID string) (*http.Response, error) {
	return uc.driveIntegration.DownloadFile(fileID)
}

// GetFileMetadata retrieves file metadata.
func (uc *DriveUseCase) GetFileMetadata(fileID string, fields string) (*gdrive.File, error) {
	return uc.driveIntegration.GetFileMetadata(fileID, fields)
}

// MoveFile moves a file to a new folder.
func (uc *DriveUseCase) MoveFile(fileID string, newParentIDs []string, oldParentIDs []string) (*gdrive.File, error) {
	return uc.driveIntegration.MoveFile(fileID, newParentIDs, oldParentIDs)
}

// TrashFile moves a file to the trash.
func (uc *DriveUseCase) TrashFile(fileID string) (*gdrive.File, error) {
	return uc.driveIntegration.TrashFile(fileID)
}

// GetOrCreateFolder checks if a folder exists with the exact same set of parents and creates it if it doesn't.
func (uc *DriveUseCase) GetOrCreateFolder(name string, parentIDs []string) (*gdrive.File, error) {
	query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and trashed = false", name)

	list, err := uc.driveIntegration.SearchFiles(query, "files(id, name, parents)")
	if err != nil {
		return nil, fmt.Errorf("failed to search for folder '%s': %w", name, err)
	}

	for _, folder := range list.Files {
		if parentsMatch(folder.Parents, parentIDs) {
			return folder, nil
		}
	}

	// Folder does not exist with all the specified parents, create it
	return uc.driveIntegration.CreateFolder(name, parentIDs)
}

// GetFileByPath finds a file or folder by a path-like string.
// The path should be absolute, starting from the root folder (e.g., "/My Documents/Reports/Q1.pdf").
func (uc *DriveUseCase) GetFileByPath(path string) (*gdrive.File, error) {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		// Return root folder if path is "/" or ""
		return uc.GetFileMetadata("root", "id,name,mimeType")
	}

	parentID := "root"
	var currentFile *gdrive.File

	for _, part := range parts {
		query := fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", part, parentID)

		list, err := uc.driveIntegration.SearchFiles(query, "files(id, name, parents, mimeType)")
		if err != nil {
			return nil, fmt.Errorf("failed to search for '%s': %w", part, err)
		}
		if len(list.Files) == 0 {
			return nil, fmt.Errorf("path not found: '%s' not found in '%s'", part, path)
		}
		currentFile = list.Files[0]
		parentID = currentFile.Id
	}

	return currentFile, nil
}

// parentsMatch checks if two slices of parent IDs are identical, regardless of order.
func parentsMatch(actualParents, expectedParents []string) bool {
	if len(actualParents) != len(expectedParents) {
		return false
	}
	parentSet := make(map[string]bool)
	for _, p := range actualParents {
		parentSet[p] = true
	}
	for _, p := range expectedParents {
		if !parentSet[p] {
			return false
		}
	}
	return true
}
