package model

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID  `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	PublishedDate CivilDate  `json:"published_date"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ImageURL      string     `json:"image_url"`
	Status        BookStatus `json:"status"`
}
