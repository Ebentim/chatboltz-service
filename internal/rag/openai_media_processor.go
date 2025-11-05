package rag

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIMediaProcessor implements MediaProcessor using OpenAI services.
// Provides AI-powered text extraction from images and audio using GPT-4 Vision
// and Whisper models with high accuracy and multilingual support.
type OpenAIMediaProcessor struct {
	// client is the OpenAI API client
	client *openai.Client
}

// NewOpenAIMediaProcessor creates a new OpenAI-based media processor.
// Uses GPT-4 Vision for image analysis and Whisper for audio transcription.
//
// Required:
//   - OpenAI API key with access to GPT-4 Vision and Whisper models
//
// Parameters:
//   - apiKey: OpenAI API key
//
// Returns:
//   - *OpenAIMediaProcessor: Configured processor ready for use
//   - error: Any error that occurred during client initialization
func NewOpenAIMediaProcessor(apiKey string) (*OpenAIMediaProcessor, error) {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIMediaProcessor{client: client}, nil
}

// ProcessImage analyzes images using GPT-4 Vision to extract text and describe content.
// Provides both OCR capabilities and intelligent content description for images
// without readable text.
//
// Features:
//   - OCR for text extraction from images
//   - Intelligent image description and analysis
//   - Multilingual text recognition
//   - Context-aware content understanding
//   - Handwriting and stylized text recognition
//
// Supported formats: JPEG, PNG, GIF, WebP
//
// Parameters:
//   - imageData: Image file data as io.Reader
//   - mimeType: MIME type of the image (e.g., "image/jpeg", "image/png")
//
// Returns:
//   - string: Extracted text or detailed image description
//   - error: Any error that occurred during image analysis
func (o *OpenAIMediaProcessor) ProcessImage(imageData io.Reader, mimeType string) (string, error) {
	// Read image data
	data, err := io.ReadAll(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Encode image to base64
	base64Image := base64.StdEncoding.EncodeToString(data)
	imageURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)

	// Create chat completion request with image
	resp, err := o.client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4Vision),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.ChatCompletionUserMessageParam{
				Role: openai.F(openai.ChatCompletionUserMessageParamRoleUser),
				Content: openai.F([]openai.ChatCompletionContentPartUnionParam{
					openai.ChatCompletionContentPartTextParam{
						Type: openai.F(openai.ChatCompletionContentPartTextTypeText),
						Text: openai.F("Extract all text from this image. If there's no readable text, provide a detailed description of what you see in the image. Focus on extracting any visible text first, then describe the visual content."),
					},
					openai.ChatCompletionContentPartImageParam{
						Type: openai.F(openai.ChatCompletionContentPartImageTypeImageURL),
						ImageURL: openai.F(openai.ChatCompletionContentPartImageImageURLParam{
							URL: openai.F(imageURL),
						}),
					},
				}),
			},
		}),
		MaxTokens: openai.Int(1000),
	})

	if err != nil {
		return "", fmt.Errorf("failed to analyze image with GPT-4 Vision: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from GPT-4 Vision")
	}

	return resp.Choices[0].Message.Content, nil
}

// ProcessAudio transcribes audio files using OpenAI's Whisper model.
// Provides state-of-the-art speech recognition with support for 99 languages
// and robust performance across various audio conditions.
//
// Features:
//   - 99 language support with automatic detection
//   - Robust to background noise and audio quality
//   - Automatic punctuation and capitalization
//   - Speaker change detection
//   - Technical and domain-specific vocabulary recognition
//
// Supported formats: MP3, MP4, MPEG, MPGA, M4A, WAV, WEBM
//
// Parameters:
//   - audioData: Audio file data as io.Reader
//   - mimeType: MIME type of the audio (e.g., "audio/wav", "audio/mp3")
//
// Returns:
//   - string: Transcribed text from the audio
//   - error: Any error that occurred during transcription
func (o *OpenAIMediaProcessor) ProcessAudio(audioData io.Reader, mimeType string) (string, error) {
	// Read audio data
	data, err := io.ReadAll(audioData)
	if err != nil {
		return "", fmt.Errorf("failed to read audio data: %w", err)
	}

	// Determine file extension from MIME type
	var extension string
	switch mimeType {
	case "audio/wav":
		extension = "wav"
	case "audio/mp3", "audio/mpeg":
		extension = "mp3"
	case "audio/mp4":
		extension = "mp4"
	case "audio/m4a":
		extension = "m4a"
	case "audio/webm":
		extension = "webm"
	default:
		extension = "wav" // Default fallback
	}

	// Create transcription request
	resp, err := o.client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		File:  openai.FileParam(bytes.NewReader(data), fmt.Sprintf("audio.%s", extension), mimeType),
		Model: openai.F(openai.AudioModelWhisper1),
	})

	if err != nil {
		return "", fmt.Errorf("failed to transcribe audio with Whisper: %w", err)
	}

	return resp.Text, nil
}

// ProcessVideo extracts audio from video and transcribes using Whisper.
// Currently supports video files that can be processed directly by Whisper.
//
// Features:
//   - Direct video processing with Whisper (for supported formats)
//   - Multilingual speech recognition
//   - Robust audio extraction and transcription
//
// Supported formats: MP4, WEBM, MOV (formats supported by Whisper)
//
// Parameters:
//   - videoData: Video file data as io.Reader
//   - mimeType: MIME type of the video (e.g., "video/mp4", "video/webm")
//
// Returns:
//   - string: Transcribed text from video's audio track
//   - error: Any error that occurred during processing
func (o *OpenAIMediaProcessor) ProcessVideo(videoData io.Reader, mimeType string) (string, error) {
	// Read video data
	data, err := io.ReadAll(videoData)
	if err != nil {
		return "", fmt.Errorf("failed to read video data: %w", err)
	}

	// Determine file extension from MIME type
	var extension string
	switch mimeType {
	case "video/mp4":
		extension = "mp4"
	case "video/webm":
		extension = "webm"
	case "video/quicktime":
		extension = "mov"
	default:
		extension = "mp4" // Default fallback
	}

	// Use Whisper to transcribe video directly (it can handle video files)
	resp, err := o.client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		File:  openai.FileParam(bytes.NewReader(data), fmt.Sprintf("video.%s", extension), mimeType),
		Model: openai.F(openai.AudioModelWhisper1),
	})

	if err != nil {
		return "", fmt.Errorf("failed to transcribe video with Whisper: %w", err)
	}

	return resp.Text, nil
}

// ProcessPDF extracts text from PDF documents.
// This is a placeholder implementation as OpenAI doesn't directly support PDF processing.
//
// Parameters:
//   - pdfData: PDF file data as io.Reader
//
// Returns:
//   - string: Extracted text content from the PDF
//   - error: Error indicating PDF processing is not implemented
//
// Note: For PDF processing with OpenAI, you would need to:
//  1. Convert PDF pages to images
//  2. Use GPT-4 Vision to extract text from each page
//  3. Combine results
func (o *OpenAIMediaProcessor) ProcessPDF(pdfData io.Reader) (string, error) {
	return "", fmt.Errorf("PDF processing not implemented for OpenAI processor - convert PDF to images and use ProcessImage")
}
