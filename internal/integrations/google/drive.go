package google

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// DriveService provides methods for interacting with the Google Drive API.
type DriveService struct {
	service *drive.Service
}

// NewDriveService creates a new DriveService.
func NewDriveService(ctx context.Context, ts oauth2.TokenSource) (*DriveService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}
	return &DriveService{service: srv}, nil
}

// ListFiles lists files in Google Drive.
func (s *DriveService) ListFiles(pageSize int64, fields string) (*drive.FileList, error) {
	return s.service.Files.List().PageSize(pageSize).Fields(googleapi.Field(fields)).Do()
}

// CreateFolder creates a new folder in Google Drive.
func (s *DriveService) CreateFolder(name string, parentIDs []string) (*drive.File, error) {
	fileMetadata := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  parentIDs,
	}
	return s.service.Files.Create(fileMetadata).Fields("id", "name", "parents").Do()
}

// SearchFiles searches for files and folders in Google Drive.
// Example query: "name = 'MyFolder' and mimeType = 'application/vnd.google-apps.folder'"
func (s *DriveService) SearchFiles(query string, fields string) (*drive.FileList, error) {
	return s.service.Files.List().Q(query).Fields(googleapi.Field(fields)).Do()
}

// UploadFile uploads a file to Google Drive, optionally into specific folders.
func (s *DriveService) UploadFile(name string, mimeType string, content io.Reader, parentIDs []string) (*drive.File, error) {
	file := &drive.File{
		Name:     name,
		MimeType: mimeType,
		Parents:  parentIDs,
	}
	return s.service.Files.Create(file).Media(content).Fields("id", "name", "parents", "mimeType").Do()
}

// DownloadFile downloads a file from Google Drive.
// The caller is responsible for closing the response body.
func (s *DriveService) DownloadFile(fileID string) (*http.Response, error) {
	return s.service.Files.Get(fileID).Download()
}

// GetFileMetadata retrieves file metadata.
func (s *DriveService) GetFileMetadata(fileID string, fields string) (*drive.File, error) {
	return s.service.Files.Get(fileID).Fields(googleapi.Field(fields)).Do()
}

// MoveFile moves a file to a new folder in Google Drive.
func (s *DriveService) MoveFile(fileID string, newParentIDs []string, oldParentIDs []string) (*drive.File, error) {
	return s.service.Files.Update(fileID, nil).
		AddParents(strings.Join(newParentIDs, ",")).
		RemoveParents(strings.Join(oldParentIDs, ",")).
		Fields("id", "parents").Do()
}

// TrashFile moves a file to the trash in Google Drive.
func (s *DriveService) TrashFile(fileID string) (*drive.File, error) {
	file := &drive.File{
		Trashed: true,
	}
	return s.service.Files.Update(fileID, file).Do()
}
