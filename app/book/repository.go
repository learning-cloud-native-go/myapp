package book

import (
	"context"

	"github.com/google/uuid"

	"myapp/model"
)

type IBookRepo interface {
	// SELECT * FROM books LIMIT @limit OFFSET @offset
	ListBooks(ctx context.Context, limit int64, offset int64) (model.Books, error)

	// INSERT INTO books (id, created_at, updated_at, title, author, published_date, image_url, description) VALUES (@data.ID, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, @data.Title, @data.Author, @data.PublishedDate, @data.ImageURL, @data.Description)
	// RETURNING *
	CreateBook(ctx context.Context, data *model.Book) (*model.Book, error)

	// SELECT * FROM books WHERE id = @id
	ReadBook(ctx context.Context, id uuid.UUID) (*model.Book, error)

	// UPDATE books
	// SET updated_at=CURRENT_TIMESTAMP, title=@data.Title, author=@data.Author, published_date=@data.PublishedDate, image_url=@data.ImageURL, description=@data.Description
	// WHERE id = @data.ID
	// RETURNING *
	UpdateBook(ctx context.Context, data *model.Book) (*model.Book, error)

	// DELETE FROM books WHERE id = @id
	// RETURNING true
	DeleteBook(ctx context.Context, id uuid.UUID) (bool, error)
}
