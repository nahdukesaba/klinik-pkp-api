package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AllowedImageTypes daftar MIME types yang diizinkan
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedImageExtensions daftar ekstensi yang diizinkan
var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// ImageUploadConfig konfigurasi untuk upload image
type ImageUploadConfig struct {
	MaxFileSize      int64  // dalam bytes
	UploadDir        string // direktori penyimpanan
	AllowedMimeTypes map[string]bool
}

// DefaultImageConfig konfigurasi default untuk upload image
func DefaultImageConfig() ImageUploadConfig {
	return ImageUploadConfig{
		MaxFileSize:      5 * 1024 * 1024, // 5MB
		UploadDir:        "./uploads/images",
		AllowedMimeTypes: AllowedImageTypes,
	}
}

// ValidateImageFile validasi file image yang diupload
func ValidateImageFile(file *multipart.FileHeader, config ImageUploadConfig) error {
	// Validasi ukuran file
	if file.Size > config.MaxFileSize {
		maxMB := float64(config.MaxFileSize) / (1024 * 1024)
		return fmt.Errorf("ukuran file terlalu besar, maksimal %.1f MB", maxMB)
	}

	// Validasi MIME type
	contentType := file.Header.Get("Content-Type")
	if !config.AllowedMimeTypes[contentType] {
		return fmt.Errorf("tipe file tidak diizinkan, hanya %v", getAllowedTypesString(config.AllowedMimeTypes))
	}

	// Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedImageExtensions[ext] {
		return fmt.Errorf("ekstensi file tidak diizinkan")
	}

	return nil
}

// GenerateUniqueFilename membuat nama file yang unik
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	// Gunakan UUID untuk nama file unik
	uniqueName := uuid.New().String()
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, uniqueName, ext)
}

// SanitizeFilename membersihkan nama file dari karakter berbahaya
func SanitizeFilename(filename string) string {
	// Hapus karakter yang tidak aman
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, " ", "_")
	return filename
}

// GetFileExtension mendapatkan ekstensi file
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// getAllowedTypesString helper untuk mendapatkan daftar tipe yang diizinkan
func getAllowedTypesString(types map[string]bool) string {
	var allowed []string
	for t := range types {
		allowed = append(allowed, t)
	}
	return strings.Join(allowed, ", ")
}

// IsImageFile cek apakah file adalah image
func IsImageFile(mimeType string) bool {
	return AllowedImageTypes[mimeType]
}
