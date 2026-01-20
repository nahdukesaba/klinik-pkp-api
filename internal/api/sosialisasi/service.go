package sosialisasi

import (
	"errors"
	"fmt"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/uploads"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"time"
	"gorm.io/gorm"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	validator     *utils.Validator
	uploadService *uploads.Service
}

// PAYLOAD FROM REQUEST BODY
type SosialisasiPayload struct {
	VillageID        string                  `form:"village_id"`
	DistrictID       string                  `form:"district_id"`
	RegionID         string                  `form:"region_id"`
	Title            string                  `form:"title"`
	Location         string                  `form:"location"`
	Description      string                  `form:"description"`
	ScheduledAtStart string                  `form:"scheduled_at_start"`
	ScheduledAtEnd   string                  `form:"scheduled_at_end"`
	Images           []*multipart.FileHeader `form:"images"`
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

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Find(&sosialisasi).Error; err != nil {
		return nil, err
	}

	return sosialisasi, nil
}

func (s *Service) GetSosialisasiById(id int) (*Sosialisasi, error) {
	var sosialisasi Sosialisasi

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		First(&sosialisasi, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	return &sosialisasi, nil
}

func (s *Service) AddSosialisasi(payload *SosialisasiPayload) (*Sosialisasi, error) {
	// VALIDATIONS
	err := s.validator.ValidateStruct(map[string]any{
		"VillageID":        payload.VillageID,
		"DistrictID":       payload.DistrictID,
		"RegionID":         payload.RegionID,
		"Title":            payload.Title,
		"Location":         payload.Location,
		"Description":      payload.Description,
		"ScheduledAtStart": payload.ScheduledAtStart,
		"ScheduledAtEnd":   payload.ScheduledAtEnd,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var village village.Village

	if err := s.db.First(&village, "id = ?", payload.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("village with id %s not found", payload.VillageID)
		}

		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district district.District

	if err := s.db.First(&district, "id = ?", payload.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("district with id %s not found", payload.DistrictID)
		}

		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region region.Region

	if err := s.db.First(&region, "id = ?", payload.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("region not found")
		}

		return nil, err
	}

	// PARSE SCHEDULED AT
	scheduledAtStart, err := time.Parse(time.RFC3339, payload.ScheduledAtStart)

	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_at_start format: %w", err)
	}

	scheduledAtEnd, err := time.Parse(time.RFC3339, payload.ScheduledAtEnd)

	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_at_end format: %w", err)
	}

	// CREATE RECORD FIRST TO GET THE ID
	sosialisasi := Sosialisasi{
		VillageID:        payload.VillageID,
		DistrictID:       payload.DistrictID,
		RegionID:         payload.RegionID,
		Title:            payload.Title,
		Location:         payload.Location,
		Description:      payload.Description,
		ScheduledAtStart: scheduledAtStart,
		ScheduledAtEnd:   scheduledAtEnd,
	}

	if err := s.db.Create(&sosialisasi).Error; err != nil {
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "sosialisasi",
			RecordID: sosialisasi.ID,
		})

		if err != nil {
			s.db.Delete(&sosialisasi)

			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		sosialisasi.ImageURLs = imageResponses
	}

	if err := s.db.Save(&sosialisasi).Error; err != nil {
		return nil, err
	}

	return &sosialisasi, nil
}

func (s *Service) EditSosialisasiById(id uint64, payload *SosialisasiPayload) (*Sosialisasi, error) {
	// VALIDATIONS
	err := s.validator.ValidateStruct(map[string]any{
		"VillageID":        payload.VillageID,
		"DistrictID":       payload.DistrictID,
		"RegionID":         payload.RegionID,
		"Title":            payload.Title,
		"Location":         payload.Location,
		"Description":      payload.Description,
		"ScheduledAtStart": payload.ScheduledAtStart,
		"ScheduledAtEnd":   payload.ScheduledAtEnd,
	})

	if err != nil {
		return nil, err
	}

	// ENSURE SOSIALISASI EXISTS
	var sosialisasi Sosialisasi

	if err := s.db.First(&sosialisasi, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	sosialisasi.VillageID = payload.VillageID
	sosialisasi.DistrictID = payload.DistrictID
	sosialisasi.RegionID = payload.RegionID
	sosialisasi.Title = payload.Title
	sosialisasi.Location = payload.Location
	sosialisasi.Description = payload.Description
	sosialisasi.ScheduledAtStart, _ = time.Parse(time.RFC3339, payload.ScheduledAtStart)
	sosialisasi.ScheduledAtEnd, _ = time.Parse(time.RFC3339, payload.ScheduledAtEnd)

	if len(payload.Images) > 0 {
		// UPLOAD NEW IMAGES
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "sosialisasi",
			RecordID: id,
		})

		// DELETE OLD IMAGES
		if err := s.uploadService.DeleteImages("sosialisasi", id); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to upload new images: %w", err)
		}

		sosialisasi.ImageURLs = imageResponses
	}

	// SAVE UPDATES
	if err := s.db.Save(&sosialisasi).Error; err != nil {
		return nil, err
	}

	return &sosialisasi, nil
}

func (s *Service) DeleteSosialisasiById(id uint64) error {
	// FIND EXISTING RECORD
	var sosialisasi Sosialisasi

	if err := s.db.First(&sosialisasi, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}

		return err
	}

	// DELETE IMAGES (IF ANY)
	if err := s.uploadService.DeleteImages("sosialisasi", id); err != nil {
		fmt.Printf("Warning: failed to delete images: %v\n", err)
	}

	// DELETE RECORD FROM DATABASE
	if err := s.db.Delete(&sosialisasi).Error; err != nil {
		return err
	}

	return nil
}
