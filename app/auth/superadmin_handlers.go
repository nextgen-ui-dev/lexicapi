package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

func superadminSignInHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body superadminSignInReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	tokens, err := superadminSignIn(ctx, body)
	if err != nil {
		switch {
		case errors.As(err, &ErrIncorrectCredentials):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, tokens)
}
