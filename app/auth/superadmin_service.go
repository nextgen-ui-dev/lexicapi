package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidJWTSigningMethod = errors.New("Invalid JWT signing method")
	ErrInvalidJWTClaims        = errors.New("Invalid JWT claims")
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	Scopes string `json:"scopes"`
}

func validateSuperadminAccessToken(tokenStr string) (token *jwt.Token, claims *accessTokenClaims, err error) {
	claims = &accessTokenClaims{}
	token, err = jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			_, isAccSigningMethod := t.Method.(*jwt.SigningMethodHMAC)
			if !isAccSigningMethod {
				return nil, ErrInvalidJWTSigningMethod
			}
			return jwtSecret, nil
		},
		jwt.WithIssuer(jwtIssuer),
		jwt.WithAudience(jwtIssuer),
		jwt.WithSubject(superadmin.Email),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return
	}

	claims, isValidClaims := token.Claims.(*accessTokenClaims)
	if !isValidClaims {
		return token, claims, ErrInvalidJWTClaims
	}

	return token, claims, nil
}

func generateSuperadminAccessToken() (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Audience:  []string{jwtIssuer},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        ulid.Make().String(),
			Subject:   superadmin.Email,
		},
		Scopes: "ROLE_SUPERADMIN",
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
