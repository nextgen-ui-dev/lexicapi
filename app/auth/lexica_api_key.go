package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/lexica-app/lexicapi/app"
	"github.com/rs/zerolog/log"
)

var (
	lexicaApiKey string

	ErrLexicaAPIKeyEmpty   = errors.New("Lexica API Key can't be empty")
	ErrInvalidLexicaAPIKey = errors.New("You are not allowed to access this resource")
)

func ConfigureLexicaAPIKey(apiKey string) {
	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) == 0 {
		log.Fatal().Err(ErrLexicaAPIKeyEmpty).Msg("Failed to configure Lexica API Key")
	}

	lexicaApiKey = apiKey
}

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
