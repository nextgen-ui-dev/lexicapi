package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"
	"gopkg.in/guregu/null.v4"
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

func extractProfileFromGoogleIdTokenPayload(p *idtoken.Payload) (accountId string, name, email, imageUrl null.String) {
	claims := p.Claims

	if nameStr, ok := claims["name"].(string); !ok {
		name = null.NewString("", false)
	} else {
		name = null.StringFrom(nameStr)
	}

	if emailStr, ok := claims["email"].(string); !ok {
		email = null.NewString("", false)
	} else {
		email = null.StringFrom(emailStr)
	}

	if pictureStr, ok := claims["picture"].(string); !ok {
		imageUrl = null.NewString("", false)
	} else {
		imageUrl = null.StringFrom(pictureStr)
	}

	return p.Subject, name, email, imageUrl
}
