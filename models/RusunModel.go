package models

import (
	"time"

	"gorm.io/gorm"
)

type Rusun struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null;index" json:"name"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	ProvinceID    uint           `gorm:"not null;index" json:"province_id"`
	Province      Province       `gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"province,omitempty"`
	RegencyID     uint           `gorm:"not null;index" json:"regency_id"`
	Regency       Regency        `gorm:"foreignKey:RegencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"regency,omitempty"`
	DistrictID    uint           `gorm:"not null;index" json:"district_id"`
	District      District       `gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"district,omitempty"`
	VillageID     uint           `gorm:"not null;index" json:"village_id"`
	Village       Village        `gorm:"foreignKey:VillageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"village,omitempty"`
	Latitude      float64        `gorm:"type:decimal(10,8);not null;index" json:"latitude"`
	Longitude     float64        `gorm:"type:decimal(11,8);not null;index" json:"longitude"`
	UnitCount     int            `gorm:"not null;default:0" json:"unit_count"`
	TowerCount    int            `gorm:"not null;default:0" json:"tower_count"`
	FloorCount    int            `gorm:"not null;default:0" json:"floor_count"`
	Type          string         `gorm:"type:varchar(255);not null" json:"type"`
	BuildYear     int            `gorm:"not null" json:"build_year"`
	HandoverYear  int            `gorm:"not null" json:"handover_year"`
	OccupiedUnits int            `gorm:"not null;default:0" json:"occupied_units"`
	Contractor    string         `gorm:"type:varchar(255);not null" json:"contractor"`
	Status        string         `gorm:"type:varchar(50);not null;default:'active';index" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Rusun) TableName() string {
	return "rusun"
}

func (r *Rusun) BeforeSave(tx *gorm.DB) error {
	if r.Latitude < -90 || r.Latitude > 90 {
		return gorm.ErrInvalidData
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return gorm.ErrInvalidData
	}
	if r.OccupiedUnits > r.UnitCount {
		return gorm.ErrInvalidData
	}
	if r.Status == "" {
		r.Status = "active"
	}
	return nil
}
