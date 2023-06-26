package article

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
)

func getArticleCategoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	category, err := getArticleCategoryById(ctx, id)
	if err != nil {
		switch err {
		case ErrInvalidArticleCategoryId:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case ErrArticleCategoryDoesNotExist:
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, category)
}

func createArticleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body createArticleCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	category, err := createArticleCategory(ctx, body.Name)
	if err != nil {
		switch err {
		case ErrArticleCategoryNameExists, ErrArticleCategoryNameTooLong, ErrArticleCategoryNameEmpty:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, category)
}

func deleteArticleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if err := deleteArticleCategory(ctx, id); err != nil {
		switch err {
		case ErrInvalidArticleCategoryId:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
