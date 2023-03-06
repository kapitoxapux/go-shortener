package models

import (
	"time"
)

type Link struct {
	ID        string    `gorm:"index" json:"id"`
	LongURL   string    `gorm:"index:idx_long;unique" json:"long_url"`
	ShortURL  string    `gorm:"not null" json:"short_url"`
	BaseURL   string    `gorm:"not null" json:"base_url"`
	SignID    uint32    `gorm:"not null;uint" json:"sign_id"`
	Sign      []byte    `gorm:"bytes;not null" json:"sign"`
	IsDeleted uint8     `gorm:"type:smallint;default:0;not null" json:"is_deleted"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
