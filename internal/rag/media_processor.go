// Package rag provides media processing services for converting various file types to text.
// This enables the RAG system to handle images, audio, video, and PDF files by extracting
// their textual content for embedding generation.
package rag

import (
	"io"
)

// MediaProcessor defines the interface for converting different media types to text.
// Each media type requires specific processing to extract meaningful textual content
// that can be embedded and searched.
type MediaProcessor interface {
	// ProcessImage extracts text from images using OCR (Optical Character Recognition).
	// Supports formats: JPEG, PNG, GIF, BMP, TIFF, WebP
	// Returns extracted text and any error that occurred during processing.
	ProcessImage(imageData io.Reader, mimeType string) (string, error)

	// ProcessImageURL extracts text from images via URL (when supported by processor).
	// Returns extracted text and any error that occurred during processing.
	ProcessImageURL(imageURL string) (string, error)

	// ProcessMultipleImages processes multiple images in a single batch request.
	// Returns combined text from all images and any error that occurred.
	ProcessMultipleImages(imageDataList []io.Reader, mimeTypes []string) (string, error)

	// ProcessAudio transcribes audio files to text using speech-to-text services.
	// Supports formats: MP3, WAV, FLAC, M4A, OGG, AAC
	// Returns transcribed text and any error that occurred during processing.
	ProcessAudio(audioData io.Reader, mimeType string) (string, error)

	// ProcessVideo extracts audio track and transcribes to text.
	// Supports formats: MP4, AVI, MOV, MKV, WebM, FLV
	// Returns transcribed text from audio track and any error that occurred.
	ProcessVideo(videoData io.Reader, mimeType string) (string, error)

	// ProcessPDF extracts text content from PDF documents.
	// Handles both text-based PDFs and scanned PDFs (with OCR fallback).
	// Returns extracted text and any error that occurred during processing.
	ProcessPDF(pdfData io.Reader) (string, error)
}
