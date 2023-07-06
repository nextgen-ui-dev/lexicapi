package auth

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	lexicaApiKey        string
	superadminJWTSecret []byte
	superadmin          *Superadmin

	ErrLexicaAPIKeyEmpty        = errors.New("Lexica API Key can't be empty")
	ErrSuperadminJWTSecretEmpty = errors.New("Superadmin JWT secret can't be empty")
	ErrSuperadminEmailEmpty     = errors.New("Superadmin email can't be empty")
	ErrSuperadminPasswordEmpty  = errors.New("Superadmin password can't be empty")
)

func ConfigureLexicaAPIKey(apiKey string) {
	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) == 0 {
		log.Fatal().Err(ErrLexicaAPIKeyEmpty).Msg("Failed to configure Lexica API Key")
	}

	lexicaApiKey = apiKey
}

func ConfigureSuperadminJWTSecret(secret string) {
	secret = strings.TrimSpace(secret)
	if len(secret) == 0 {
		log.Fatal().Err(ErrSuperadminJWTSecretEmpty).Msg("Failed to configure superadminJWTSecret")
	}

	superadminJWTSecret = []byte(secret)
}

func ConfigureSuperadmin(email, password string) {
	email, password = strings.TrimSpace(email), strings.TrimSpace(password)
	if len(email) == 0 {
		log.Fatal().Err(ErrSuperadminEmailEmpty).Msg("Failed to configure superadmin")
	}
	if len(password) == 0 {
		log.Fatal().Err(ErrSuperadminPasswordEmpty).Msg("Failed to configure superadmin")
	}

	superadmin = &Superadmin{
		Email:    email,
		Password: password,
	}
}
