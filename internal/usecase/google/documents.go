package usecase

import (
	gdocs "google.golang.org/api/docs/v1"
)

// DocsIntegration defines the interface for docs operations.
type DocsIntegration interface {
	CreateDocument(title string) (*gdocs.Document, error)
	GetDocument(documentID string) (*gdocs.Document, error)
	BatchUpdate(documentID string, requests []*gdocs.Request) (*gdocs.BatchUpdateDocumentResponse, error)
}

// DocsUseCase handles the business logic for docs operations.
type DocsUseCase struct {
	docsIntegration  DocsIntegration
	driveIntegration DriveIntegration
}

// NewDocsUseCase creates a new DocsUseCase.
func NewDocsUseCase(di DocsIntegration, driveIntegration DriveIntegration) *DocsUseCase {
	return &DocsUseCase{
		docsIntegration:  di,
		driveIntegration: driveIntegration,
	}
}

// CreateDocument creates a new Google Doc.
func (uc *DocsUseCase) CreateDocument(title string) (*gdocs.Document, error) {
	return uc.docsIntegration.CreateDocument(title)
}

// GetDocument retrieves a Google Doc.
func (uc *DocsUseCase) GetDocument(documentID string) (*gdocs.Document, error) {
	return uc.docsIntegration.GetDocument(documentID)
}

// DeleteDocument moves a Google Doc to the trash.
func (uc *DocsUseCase) DeleteDocument(documentID string) error {
	_, err := uc.driveIntegration.TrashFile(documentID)
	return err
}

// InsertText inserts text into a Google Doc at a specific index.
func (uc *DocsUseCase) InsertText(documentID string, text string, index int64) (*gdocs.BatchUpdateDocumentResponse, error) {
	requests := []*gdocs.Request{
		{
			InsertText: &gdocs.InsertTextRequest{
				Text: text,
				Location: &gdocs.Location{
					Index: index,
				},
			},
		},
	}
	return uc.docsIntegration.BatchUpdate(documentID, requests)
}

// AppendText appends text to the end of a Google Doc.
func (uc *DocsUseCase) AppendText(documentID string, text string) (*gdocs.BatchUpdateDocumentResponse, error) {
	doc, err := uc.GetDocument(documentID)
	if err != nil {
		return nil, err
	}
	if doc.Body == nil || len(doc.Body.Content) == 0 {
		// If the document body is empty, insert at index 1.
		return uc.InsertText(documentID, text, 1)
	}

	// The body content ends at the last index of the last content element.
	// The actual end index for new content is the one after the last character.
	// The document's body has a content array, and the last element has an endIndex.
	endIndex := doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1

	requests := []*gdocs.Request{
		{
			InsertText: &gdocs.InsertTextRequest{
				Text: text,
				Location: &gdocs.Location{
					Index: endIndex,
				},
			},
		},
	}
	return uc.docsIntegration.BatchUpdate(documentID, requests)
}
