package book

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	e "myapp/api/resource/common/err"
	"myapp/util/validator"
)

// List godoc
//
//	@summary		List books
//	@description	List books
//	@tags			books
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	err.Error
//	@router			/books [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	books, err := a.repository.ListBooks()
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(books.ToDto()); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.JsonEncodingFailure)
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
//	@param			body	body	Form	true	"Book form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/books [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.BadRequest(w, e.JsonDecodingFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.ServerError(w, e.FormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			e.ServerError(w, e.JsonEncodingFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	book, err := a.repository.CreateBook(form.ToModel())
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataCreationFailure)
		return
	}

	a.logger.Info().Str("id", book.ID.String()).Msg("new book created")
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read book
//	@description	Read book
//	@tags			books
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Book ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	book, err := a.repository.ReadBook(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.JsonEncodingFailure)
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
//	@param			id		path	string	true	"Book ID"
//	@param			body	body	Form	true	"Book form"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.BadRequest(w, e.JsonDecodingFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.ServerError(w, e.FormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			e.ServerError(w, e.JsonEncodingFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	bookModel := form.ToModel()
	bookModel.ID = id

	if err := a.repository.UpdateBook(bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataUpdateFailure)
		return
	}

	a.logger.Info().Str("id", id.String()).Msg("book updated")
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
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/books/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		e.BadRequest(w, e.InvalidIdInUrlParam)
		return
	}

	if err := a.repository.DeleteBook(id); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.DataDeletionFailure)
		return
	}

	a.logger.Info().Str("id", id.String()).Msg("book deleted")
}
