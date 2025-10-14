package aiprovider

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// ConvertImageDataToBase64 converts raw image bytes to base64 string
func ConvertImageDataToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ExtractBase64FromDataURI extracts base64 data from data URI format
func ExtractBase64FromDataURI(dataURI string) (string, error) {
	if !strings.HasPrefix(dataURI, "data:") {
		return "", fmt.Errorf("invalid data URI format")
	}

	parts := strings.Split(dataURI, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URI format")
	}

	return parts[1], nil
}

// DetectImageMimeType attempts to detect MIME type from base64 data
func DetectImageMimeType(base64Data string) string {
	// Decode first few bytes to detect format
	data, err := base64.StdEncoding.DecodeString(base64Data[:min(100, len(base64Data))])
	if err != nil {
		return "image/jpeg" // default fallback
	}

	// Check magic bytes
	if len(data) >= 4 {
		if data[0] == 0xFF && data[1] == 0xD8 {
			return "image/jpeg"
		}
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png"
		}
		if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
			return "image/gif"
		}
		if string(data[0:4]) == "RIFF" && len(data) >= 12 && string(data[8:12]) == "WEBP" {
			return "image/webp"
		}
	}

	return "image/jpeg" // default fallback
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
