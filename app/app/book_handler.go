package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"

	"myapp/model"
	"myapp/repository"
	"myapp/util/validator"
)

func (app *App) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.ListBooks(app.db)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}

	dtos := books.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	form := &model.BookForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	book, err := repository.CreateBook(app.db, bookModel)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataCreationFailure)
		return
	}

	app.logger.Info().Msgf("New book created: %d", book.ID)
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	book, err := repository.ReadBook(app.db, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	form := &model.BookForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		app.logger.Warn().Err(err).Msg("")

		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	bookModel, err := form.ToModel()
	if err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	bookModel.ID = uint(id)
	if err := repository.UpdateBook(app.db, bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataUpdateFailure)
		return
	}

	app.logger.Info().Msgf("Book updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		app.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := repository.DeleteBook(app.db, uint(id)); err != nil {
		app.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}

	app.logger.Info().Msgf("Book deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
