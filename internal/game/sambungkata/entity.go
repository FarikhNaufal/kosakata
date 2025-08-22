package sambungkata

import (
	"time"
	"gorm.io/datatypes"
)

type Word struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	Start     string         `json:"start" binding:"required"`
	End       string         `json:"end"`
	List      datatypes.JSON `json:"list" gorm:"type:json"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:date"`
	ReleaseAt string    `json:"release_at" gorm:"type:date"`
}
