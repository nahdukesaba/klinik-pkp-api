package uploads

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CURRENT INSTANCE
type Service struct {
	basePath string
}

// PAYLOAD FROM REQUEST BODY
type FilePayload struct {
	Files    []*multipart.FileHeader
	MaxCount int
	Category string
	RecordID uint64
}

// CONSTRUCTOR
func NewService(basePath string) *Service {
	return &Service{
		basePath: basePath,
	}
}

func (s *Service) SaveImages(payload *FilePayload) ([]string, error) {
	// VALIDATIONS
	if len(payload.Files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	if len(payload.Files) > payload.MaxCount {
		return nil, fmt.Errorf("maximum %d images allowed, got %d", payload.MaxCount, len(payload.Files))
	}

	category, err := utils.ValidateImageCategory(payload.Category)
	if err != nil {
		return nil, err
	}

	recordID := strconv.FormatUint(uint64(payload.RecordID), 10)
	if recordID == "0" {
		return nil, fmt.Errorf("record ID is required")
	}

	var responses []string

	// VALIDATE ALL FILES BEFORE CREATING DIRECTORY
	for _, fileHeader := range payload.Files {
		if err := utils.ValidateImage(fileHeader); err != nil {
			return nil, fmt.Errorf("invalid file %s: %w", fileHeader.Filename, err)
		}
	}

	// DIRECTORY PATH: storage/{{category}}/{{year}}/{{month}}}/{{day}}/{{record_id}}/
	fileDirectoryPath := filepath.Join(s.basePath, category, recordID)

	// CREATE DIRECTORY IF NOT EXISTS
	if err := os.MkdirAll(fileDirectoryPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// PROCESS EACH FILE
	for _, fileHeader := range payload.Files {
		extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
		newFilename := fmt.Sprintf("%s_%s%s", category, uuid.New().String(), extension)
		fileFullPath := filepath.Join(fileDirectoryPath, newFilename)

		// OPEN UPLOADED FILE
		file, err := fileHeader.Open()

		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}

		// CREATE DESTINATION FILE
		destFile, err := os.Create(fileFullPath)

		if err != nil {
			file.Close()

			return nil, fmt.Errorf("failed to create file: %w", err)
		}

		// COPY CONTENTS
		if _, err := io.Copy(destFile, file); err != nil {
			os.Remove(fileFullPath)

			return nil, fmt.Errorf("failed to save file: %w", err)
		}

		file.Close()
		destFile.Close()

		responses = append(responses, fmt.Sprintf("uploads/%s", strings.ReplaceAll(filepath.Join(category, recordID, newFilename), "\\", "/")))
	}

	return responses, nil
}

func (s *Service) DeleteImages(paths []string) error {
	basePath := filepath.Clean(s.basePath)
	basePrefix := basePath + string(os.PathSeparator)

	for _, rawPath := range paths {
		relPath := strings.TrimPrefix(rawPath, "uploads/")
		relPath = filepath.Clean(relPath)

		if relPath == "." || relPath == "" {
			continue
		}

		fullPath := filepath.Clean(filepath.Join(basePath, relPath))

		// PREVENT PATH TRAVERSAL OUTSIDE STORAGE DIRECTORY
		if fullPath != basePath && !strings.HasPrefix(fullPath, basePrefix) {
			return fmt.Errorf("invalid path: %s", rawPath)
		}

		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", rawPath, err)
		}
	}

	return nil
}

func (s *Service) SaveFiles(payload *FilePayload) ([]string, error) {
	// VALIDATIONS
	if len(payload.Files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	if len(payload.Files) > payload.MaxCount {
		return nil, fmt.Errorf("maximum %d files allowed, got %d", payload.MaxCount, len(payload.Files))
	}

	category, err := utils.ValidateFileCategory(payload.Category)
	if err != nil {
		return nil, err
	}

	recordID := strconv.FormatUint(uint64(payload.RecordID), 10)
	if recordID == "0" {
		return nil, fmt.Errorf("record ID is required")
	}

	var responses []string

	// VALIDATE ALL FILES BEFORE CREATING DIRECTORY
	for _, fileHeader := range payload.Files {
		if err := utils.ValidatePDF(fileHeader); err != nil {
			return nil, fmt.Errorf("invalid file %s: %w", fileHeader.Filename, err)
		}
	}

	// DIRECTORY PATH: storage/{{category}}/{{record_id}}/
	fileDirectoryPath := filepath.Join(s.basePath, category, recordID)

	// CREATE DIRECTORY IF NOT EXISTS
	if err := os.MkdirAll(fileDirectoryPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// PROCESS EACH FILE
	for _, fileHeader := range payload.Files {

		extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
		newFilename := fmt.Sprintf("%s_%s%s", category, uuid.New().String(), extension)
		fileFullPath := filepath.Join(fileDirectoryPath, newFilename)

		// OPEN UPLOADED FILE
		file, err := fileHeader.Open()

		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
		}

		// CREATE DESTINATION FILE
		destFile, err := os.Create(fileFullPath)

		if err != nil {
			file.Close()

			return nil, fmt.Errorf("failed to create file: %w", err)
		}

		// COPY CONTENTS
		if _, err := io.Copy(destFile, file); err != nil {
			os.Remove(fileFullPath)

			return nil, fmt.Errorf("failed to save file: %w", err)
		}

		file.Close()
		destFile.Close()

		responses = append(responses, fmt.Sprintf("uploads/%s", strings.ReplaceAll(filepath.Join(category, recordID, newFilename), "\\", "/")))
	}

	return responses, nil
}

func (s *Service) DeleteFiles(paths []string) error {
	basePath := filepath.Clean(s.basePath)
	basePrefix := basePath + string(os.PathSeparator)

	for _, rawPath := range paths {
		relPath := strings.TrimPrefix(rawPath, "uploads/")
		relPath = filepath.Clean(relPath)

		if relPath == "." || relPath == "" {
			continue
		}

		fullPath := filepath.Clean(filepath.Join(basePath, relPath))

		// PREVENT PATH TRAVERSAL OUTSIDE STORAGE DIRECTORY
		if fullPath != basePath && !strings.HasPrefix(fullPath, basePrefix) {
			return fmt.Errorf("invalid path: %s", rawPath)
		}

		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", rawPath, err)
		}
	}

	return nil
}
