package models

import (
	"time"

	"gorm.io/gorm"
)

// Province represents province data
type Province struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Code      string         `gorm:"size:10;uniqueIndex" json:"code"`
	
	// Audit fields
	CreatedBy *uint          `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *uint          `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy *uint          `gorm:"index" json:"deleted_by,omitempty"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Regencies []Regency `gorm:"foreignKey:ProvinceID" json:"regencies,omitempty"`
}

// TableName specifies the table name for Province model
func (Province) TableName() string {
	return "provinces"
}
