package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents category data
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	
	// Audit fields
	CreatedBy   *uint          `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy   *uint          `gorm:"index" json:"updated_by,omitempty"`
	DeletedBy   *uint          `gorm:"index" json:"deleted_by,omitempty"`
	
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}
