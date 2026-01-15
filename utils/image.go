package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const MAX_FILE_SIZE = 10 * 1024 * 1024

var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"image/x-icon":    true,
	"image/svg+xml":   true,
	"application/pdf": true,
}

var AllowedImageExtensions = map[string]bool{
	".jpeg": true,
	".jpg":  true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

func ValidateImage(file *multipart.FileHeader) error {
	if file.Size > MAX_FILE_SIZE {
		return fmt.Errorf("ukuran file terlalu besar, maksimal %.1f MB", float64(MAX_FILE_SIZE)/(1024*1024))
	}

	if !AllowedMimeTypes[file.Header.Get("Content-Type")] {
		return fmt.Errorf("tipe file tidak diizinkan, hanya %v", getAllowedTypesString(AllowedMimeTypes))
	}

	if !AllowedImageExtensions[strings.ToLower(filepath.Ext(file.Filename))] {
		return fmt.Errorf("ekstensi file tidak diizinkan")
	}

	return nil
}

func ValidateImageCategory(category string) (string, error) {
	var formattedCategory string = strings.ToLower(strings.TrimSpace(category))

	if formattedCategory == "" {
		return "", fmt.Errorf("category is required")
	}

	if !isValidCategory(formattedCategory) {
		return "", fmt.Errorf("category contains invalid characters")
	}

	return formattedCategory, nil
}

// HELPER TO GET ALLOWED MIME TYPES AS STRING
func getAllowedTypesString(types map[string]bool) string {
	var allowed []string

	for t := range types {
		allowed = append(allowed, t)
	}

	return strings.Join(allowed, ", ")
}

// CHECK IF CATEGORY ONLY CONTAINS ALPHANUMERIC CHARACTERS AND UNDERSCORES
func isValidCategory(category string) bool {
	for _, char := range category {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return len(category) > 0
}
