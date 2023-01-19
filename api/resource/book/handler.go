package book

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	e "myapp/api/resource/common/err"
	"myapp/util/validator"
)

func (a *API) List(w http.ResponseWriter, r *http.Request) {
	books, err := a.repository.ListBooks()
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.DataAccessFailure)
		return
	}

	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(books.ToDto()); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.JsonEncodingFailure)
		return
	}
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &FormBook{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		e.ValError(w, e.JsonDecodingFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.AppError(w, e.FormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			e.AppError(w, e.JsonEncodingFailure)
			return
		}

		e.FormValErrors(w, respBody)
		return
	}

	book, err := a.repository.CreateBook(form.ToModel())
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.DataCreationFailure)
		return
	}

	a.logger.Info().Uint("id", book.ID).Msg("new book created")
	w.WriteHeader(http.StatusCreated)
}

func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.InvalidIdInUrlParam)
		return
	}

	book, err := a.repository.ReadBook(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.DataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.JsonEncodingFailure)
		return
	}
}

func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.InvalidIdInUrlParam)
		return
	}

	form := &FormBook{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ValError(w, e.JsonDecodingFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.AppError(w, e.FormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			e.AppError(w, e.JsonEncodingFailure)
			return
		}

		e.FormValErrors(w, respBody)
		return
	}

	bookModel := form.ToModel()
	bookModel.ID = uint(id)

	if err := a.repository.UpdateBook(bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.DataUpdateFailure)
		return
	}

	a.logger.Info().Uint64("id", id).Msg("book updated")
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.InvalidIdInUrlParam)
		return
	}

	if err := a.repository.DeleteBook(uint(id)); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.AppError(w, e.DataDeletionFailure)
		return
	}

	a.logger.Info().Uint64("id", id).Msg("book deleted")
}
