package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `json:"id" gorm:"primarykey"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"published_date,format:'2006-01-02'"`
	ImageURL      string    `json:"image_url"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at,format:RFC3339" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at,format:RFC3339" gorm:"autoUpdateTime"`
}
