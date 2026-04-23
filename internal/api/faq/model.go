package faq

import (
	"time"

	"gorm.io/gorm"
)

type FAQ struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Question string `gorm:"type:text;not null" json:"question"`
	Answer   string `gorm:"type:text;not null" json:"answer"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (FAQ) TableName() string {
	return "faqs"
}
