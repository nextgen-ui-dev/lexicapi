package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func generateSuperadminAccessToken() (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    jwtIssuer,
		"aud":    []string{jwtIssuer},
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour).Unix(),
		"jti":    ulid.Make(),
		"sub":    superadmin.Email,
		"scopes": "ROLE_SUPERADMIN",
	})

	token, err = tokenObj.SignedString(jwtSecret)
	if err != nil {
		log.Err(err).Msg("Failed to generate superadmin access token")
		return
	}

	return token, nil
}

func superadminSignIn(ctx context.Context, body superadminSignInReq) (tokens SuperadminTokens, err error) {
	if err = superadmin.ValidateCredentials(body.Email, body.Password); err != nil {
		return
	}

	accessToken, err := generateSuperadminAccessToken()
	if err != nil {
		return
	}

	return SuperadminTokens{
		AccessToken: accessToken,
	}, nil
}
