package service

import (
	"errors"
	"gorm.io/gorm"
	"klinik-pkp-api/models"
	"klinik-pkp-api/utils"
)

// RUSUN SERVICE STRUCT
type RusunService struct {
	db        *gorm.DB
	validator *utils.Validator
}

// REPRESENTS DATA FROM INPUT
type RusunPayload struct {
	VillageID  string              `json:"village_id"`
	DistrictID string              `json:"district_id"`
	RegionID   string              `json:"region_id"`
	ProvinceID string              `json:"province_id"`
	Name       string              `json:"name"`
	Address    string              `json:"address"`
	Coordinate []models.Coordinate `json:"coordinate"`
	Tower      string              `json:"tower"`
	UnitType   string              `json:"unit_type"`
	Floor      int                 `json:"floor"`
	UnitCount  int                 `json:"unit_count"`
}

// RUSUN SERVICE CONSTRUCTOR
func NewRusunService(db *gorm.DB) *RusunService {
	return &RusunService{db: db, validator: &utils.Validator{}}
}

func (s *RusunService) GetAll() ([]models.Rusun, error) {
	var rusuns []models.Rusun

	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Preload("Province").
		Find(&rusuns).Error; err != nil {
		return nil, err
	}

	return rusuns, nil
}

func (s *RusunService) GetByID(id int) (*models.Rusun, error) {
	var rusun models.Rusun
	if err := s.db.
		Preload("Village").
		Preload("District").
		Preload("Region").
		Preload("Province").
		First(&rusun, id).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *RusunService) Create(input RusunPayload) (*models.Rusun, error) {
	if ok, msginput := s.validator.ValidateRequired(input.VillageID, "VillageID"); !ok {
		return nil, errors.New(msginput)
	}

	if ok, msginput := s.validator.ValidateRequired(input.DistrictID, "DistrictID"); !ok {
		return nil, errors.New(msginput)
	}

	if ok, msginput := s.validator.ValidateRequired(input.RegionID, "RegionID"); !ok {
		return nil, errors.New(msginput)
	}

	if ok, msginput := s.validator.ValidateRequired(input.ProvinceID, "ProvinceID"); !ok {
		return nil, errors.New(msginput)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Address, "Address"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Coordinate, "Coordinate"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Tower, "Tower"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.UnitType, "UnitType"); !ok {
		return nil, errors.New(msg)
	}

	// ENSURE RELATIONS EXIST
	var village models.Village

	if err := s.db.First(&village, input.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district models.District

	if err := s.db.First(&district, input.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region models.Region

	if err := s.db.First(&region, input.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}
		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var province models.Province

	if err := s.db.First(&province, input.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}
		return nil, err
	}

	rusun := models.Rusun{
		VillageID:  input.VillageID,
		DistrictID: input.DistrictID,
		RegionID:   input.RegionID,
		ProvinceID: input.ProvinceID,
		Name:       input.Name,
		Address:    input.Address,
		Coordinate: input.Coordinate,
		Tower:      input.Tower,
		UnitType:   input.UnitType,
		Floor:      input.Floor,
		UnitCount:  input.UnitCount,
	}

	if err := s.db.Create(&rusun).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *RusunService) Update(id int, input RusunPayload) (*models.Rusun, error) {
	// VALIDATIONS
	if ok, msg := s.validator.ValidateRequired(input.VillageID, "VillageID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.DistrictID, "DistrictID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.RegionID, "RegionID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.ProvinceID, "ProvinceID"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Name, "Name"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Address, "Address"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Coordinate, "Coordinate"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.Tower, "Tower"); !ok {
		return nil, errors.New(msg)
	}

	if ok, msg := s.validator.ValidateRequired(input.UnitType, "UnitType"); !ok {
		return nil, errors.New(msg)
	}

	var rusun models.Rusun

	// ENSURE RUSUN EXISTS
	if err := s.db.First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Rusun not found")
		}

		return nil, err
	}

	// ENSURE RELATED VILLAGE EXISTS
	var village models.Village

	if err := s.db.First(&village, input.VillageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Village not found")
		}
		return nil, err
	}

	// ENSURE RELATED DISTRICT EXISTS
	var district models.District

	if err := s.db.First(&district, input.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("District not found")
		}
		return nil, err
	}

	// ENSURE RELATED REGION EXISTS
	var region models.Region

	if err := s.db.First(&region, input.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Region not found")
		}

		return nil, err
	}

	// ENSURE RELATED PROVINCE EXISTS
	var province models.Province

	if err := s.db.First(&province, input.ProvinceID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Province not found")
		}

		return nil, err
	}

	rusun.VillageID = input.VillageID
	rusun.DistrictID = input.DistrictID
	rusun.RegionID = input.RegionID
	rusun.ProvinceID = input.ProvinceID
	rusun.Name = input.Name
	rusun.Address = input.Address
	rusun.Coordinate = input.Coordinate
	rusun.Tower = input.Tower
	rusun.UnitType = input.UnitType
	rusun.Floor = input.Floor
	rusun.UnitCount = input.UnitCount

	if err := s.db.Save(&rusun).Error; err != nil {
		return nil, err
	}

	return &rusun, nil
}

func (s *RusunService) Delete(id int) error {
	var rusun models.Rusun

	if err := s.db.First(&rusun, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Rusun not found")
		}

		return err
	}

	if err := s.db.Delete(&rusun).Error; err != nil {
		return err
	}

	return nil
}
