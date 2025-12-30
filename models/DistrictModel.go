package models

import (
	"gorm.io/gorm"
	"time"
)

type District struct {
	ID       string `gorm:"primaryKey;autoIncrement:false" json:"id"`
	RegionID string `gorm:"column:region_id" json:"region_id"`
	Name     string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// RELATIONSHIPS
	Region *Region `gorm:"foreignKey:RegionID" json:"region,omitempty"`
}

func (District) TableName() string {
	return "district"
}
