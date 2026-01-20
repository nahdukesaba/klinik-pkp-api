package region

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/province"
	"time"
)

type Region struct {
	ID         string `gorm:"primaryKey;autoIncrement:false" json:"id"`
	ProvinceID string `gorm:"column:province_id" json:"province_id"`
	Name       string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// RELATIONSHIPS
	Province *province.Province `gorm:"foreignKey:ProvinceID;constraint:OnDelete:CASCADE;" json:"province,omitempty"`
}

// TABLE NAME
func (Region) TableName() string {
	return "region"
}
