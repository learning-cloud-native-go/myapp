package book

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/hlog"
	"gorm.io/gorm"

	"myapp/app/book/bookrepo"
	"myapp/model"
	"myapp/pkg/ctxutil"
	e "myapp/pkg/errors"
	m "myapp/pkg/middleware"
	"myapp/pkg/paramsutil"
)

type Handler struct {
	validator *validator.Validate
	bookRepo  IBookRepo
}

func New(validator *validator.Validate, db *gorm.DB) *Handler {
	return &Handler{
		validator: validator,
		bookRepo:  bookrepo.IBookRepo[model.Book](db),
	}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/", h.list)
	r.With(m.Validate[Form](h.validator)).Post("/", h.create)
	r.Get("/{id}", h.read)
	r.With(m.Validate[Form](h.validator)).Put("/{id}", h.update)
	r.Delete("/{id}", h.delete)
}

// list godoc
//
//	@summary		List books
//	@description	List books
//	@tags			books
//	@accept			json
//	@produce		json
//	@Param			page		query		int64	false	"Page number for pagination"	default(1)	minimum(1)
//	@Param			pageSize	query		int64	false	"Number of items per page"		default(10)	minimum(1)	maximum(100)
//	@success		200			{array}		model.Book
//	@failure		500			{object}	e.Error
//	@router			/books [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := hlog.FromRequest(r)

	limit, offset := paramsutil.LimitOffset(r)
	books, err := h.bookRepo.ListBooks(ctx, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(books) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.MarshalWrite(w, books); err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// create godoc
//
//	@summary		Create book
//	@description	Create book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			body	body		Form	true	"Book form"
//	@success		201		{object}	model.Book
//	@failure		400		{object}	e.Error
//	@failure		422		{object}	e.Errors
//	@failure		500		{object}	e.Error
//	@router			/books [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := hlog.FromRequest(r)

	ctxForm, ok := ctxutil.ValidatedForm[Form](ctx)
	if !ok {
		logger.Error().Msg("validated form not found")
		e.ServerError(w, e.RespValidatedFormNotFound)
		return
	}

	book, err := h.bookRepo.CreateBook(ctx, ctxForm.ToCreateModel())
	if err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.MarshalWrite(w, book); err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	logger.Info().Str("id", book.ID.String()).Msg("new book created")
}

// read godoc
//
//	@summary		Read book
//	@description	Read book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Book ID"
//	@success		200	{object}	model.Book
//	@failure		400	{object}	e.Error
//	@failure		404
//	@failure		500	{object}	e.Error
//	@router			/books/{id} [get]
func (h *Handler) read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := hlog.FromRequest(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	book, err := h.bookRepo.ReadBook(ctx, id)
	if err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if book.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.MarshalWrite(w, book); err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// update godoc
//
//	@summary		Update book
//	@description	Update book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id		path		string	true	"Book ID"
//	@param			body	body		Form	true	"Book form"
//	@success		200		{object}	model.Book
//	@failure		400		{object}	e.Error
//	@failure		404
//	@failure		422	{object}	e.Errors
//	@failure		500	{object}	e.Error
//	@router			/books/{id} [put]
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := hlog.FromRequest(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	ctxForm, ok := ctxutil.ValidatedForm[Form](ctx)
	if !ok {
		logger.Error().Msg("validated form not found")
		e.ServerError(w, e.RespValidatedFormNotFound)
		return
	}

	book, err := h.bookRepo.UpdateBook(ctx, ctxForm.ToUpdateModel(id))
	if err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}

	if book.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.MarshalWrite(w, book); err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	logger.Info().Str("id", id.String()).Msg("book updated")
}

// delete godoc
//
//	@summary		Delete book
//	@description	Delete book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Book ID"
//	@success		200
//	@failure		400	{object}	e.Error
//	@failure		404
//	@failure		500	{object}	e.Error
//	@router			/books/{id} [delete]
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := hlog.FromRequest(r)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	isDeleted, err := h.bookRepo.DeleteBook(ctx, id)
	if err != nil {
		logger.Error().Err(err).Msg("")
		e.ServerError(w, e.RespDBDataRemoveFailure)
		return
	}
	if !isDeleted {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger.Info().Str("id", id.String()).Msg("book deleted")
}
