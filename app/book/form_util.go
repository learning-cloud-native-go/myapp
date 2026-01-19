package book

import (
	"time"

	"github.com/google/uuid"

	"myapp/form"
	"myapp/model"
)

func CreateFormToModel(f *form.BookForm) *model.Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)

	return &model.Book{
		ID:            uuid.New(),
		Title:         f.Title,
		Author:        f.Author,
		PublishedDate: pubDate,
		ImageURL:      f.ImageURL,
		Description:   f.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func UpdateFormToModel(f *form.BookForm, id uuid.UUID) *model.Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)

	return &model.Book{
		ID:            id,
		Title:         f.Title,
		Author:        f.Author,
		PublishedDate: pubDate,
		ImageURL:      f.ImageURL,
		Description:   f.Description,
		UpdatedAt:     time.Now(),
	}
}
