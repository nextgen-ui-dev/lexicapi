package auth

import (
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidRefreshToken)
		return
	}

	signIn, err := refreshToken(ctx, user)
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, signIn)
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, user)
}

func signInWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idToken := r.Header.Get("X-Google-Id-Token")

	signIn, errs, err := signInWithGoogle(ctx, idToken)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		case ErrInvalidGoogleIdToken:
			app.WriteHttpError(w, http.StatusUnauthorized, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, signIn)
}
