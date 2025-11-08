package rag

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// OpenAIMediaProcessor implements MediaProcessor using OpenAI APIs
type OpenAIMediaProcessor struct {
	apiKey string
	client *http.Client
}

// NewOpenAIMediaProcessor creates a new OpenAI-based media processor
func NewOpenAIMediaProcessor(apiKey string) *OpenAIMediaProcessor {
	return &OpenAIMediaProcessor{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// ProcessImage uses GPT-4 Vision to extract text from images or URLs
func (o *OpenAIMediaProcessor) ProcessImage(imageData io.Reader, mimeType string) (string, error) {
	return o.processImageWithURL("", imageData, mimeType)
}

// ProcessImageURL processes an image from a URL directly
func (o *OpenAIMediaProcessor) ProcessImageURL(imageURL string) (string, error) {
	return o.processImageWithURL(imageURL, nil, "")
}

func (o *OpenAIMediaProcessor) processImageWithURL(imageURL string, imageData io.Reader, mimeType string) (string, error) {
	var imageURLValue string

	if imageURL != "" {
		// Use direct URL
		imageURLValue = imageURL
	} else {
		// Convert binary data to base64
		data, err := io.ReadAll(imageData)
		if err != nil {
			return "", fmt.Errorf("failed to read image data: %w", err)
		}
		base64Image := base64.StdEncoding.EncodeToString(data)
		imageURLValue = fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image)
	}

	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "Extract all text content from this image. Return only the text, no explanations.",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageURLValue,
						},
					},
				},
			},
		},
		"max_tokens": 1000,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from GPT-4 Vision")
	}

	return result.Choices[0].Message.Content, nil
}

// ProcessAudio uses Whisper API to transcribe audio
func (o *OpenAIMediaProcessor) ProcessAudio(audioData io.Reader, mimeType string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the audio file
	part, err := writer.CreateFormFile("file", "audio")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, audioData); err != nil {
		return "", fmt.Errorf("failed to copy audio data: %w", err)
	}

	// Add model parameter
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return "", fmt.Errorf("failed to write model field: %w", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Text, nil
}

// ProcessVideo extracts audio and transcribes using Whisper
func (o *OpenAIMediaProcessor) ProcessVideo(videoData io.Reader, mimeType string) (string, error) {
	// Simplified approach: treat video as audio for transcription
	// In production, use FFmpeg to extract audio track first
	return o.ProcessAudio(videoData, "audio/mp3")
}

// ProcessMultipleImages processes multiple images sequentially (OpenAI doesn't support batch)
func (o *OpenAIMediaProcessor) ProcessMultipleImages(imageDataList []io.Reader, mimeTypes []string) (string, error) {
	var allText []string
	for i, imageData := range imageDataList {
		mimeType := ""
		if i < len(mimeTypes) {
			mimeType = mimeTypes[i]
		}
		text, err := o.ProcessImage(imageData, mimeType)
		if err != nil {
			return "", fmt.Errorf("failed to process image %d: %w", i+1, err)
		}
		if text != "" {
			allText = append(allText, text)
		}
	}
	return strings.Join(allText, "\n\n"), nil
}

// ProcessPDF placeholder - not supported by OpenAI APIs directly
func (o *OpenAIMediaProcessor) ProcessPDF(pdfData io.Reader) (string, error) {
	return "", fmt.Errorf("PDF processing not supported by OpenAI APIs")
}
