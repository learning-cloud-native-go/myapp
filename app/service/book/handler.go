package book

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	e "myapp/app/service/error"
	"myapp/util/validator"
)

func (app *App) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := app.repository.ListBooks()
	if err != nil {
		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrDataAccessFailure)
		return
	}
	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(books.ToDto()); err != nil {
		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	form := &FormBook{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		e.ValError(w, e.ErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.AppError(w, e.ErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Error().Err(err).Msg("")
			e.AppError(w, e.ErrJsonCreationFailure)
			return
		}

		e.FormValErrors(w, respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Error().Err(err).Msg("")
		e.ValError(w, e.ErrFormDecodingFailure)
		return
	}

	book, err := app.repository.CreateBook(bookModel)
	if err != nil {
		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrDataCreationFailure)
		return
	}

	app.logger.Info().Uint("id", book.ID).Msg("new book created")
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.ErrInvalidIdInUrlParam)
		return
	}

	book, err := app.repository.ReadBook(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrDataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.ErrInvalidIdInUrlParam)
		return
	}

	form := &FormBook{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.logger.Error().Err(err).Msg("")
		e.ValError(w, e.ErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.AppError(w, e.ErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Error().Err(err).Msg("")
			e.AppError(w, e.ErrJsonCreationFailure)
			return
		}

		e.FormValErrors(w, respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Error().Err(err).Msg("")
		e.ValError(w, e.ErrFormDecodingFailure)
		return
	}

	bookModel.ID = uint(id)
	if err := app.repository.UpdateBook(bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrDataUpdateFailure)
		return
	}

	app.logger.Info().Uint64("id", id).Msg("book updated")
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		e.ValError(w, e.ErrInvalidIdInUrlParam)
		return
	}

	if err := app.repository.DeleteBook(uint(id)); err != nil {
		app.logger.Error().Err(err).Msg("")
		e.AppError(w, e.ErrDataDeletionFailure)
		return
	}

	app.logger.Info().Uint64("id", id).Msg("book deleted")
	w.WriteHeader(http.StatusAccepted)
}
