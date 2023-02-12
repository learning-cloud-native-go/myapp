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

func (r *Repository) ListBooks() (Books, error) {
	books := make([]*Book, 0)
	if err := r.db.Find(&books).Error; err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, nil
	}

	return books, nil
}

func (r *Repository) CreateBook(book *Book) (*Book, error) {
	if err := r.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) ReadBook(id uuid.UUID) (*Book, error) {
	book := &Book{}
	if err := r.db.Where("id = ?", id).First(&book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) UpdateBook(book *Book) error {
	if err := r.db.First(&Book{}, book.ID).Updates(book).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteBook(id uuid.UUID) error {
	book := &Book{}
	if err := r.db.Where("id = ?", id).Delete(&book).Error; err != nil {
		return err
	}

	return nil
}
