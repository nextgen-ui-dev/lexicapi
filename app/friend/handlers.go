package friend

import (
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
		case ErrCantSendFriendRequestToSelf:
			app.WriteHttpError(w, http.StatusBadRequest, err)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, friend)
}
