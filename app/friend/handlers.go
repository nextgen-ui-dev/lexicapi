package friend

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lexica-app/lexicapi/app"
	"github.com/lexica-app/lexicapi/app/auth"
)

func sendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	requesteeId := chi.URLParam(r, "requesteeId")
	friend, errs, err := sendFriendRequest(ctx, user.Id.String(), requesteeId)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		case ErrCantSendFriendRequestToSelf, ErrAlreadyFriends, ErrFriendRequestAlreadyExist:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, friend)
}

func acceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestee, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	friendId := chi.URLParam(r, "friendId")

	friend, err := acceptFriendRequest(ctx, requestee.Id.String(), friendId)
	if err != nil {
		switch {
		case
			errors.As(err, &ErrInvalidFriendId),
			errors.As(err, &ErrInvalidRequesteeId),
			errors.Is(err, ErrFriendRequestAlreadyAccepted),
			errors.Is(err, ErrFriendRequestAlreadyRejected),
			errors.Is(err, ErrCantAcceptFriendRequestOfOtherRequestee):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrFriendDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, friend)
}

func rejectFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestee, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	friendId := chi.URLParam(r, "friendId")

	friend, err := rejectFriendRequest(ctx, requestee.Id.String(), friendId)
	if err != nil {
		switch {
		case
			errors.As(err, &ErrInvalidFriendId),
			errors.As(err, &ErrInvalidRequesteeId),
			errors.Is(err, ErrFriendRequestAlreadyAccepted),
			errors.Is(err, ErrFriendRequestAlreadyRejected),
			errors.Is(err, ErrCantRejectFriendRequestOfOtherRequestee):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrFriendDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, friend)
}

func unfriendHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	friendId := chi.URLParam(r, "friendId")

	friend, err := unfriend(ctx, user.Id.String(), friendId)
	if err != nil {
		switch {
		case
			errors.As(err, &ErrInvalidFriendId),
			errors.As(err, &ErrInvalidRequesteeId),
			errors.Is(err, ErrFriendRequestPending),
			errors.Is(err, ErrFriendRequestAlreadyRejected),
			errors.Is(err, ErrCantUnfriendOtherUserFriend):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrFriendDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, friend)
}

func getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}
	
	friends, err := getFriends(ctx, user.Id)
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, friends)
}