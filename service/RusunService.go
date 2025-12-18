package service

import (
	"errors"
	"html"
	"klinik-api/models"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RusunService struct {
	DB *gorm.DB
}

func NewRusunService(db *gorm.DB) *RusunService {
	return &RusunService{DB: db}
}

// sanitize cleans input from XSS attacks
func (s *RusunService) sanitize(input string) string {
	input = html.EscapeString(strings.TrimSpace(input))
	input = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`).ReplaceAllString(input, "")
	return input
}

// validate validates rusun data
func (s *RusunService) validate(r *models.Rusun) error {
	r.Name = s.sanitize(r.Name)
	r.Address = s.sanitize(r.Address)
	r.Type = s.sanitize(r.Type)
	r.Contractor = s.sanitize(r.Contractor)
	r.Status = s.sanitize(r.Status)

	if r.Name == "" || len(r.Name) > 255 {
		return errors.New("nama tidak valid")
	}
	if r.Address == "" {
		return errors.New("alamat tidak boleh kosong")
	}
	if r.Latitude < -90 || r.Latitude > 90 {
		return errors.New("latitude tidak valid")
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return errors.New("longitude tidak valid")
	}
	if r.OccupiedUnits > r.UnitCount {
		return errors.New("unit terisi melebihi total unit")
	}
	if r.HandoverYear < r.BuildYear {
		return errors.New("tahun serah terima tidak valid")
	}
	
	// Validate foreign keys
	if r.ProvinceID == 0 {
		return errors.New("provinsi harus dipilih")
	}
	if r.RegencyID == 0 {
		return errors.New("kabupaten/kota harus dipilih")
	}
	
	var province models.Province
	if err := s.DB.First(&province, r.ProvinceID).Error; err != nil {
		return errors.New("provinsi tidak ditemukan")
	}
	
	var regency models.Regency
	if err := s.DB.Where("id = ? AND province_id = ?", r.RegencyID, r.ProvinceID).First(&regency).Error; err != nil {
		return errors.New("kabupaten/kota tidak sesuai dengan provinsi")
	}

	if r.Status == "" {
		r.Status = "active"
	}
	validStatus := map[string]bool{"active": true, "inactive": true, "maintenance": true}
	if !validStatus[r.Status] {
		return errors.New("status tidak valid")
	}

	return nil
}

// GetAll retrieves all rusun with filters
func (s *RusunService) GetAll(page, limit int, provinceID, regencyID uint, status, search string) ([]models.Rusun, int64, error) {
	var rusun []models.Rusun
	var total int64

	query := s.DB.Model(&models.Rusun{}).Preload("Province").Preload("Regency")

	if provinceID > 0 {
		query = query.Where("province_id = ?", provinceID)
	}
	if regencyID > 0 {
		query = query.Where("regency_id = ?", regencyID)
	}
	if status != "" && (status == "active" || status == "inactive" || status == "maintenance") {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		search = s.sanitize(search)
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR address ILIKE ? OR contractor ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Order("created_at DESC").Find(&rusun).Error; err != nil {
		return nil, 0, err
	}

	return rusun, total, nil
}

// GetByID retrieves rusun by ID
func (s *RusunService) GetByID(id uint) (*models.Rusun, error) {
	if id == 0 {
		return nil, errors.New("ID tidak valid")
	}

	var rusun models.Rusun
	if err := s.DB.Preload("Province").Preload("Regency").First(&rusun, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("rusun tidak ditemukan")
		}
		return nil, err
	}

	return &rusun, nil
}

// GetByProvince retrieves rusun by province
func (s *RusunService) GetByProvince(provinceID uint) ([]models.Rusun, error) {
	if provinceID == 0 {
		return nil, errors.New("province ID tidak valid")
	}

	var rusun []models.Rusun
	if err := s.DB.Preload("Province").Preload("Regency").
		Where("province_id = ? AND status = ?", provinceID, "active").
		Order("name ASC").
		Find(&rusun).Error; err != nil {
		return nil, err
	}

	return rusun, nil
}

// GetMapData retrieves coordinate data for map
func (s *RusunService) GetMapData(provinceID uint) ([]map[string]interface{}, error) {
	var rusun []models.Rusun
	
	query := s.DB.Select("id, name, address, latitude, longitude, province_id, regency_id, unit_count, occupied_units").
		Where("status = ?", "active")
	
	if provinceID > 0 {
		query = query.Where("province_id = ?", provinceID)
	}
	
	if err := query.Find(&rusun).Error; err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rusun))
	for i, r := range rusun {
		occupancyRate := float64(0)
		if r.UnitCount > 0 {
			occupancyRate = float64(r.OccupiedUnits) / float64(r.UnitCount) * 100
		}
		
		result[i] = map[string]interface{}{
			"id":             r.ID,
			"name":           r.Name,
			"address":        r.Address,
			"latitude":       r.Latitude,
			"longitude":      r.Longitude,
			"province_id":    r.ProvinceID,
			"regency_id":     r.RegencyID,
			"unit_count":     r.UnitCount,
			"occupied_units": r.OccupiedUnits,
			"occupancy_rate": occupancyRate,
		}
	}

	return result, nil
}

// Create creates new rusun
func (s *RusunService) Create(rusun *models.Rusun) error {
	if err := s.validate(rusun); err != nil {
		return err
	}

	var count int64
	s.DB.Model(&models.Rusun{}).Where("name = ? AND province_id = ?", rusun.Name, rusun.ProvinceID).Count(&count)
	if count > 0 {
		return errors.New("rusun dengan nama yang sama sudah ada")
	}

	return s.DB.Create(rusun).Error
}

// Update updates rusun
func (s *RusunService) Update(id uint, rusun *models.Rusun) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	if err := s.validate(rusun); err != nil {
		return err
	}

	var existing models.Rusun
	if err := s.DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rusun tidak ditemukan")
		}
		return err
	}

	var count int64
	s.DB.Model(&models.Rusun{}).Where("name = ? AND province_id = ? AND id != ?", rusun.Name, rusun.ProvinceID, id).Count(&count)
	if count > 0 {
		return errors.New("rusun dengan nama yang sama sudah ada")
	}

	rusun.ID = id
	rusun.CreatedAt = existing.CreatedAt
	rusun.UpdatedAt = time.Now()
	
	return s.DB.Save(rusun).Error
}

// Delete soft deletes rusun
func (s *RusunService) Delete(id uint) error {
	if id == 0 {
		return errors.New("ID tidak valid")
	}

	var rusun models.Rusun
	if err := s.DB.First(&rusun, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("rusun tidak ditemukan")
		}
		return err
	}

	return s.DB.Delete(&rusun).Error
}

// GetStatistics retrieves rusun statistics
func (s *RusunService) GetStatistics(provinceID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := s.DB.Model(&models.Rusun{}).Where("status = ?", "active")
	if provinceID > 0 {
		query = query.Where("province_id = ?", provinceID)
	}

	var totalRusun int64
	query.Count(&totalRusun)
	stats["total_rusun"] = totalRusun

	var result struct {
		TotalUnits    int64
		TotalOccupied int64
		TotalTowers   int64
	}
	
	query.Select("COALESCE(SUM(unit_count), 0) as total_units, COALESCE(SUM(occupied_units), 0) as total_occupied, COALESCE(SUM(tower_count), 0) as total_towers").Scan(&result)
	
	stats["total_units"] = result.TotalUnits
	stats["total_occupied"] = result.TotalOccupied
	stats["total_towers"] = result.TotalTowers
	
	if result.TotalUnits > 0 {
		stats["occupancy_rate"] = float64(result.TotalOccupied) / float64(result.TotalUnits) * 100
	} else {
		stats["occupancy_rate"] = 0
	}

	return stats, nil
}

