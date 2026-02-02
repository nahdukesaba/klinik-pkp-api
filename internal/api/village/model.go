package village

import (
	"gorm.io/gorm"
	"klinik-pkp-api/internal/api/district"
	"time"
)

type Village struct {
	ID         uint64 `gorm:"primaryKey" json:"id,string"`
	DistrictID uint64 `gorm:"column:district_id" json:"district_id,string"`
	Name       string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`

	// RELATIONSHIPS
	District *district.District `gorm:"foreignKey:DistrictID;constraint:OnDelete:CASCADE;" json:"district,omitempty"`
}

func (Village) TableName() string {
	return "village"
}
