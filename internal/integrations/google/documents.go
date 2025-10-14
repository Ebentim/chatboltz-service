package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

// DocsService provides methods for interacting with the Google Docs API.
type DocsService struct {
	service *docs.Service
}

// NewDocsService creates a new DocsService.
func NewDocsService(ctx context.Context, ts oauth2.TokenSource) (*DocsService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Docs client: %v", err)
	}
	return &DocsService{service: srv}, nil
}

// CreateDocument creates a new Google Doc.
func (s *DocsService) CreateDocument(title string) (*docs.Document, error) {
	doc := &docs.Document{
		Title: title,
	}
	return s.service.Documents.Create(doc).Do()
}

// GetDocument retrieves a Google Doc.
func (s *DocsService) GetDocument(documentID string) (*docs.Document, error) {
	return s.service.Documents.Get(documentID).Do()
}

// BatchUpdate performs a batch update on a Google Doc.
func (s *DocsService) BatchUpdate(documentID string, requests []*docs.Request) (*docs.BatchUpdateDocumentResponse, error) {
	batchUpdate := &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}
	return s.service.Documents.BatchUpdate(documentID, batchUpdate).Do()
}
