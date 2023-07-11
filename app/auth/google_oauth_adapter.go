package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
)

var (
	ErrInvalidGoogleIdToken = errors.New("Invalid google id token")
)

func validateUserGoogleIdToken(ctx context.Context, idToken string) (payload *idtoken.Payload, err error) {
	payload, err = idtoken.Validate(ctx, idToken, googleOAuthClientId)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to validate user google id token")
		return payload, ErrInvalidGoogleIdToken
	}

	return
}
