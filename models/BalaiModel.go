package models

import (
	"time"

	"gorm.io/gorm"
)

// Balai represents balai data
type Balai struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProvinceID uint           `gorm:"not null;index" json:"province_id"`
	CategoryID uint           `gorm:"not null;index" json:"category_id"`
	Name       string         `gorm:"size:255;not null" json:"name"`
	Address    string         `gorm:"type:text" json:"address"`
	
	// Audit fields
	CreatedBy  *uint          `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy  *uint          `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy  *uint          `gorm:"index" json:"deleted_by,omitempty"`
	
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Province Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName specifies the table name for Balai model
func (Balai) TableName() string {
	return "balais"
}
