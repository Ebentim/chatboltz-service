package utils

import (
	"net/http"
	"strings"
)

// DetectMimeType detects MIME type from file data
func DetectMimeType(data []byte) string {
	return http.DetectContentType(data)
}

// ValidateMimeType validates if the MIME type is supported for processing
func ValidateMimeType(mimeType string) bool {
	supportedTypes := map[string]bool{
		// Images
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/bmp":  true,
		"image/tiff": true,
		"image/webp": true,
		// Audio
		"audio/mpeg": true, // MP3
		"audio/wav":  true,
		"audio/flac": true,
		"audio/mp4":  true, // M4A
		"audio/ogg":  true,
		"audio/aac":  true,
		// Video
		"video/mp4":       true,
		"video/avi":       true,
		"video/quicktime": true, // MOV
		"video/x-msvideo": true, // AVI
		"video/webm":      true,
		"video/x-flv":     true,
		// Documents
		"application/pdf": true,
		"text/plain":      true,
	}

	return supportedTypes[strings.ToLower(mimeType)]
}

// GetDocumentTypeFromMime maps MIME type to document type
func GetDocumentTypeFromMime(mimeType string) string {
	mimeType = strings.ToLower(mimeType)

	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "video"
	}
	if mimeType == "application/pdf" {
		return "pdf"
	}
	if strings.HasPrefix(mimeType, "text/") {
		return "text"
	}

	return "unknown"
}
