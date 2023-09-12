package article

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
	"github.com/lexica-app/lexicapi/app/auth"
)

func createCollectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	var body createCollectionReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	collection, errs, err := createCollection(ctx, user.Id.String(), body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, collection)
}

func updateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	collectionId := chi.URLParam(r, "collectionId")

	var body updateCollectionReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	collection, errs, err := updateCollection(ctx, collectionId, user.Id.String(), body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidCollectionId), errors.As(err, &ErrInvalidCollectionCreatorId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrCollectionDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		case errors.Is(err, ErrNotAllowedToUpdateCollection):
			app.WriteHttpError(w, http.StatusForbidden, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, collection)
}

func deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	collectionId := chi.URLParam(r, "collectionId")

	collection, err := deleteCollection(ctx, collectionId, user.Id.String())
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidCollectionId), errors.As(err, &ErrInvalidCollectionCreatorId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrCollectionDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		case errors.Is(err, ErrNotAllowedToDeleteCollection):
			app.WriteHttpError(w, http.StatusForbidden, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, collection)
}

func getOwnCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creator, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	collections, err := getOwnCollections(ctx, creator.Id)
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, collections)
}

func getAddedCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creator, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	articleId := chi.URLParam(r, "articleId")

	collections, err := getAddedCollections(ctx, creator.Id, articleId)
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

	app.WriteHttpBodyJson(w, http.StatusOK, collections)
}

func addArticleToCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	creator, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	articleId := chi.URLParam(r, "articleId")

	var body addArticleToCollectionsReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	collections, err := addArticleToCollections(ctx, creator.Id, articleId, body.CollectionIds)
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

	app.WriteHttpBodyJson(w, http.StatusOK, collections)
}

func getPublicCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	collections, err := getPublicCollections(ctx)
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, collections)
}
