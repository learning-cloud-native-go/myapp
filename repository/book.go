package repository

import (
	"github.com/jinzhu/gorm"

	"myapp/model"
)

func ListBooks(db *gorm.DB) (model.Books, error) {
	books := make([]*model.Book, 0)
	if err := db.Find(&books).Error; err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, nil
	}

	return books, nil
}
