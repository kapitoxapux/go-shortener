package models

import (
	"time"
)

type Shortener struct {
	ID        string    `gorm:"index", json:"id"`
	LongURL   string    `gorm:"not null", json:"long_url"`
	ShortURL  string    `gorm:"not null", json:"short_url"`
	BaseURL   string    `gorm:"not null", json:"base_url"`
	SignID    uint32    `gorm:"not null,uint", json:"sign_id"`
	Sign      []byte    `gorm:"bytes,not null", json:"sign"`
	CreatedAt time.Time `gorm:"autoCreateTime", json:"created_at"`
}
