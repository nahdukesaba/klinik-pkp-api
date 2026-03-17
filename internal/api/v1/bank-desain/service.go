package bank_desain

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/v1/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"
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
	// CREATE RECORD FIRST TO GET ID
	newBankDesain := BankDesain{
		Name:          payload.Name,
		Type:          payload.Type,
		BedroomCount:  payload.BedroomCount,
		BathroomCount: payload.BathroomCount,
		TotalArea:     payload.TotalArea,
		HasGarage:     payload.HasGarage,
	}

	if err := s.db.Create(&newBankDesain).Error; err != nil {
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 4,
			Category: "bank_desain",
			RecordID: newBankDesain.ID,
		})

		if err != nil {
			s.db.Delete(&newBankDesain)

			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newBankDesain.ImageURLs = imageResponses
	}

	if len(payload.Files) > 0 {
		fileResponses, err := s.uploadService.SaveFiles(&uploads.FilePayload{
			Files:    payload.Files,
			MaxCount: 8,
			Category: "bank_desain",
			RecordID: newBankDesain.ID,
		})

		if err != nil {
			s.db.Delete(&newBankDesain)

			return nil, fmt.Errorf("failed to upload files: %w", err)
		}

		newBankDesain.FileURLs = fileResponses
	}

	if err := s.db.Save(&newBankDesain).Error; err != nil {
		return nil, err
	}

	return &newBankDesain, nil
}

func (s *Service) EditBankDesainById(id uint64, payload *BankDesainPayload) error {
	// CHECK IF RECORD EXISTS
	existingBankDesain := BankDesain{}

	if err := s.db.First(&existingBankDesain, id).Error; err != nil {
		return err
	}

	newBankDesain := BankDesain{
		Name:          payload.Name,
		Type:          payload.Type,
		BedroomCount:  payload.BedroomCount,
		BathroomCount: payload.BathroomCount,
		TotalArea:     payload.TotalArea,
		HasGarage:     payload.HasGarage,
	}

	if len(payload.Images) > 0 {
		// DELETE OLD IMAGES
		if err := s.uploadService.DeleteImages("bank_desain", id); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}

		// UPLOAD NEW IMAGES
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 8,
			Category: "bank_desain",
			RecordID: id,
		})

		if err != nil {
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		println(imageResponses)

		newBankDesain.ImageURLs = imageResponses
	}

	if len(payload.Files) > 0 {
		// DELETE OLD FILES
		if err := s.uploadService.DeleteFiles("bank_desain", id); err != nil {
			fmt.Printf("Warning: failed to delete old files: %v\n", err)
		}

		// UPLOAD NEW FILES
		fileResponses, err := s.uploadService.SaveFiles(&uploads.FilePayload{
			Files:    payload.Files,
			MaxCount: 1,
			Category: "bank_desain",
			RecordID: id,
		})

		if err != nil {
			return fmt.Errorf("failed to upload new files: %w", err)
		}

		newBankDesain.FileURLs = fileResponses
	}

	result := s.db.Model(&BankDesain{}).Where("id = ?", id).Updates(newBankDesain)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Service) DeleteBankDesainById(id uint64) error {
	// DELETE IMAGES (IF ANY)
	if err := s.uploadService.DeleteImages("bank_desain", id); err != nil {
		fmt.Printf("Warning: failed to delete images: %v\n", err)
	}

	result := s.db.Delete(&BankDesain{}, id)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	// NOT FOUND ERROR
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
