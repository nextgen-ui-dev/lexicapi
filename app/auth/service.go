package auth

import (
	"context"

	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

func refreshToken(ctx context.Context, user User) (signIn UserSignIn, err error) {
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
	}, nil
}

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

		_, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
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

			_, err = saveAccount(ctx, tx, account)
			if err != nil {
				return
			}
		}
	} else {
		_, err = findAccountByProviderAndProviderAccountId(ctx, tx, GOOGLE, accountId)
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

			_, err = saveAccount(ctx, tx, account)
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

func onboardUser(ctx context.Context, user User, body onboardReq) (User, map[string]error, error) {
	errs := user.Onboard(body.Role, body.EducationLevel)
	if errs != nil {
		return User{}, errs, nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to onboard user")
		return User{}, nil, err
	}

	defer tx.Rollback(ctx)

	user, err = updateUserForOnboarding(ctx, tx, user, body.InterestIds)
	if err != nil {
		return User{}, nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to onboard user")
		return User{}, nil, err
	}

	return user, nil, nil
}
