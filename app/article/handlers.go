package article

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
)

func getArticleCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	// Don't throw error to client just because of misinputs
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	categories, err := getArticleCategories(ctx, query, uint(limit))
	if err != nil {
		switch err {
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, categories)
}

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

func updateArticleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	var body updateArticleCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	category, err := updateArticleCategory(ctx, id, body.Name)
	if err != nil {
		switch err {
		case ErrArticleCategoryNameTooLong, ErrArticleCategoryNameEmpty,
			ErrInvalidArticleCategoryId, ErrArticleCategoryNameExists:
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

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body createArticleReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	article, errs, err := createArticle(ctx, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		case ErrArticleCategoryDoesNotExist:
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpError(w, http.StatusInternalServerError, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, article)
}
