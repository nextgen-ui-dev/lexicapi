package auth

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	lexicaApiKey          string
	jwtIssuer             string
	jwtSuperadminSecret   []byte
	jwtAccessTokenSecret  []byte
	jwtRefreshTokenSecret []byte
	googleOAuthClientId   string
	superadmin            *Superadmin

	ErrGoogleOAuthClientIdEmpty   = errors.New("Google OAuth client id can't be empty")
	ErrLexicaAPIKeyEmpty          = errors.New("Lexica API Key can't be empty")
	ErrJwtIssuerEmpty             = errors.New("JWT issuer can't be empty")
	ErrJWTSuperadminSecretEmpty   = errors.New("JWT superadmin secret can't be empty")
	ErrJWTAccessTokenSecretEmpty  = errors.New("JWT access token secret can't be empty")
	ErrJWTRefreshTokenSecretEmpty = errors.New("JWT refresh token secret can't be empty")
	ErrSuperadminEmailEmpty       = errors.New("Superadmin email can't be empty")
	ErrSuperadminPasswordEmpty    = errors.New("Superadmin password can't be empty")
)

func ConfigureGoogleOAuth(clientId string) {
	clientId = strings.TrimSpace(clientId)
	if len(clientId) == 0 {
		log.Fatal().Err(ErrGoogleOAuthClientIdEmpty).Msg("Failed to configure Google OAuth")
	}

	googleOAuthClientId = clientId
}

func ConfigureLexicaAPIKey(apiKey string) {
	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) == 0 {
		log.Fatal().Err(ErrLexicaAPIKeyEmpty).Msg("Failed to configure Lexica API Key")
	}

	lexicaApiKey = apiKey
}

func ConfigureJWTProperties(
	issuer,
	superadminSecret,
	accessTokenSecret,
	refreshTokenSecret string,
) {
	issuer = strings.TrimSpace(issuer)
	if len(issuer) == 0 {
		log.Fatal().Err(ErrJwtIssuerEmpty).Msg("Failed to configure JWT properties")
	}
	superadminSecret = strings.TrimSpace(superadminSecret)
	if len(superadminSecret) == 0 {
		log.Fatal().Err(ErrJWTSuperadminSecretEmpty).Msg("Failed to configure JWT properties")
	}
	accessTokenSecret = strings.TrimSpace(accessTokenSecret)
	if len(accessTokenSecret) == 0 {
		log.Fatal().Err(ErrJWTAccessTokenSecretEmpty).Msg("Failed to configure JWT properties")
	}
	refreshTokenSecret = strings.TrimSpace(refreshTokenSecret)
	if len(refreshTokenSecret) == 0 {
		log.Fatal().Err(ErrJWTRefreshTokenSecretEmpty).Msg("Failed to configure JWT properties")
	}

	jwtIssuer = issuer
	jwtSuperadminSecret = []byte(superadminSecret)
	jwtAccessTokenSecret = []byte(accessTokenSecret)
	jwtRefreshTokenSecret = []byte(refreshTokenSecret)
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
