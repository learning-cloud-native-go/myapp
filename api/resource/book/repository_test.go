package book_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"myapp/api/resource/book"
	mockDB "myapp/mock/db"
	testUtil "myapp/util/test"
)

func TestListBooks(t *testing.T) {
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(uuid.New(), "Book1", "Author1").
		AddRow(uuid.New(), "Book2", "Author2")

	mock.ExpectQuery("^SELECT (.+) FROM \"books\"").WillReturnRows(rows)

	books, err := repo.ListBooks()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(books), 2)
}

func TestCreateBook(t *testing.T) {
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"books\" ").
		WithArgs(id, "Title", "Author", testUtil.AnyTime{}, "", "", testUtil.AnyTime{}, testUtil.AnyTime{}, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	book := &book.Book{ID: id, Title: "Title", Author: "Author", PublishedDate: time.Now()}
	_, err = repo.CreateBook(book)
	testUtil.NoError(t, err)
}

func TestReadBook(t *testing.T) {
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(id, "Book1", "Author1")

	mock.ExpectQuery("^SELECT (.+) FROM \"books\" WHERE (.+)").
		WithArgs(id).
		WillReturnRows(rows)

	book, err := repo.ReadBook(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, "Book1", book.Title)
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"books\" SET").
		WithArgs("Title", "Author", testUtil.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	book := &book.Book{ID: id, Title: "Title", Author: "Author"}
	err = repo.UpdateBook(book)
	testUtil.NoError(t, err)
}

func TestDeleteBook(t *testing.T) {
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"books\" SET \"deleted_at\"").
		WithArgs(testUtil.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteBook(id)
	testUtil.NoError(t, err)
}
