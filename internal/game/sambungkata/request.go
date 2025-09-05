package sambungkata

import (
	"gorm.io/datatypes"
)

type WordRequest struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	Start     string         `json:"start" binding:"required"`
	End       string         `json:"end"`
	List      datatypes.JSON `json:"list" gorm:"type:json"`
	ReleaseAt string         `json:"release_at" gorm:"type:date"`
}

type NextWordRequest struct {
	PrevWord *string `json:"prev_word"`
	NextWord string  `json:"next_word" binding:"required"`
	ID string `json:"id"`
}
