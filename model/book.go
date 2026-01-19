package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	Books []*Book
	Book  struct {
		ID            uuid.UUID `gorm:"primarykey"`
		Title         string
		Author        string
		PublishedDate time.Time
		ImageURL      string
		Description   string
		CreatedAt     time.Time `gorm:"autoCreateTime"`
		UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	}
)

type BookDTO struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ImageURL      string `json:"image_url"`
	Description   string `json:"description"`
}

func (b *Book) ToDTO() *BookDTO {
	return &BookDTO{
		ID:            b.ID.String(),
		Title:         b.Title,
		Author:        b.Author,
		PublishedDate: b.PublishedDate.Format("2006-01-02"),
		ImageURL:      b.ImageURL,
		Description:   b.Description,
	}
}

func (bs Books) ToDTO() []*BookDTO {
	resp := make([]*BookDTO, len(bs))
	for i, v := range bs {
		resp[i] = v.ToDTO()
	}

	return resp
}
