package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null;index" json:"filename"`
	OriginalName string         `gorm:"type:varchar(255);not null" json:"original_name"`
	FilePath     string         `gorm:"type:varchar(500);not null" json:"file_path"`
	FileURL      string         `gorm:"type:varchar(500);not null" json:"file_url"`
	MimeType     string         `gorm:"type:varchar(100);not null" json:"mime_type"`
	FileSize     int64          `gorm:"not null" json:"file_size"`
	Width        int            `gorm:"default:0" json:"width,omitempty"`
	Height       int            `gorm:"default:0" json:"height,omitempty"`
	EntityType   string         `gorm:"type:varchar(50);index" json:"entity_type,omitempty"` // rusun, balai, etc
	EntityID     uint           `gorm:"index" json:"entity_id,omitempty"`                     // related entity ID
	UploadedBy   uint           `gorm:"index" json:"uploaded_by,omitempty"`                   // user ID
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Image) TableName() string {
	return "images"
}

func (i *Image) BeforeSave(tx *gorm.DB) error {
	if i.FileSize <= 0 {
		return gorm.ErrInvalidData
	}
	if i.Filename == "" || i.FilePath == "" {
		return gorm.ErrInvalidData
	}
	return nil
}
