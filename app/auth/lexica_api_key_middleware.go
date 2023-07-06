package auth

import (
	"errors"
	"net/http"

	"github.com/lexica-app/lexicapi/app"
)

var ErrInvalidLexicaAPIKey = errors.New("You are not allowed to access this resource")

func LexicaAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Lexica-Api-Key")
		if apiKey != lexicaApiKey {
			app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidLexicaAPIKey)
			return
		}

		next.ServeHTTP(w, r)
	})
}
