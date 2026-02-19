package sosialisasi

import (
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"time"

	"gorm.io/gorm"
)

type Sosialisasi struct {
	ID          uint64           `gorm:"primaryKey" json:"id,string"`
	VillageID   uint64           `gorm:"column:village_id;not null" json:"village_id,string"`
	DistrictID  uint64           `gorm:"column:district_id;not null" json:"district_id,string"`
	RegionID    uint64           `gorm:"column:region_id;not null" json:"region_id,string"`
	Title       string           `gorm:"column:title" json:"title"`
	Location    string           `gorm:"column:location" json:"location"`
	Description string           `gorm:"column:description" json:"description"`
	ImageURLs   []string         `gorm:"column:image_urls;serializer:json" json:"image_urls"`
	Coordinate  utils.Coordinate `gorm:"serializer:json" json:"coordinate"`

	// TIMESTAMPS
	ScheduledAtStart time.Time      `gorm:"column:scheduled_at_start" json:"scheduled_at_start"`
	ScheduledAtEnd   time.Time      `gorm:"column:scheduled_at_end" json:"scheduled_at_end"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (Sosialisasi) TableName() string {
	return "sosialisasi"
}
