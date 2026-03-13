package rusun

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"time"
)

type Rusun struct {
	ID         uint64           `gorm:"primaryKey" json:"id,string"`
	VillageID  uint64           `gorm:"column:village_id" json:"village_id,string"`
	DistrictID uint64           `gorm:"column:district_id" json:"district_id,string"`
	RegionID   uint64           `gorm:"column:region_id" json:"region_id,string"`
	Name       string           `gorm:"size:255;uniqueIndex" json:"name"`
	Address    string           `gorm:"size:500" json:"address"`
	Tower      uint64           `gorm:"size:100" json:"tower"`
	UnitType   string           `gorm:"size:150" json:"unit_type"`
	Floor      uint64           `json:"floor"`
	UnitCount  uint64           `json:"unit_count"`
	YearGiven  uint64           `json:"year_given"`
	ImageURLs  []string         `gorm:"column:image_urls;serializer:json" json:"image_urls"`
	Coordinate utils.Coordinate `gorm:"serializer:json" json:"coordinate"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

// TABLE NAME
func (Rusun) TableName() string {
	return "rusun"
}
