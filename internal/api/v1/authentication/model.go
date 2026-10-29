package authentication

import (
	"time"
)

type Authentication struct {
	Token     string    `gorm:"primaryKey;size:500;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
