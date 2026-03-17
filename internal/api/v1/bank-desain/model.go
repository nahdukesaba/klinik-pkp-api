package bank_desain

import (
	"time"

	"gorm.io/gorm"
)

type BankDesain struct {
	ID            uint64   `gorm:"primaryKey" json:"id,string"`
	Name          string   `gorm:"type:varchar(255);not null" json:"name"`
	Type          string   `gorm:"type:varchar(255);not null;default:'Tipe 36';check:type IN ('Tipe 36', 'Tipe 45', 'Tipe 54', 'Rusun')" json:"type"`
	BedroomCount  uint64   `json:"bedroom_count"`
	BathroomCount uint64   `json:"bathroom_count"`
	TotalArea     uint64   `json:"total_area"`
	HasGarage     *bool    `json:"has_garage"`
	ImageURLs     []string `gorm:"serializer:json" json:"image_urls"`
	FileURLs      []string `gorm:"serializer:json" json:"file_urls"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

func (BankDesain) TableName() string {
	return "bank_desain"
}
