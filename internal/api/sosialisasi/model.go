package sosialisasi

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"time"
)

type Sosialisasi struct {
	ID          uint64   `gorm:"primaryKey" json:"id"`
	VillageID   string   `gorm:"column:village_id;not null" json:"village_id"`
	DistrictID  string   `gorm:"column:district_id;not null" json:"district_id"`
	RegionID    string   `gorm:"column:region_id;not null" json:"region_id"`
	Title       string   `gorm:"column:title" json:"title"`
	Location    string   `gorm:"column:location" json:"location"`
	Description string   `gorm:"column:description" json:"description"`
	ImageURLs   []string `gorm:"column:image_urls;serializer:json" json:"image_urls"`
	Coordinates []utils.Coordinate `gorm:"serializer:json" json:"coordinates"`

	// TIMESTAMPS
	ScheduledAtStart time.Time      `gorm:"column:scheduled_at_start" json:"scheduled_at_start"`
	ScheduledAtEnd   time.Time      `gorm:"column:scheduled_at_end" json:"scheduled_at_end"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (Sosialisasi) TableName() string {
	return "sosialisasi"
}
