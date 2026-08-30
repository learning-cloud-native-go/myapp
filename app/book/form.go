package book

import (
	"time"

	"github.com/google/uuid"

	"myapp/model"
)

type Form struct {
	Title         string `json:"title" validate:"required,max=255"`
	PublishedDate string `json:"published_date" validate:"required,datetime=2006-01-02"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url" validate:"url"`
	Status        string `json:"status" validate:"required,oneof=pending verified"`
}

func (f *Form) ToCreateModel() *model.Book {
	return f.ToUpdateModel(uuid.New())
}

func (f *Form) ToUpdateModel(id uuid.UUID) *model.Book {
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
