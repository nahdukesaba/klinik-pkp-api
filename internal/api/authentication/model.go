package authentication

import (
	"time"
)

type Authentication struct {
	Token     string    `gorm:"primaryKey;size:500;not null" json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// TABLE NAME
func (Authentication) TableName() string {
	return "authentication"
}
