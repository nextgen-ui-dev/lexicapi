package auth

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	lexicaApiKey string
	jwtIssuer    string
	jwtSecret    []byte
	superadmin   *Superadmin

	ErrLexicaAPIKeyEmpty       = errors.New("Lexica API Key can't be empty")
	ErrJwtIssuerEmpty          = errors.New("JWT issuer can't be empty")
	ErrJWTSecretEmpty          = errors.New("JWT secret can't be empty")
	ErrSuperadminEmailEmpty    = errors.New("Superadmin email can't be empty")
	ErrSuperadminPasswordEmpty = errors.New("Superadmin password can't be empty")
)

func ConfigureLexicaAPIKey(apiKey string) {
	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) == 0 {
		log.Fatal().Err(ErrLexicaAPIKeyEmpty).Msg("Failed to configure Lexica API Key")
	}

	lexicaApiKey = apiKey
}

func ConfigureJWTProperties(issuer, secret string) {
	issuer, secret = strings.TrimSpace(issuer), strings.TrimSpace(secret)
	if len(issuer) == 0 {
		log.Fatal().Err(ErrJwtIssuerEmpty).Msg("Failed to configure JWT properties")
	}
	if len(secret) == 0 {
		log.Fatal().Err(ErrJWTSecretEmpty).Msg("Failed to configure JWT properties")
	}

	jwtIssuer = issuer
	jwtSecret = []byte(secret)
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
