package book

import (
	"time"

	"github.com/google/uuid"

	"myapp/form"
	"myapp/model"
)

func CreateFormToModel(f *form.BookForm) *model.Book {
	return UpdateFormToModel(f, uuid.New())
}

func UpdateFormToModel(f *form.BookForm, id uuid.UUID) *model.Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)
	status, _ := model.ParseBookStatus(f.Status)

	return &model.Book{
		ID:            id,
		PublishedDate: model.CivilDate(pubDate),
		Title:         f.Title,
		ImageURL:      f.ImageURL,
		Description:   f.Description,
		Status:        *status,
	}
}
