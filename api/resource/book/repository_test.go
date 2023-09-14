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

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	mockRows := sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(uuid.New(), "Book1", "Author1").
		AddRow(uuid.New(), "Book2", "Author2")

	mock.ExpectQuery("^SELECT (.+) FROM \"books\"").WillReturnRows(mockRows)

	books, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(books), 2)
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO \"books\" ").
		WithArgs(id, "Title", "Author", mockDB.AnyTime{}, "", "", mockDB.AnyTime{}, mockDB.AnyTime{}, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	book := &book.Book{ID: id, Title: "Title", Author: "Author", PublishedDate: time.Now()}
	_, err = repo.Create(book)
	testUtil.NoError(t, err)
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	mockRows := sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(id, "Book1", "Author1")

	mock.ExpectQuery("^SELECT (.+) FROM \"books\" WHERE (.+)").
		WithArgs(id).
		WillReturnRows(mockRows)

	book, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, "Book1", book.Title)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(id, "Book1", "Author1")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"books\" SET").
		WithArgs("Title", "Author", mockDB.AnyTime{}, "", "", mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	book := &book.Book{ID: id, Title: "Title", Author: "Author"}
	rows, err := repo.Update(book)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := book.NewRepository(db)

	id := uuid.New()
	_ = sqlmock.NewRows([]string{"id", "title", "author"}).
		AddRow(id, "Book1", "Author1")

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"books\" SET \"deleted_at\"").
		WithArgs(mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, 1, rows)
}
