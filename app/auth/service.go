package auth

import (
	"context"

	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func signInWithGoogle(ctx context.Context, idToken string) (signIn UserSignIn, errs map[string]error, err error) {
	payload, err := validateUserGoogleIdToken(ctx, idToken)
	if err != nil {
		return
	}

	accountId, name, email, imageUrl := extractProfileFromGoogleIdTokenPayload(payload)

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to sign in with Google")
		return
	}

	defer tx.Rollback(ctx)

	var user User
	var account Account

	user, err = findUserByEmail(ctx, tx, email.String)
	if err != nil {
		if err != ErrUserDoesNotExist {
			return
		}

		user, errs = NewUserWithOAuth(name, email, imageUrl)
		if errs != nil {
			return
		}

		user, err = saveUser(ctx, tx, user)
		if err != nil {
			return
		}

		account, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
		if err != nil {
			if err != ErrAccountDoesNotExist {
				return
			}

			account, errs = NewAccount(
				user.Id.String(),
				string(OAUTH),
				null.NewString("", false),
				string(GOOGLE),
				accountId,
			)
			if errs != nil {
				return
			}

			account, err = saveAccount(ctx, tx, account)
			if err != nil {
				return
			}
		}
	} else {
		account, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
		if err != nil {
			if err != ErrAccountDoesNotExist {
				return
			}

			account, errs = NewAccount(
				user.Id.String(),
				string(OAUTH),
				null.NewString("", false),
				string(GOOGLE),
				accountId,
			)
			if errs != nil {
				return
			}

			account, err = saveAccount(ctx, tx, account)
			if err != nil {
				return
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to sign in with Google")
		return
	}

	accessToken, err := generateUserAccessToken(user.Id.String())
	if err != nil {
		return
	}

	refreshToken, err := generateUserRefreshToken(user.Id.String())
	if err != nil {
		return
	}

	return UserSignIn{
		UserId:       user.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil, nil
}
