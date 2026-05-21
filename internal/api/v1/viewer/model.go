package viewer

import "time"

type Viewer struct {
	ID        uint64 `gorm:"primaryKey"`
	Path      string `gorm:"type:text"`
	Method    string `gorm:"type:varchar(10)"`
	CreatedAt time.Time
}

func (Viewer) TableName() string {
	return "viewers"
}
