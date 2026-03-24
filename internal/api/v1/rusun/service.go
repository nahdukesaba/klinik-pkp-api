package rusun

import (
	"fmt"
	"klinik-pkp-api/internal/api/v1/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"strconv"

	"gorm.io/gorm"
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

// FILTER FOR QUERY PARAMETERS
type RusunFilter struct {
	VillageID  *uint64
	DistrictID *uint64
	RegionID   *uint64
}

// CONSTRUCTOR
func NewService(db *gorm.DB, uploadService *uploads.Service) *Service {
	return &Service{
		db:            db,
		validator:     utils.NewValidator(),
		uploadService: uploadService,
	}
}

func (s *Service) GetRusun(filter RusunFilter) ([]Rusun, error) {
	var rusuns []Rusun

	query := s.db.Preload("Village").Preload("District").Preload("Region").Model(&Rusun{})

	// FILTER BY VILLAGE
	if filter.VillageID != nil {
		query = query.Where("rusun.village_id = ?", *filter.VillageID)
	}

	// FILTER BY DISTRICT
	if filter.DistrictID != nil {
		query = query.Where("rusun.district_id = ?", *filter.DistrictID)
	}

	// FILTER BY REGION
	if filter.RegionID != nil {
		query = query.Where("rusun.region_id = ?", *filter.RegionID)
	}

	if err := query.Find(&rusuns).Error; err != nil {
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
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

	if err := tx.Create(&newRusun).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 1,
			Category: "rusun",
			RecordID: newRusun.ID,
		})

		if err != nil {
			tx.Rollback()
			// CLEANUP: IMAGES UPLOADED, DELETE THEM ASYNCHRONOUSLY
			go s.uploadService.DeleteImages(newRusun.ImageURLs)
			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newRusun.ImageURLs = imageResponses
	}

	if err := tx.Save(&newRusun).Error; err != nil {
		tx.Rollback()
		// CLEANUP: DELETE UPLOADED IMAGES ASYNCHRONOUSLY
		go s.uploadService.DeleteImages(newRusun.ImageURLs)
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		// CLEANUP: DELETE UPLOADED IMAGES ASYNCHRONOUSLY
		go s.uploadService.DeleteImages(newRusun.ImageURLs)
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// CHECK IF RECORD EXISTS
	existingRusun := Rusun{}

	if err := tx.First(&existingRusun, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	existingImageURLs := append([]string(nil), existingRusun.ImageURLs...)
	isImagesUpdated := false

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
		// UPLOAD NEW IMAGES FIRST; OLD IMAGES STAY UNTIL TRANSACTION COMMITS
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 1,
			Category: "rusun",
			RecordID: id,
		})

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		newRusun.ImageURLs = imageResponses
		isImagesUpdated = true
	}

	result := tx.Model(&Rusun{}).Where("id = ?", id).Updates(newRusun)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newRusun.ImageURLs)
		}
		return result.Error
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newRusun.ImageURLs)
		}
		return err
	}

	if isImagesUpdated && len(existingImageURLs) > 0 {
		if err := s.uploadService.DeleteImages(existingImageURLs); err != nil {
			fmt.Printf("Warning: failed to delete old images: %v\n", err)
		}
	}

	return nil
}

func (s *Service) DeleteRusunById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingRusun Rusun
	if err := tx.First(&existingRusun, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// DELETE IMAGES (IF ANY) - NON-BLOCKING, ASYNC CLEANUP
	go s.uploadService.DeleteImages(existingRusun.ImageURLs)

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&Rusun{}, id)

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
