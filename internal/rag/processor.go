package rag

import (
	"fmt"
	"io"
	"strings"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

// ContentProcessor handles the processing of different document types into chunks with embeddings.
// It applies type-specific chunking strategies, converts media to text, and generates multilingual embeddings using Cohere.
// Supports content from images (OCR/descriptions), audio (transcripts), video (transcripts), PDFs, and text.
type ContentProcessor struct {
	// cohere is the client used to generate multilingual embeddings for text chunks
	cohere *CohereClient
	// mediaProcessor handles conversion of media files to text
	mediaProcessor MediaProcessor
}

// NewContentProcessor creates a new content processor with the provided clients.
//
// Parameters:
//   - cohere: Configured Cohere client for generating multilingual embeddings
//   - mediaProcessor: Media processor for converting files to text (optional, can be nil for text-only processing)
//
// Returns:
//   - *ContentProcessor: Ready-to-use content processor
func NewContentProcessor(cohere *CohereClient, mediaProcessor MediaProcessor) *ContentProcessor {
	return &ContentProcessor{
		cohere:         cohere,
		mediaProcessor: mediaProcessor,
	}
}

// ProcessDocument processes a training document into chunks with multilingual embeddings.
// It applies document-type-specific chunking strategies and generates language-agnostic embeddings for each chunk.
// Supports content in 100+ languages from various sources:
//
// Chunking strategies by document type:
//   - Text: 500 character chunks with word boundaries (any language)
//   - PDF: 800 character chunks (larger for better context, multilingual)
//   - FAQ: Split by Q&A pairs (double newlines, language-agnostic)
//   - Audio: 600 character chunks (transcribed content, any language)
//   - Video: 600 character chunks (transcribed content, any language)
//   - Image: Single chunk (OCR text or descriptions, multilingual)
//
// Parameters:
//   - doc: The training document metadata
//   - content: The raw text content to process (any language)
//
// Returns:
//   - []entity.DocumentChunk: Array of processed chunks with multilingual embeddings
//   - error: Any error that occurred during processing
func (p *ContentProcessor) ProcessDocument(doc *entity.TrainingDocument, content string) ([]entity.DocumentChunk, error) {
	var chunks []string
	var metadata map[string]string

	switch doc.DocumentType {
	case entity.DocumentTypeText:
		chunks = p.chunkText(content, 500)
		metadata = map[string]string{"type": "text", "title": doc.Title}
	case entity.DocumentTypePDF:
		chunks = p.chunkText(content, 800) // Larger chunks for PDFs
		metadata = map[string]string{"type": "pdf", "title": doc.Title, "source": *doc.SourceURL}
	case entity.DocumentTypeFAQ:
		chunks = p.processFAQ(content)
		metadata = map[string]string{"type": "faq", "title": doc.Title}
	case entity.DocumentTypeAudio:
		chunks = p.chunkText(content, 600) // Transcribed audio
		metadata = map[string]string{"type": "audio_transcript", "title": doc.Title, "source": *doc.SourceURL}
	case entity.DocumentTypeVideo:
		chunks = p.chunkText(content, 600) // Transcribed video
		metadata = map[string]string{"type": "video_transcript", "title": doc.Title, "source": *doc.SourceURL}
	case entity.DocumentTypeImage:
		chunks = []string{content} // Image description/OCR text
		metadata = map[string]string{"type": "image_description", "title": doc.Title, "source": *doc.SourceURL}
	default:
		return nil, fmt.Errorf("unsupported document type: %s", doc.DocumentType)
	}

	// Generate embeddings
	embeddings, err := p.cohere.Embed(chunks, "search_document")
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Create document chunks
	var docChunks []entity.DocumentChunk
	for i, chunk := range chunks {
		docChunks = append(docChunks, entity.DocumentChunk{
			DocumentID: doc.ID,
			AgentID:    doc.AgentID,
			Content:    chunk,
			ChunkIndex: i,
			Metadata:   metadata,
			Embedding:  embeddings[i],
		})
	}

	return docChunks, nil
}

// chunkText splits text into chunks of approximately maxChunkSize characters.
// It respects word boundaries to avoid splitting words in the middle.
//
// Parameters:
//   - text: The input text to chunk
//   - maxChunkSize: Maximum size of each chunk in characters
//
// Returns:
//   - []string: Array of text chunks
func (p *ContentProcessor) chunkText(text string, maxChunkSize int) []string {
	words := strings.Fields(text)
	var chunks []string
	var currentChunk []string
	currentSize := 0

	for _, word := range words {
		wordSize := len(word) + 1
		if currentSize+wordSize > maxChunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{word}
			currentSize = len(word)
		} else {
			currentChunk = append(currentChunk, word)
			currentSize += wordSize
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}

// processFAQ processes FAQ content by splitting it into individual Q&A pairs.
// It expects FAQ content to be formatted with double newlines separating each Q&A pair.
//
// Parameters:
//   - content: FAQ content with Q&A pairs separated by double newlines
//
// Returns:
//   - []string: Array of individual Q&A pair chunks
func (p *ContentProcessor) processFAQ(content string) []string {
	// Split FAQ into Q&A pairs
	pairs := strings.Split(content, "\n\n")
	var chunks []string

	for _, pair := range pairs {
		if strings.TrimSpace(pair) != "" {
			chunks = append(chunks, strings.TrimSpace(pair))
		}
	}

	return chunks
}

// ProcessMediaToText converts media files to text using the configured media processor.
// This method handles the conversion of images, audio, video, and PDF files to text
// before they can be processed into chunks and embeddings.
//
// Parameters:
//   - mediaData: The raw media file data
//   - docType: Type of document being processed
//   - mimeType: MIME type of the media file (optional for PDFs)
//
// Returns:
//   - string: Extracted text content from the media file
//   - error: Any error that occurred during media processing
func (p *ContentProcessor) ProcessMediaToText(mediaData io.Reader, docType entity.DocumentType, mimeType string) (string, error) {
	if p.mediaProcessor == nil {
		return "", fmt.Errorf("media processor not configured")
	}

	switch docType {
	case entity.DocumentTypeImage:
		return p.mediaProcessor.ProcessImage(mediaData, mimeType)
	case entity.DocumentTypeAudio:
		return p.mediaProcessor.ProcessAudio(mediaData, mimeType)
	case entity.DocumentTypeVideo:
		return p.mediaProcessor.ProcessVideo(mediaData, mimeType)
	case entity.DocumentTypePDF:
		return p.mediaProcessor.ProcessPDF(mediaData)
	default:
		return "", fmt.Errorf("unsupported media type for processing: %s", docType)
	}
}
