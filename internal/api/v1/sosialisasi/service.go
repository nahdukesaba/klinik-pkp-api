package sosialisasi

import (
	"fmt"
	"klinik-pkp-api/internal/api/v1/uploads"
	"klinik-pkp-api/utils"
	"mime/multipart"
	"strconv"
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

// FILTER FOR QUERY PARAMETERS
type SosialisasiFilter struct {
	VillageID  *uint64
	DistrictID *uint64
	RegionID   *uint64
	Location   *string
	Title      *string
}

// CONSTRUCTOR
func NewService(db *gorm.DB, uploadService *uploads.Service) *Service {
	return &Service{
		db:            db,
		validator:     utils.NewValidator(),
		uploadService: uploadService,
	}
}

func (s *Service) GetSosialisasi(filter SosialisasiFilter, pagination utils.PaginationQuery) ([]Sosialisasi, int64, error) {
	var sosialisasi []Sosialisasi
	var totalRecords int64

	query := s.db.Preload("Village").Preload("District").Preload("Region.Province").Model(&Sosialisasi{})

	// FILTER BY VILLAGE
	if filter.VillageID != nil {
		query = query.Where("sosialisasi.village_id = ?", *filter.VillageID)
	}

	// FILTER BY DISTRICT
	if filter.DistrictID != nil {
		query = query.Where("sosialisasi.district_id = ?", *filter.DistrictID)
	}

	// FILTER BY REGION
	if filter.RegionID != nil {
		query = query.Where("sosialisasi.region_id = ?", *filter.RegionID)
	}

	// FILTER BY LOCATION (SEARCH)
	if filter.Location != nil {
		query = query.Where("sosialisasi.location LIKE ?", "%"+*filter.Location+"%")
	}

	// FILTER BY TITLE (SEARCH)
	if filter.Title != nil {
		query = query.Where("sosialisasi.title LIKE ?", "%"+*filter.Title+"%")
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset()).Find(&sosialisasi).Error; err != nil {
		return nil, 0, err
	}

	return sosialisasi, totalRecords, nil
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

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
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

	if err := tx.Create(&newSosialisasi).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(payload.Images) > 0 {
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 4,
			Category: "sosialisasi",
			RecordID: newSosialisasi.ID,
		})

		if err != nil {
			tx.Rollback()
			// CLEANUP: IMAGES UPLOADED, DELETE THEM ASYNCHRONOUSLY
			go s.uploadService.DeleteImages(newSosialisasi.ImageURLs)
			return nil, fmt.Errorf("failed to upload images: %w", err)
		}

		newSosialisasi.ImageURLs = imageResponses
	}

	if err := tx.Save(&newSosialisasi).Error; err != nil {
		tx.Rollback()
		// CLEANUP: DELETE UPLOADED IMAGES ASYNCHRONOUSLY
		go s.uploadService.DeleteImages(newSosialisasi.ImageURLs)
		return nil, err
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		// CLEANUP: DELETE UPLOADED IMAGES ASYNCHRONOUSLY
		go s.uploadService.DeleteImages(newSosialisasi.ImageURLs)
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

	// PARSE TIME
	scheduledAtStart, err := time.Parse(time.RFC3339, payload.ScheduledAtStart)

	if err != nil {
		return err
	}

	scheduledAtEnd, err := time.Parse(time.RFC3339, payload.ScheduledAtEnd)

	if err != nil {
		return err
	}

	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// CHECK IF RECORD EXISTS
	existingSosialisasi := Sosialisasi{}

	if err := tx.First(&existingSosialisasi, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	existingImageURLs := append([]string(nil), existingSosialisasi.ImageURLs...)
	isImagesUpdated := false

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
		// UPLOAD NEW IMAGES FIRST; OLD IMAGES STAY UNTIL TRANSACTION COMMITS
		imageResponses, err := s.uploadService.SaveImages(&uploads.FilePayload{
			Files:    payload.Images,
			MaxCount: 4,
			Category: "sosialisasi",
			RecordID: id,
		})

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		newSosialisasi.ImageURLs = imageResponses
		isImagesUpdated = true
	}

	result := tx.Model(&Sosialisasi{}).Where("id = ?", id).Updates(newSosialisasi)

	// SERVER ERRORS
	if result.Error != nil {
		tx.Rollback()
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newSosialisasi.ImageURLs)
		}
		return result.Error
	}

	// COMMIT TRANSACTION
	if err := tx.Commit().Error; err != nil {
		if isImagesUpdated {
			go s.uploadService.DeleteImages(newSosialisasi.ImageURLs)
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

func (s *Service) DeleteSosialisasiById(id uint64) error {
	// START TRANSACTION
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingSosialisasi Sosialisasi
	if err := tx.First(&existingSosialisasi, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// DELETE IMAGES (IF ANY) - NON-BLOCKING, ASYNC CLEANUP
	go s.uploadService.DeleteImages(existingSosialisasi.ImageURLs)

	// DELETE RECORD FROM DATABASE
	result := tx.Delete(&Sosialisasi{}, id)

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
