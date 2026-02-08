package bsps

import (
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/internal/api/village"
	"klinik-pkp-api/utils"
	"time"

	"gorm.io/gorm"
)

type BSPS struct {
	ID          uint64             `gorm:"primaryKey" json:"id,string"`
	VillageID   uint64             `gorm:"column:village_id;not null" json:"village_id"`
	DistrictID  uint64             `gorm:"column:district_id;not null" json:"district_id"`
	RegionID    uint64             `gorm:"column:region_id;not null" json:"region_id"`
	UnitCount   uint64             `json:"unit_count"`
	YearGiven   uint64             `json:"year_given"`
	Status      string             `gorm:"type:varchar(255);not null;default:'Rencana';check:status IN ('Rencana', 'Dalam Proses', 'Selesai')" json:"status"`
	Coordinates []utils.Coordinate `gorm:"serializer:json" json:"coordinates"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (BSPS) TableName() string {
	return "penerimaan_bsps"
}
