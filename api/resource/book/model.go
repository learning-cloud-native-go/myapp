package book

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DTO struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ImageUrl      string `json:"image_url"`
	Description   string `json:"description"`
}

type Form struct {
	Title         string `json:"title" form:"required,max=255"`
	Author        string `json:"author" form:"required,alpha_space,max=255"`
	PublishedDate string `json:"published_date" form:"required,datetime=2006-01-02"`
	ImageURL      string `json:"image_url" form:"url"`
	Description   string `json:"description"`
}

type Book struct {
	ID            uuid.UUID `gorm:"primarykey"`
	Title         string
	Author        string
	PublishedDate time.Time
	ImageUrl      string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Books []*Book

func (b *Book) ToDto() *DTO {
	return &DTO{
		ID:            b.ID.String(),
		Title:         b.Title,
		Author:        b.Author,
		PublishedDate: b.PublishedDate.Format("2006-01-02"),
		ImageUrl:      b.ImageUrl,
		Description:   b.Description,
	}
}

func (bs Books) ToDto() []*DTO {
	dtos := make([]*DTO, len(bs))
	for i, v := range bs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() *Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)

	return &Book{
		ID:            uuid.New(),
		Title:         f.Title,
		Author:        f.Author,
		PublishedDate: pubDate,
		ImageUrl:      f.ImageURL,
		Description:   f.Description,
	}
}
