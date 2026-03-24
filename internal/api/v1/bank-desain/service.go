package bank_desain

import (
	"fmt"
	"klinik-pkp-api/internal/api/v1/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"

	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	validator     *utils.Validator
	uploadService *uploads.Service
}

// PAYLOAD FROM REQUEST BODY
type BankDesainPayload struct {
	Name          string                  `form:"name" validate:"required,ne=Null,ne=null,ne=NULL"`
	Type          string                  `form:"type" validate:"required,ne=Null,ne=null,ne=NULL"`
	BedroomCount  uint64                  `form:"bedroom_count" validate:"required"`
	BathroomCount uint64                  `form:"bathroom_count" validate:"required"`
	TotalArea     uint64                  `form:"total_area" validate:"required"`
	HasGarage     *bool                   `form:"has_garage" validate:"required"`
	Images        []*multipart.FileHeader `form:"images" validate:"required"`
	Files         []*multipart.FileHeader `form:"files" validate:"required"`
}

// FILTER FOR QUERY PARAMETERS
type BankDesainFilter struct {
	Type          *string
	BedroomCount  *uint64
	BathroomCount *uint64
	HasGarage     *bool
	Name          *string
}

// CONSTRUCTOR
func NewService(db *gorm.DB, uploadService *uploads.Service) *Service {
	return &Service{
		db:            db,
		validator:     utils.NewValidator(),
		uploadService: uploadService,
	}
}

func (s *Service) GetBankDesain(filter BankDesainFilter) ([]BankDesain, error) {
	var bankDesains []BankDesain

	query := s.db.Model(&BankDesain{})

	// FILTER BY TYPE
	if filter.Type != nil {
		query = query.Where("bank_desain.type = ?", *filter.Type)
	}

	// FILTER BY BEDROOM COUNT
	if filter.BedroomCount != nil {
		query = query.Where("bank_desain.bedroom_count = ?", *filter.BedroomCount)
	}

	// FILTER BY BATHROOM COUNT
	if filter.BathroomCount != nil {
		query = query.Where("bank_desain.bathroom_count = ?", *filter.BathroomCount)
	}

	// FILTER BY HAS GARAGE
	if filter.HasGarage != nil {
		query = query.Where("bank_desain.has_garage = ?", *filter.HasGarage)
	}

	// FILTER BY NAME (SEARCH)
	if filter.Name != nil {
		query = query.Where("bank_desain.name LIKE ?", "%"+*filter.Name+"%")
	}

	if err := query.Find(&bankDesains).Error; err != nil {
		return nil, err
	}

	return bankDesains, nil
}

func (s *Service) GetBankDesainById(id uint64) (*BankDesain, error) {
	var bankDesain BankDesain

	if err := s.db.First(&bankDesain, id).Error; err != nil {
		return nil, err
	}

	return &bankDesain, nil
}

func (s *Service) AddBankDesain(payload *BankDesainPayload) (*BankDesain, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// CREATE RECORD FIRST TO GET ID
	newBankDesain := BankDesain{
		Name:          payload.Name,
		Type:          payload.Type,
		BedroomCount:  payload.BedroomCount,
		BathroomCount: payload.BathroomCount,
		TotalArea:     payload.TotalArea,
		HasGarage:     payload.HasGarage,
	}

	if err := tx.Create(&newBankDesain).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 8,
			Category: "bank_desain",
			RecordID: newBankDesain.ID,
		})

		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newBankDesain.ImageURLs = imageResponses
	}

	if len(payload.Files) > 0 {
		fileResponses, err := s.uploadService.SaveFiles(&uploads.FilePayload{
			Files:    payload.Files,
			MaxCount: 1,
			Category: "bank_desain",
			RecordID: newBankDesain.ID,
		})

		if err != nil {
			// CLEANUP: DELETE THE CREATED RECORD AND UPLOADED IMAGES
			tx.Rollback()
			// DELETE UPLOADED IMAGES FROM STORAGE (ASYNC, NON-BLOCKING)
			go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
			return nil, fmt.Errorf("failed to upload files: %w", err)
		}

		newBankDesain.FileURLs = fileResponses
	}

	if err := tx.Save(&newBankDesain).Error; err != nil {
		tx.Rollback()
		// CLEANUP: DELETE UPLOADED FILES AND IMAGES
		go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
		go s.uploadService.DeleteFiles(newBankDesain.FileURLs)
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		// CLEANUP: DELETE UPLOADED FILES AND IMAGES
		go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
		go s.uploadService.DeleteFiles(newBankDesain.FileURLs)
		return nil, err
	}

	return &newBankDesain, nil
}

func (s *Service) EditBankDesainById(id uint64, payload *BankDesainPayload) error {
	// BEGIN TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// CHECK IF RECORD EXISTS
	existingBankDesain := BankDesain{}
	if err := tx.First(&existingBankDesain, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	existingImageURLs := append([]string(nil), existingBankDesain.ImageURLs...)
	existingFileURLs := append([]string(nil), existingBankDesain.FileURLs...)
	isImagesUpdated := false
	isFilesUpdated := false
	newBankDesain := BankDesain{
		Name:          payload.Name,
		Type:          payload.Type,
		BedroomCount:  payload.BedroomCount,
		BathroomCount: payload.BathroomCount,
		TotalArea:     payload.TotalArea,
		HasGarage:     payload.HasGarage,
	}

	if len(payload.Images) > 0 {
		// UPLOAD NEW IMAGES FIRST; OLD IMAGES STAY UNTIL TRANSACTION COMMITS
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 8,
			Category: "bank_desain",
			RecordID: id,
		})

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		newBankDesain.ImageURLs = imageResponses
		isImagesUpdated = true
	}

	if len(payload.Files) > 0 {
		// UPLOAD NEW FILES FIRST; OLD FILES STAY UNTIL TRANSACTION COMMITS
		fileResponses, err := s.uploadService.SaveFiles(&uploads.FilePayload{
			Files:    payload.Files,
			MaxCount: 1,
			Category: "bank_desain",
			RecordID: id,
		})

		if err != nil {
			tx.Rollback()
			if isImagesUpdated {
				go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
			}
			return fmt.Errorf("failed to upload new files: %w", err)
		}

		newBankDesain.FileURLs = fileResponses
		isFilesUpdated = true
	}

	result := tx.Model(&BankDesain{}).Where("id = ?", id).Updates(newBankDesain)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
		}
		if isFilesUpdated {
			go s.uploadService.DeleteFiles(newBankDesain.FileURLs)
		}
		return result.Error
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newBankDesain.ImageURLs)
		}
		if isFilesUpdated {
			go s.uploadService.DeleteFiles(newBankDesain.FileURLs)
		}
		return err
	}

	if isImagesUpdated && len(existingImageURLs) > 0 {
		if err := s.uploadService.DeleteImages(existingImageURLs); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}
	}

	if isFilesUpdated && len(existingFileURLs) > 0 {
		if err := s.uploadService.DeleteFiles(existingFileURLs); err != nil {
			fmt.Printf("Warning: failed to delete old files: %v\n", err)
		}
	}

	return nil
}

func (s *Service) DeleteBankDesainById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingBankDesain BankDesain
	if err := tx.First(&existingBankDesain, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// DELETE IMAGES AND FILES (IF ANY) - NON-BLOCKING, ASYNC CLEANUP
	go s.uploadService.DeleteImages(existingBankDesain.ImageURLs)
	go s.uploadService.DeleteFiles(existingBankDesain.FileURLs)

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&BankDesain{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
