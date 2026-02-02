package region

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/province"
	"time"
)

type Region struct {
	ID         uint64 `gorm:"primaryKey" json:"id,string"`
	ProvinceID uint64 `gorm:"column:province_id" json:"province_id,string"`
	Name       string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	Province *province.Province `gorm:"foreignKey:ProvinceID;constraint:OnDelete:CASCADE;" json:"province,omitempty"`
}

// TABLE NAME
func (Region) TableName() string {
	return "region"
}
