package rusun

import (
	"fmt"
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"strconv"
)

// CURRENT INSTANCE
type Service struct {
	db            *gorm.DB
	validator     *utils.Validator
	uploadService *uploads.Service
}

// PAYLOAD FROM REQUEST BODY
type RusunPayload struct {
	VillageID  string                  `form:"village_id" validate:"required"`
	DistrictID string                  `form:"district_id" validate:"required"`
	RegionID   string                  `form:"region_id" validate:"required"`
	Name       string                  `form:"name" validate:"required,ne=Null,ne=null,ne=NULL"`
	Address    string                  `form:"address" validate:"required,ne=null,ne=NULL,ne=Null"`
	Tower      uint64                  `form:"tower" validate:"required"`
	UnitType   string                  `form:"unit_type" validate:"required,ne=Null,ne=null,ne=NULL"`
	Floor      uint64                  `form:"floor" validate:"required"`
	UnitCount  uint64                  `form:"unit_count" validate:"required"`
	YearGiven  uint64                  `form:"year_given" validate:"required"`
	Images     []*multipart.FileHeader `form:"images" validate:"required"`

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

func (s *Service) GetRusun() ([]Rusun, error) {
	var rusuns []Rusun

	if err := s.db.Preload("Village").Preload("District").Preload("Region").Find(&rusuns).Error; err != nil {
		return nil, err
	}

	return rusuns, nil
}

func (s *Service) GetRusunById(id uint64) (*Rusun, error) {
	var rusun Rusun

	if err := s.db.Preload("Village").Preload("District").Preload("Region").First(&rusun, id).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *Service) AddRusun(payload *RusunPayload) (*Rusun, error) {
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

	// CREATE RECORD FIRST TO GET ID
	newRusun := Rusun{
		VillageID:  villageID,
		DistrictID: districtID,
		RegionID:   regionID,
		Name:       payload.Name,
		Address:    payload.Address,
		Tower:      payload.Tower,
		UnitType:   payload.UnitType,
		Floor:      payload.Floor,
		UnitCount:  payload.UnitCount,
		YearGiven:  payload.YearGiven,
		Coordinate: payload.Coordinate,
	}

	if err := s.db.Create(&newRusun).Error; err != nil {
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "rusun",
			RecordID: newRusun.ID,
		})

		if err != nil {
			s.db.Delete(&newRusun)

			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newRusun.ImageURLs = imageResponses
	}

	if err := s.db.Save(&newRusun).Error; err != nil {
		return nil, err
	}

	return &newRusun, nil
}

func (s *Service) EditRusunById(id uint64, payload *RusunPayload) error {
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
	existingRusun := Rusun{}

	if err := s.db.First(&existingRusun, id).Error; err != nil {
		return err
	}

	newRusun := Rusun{
		VillageID:  villageID,
		DistrictID: districtID,
		RegionID:   regionID,
		Name:       payload.Name,
		Address:    payload.Address,
		Tower:      payload.Tower,
		UnitType:   payload.UnitType,
		Floor:      payload.Floor,
		UnitCount:  payload.UnitCount,
		YearGiven:  payload.YearGiven,
		Coordinate: payload.Coordinate,
	}

	if len(payload.Images) > 0 {
		// DELETE OLD IMAGES
		if err := s.uploadService.DeleteImages("rusun", id); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}

		// UPLOAD NEW IMAGES
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			Category: "rusun",
			RecordID: id,
		})

		if err != nil {
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		newRusun.ImageURLs = imageResponses
	}

	result := s.db.Model(&Rusun{}).Where("id = ?", id).Updates(newRusun)

	// SERVER ERRORS
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Service) DeleteRusunById(id uint64) error {
	// DELETE IMAGES (IF ANY)
	if err := s.uploadService.DeleteImages("rusun", id); err != nil {
		fmt.Printf("Warning: failed to delete images: %v\n", err)
	}

	result := s.db.Delete(&Rusun{}, id)

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
