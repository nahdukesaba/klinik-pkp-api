package sosialisasi

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"strconv"
	"time"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	validator     *utils.Validator
	uploadService *uploads.Service
}

// PAYLOAD FROM REQUEST BODY
type SosialisasiPayload struct {
	VillageID        string                  `form:"village_id" validate:"required"`
	DistrictID       string                  `form:"district_id" validate:"required"`
	RegionID         string                  `form:"region_id" validate:"required"`
	Title            string                  `form:"title" validate:"required,ne=Null,ne=null,ne=NULL"`
	Location         string                  `form:"location" validate:"required,ne=Null,ne=null,ne=NULL"`
	Description      string                  `form:"description" validate:"required,ne=Null,ne=null,ne=NULL"`
	ScheduledAtStart string                  `form:"scheduled_at_start" validate:"required"`
	ScheduledAtEnd   string                  `form:"scheduled_at_end" validate:"required"`
	Images           []*multipart.FileHeader `form:"images" validate:"-"`

	// FIBER CANNOT PROCESS JSON ARRAYS IN FORM DATA, NEED TO UNMARSHALL MANUALLY
	CoordinateRaw string           `form:"coordinate" validate:"required"`
	Coordinate    utils.Coordinate `gorm:"-"`
}

// CONSTRUCTOR
func NewService(db *gorm.DB, uploadService *uploads.Service) *Service {
	return &Service{
		db:            db,
		validator:     &utils.Validator{},
		uploadService: uploadService,
	}
}

func (s *Service) GetSosialisasi() ([]Sosialisasi, error) {
	var sosialisasi []Sosialisasi

	if err := s.db.Preload("Village").Preload("District").Preload("Region.Province").Find(&sosialisasi).Error; err != nil {
		return nil, err
	}

	return sosialisasi, nil
}

func (s *Service) GetSosialisasiById(id uint64) (*Sosialisasi, error) {
	var sosialisasi Sosialisasi

	if err := s.db.Preload("Village").Preload("District").Preload("Region.Province").First(&sosialisasi, id).Error; err != nil {
		return nil, err
	}

	return &sosialisasi, nil
}

func (s *Service) AddSosialisasi(payload *SosialisasiPayload) (*Sosialisasi, error) {
	villageID, err := strconv.ParseUint(payload.VillageID, 10, 64)

	if err != nil {
		return nil, err
	}

	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return nil, err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return nil, err
	}

	// PARSE TIME
	scheduledAtStart, err := time.Parse(time.RFC3339, payload.ScheduledAtStart)

	if err != nil {
		return nil, err
	}

	scheduledAtEnd, err := time.Parse(time.RFC3339, payload.ScheduledAtEnd)

	if err != nil {
		return nil, err
	}

	// CREATE RECORD FIRST TO GET ID
	newSosialisasi := Sosialisasi{
		VillageID:        villageID,
		DistrictID:       districtID,
		RegionID:         regionID,
		Title:            payload.Title,
		Location:         payload.Location,
		Description:      payload.Description,
		Coordinate:       payload.Coordinate,
		ScheduledAtStart: scheduledAtStart,
		ScheduledAtEnd:   scheduledAtEnd,
	}

	if err := s.db.Create(&newSosialisasi).Error; err != nil {
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "sosialisasi",
			RecordID: newSosialisasi.ID,
		})

		if err != nil {
			s.db.Delete(&newSosialisasi)

			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newSosialisasi.ImageURLs = imageResponses
	}

	if err := s.db.Save(&newSosialisasi).Error; err != nil {
		return nil, err
	}

	return &newSosialisasi, nil
}

func (s *Service) EditSosialisasiById(id uint64, payload *SosialisasiPayload) error {
	villageID, err := strconv.ParseUint(payload.VillageID, 10, 64)

	if err != nil {
		return err
	}

	districtID, err := strconv.ParseUint(payload.DistrictID, 10, 64)

	if err != nil {
		return err
	}

	regionID, err := strconv.ParseUint(payload.RegionID, 10, 64)

	if err != nil {
		return err
	}

	// CHECK IF RECORD EXISTS
	existingSosialisasi := Sosialisasi{}

	if err := s.db.First(&existingSosialisasi, id).Error; err != nil {
		return err
	}

	// PARSE TIME
	scheduledAtStart, err := time.Parse(time.RFC3339, payload.ScheduledAtStart)

	if err != nil {
		return err
	}

	scheduledAtEnd, err := time.Parse(time.RFC3339, payload.ScheduledAtEnd)

	if err != nil {
		return err
	}

	// CREATE UPDATE STRUCT (TO STORE IMAGE URLS LATER)
	newSosialisasi := Sosialisasi{
		VillageID:        villageID,
		DistrictID:       districtID,
		RegionID:         regionID,
		Title:            payload.Title,
		Location:         payload.Location,
		Description:      payload.Description,
		Coordinate:       payload.Coordinate,
		ScheduledAtStart: scheduledAtStart,
		ScheduledAtEnd:   scheduledAtEnd,
	}

	if len(payload.Images) > 0 {
		// DELETE OLD IMAGES
		if err := s.uploadService.DeleteImages("sosialisasi", id); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}

		// UPLOAD NEW IMAGES
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "sosialisasi",
			RecordID: id,
		})

		if err != nil {
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		newSosialisasi.ImageURLs = imageResponses
	}

	result := s.db.Model(&Sosialisasi{}).Where("id = ?", id).Updates(newSosialisasi)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Service) DeleteSosialisasiById(id uint64) error {
	// DELETE IMAGES (IF ANY)
	if err := s.uploadService.DeleteImages("sosialisasi", id); err != nil {
		fmt.Printf("Warning: failed to delete images: %v\n", err)
	}

	result := s.db.Delete(&Sosialisasi{}, id)

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
