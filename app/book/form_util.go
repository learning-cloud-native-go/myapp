package book

import (
	"time"

	"github.com/google/uuid"

	"myapp/form"
	"myapp/model"
)

func createFormToModel(f *form.BookForm) *model.Book {
	return updateFormToModel(f, uuid.New())
}

func updateFormToModel(f *form.BookForm, id uuid.UUID) *model.Book {
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
