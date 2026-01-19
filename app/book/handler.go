package book

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"myapp/app/book/bookrepo"
	"myapp/form"
	"myapp/model"
	"myapp/pkg/ctxutil"
	e "myapp/pkg/errors"
	l "myapp/pkg/logger"
)

type API struct {
	logger    *l.Logger
	validator *validator.Validate
	bookRepo  IBookRepo
}

func New(logger *l.Logger, validator *validator.Validate, db *gorm.DB) *API {
	return &API{
		logger:    logger,
		validator: validator,
		bookRepo:  bookrepo.IBookRepo[model.Book](db),
	}
}

// List godoc
//
//	@summary		List books
//	@description	List books
//	@tags			books
//	@accept			json
//	@produce		json
//	@success		200	{array}		model.BookDTO
//	@failure		500	{object}	e.Error
//	@router			/books [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	books, err := a.bookRepo.ListBooks(ctx, 10, 0)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(books) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.MarshalWrite(w, books.ToDTO()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
//	@summary		Create book
//	@description	Create book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			body	body		form.BookForm	true	"Book form"
//	@success		201		{object}	model.BookDTO
//	@failure		400		{object}	e.Error
//	@failure		422		{object}	e.Errors
//	@failure		500		{object}	e.Error
//	@router			/books [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	ctxForm, ok := ctxutil.ValidatedForm[form.BookForm](ctx)
	if !ok {
		a.logger.Error().Str(l.KeyReqID, reqID).Msg(l.MsgErrCTXValidatedFormNotFound)
		e.ServerError(w, e.RespValidatedFormNotFound)
		return
	}

	book, err := a.bookRepo.CreateBook(ctx, CreateFormToModel(&ctxForm))
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.MarshalWrite(w, book.ToDTO()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", book.ID.String()).Msg("new book created")
}

// Read godoc
//
//	@summary		Read book
//	@description	Read book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Book ID"
//	@success		200	{object}	model.BookDTO
//	@failure		400	{object}	e.Error
//	@failure		404
//	@failure		500	{object}	e.Error
//	@router			/books/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	book, err := a.bookRepo.ReadBook(ctx, id)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if book.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.MarshalWrite(w, book.ToDTO()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Update godoc
//
//	@summary		Update book
//	@description	Update book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id		path		string			true	"Book ID"
//	@param			body	body		form.BookForm	true	"Book form"
//	@success		200		{object}	model.BookDTO
//	@failure		400		{object}	e.Error
//	@failure		404
//	@failure		422	{object}	e.Errors
//	@failure		500	{object}	e.Error
//	@router			/books/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	ctxForm, ok := ctxutil.ValidatedForm[form.BookForm](ctx)
	if !ok {
		a.logger.Error().Str(l.KeyReqID, reqID).Msg(l.MsgErrCTXValidatedFormNotFound)
		e.ServerError(w, e.RespValidatedFormNotFound)
		return
	}

	book, err := a.bookRepo.UpdateBook(ctx, UpdateFormToModel(&ctxForm, id))
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}

	if book.ID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.MarshalWrite(w, book.ToDTO()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", id.String()).Msg("book updated")
}

// Delete godoc
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
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := ctxutil.RequestID(ctx)

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	isDeleted, err := a.bookRepo.DeleteBook(ctx, id)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataRemoveFailure)
		return
	}
	if !isDeleted {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Str("id", id.String()).Msg("book deleted")
}
