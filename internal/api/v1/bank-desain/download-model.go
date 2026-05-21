package bank_desain

import "time"

// /internal/api/bank-desain/download_model.go
type BankDesainDownload struct {
	ID           uint64 `gorm:"primaryKey"`
	BankDesainID uint64
	Name         string `gorm:"type:varchar(255);not null"`
	Address      string `gorm:"type:text;not null"`

	CreatedAt time.Time
}

func (BankDesainDownload) TableName() string {
	return "bank_desain_downloads"
}
