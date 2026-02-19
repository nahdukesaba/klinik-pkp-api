package province

import (
	"gorm.io/gorm"
	"time"
)

type Province struct {
	ID   uint64 `gorm:"primaryKey" json:"id,string"`
	Name string `gorm:"column:name" json:"name"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// USE omitzero TO AVOID "null" IN JSON RESPONSE
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

func (Province) TableName() string {
	return "province"
}
