package models

import (
	"time"

	"gorm.io/gorm"
)

// Regency represents regency/city data
type Regency struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProvinceID uint           `gorm:"not null;index" json:"province_id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Code       string         `gorm:"size:10;uniqueIndex" json:"code"`
	
	// Audit fields
	CreatedBy  *uint          `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy  *uint          `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy  *uint          `gorm:"index" json:"deleted_by,omitempty"`
	
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Province Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
}

// TableName specifies the table name for Regency model
func (Regency) TableName() string {
	return "regencies"
}
