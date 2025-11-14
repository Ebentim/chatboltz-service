package rag

import (
	"fmt"
	"io"
)

// MediaProcessorFactory creates media processors with fallback hierarchy
type MediaProcessorFactory struct {
	openaiKey  string
	googleKey  string
	cohereKey  string
	processors []MediaProcessor
}

// NewMediaProcessorFactory creates a factory with provider hierarchy: OpenAI -> Google -> Cohere
func NewMediaProcessorFactory(openaiKey, googleKey, cohereKey string) *MediaProcessorFactory {
	factory := &MediaProcessorFactory{
		openaiKey: openaiKey,
		googleKey: googleKey,
		cohereKey: cohereKey,
	}

	// Initialize processors in order of preference
	if openaiKey != "" {
		factory.processors = append(factory.processors, NewOpenAIMediaProcessor(openaiKey))
	}
	if googleKey != "" {
		processor, _ := NewGoogleMediaProcessor()
		factory.processors = append(factory.processors, processor)

	}

	return factory
}

// ProcessImage processes image with fallback
func (f *MediaProcessorFactory) ProcessImage(imageData io.Reader, mimeType string) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessImage(imageData, mimeType)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}

// ProcessImageURL processes image URL with fallback
func (f *MediaProcessorFactory) ProcessImageURL(imageURL string) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessImageURL(imageURL)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}

// ProcessMultipleImages processes multiple images with fallback
func (f *MediaProcessorFactory) ProcessMultipleImages(imageDataList []io.Reader, mimeTypes []string) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessMultipleImages(imageDataList, mimeTypes)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}

// ProcessAudio processes audio with fallback (OpenAI Whisper preferred)
func (f *MediaProcessorFactory) ProcessAudio(audioData io.Reader, mimeType string) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessAudio(audioData, mimeType)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}

// ProcessVideo processes video with fallback
func (f *MediaProcessorFactory) ProcessVideo(videoData io.Reader, mimeType string) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessVideo(videoData, mimeType)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}

// ProcessPDF processes PDF with fallback
func (f *MediaProcessorFactory) ProcessPDF(pdfData io.Reader) (string, error) {
	for i, processor := range f.processors {
		result, err := processor.ProcessPDF(pdfData)
		if err == nil {
			return result, nil
		}
		if i == len(f.processors)-1 {
			return "", fmt.Errorf("all processors failed: %w", err)
		}
	}
	return "", fmt.Errorf("no processors available")
}
