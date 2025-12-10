// Package usecase provides business logic for RAG operations.
// This package contains use cases for training agents and retrieving context.
package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/rag"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
	"github.com/alpinesboltltd/boltz-ai/internal/scraper"
	"github.com/alpinesboltltd/boltz-ai/internal/utils"
	"gorm.io/gorm"
)

// TrainingUseCase handles agent training operations including document processing.
// It orchestrates the training workflow from document ingestion to knowledge base creation.
type TrainingUseCase struct {
	// ragService provides RAG operations for document processing and retrieval
	ragService *rag.RAGService
	// agentRepo provides access to agent data
	agentRepo repository.AgentRepositoryInterface
	// scraperService provides web scraping capabilities
	scraperService *scraper.Service
}

// NewTrainingUseCase creates a new training use case with the required dependencies.
// It initializes the RAG service with Cohere client and media processor factory.
//
// Parameters:
//   - cohereKey: Cohere API key for embeddings
//   - openaiKey: OpenAI API key for media processing (preferred)
//   - googleKey: Google API key for media processing (fallback)
//   - db: Database connection for repository operations
//   - agentRepo: Repository for agent operations
//
// Returns:
//   - *TrainingUseCase: Configured training use case
func NewTrainingUseCase(cohereKey, openaiKey, googleKey, pineconeKey, pineconeIndex, vectorDBType string, db *gorm.DB, agentRepo repository.AgentRepositoryInterface) (*TrainingUseCase, error) {
	// Trim whitespace from API keys (handles trailing newlines from .env)
	cohereKey = strings.TrimSpace(cohereKey)
	openaiKey = strings.TrimSpace(openaiKey)
	googleKey = strings.TrimSpace(googleKey)
	pineconeKey = strings.TrimSpace(pineconeKey)

	if cohereKey == "" {
		return nil, fmt.Errorf("COHERE_API_KEY is empty or not set")
	}

	cohere, err := rag.NewCohereClient(cohereKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cohere client: %w", err)
	}

	// Initialize media processor factory with fallback hierarchy
	mediaProcessor := rag.NewMediaProcessorFactory(openaiKey, googleKey, cohereKey)

	ragRepo := repository.NewRAGRepository(db)

	// Initialize vector DB based on type
	var vectorDB rag.VectorDB
	if vectorDBType == "pinecone" && pineconeKey != "" {
		vectorDB, err = rag.NewPineconeDB(pineconeKey, pineconeIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
		}
	}

	ragService := rag.NewRAGService(cohere, ragRepo, mediaProcessor, vectorDB, vectorDBType)

	return &TrainingUseCase{
		ragService:     ragService,
		agentRepo:      agentRepo,
		scraperService: scraper.NewService(nil),
	}, nil
}

// ProcessDocument processes a single document for an agent's knowledge base.
// This is the primary method for adding new training content to an agent.
//
// Supported document types:
//   - text: Plain text content
//   - pdf: Extracted PDF text
//   - audio: Transcribed audio content
//   - video: Transcribed video content
//   - faq: FAQ content in Q&A format
//   - image: OCR text or image descriptions
//
// Parameters:
//   - agentID: ID of the agent to train
//   - title: Human-readable title for the document
//   - content: Raw text content to process
//   - docType: Type of document for appropriate processing
//   - sourceURL: Optional URL of the original document
//
// Returns:
//   - error: Any error that occurred during processing
func (t *TrainingUseCase) ProcessDocument(agentID, title, content string, docType entity.DocumentType, sourceURL *string) error {
	return t.ragService.ProcessDocument(agentID, title, docType, content, sourceURL)
}

// ProcessFileWithMimeDetection processes a file with automatic MIME type detection.
// This method detects the file type and validates it before processing.
//
// Parameters:
//   - agentID: ID of the agent to train
//   - title: Human-readable title for the document
//   - fileData: Raw file data
//   - mimeType: Optional MIME type (if empty, will be detected)
//   - sourceURL: Optional URL of the original document
//
// Returns:
//   - error: Any error that occurred during processing
func (t *TrainingUseCase) ProcessFileWithMimeDetection(agentID, title string, fileData []byte, mimeType string, sourceURL *string) error {
	// Detect MIME type if not provided
	if mimeType == "" {
		mimeType = utils.DetectMimeType(fileData)
	}

	// Validate MIME type
	if !utils.ValidateMimeType(mimeType) {
		return fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Map MIME type to document type
	docTypeStr := utils.GetDocumentTypeFromMime(mimeType)
	docType := entity.DocumentType(docTypeStr)

	// For text files, convert bytes to string
	if docType == entity.DocumentTypeText {
		content := string(fileData)
		return t.ragService.ProcessDocument(agentID, title, docType, content, sourceURL)
	}

	// Process media files through MediaProcessor
	return t.ragService.ProcessMediaFile(agentID, title, docType, fileData, mimeType, sourceURL)
}

// TrainAgentFromLegacyData migrates legacy training data to the new RAG system.
// This method processes existing TrainingData records and converts them to the new format.
//
// Parameters:
//   - agentID: ID of the agent whose legacy data to migrate
//
// Returns:
//   - error: Any error that occurred during migration
func (t *TrainingUseCase) TrainAgentFromLegacyData(agentID string) error {
	agent, err := t.agentRepo.GetAgent(agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent: %w", err)
	}

	if len(agent.TrainingData) == 0 {
		return fmt.Errorf("no training data found for agent %s", agentID)
	}

	// Process legacy training data
	for _, td := range agent.TrainingData {
		if td.IsActive {
			for _, text := range td.Content {
				if err := t.ragService.ProcessDocument(agentID, text.Title, entity.DocumentTypeText, text.Content, nil); err != nil {
					return fmt.Errorf("failed to process legacy text: %w", err)
				}
			}
		}
	}

	return nil
}

// GetAgentDocuments retrieves all training documents for an agent.
// This includes document metadata and processing status.
//
// Parameters:
//   - agentID: ID of the agent whose documents to retrieve
//
// Returns:
//   - []entity.TrainingDocument: Array of training documents
//   - error: Any error that occurred during retrieval
func (t *TrainingUseCase) GetAgentDocuments(agentID string) ([]entity.TrainingDocument, error) {
	return t.ragService.GetAgentDocuments(agentID)
}

// DeleteAgentTraining removes all training data for an agent.
// This clears the agent's knowledge base completely.
//
// Parameters:
//   - agentID: ID of the agent whose training data to delete
//
// Returns:
//   - error: Any error that occurred during deletion
func (t *TrainingUseCase) DeleteAgentTraining(agentID string) error {
	return t.ragService.DeleteAgentDocuments(agentID)
}

// QueryKnowledgeBase performs RAG query on agent's knowledge base
func (t *TrainingUseCase) QueryKnowledgeBase(ragQuery entity.RAGQuery) (*entity.RAGResponse, error) {
	return t.ragService.Query(ragQuery)
}

// ProcessURL scrapes a URL and processes the content for training
func (t *TrainingUseCase) ProcessURL(agentID, url, title string, trace bool, maxPages int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	opts := scraper.ScrapeOptions{
		Trace:    trace,
		MaxPages: maxPages,
	}

	if maxPages == 0 {
		opts.MaxPages = 5 // Default to 5 pages
	}

	result, err := t.scraperService.Scrape(ctx, url, opts)
	if err != nil {
		return fmt.Errorf("failed to scrape URL: %w", err)
	}

	// Process each scraped page
	for i, page := range result.Pages {
		pageTitle := title
		if len(result.Pages) > 1 {
			pageTitle = fmt.Sprintf("%s - Page %d", title, i+1)
		}
		if page.Title != "" {
			pageTitle = page.Title
		}

		// Combine all text content from the page
		content := t.extractTextFromPage(page)
		if content == "" {
			continue // Skip empty pages
		}

		// Process the page content
		if err := t.ragService.ProcessDocument(agentID, pageTitle, entity.DocumentTypeText, content, &page.URL); err != nil {
			return fmt.Errorf("failed to process page %s: %w", page.URL, err)
		}
	}

	return nil
}

// extractTextFromPage extracts meaningful text content from scraped page data
func (t *TrainingUseCase) extractTextFromPage(page scraper.PageData) string {
	var parts []string

	// Add title if available
	if page.Title != "" {
		parts = append(parts, "Title: "+page.Title)
	}

	// Add sections (headings and paragraphs)
	for _, section := range page.Sections {
		if section.Text != "" {
			// Format headings differently
			if strings.HasPrefix(section.Tag, "h") {
				parts = append(parts, "\n"+section.Text+"\n")
			} else {
				parts = append(parts, section.Text)
			}
		}
	}

	return strings.Join(parts, "\n")
}
