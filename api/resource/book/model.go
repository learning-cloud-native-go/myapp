package book

import (
	"time"

	"gorm.io/gorm"
)

type FormBook struct {
	Title         string `json:"title" form:"required,max=255"`
	Author        string `json:"author" form:"required,alpha_space,max=255"`
	PublishedDate string `json:"published_date" form:"required,datetime=2006-01-02"`
	ImageUrl      string `json:"image_url" form:"url"`
	Description   string `json:"description"`
}

type Book struct {
	gorm.Model
	Title         string
	Author        string
	PublishedDate time.Time
	ImageUrl      string
	Description   string
}

type Books []*Book

type BookDto struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ImageUrl      string `json:"image_url"`
	Description   string `json:"description"`
}

func (b Book) ToDto() *BookDto {
	return &BookDto{
		ID:            b.ID,
		Title:         b.Title,
		Author:        b.Author,
		PublishedDate: b.PublishedDate.Format("2006-01-02"),
		ImageUrl:      b.ImageUrl,
		Description:   b.Description,
	}
}

func (bs Books) ToDto() []*BookDto {
	dtos := make([]*BookDto, len(bs))
	for i, v := range bs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *FormBook) ToModel() *Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)

	return &Book{
		Title:         f.Title,
		Author:        f.Author,
		PublishedDate: pubDate,
		ImageUrl:      f.ImageUrl,
		Description:   f.Description,
	}
}
