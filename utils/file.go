package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

const MAX_FILE_SIZE_PDF = 10 * 1024 * 1024 // 10 MB for PDF

var AllowedPDFMimeTypes = map[string]bool{
	"application/pdf": true,
}

var AllowedPDFExtensions = map[string]bool{
	".pdf": true,
}

func ValidatePDF(file *multipart.FileHeader) error {
	if file.Size > MAX_FILE_SIZE_PDF {
		return fmt.Errorf("ukuran file terlalu besar, maksimal %.1f MB", float64(MAX_FILE_SIZE_PDF)/(1024*1024))
	}

	contentType := file.Header.Get("Content-Type")
	if !AllowedPDFMimeTypes[contentType] {
		return fmt.Errorf("tipe file tidak diizinkan, hanya application/pdf yang diterima")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedPDFExtensions[ext] {
		return fmt.Errorf("ekstensi file tidak diizinkan, hanya .pdf yang diterima")
	}

	return nil
}

func ValidateFileCategory(category string) (string, error) {
	var formattedCategory string = strings.ToLower(strings.TrimSpace(category))

	if formattedCategory == "" {
		return "", fmt.Errorf("category is required")
	}

	if !isValidCategory(formattedCategory) {
		return "", fmt.Errorf("category contains invalid characters")
	}

	return formattedCategory, nil
}

// HELPER TO PARSE STRING TO UINT64
func ParseUint64(value string) (uint64, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint64 format: %w", err)
	}
	return id, nil
}
