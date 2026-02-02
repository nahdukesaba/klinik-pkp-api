package district

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/region"
	"time"
)

type District struct {
	ID       uint64 `gorm:"primaryKey" json:"id,string"`
	RegionID uint64 `gorm:"column:region_id" json:"region_id,string"`
	Name     string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Region *region.Region `gorm:"foreignKey:RegionID;constraint:OnDelete:CASCADE;" json:"region,omitempty"`
}

func (District) TableName() string {
	return "district"
}
