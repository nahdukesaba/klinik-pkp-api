package province

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID   string `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name string `gorm:"column:name" json:"name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	// USE omitzero TO AVOID "null" IN JSON RESPONSE
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Province) TableName() string {
	return "province"
}
