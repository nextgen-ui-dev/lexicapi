package article

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
)

func regenerateOpenAIArticleTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	articleId := chi.URLParam(r, "articleId")

	var body regenerateOpenAIArticleTextReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	text, errs, err := regenerateOpenAIArticleText(ctx, id, articleId, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleId),
			errors.As(err, &ErrInvalidArticleTextId),
			errors.As(err, &ErrInvalidArticleTextDifficulty),
			errors.Is(err, ErrArticleTextDifficultyExist):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrInvalidOpenAIAPIKey):
			app.WriteHttpError(w, http.StatusUnauthorized, err)
		case errors.Is(err, ErrArticleDoesNotExist), errors.Is(err, ErrArticleTextDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		case errors.Is(err, ErrOpenAIRateLimited):
			app.WriteHttpError(w, http.StatusTooManyRequests, err)
		case errors.Is(err, ErrOpenAIServiceError):
			app.WriteHttpError(w, http.StatusServiceUnavailable, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, text)
}

func generateOpenAIArticleTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articleId := chi.URLParam(r, "articleId")
	var body generateOpenAIArticleTextReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	text, errs, err := generateOpenAIArticleText(ctx, articleId, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleId), errors.As(err, &ErrInvalidArticleTextDifficulty), errors.Is(err, ErrArticleTextDifficultyExist):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrInvalidOpenAIAPIKey):
			app.WriteHttpError(w, http.StatusUnauthorized, err)
		case errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		case errors.Is(err, ErrOpenAIRateLimited):
			app.WriteHttpError(w, http.StatusTooManyRequests, err)
		case errors.Is(err, ErrOpenAIServiceError):
			app.WriteHttpError(w, http.StatusServiceUnavailable, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, text)
}

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
			app.WriteHttpInternalServerError(w)
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
		switch {
		case errors.As(err, &ErrInvalidArticleCategoryId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleCategoryDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
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
			app.WriteHttpInternalServerError(w)
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
			app.WriteHttpInternalServerError(w)
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
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, category)
}

func getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query().Get("q")
	categoryId := r.URL.Query().Get("category_id")
	pageSizeStr := r.URL.Query().Get("page_size")
	direction := r.URL.Query().Get("direction")
	cursor := r.URL.Query().Get("cursor")

	includeUnpublished := strings.HasPrefix(r.URL.Path, "/admin")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 100
	}

	articles, err := getArticles(ctx, q, categoryId, uint(pageSize), direction, cursor, includeUnpublished)
	if err != nil {
		switch {
		default:
			app.WriteHttpInternalServerError(w)
		}
	}

	app.WriteHttpBodyJson(w, http.StatusOK, articles)
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
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, article)
}

func getArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	article, err := getArticleById(ctx, id)
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, article)
}

func updateArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	var body updateArticleReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	article, errs, err := updateArticle(ctx, id, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleCategoryDoesNotExist), errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, article)
}

func removeArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if err := removeArticle(ctx, id); err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createArticleTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	articleId := chi.URLParam(r, "articleId")
	var body createArticleTextReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	text, errs, err := createArticleText(ctx, articleId, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		case ErrArticleTextDifficultyExist:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case ErrArticleDoesNotExist:
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, text)
}

func updateArticleTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	articleId := chi.URLParam(r, "articleId")

	var body updateArticleTextReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	text, errs, err := updateArticleText(ctx, id, articleId, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleTextId), errors.As(err, &ErrInvalidArticleId), errors.Is(err, ErrArticleTextDifficultyExist):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleTextDoesNotExist), errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, text)
}

func removeArticleTextHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	articleId := chi.URLParam(r, "articleId")

	if err := removeArticleText(ctx, id, articleId); err != nil {
		switch {
		case errors.As(err, &ErrInvalidArticleTextId), errors.As(err, &ErrInvalidArticleId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrArticleTextDoesNotExist), errors.Is(err, ErrArticleDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
