package auth

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrAccountDoesNotExist = errors.New("Account does not exist")

func findUserByEmail(ctx context.Context, tx pgx.Tx, email string) (user User, err error) {
	q := "SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &user, q, email); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return user, ErrUserDoesNotExist
		}

		log.Err(err).Msg("Failed to find user by email")
		return user, err
	}

	return user, nil
}

func findAccountByProviderAndProviderAccountId(ctx context.Context, tx pgx.Tx, provider AccountProvider, providerAccountId string) (account Account, err error) {
	q := "SELECT * FROM accounts WHERE provider = $1 AND provider_account_id = $2 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &account, q, provider, providerAccountId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return account, ErrAccountDoesNotExist
		}

		log.Err(err).Msg("Failed to find account by provider and provider account id")
		return account, err
	}

	return account, nil
}
