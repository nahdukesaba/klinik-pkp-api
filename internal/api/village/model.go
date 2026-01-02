package village

import (
	"klinik-pkp-api/internal/api/district"
	"time"

	"gorm.io/gorm"
)

type Village struct {
	ID         string `gorm:"primaryKey;autoIncrement:false" json:"id"`
	DistrictID string `gorm:"column:district_id" json:"district_id"`
	Name       string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// RELATIONSHIPS
	District *district.District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

func (Village) TableName() string {
	return "village"
}
