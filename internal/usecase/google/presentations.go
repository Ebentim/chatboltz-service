package usecase

import (
	"fmt"

	gslides "google.golang.org/api/slides/v1"
)

// SlidesIntegration defines the interface for slides operations.
type SlidesIntegration interface {
	CreatePresentation(title string) (*gslides.Presentation, error)
	GetPresentation(presentationID string) (*gslides.Presentation, error)
	BatchUpdate(presentationID string, requests []*gslides.Request) (*gslides.BatchUpdatePresentationResponse, error)
}

// SlidesUseCase handles the business logic for slides operations.
type SlidesUseCase struct {
	slidesIntegration SlidesIntegration
	driveIntegration  DriveIntegration
}

// NewSlidesUseCase creates a new SlidesUseCase.
func NewSlidesUseCase(si SlidesIntegration, driveIntegration DriveIntegration) *SlidesUseCase {
	return &SlidesUseCase{
		slidesIntegration: si,
		driveIntegration:  driveIntegration,
	}
}

// CreatePresentation creates a new Google Slides presentation.
func (uc *SlidesUseCase) CreatePresentation(title string) (*gslides.Presentation, error) {
	return uc.slidesIntegration.CreatePresentation(title)
}

// GetPresentation retrieves a Google Slides presentation.
func (uc *SlidesUseCase) GetPresentation(presentationID string) (*gslides.Presentation, error) {
	return uc.slidesIntegration.GetPresentation(presentationID)
}

// CreateSlide creates a new slide in a presentation.
func (uc *SlidesUseCase) CreateSlide(presentationID string, slideID string, layoutID string) (*gslides.BatchUpdatePresentationResponse, error) {
	requests := []*gslides.Request{
		{
			CreateSlide: &gslides.CreateSlideRequest{
				ObjectId: slideID,
				SlideLayoutReference: &gslides.LayoutReference{
					PredefinedLayout: layoutID,
				},
			},
		},
	}
	return uc.slidesIntegration.BatchUpdate(presentationID, requests)
}

// InsertText inserts text into a shape on a slide.
func (uc *SlidesUseCase) InsertText(presentationID string, shapeID string, text string) (*gslides.BatchUpdatePresentationResponse, error) {
	requests := []*gslides.Request{
		{
			InsertText: &gslides.InsertTextRequest{
				ObjectId: shapeID,
				Text:     text,
			},
		},
	}
	return uc.slidesIntegration.BatchUpdate(presentationID, requests)
}

// InsertImage inserts an image into a slide.
func (uc *SlidesUseCase) InsertImage(presentationID string, slideID string, imageURL string, objectID string, height float64, width float64) (*gslides.BatchUpdatePresentationResponse, error) {
	requests := []*gslides.Request{
		{
			CreateImage: &gslides.CreateImageRequest{
				ObjectId: objectID,
				Url:      imageURL,
				ElementProperties: &gslides.PageElementProperties{
					PageObjectId: slideID,
					Size: &gslides.Size{
						Height: &gslides.Dimension{Magnitude: height, Unit: "PT"},
						Width:  &gslides.Dimension{Magnitude: width, Unit: "PT"},
					},
				},
			},
		},
	}
	return uc.slidesIntegration.BatchUpdate(presentationID, requests)
}

// CreatePresentationWithSlides creates a new presentation with a specified number of slides.
// If it fails to create the slides, it will delete the partially created presentation.
func (uc *SlidesUseCase) CreatePresentationWithSlides(title string, numSlides int, layoutID string) (*gslides.Presentation, error) {
	presentation, err := uc.CreatePresentation(title)
	if err != nil {
		return nil, err
	}

	if numSlides <= 1 {
		return presentation, nil
	}

	var requests []*gslides.Request
	for i := 1; i < numSlides; i++ {
		slideID := fmt.Sprintf("%s_slide_%d", title, i)
		requests = append(requests, &gslides.Request{
			CreateSlide: &gslides.CreateSlideRequest{
				ObjectId: slideID,
				SlideLayoutReference: &gslides.LayoutReference{
					PredefinedLayout: layoutID,
				},
			},
		})
	}

	_, err = uc.slidesIntegration.BatchUpdate(presentation.PresentationId, requests)
	if err != nil {
		// Attempt to delete the presentation if slide creation fails.
		_, _ = uc.driveIntegration.TrashFile(presentation.PresentationId)
		return nil, fmt.Errorf("failed to create slides, presentation has been deleted: %w", err)
	}

	// We need to get the presentation again to have the updated list of slides.
	return uc.GetPresentation(presentation.PresentationId)
}
