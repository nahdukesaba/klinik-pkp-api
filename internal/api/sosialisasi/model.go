package sosialisasi

import (
	"gorm.io/gorm"
	"time"
)

type Sosialisasi struct {
	ID          string `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Title       string `gorm:"column:title" json:"title"`
	Location    string `gorm:"column:location" json:"location"`
	Description string `gorm:"column:description" json:"description"`
	ImageURL    string `gorm:"column:image_url" json:"image_url"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
