package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

// FormsService provides methods for interacting with the Google Forms API.
type FormsService struct {
	service *forms.Service
}

// NewFormsService creates a new FormsService.
func NewFormsService(ctx context.Context, ts oauth2.TokenSource) (*FormsService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := forms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Forms client: %v", err)
	}
	return &FormsService{service: srv}, nil
}

// CreateForm creates a new Google Form.
func (s *FormsService) CreateForm(form *forms.Form) (*forms.Form, error) {
	return s.service.Forms.Create(form).Do()
}

// GetForm retrieves a Google Form.
func (s *FormsService) GetForm(formID string) (*forms.Form, error) {
	return s.service.Forms.Get(formID).Do()
}

// GetFormResponses retrieves all responses from a Google Form.
func (s *FormsService) GetFormResponses(formID string) (*forms.ListFormResponsesResponse, error) {
	return s.service.Forms.Responses.List(formID).Do()
}

// BatchUpdate performs a batch update on a Google Form.
func (s *FormsService) BatchUpdate(formID string, requests []*forms.Request) (*forms.BatchUpdateFormResponse, error) {
	batchUpdate := &forms.BatchUpdateFormRequest{
		Requests: requests,
	}
	return s.service.Forms.BatchUpdate(formID, batchUpdate).Do()
}
