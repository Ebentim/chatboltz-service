package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

// SlidesService provides methods for interacting with the Google Slides API.
type SlidesService struct {
	service *slides.Service
}

// NewSlidesService creates a new SlidesService.
func NewSlidesService(ctx context.Context, ts oauth2.TokenSource) (*SlidesService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := slides.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Slides client: %v", err)
	}
	return &SlidesService{service: srv}, nil
}

// CreatePresentation creates a new Google Slides presentation.
func (s *SlidesService) CreatePresentation(title string) (*slides.Presentation, error) {
	presentation := &slides.Presentation{
		Title: title,
	}
	return s.service.Presentations.Create(presentation).Do()
}

// GetPresentation retrieves a Google Slides presentation.
func (s *SlidesService) GetPresentation(presentationID string) (*slides.Presentation, error) {
	return s.service.Presentations.Get(presentationID).Do()
}

// BatchUpdate performs a batch update on a Google Slides presentation.
func (s *SlidesService) BatchUpdate(presentationID string, requests []*slides.Request) (*slides.BatchUpdatePresentationResponse, error) {
	batchUpdate := &slides.BatchUpdatePresentationRequest{
		Requests: requests,
	}
	return s.service.Presentations.BatchUpdate(presentationID, batchUpdate).Do()
}
