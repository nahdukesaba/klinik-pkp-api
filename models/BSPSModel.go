package models

import (
	"gorm.io/gorm"
	"time"
)

type BSPS struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	VillageID  string       `gorm:"column:village_id;not null" json:"village_id"`
	DistrictID string       `gorm:"column:district_id;not null" json:"district_id"`
	RegionID   string       `gorm:"column:region_id;not null" json:"region_id"`
	ProvinceID string       `gorm:"column:province_id;not null" json:"province_id"`
	Village    *Village     `gorm:"foreignKey:VillageID" json:"village,omitempty"`
	District   *District    `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Region     *Region      `gorm:"foreignKey:RegionID" json:"region,omitempty"`
	Province   *Province    `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	Coordinate []Coordinate `gorm:"serializer:json" json:"coordinate"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (BSPS) TableName() string {
	return "penerimaan_bsps"
}
