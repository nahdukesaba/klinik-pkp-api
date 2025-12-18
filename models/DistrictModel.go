package models

import (
	"time"

	"gorm.io/gorm"
)

// District represents Kecamatan
type District struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RegencyID  uint           `gorm:"not null;index" json:"regency_id"`
	Regency    Regency        `gorm:"foreignKey:RegencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"regency,omitempty"`
	Name       string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Code       string         `gorm:"type:varchar(20);uniqueIndex" json:"code"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (District) TableName() string {
	return "districts"
}
