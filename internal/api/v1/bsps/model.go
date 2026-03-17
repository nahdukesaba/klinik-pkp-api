package bsps

import (
	"klinik-pkp-api/internal/api/v1/district"
	"klinik-pkp-api/internal/api/v1/region"
	"klinik-pkp-api/internal/api/v1/village"
	"klinik-pkp-api/utils"
	"time"

	"gorm.io/gorm"
)

type BSPS struct {
	ID         uint64           `gorm:"primaryKey" json:"id,string"`
	VillageID  uint64           `gorm:"column:village_id;not null" json:"village_id,string"`
	DistrictID uint64           `gorm:"column:district_id;not null" json:"district_id,string"`
	RegionID   uint64           `gorm:"column:region_id;not null" json:"region_id,string"`
	UnitCount  uint64           `json:"unit_count"`
	YearGiven  uint64           `json:"year_given"`
	Status     string           `gorm:"type:varchar(255);not null;default:'Rencana';check:status IN ('Rencana', 'Dalam Proses', 'Selesai')" json:"status"`
	Coordinate utils.Coordinate `gorm:"serializer:json" json:"coordinate"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (BSPS) TableName() string {
	return "bsps"
}
