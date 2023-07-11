package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/lexica-app/lexicapi/app"
	"github.com/rs/zerolog/log"
)

type contextkey string

const (
	UserInfoCtx contextkey = "auth.userinfo"
)

var (
	ErrInvalidBearerAuthHeader = errors.New("Invalid bearer authorization header")
	ErrInvalidAccessToken      = errors.New("Invalid access token")
	ErrInvalidRefreshToken     = errors.New("Invalid refresh token")
)

func extractBearerTokenFromAuthorizationHeader(authHeader string) (token string, err error) {
	if len(authHeader) < 7 {
		return token, ErrInvalidBearerAuthHeader
	}

	authType := authHeader[:6]
	if authType != "Bearer" {
		return token, ErrInvalidBearerAuthHeader
	}

	token = authHeader[7:]
	return token, nil
}

func UserRefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("X-Forwarded-Authorization")
			if authHeader == "" {
				app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidBearerAuthHeader)
				return
			}
		}

		tokenStr, err := extractBearerTokenFromAuthorizationHeader(authHeader)
		if err != nil {
			app.WriteHttpError(w, http.StatusUnauthorized, err)
			return
		}

		user, err := validateUserRefreshToken(ctx, tokenStr)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to validate user refresh token")
			app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidRefreshToken)
			return
		}

		ctx = context.WithValue(ctx, UserInfoCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("X-Forwarded-Authorization")
			if authHeader == "" {
				app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidBearerAuthHeader)
				return
			}
		}

		tokenStr, err := extractBearerTokenFromAuthorizationHeader(authHeader)
		if err != nil {
			app.WriteHttpError(w, http.StatusUnauthorized, err)
			return
		}

		user, err := validateUserAccessToken(ctx, tokenStr)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to validate user access token")
			app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		ctx = context.WithValue(ctx, UserInfoCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SuperadminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("X-Forwarded-Authorization")
			if authHeader == "" {
				app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidBearerAuthHeader)
				return
			}
		}

		tokenStr, err := extractBearerTokenFromAuthorizationHeader(authHeader)
		if err != nil {
			app.WriteHttpError(w, http.StatusUnauthorized, err)
			return
		}

		_, _, err = validateSuperadminAccessToken(tokenStr)
		if err != nil {
			log.Debug().Err(err).Msg("Failed to validate superadmin access token")
			app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}
