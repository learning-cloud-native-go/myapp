package book

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Books, error) {
	books := make([]*Book, 0)
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, nil
	}

	return books, nil
}

func (r *Repository) Create(book *Book) (*Book, error) {
	if err := r.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) Read(id uuid.UUID) (*Book, error) {
	book := &Book{}
	if err := r.db.Where("id = ?", id).First(&book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) Update(book *Book) error {
	return r.db.Updates(book).Where("id = %s", book.ID).Error
}

func (r *Repository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&Book{}).Error
}
