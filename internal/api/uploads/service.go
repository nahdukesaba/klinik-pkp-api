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

	if len(payload.Files) > 4 {
		return nil, fmt.Errorf("maximum 4 images allowed, got %d", len(payload.Files))
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

	// DIRECTORY PATH: storage/{{category}}/{{year}}/{{month}}}/{{day}}/{{record_id}}/
	fileDirectoryPath := filepath.Join(s.basePath, category, recordID)

	// CREATE DIRECTORY IF NOT EXISTS
	if err := os.MkdirAll(fileDirectoryPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// PROCESS EACH FILE
	for _, fileHeader := range payload.Files {
		if err := utils.ValidateImage(fileHeader); err != nil {
			return nil, fmt.Errorf("invalid file %s: %w", fileHeader.Filename, err)
		}

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

func (s *Service) DeleteImages(category string, recordID uint64) error {
	if strings.Contains(category, "..") || strings.Contains(strconv.FormatUint(recordID, 10), "..") {
		return fmt.Errorf("invalid path")
	}

	// FIND ALL FILES IN THE RECORD DIRECTORY (storage/{{category}}/{{record_id}}/)
	matches, err := filepath.Glob(filepath.Join(s.basePath, category, strconv.FormatUint(recordID, 10)))

	if err != nil {
		return fmt.Errorf("failed to find record directories: %w", err)
	}

	for _, dir := range matches {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to delete directory %s: %w", dir, err)
		}
	}

	return nil
}
