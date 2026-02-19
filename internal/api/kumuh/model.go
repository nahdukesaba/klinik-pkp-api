package kumuh

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"klinik-pkp-api/internal/api/region"
	"klinik-pkp-api/utils"
	"time"
)

type KawasanKumuh struct {
	ID              uint64           `gorm:"primaryKey" json:"id,string"`
	DistrictID      uint64           `gorm:"column:district_id" json:"district_id,string"`
	RegionID        uint64           `gorm:"column:region_id" json:"region_id,string"`
	AreaName        string           `gorm:"column:area_name;size:255" json:"area_name"`
	Environments    string           `gorm:"column:environments;size:255" json:"environments"`
	Villages        string           `gorm:"column:villages;type:text" json:"villages"`
	TotalArea       float64          `gorm:"column:total_area" json:"total_area"`
	TotalPopulation uint64           `gorm:"column:total_population" json:"total_population"`
	SlumValue       uint64           `gorm:"column:slum_value" json:"slum_value"`
	Coordinate      utils.Coordinate `gorm:"serializer:json" json:"coordinate"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region   *region.Region     `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

// TABLE NAME
func (KawasanKumuh) TableName() string {
	return "kawasan_kumuh"
}
