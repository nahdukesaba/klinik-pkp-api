package models

import (
	"time"

	"gorm.io/gorm"
)

// Village represents Kelurahan/Desa
type Village struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	DistrictID uint           `gorm:"not null;index" json:"district_id"`
	District   District       `gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"district,omitempty"`
	Name       string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Code       string         `gorm:"type:varchar(20);uniqueIndex" json:"code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Village) TableName() string {
	return "villages"
}
