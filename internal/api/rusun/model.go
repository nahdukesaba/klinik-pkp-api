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
	ID          uint               `gorm:"primaryKey" json:"id"`
	VillageID   string             `gorm:"column:village_id;not null" json:"village_id"`
	DistrictID  string             `gorm:"column:district_id;not null" json:"district_id"`
	RegionID    string             `gorm:"column:region_id;not null" json:"region_id"`
	Name        string             `gorm:"size:255;not null" json:"name"`
	Address     string             `gorm:"size:500;not null" json:"address"`
	Tower       int                `gorm:"size:100" json:"tower"`
	UnitType    string             `gorm:"size:150" json:"unit_type"`
	Floor       int                `json:"floor"`
	UnitCount   int                `json:"unit_count"`
	YearGiven   int                `json:"year_given"`
	Coordinates []utils.Coordinate `gorm:"serializer:json" json:"coordinates"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// RELATIONSHIPS
	Village  *village.Village   `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

// TABLE NAME
func (Rusun) TableName() string {
	return "rusun"
}
