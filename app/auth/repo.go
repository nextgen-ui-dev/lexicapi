package auth

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrAccountDoesNotExist = errors.New("Account does not exist")

func findUserById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (user User, err error) {
	q := "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &user, q, id); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return user, ErrUserDoesNotExist
		}

		log.Err(err).Msg("Failed to find user by id")
		return user, err
	}

	return user, nil
}

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

func saveUser(ctx context.Context, tx pgx.Tx, user User) (newUser User, err error) {
	q := `
  INSERT INTO users (id, name, email, image_url, status, created_at) VALUES
  ($1, $2, $3, $4, $5, $6)
  ON CONFLICT ON CONSTRAINT users_email_unique
  DO NOTHING
  RETURNING *
  `

	if err = pgxscan.Get(
		ctx,
		tx,
		&newUser,
		q,
		user.Id,
		user.Name,
		user.Email,
		user.ImageUrl,
		user.Status,
		user.CreatedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			newUser, err = findUserByEmail(ctx, tx, user.Email.String)
			if err != nil {
				return
			}

			return newUser, nil
		}

		log.Err(err).Msg("Failed to save user")
		return
	}

	return newUser, nil
}

func saveAccount(ctx context.Context, tx pgx.Tx, account Account) (newAccount Account, err error) {
	if _, err = findUserById(ctx, tx, account.UserId); err != nil {
		return
	}

	q := `
  INSERT INTO accounts (id, user_id, type, password_hash, provider, provider_account_id, created_at) VALUES
  ($1, $2, $3, $4, $5, $6, $7)
  ON CONFLICT ON CONSTRAINT accounts_provider_provider_account_id_unique
  DO NOTHING
  RETURNING *
  `

	if err = pgxscan.Get(
		ctx,
		tx,
		&newAccount,
		q,
		account.Id,
		account.UserId,
		account.Type,
		account.PasswordHash,
		account.Provider,
		account.ProviderAccountId,
		account.CreatedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			newAccount, err = findAccountByProviderAndProviderAccountId(ctx, tx, account.Provider, account.ProviderAccountId)
			if err != nil {
				return
			}

			return newAccount, nil
		}

		log.Err(err).Msg("Failed to save account")
		return
	}

	return newAccount, nil
}

func updateUserForOnboarding(ctx context.Context, tx pgx.Tx, user User, interests []ulid.ULID) (onboardedUser User, err error) {
	if _, err = findUserById(ctx, tx, user.Id); err != nil {
		return
	}

	q := `
  UPDATE users
  SET role = $1, education_level = $2, status = $3, updated_at = $4
  WHERE id = $5 AND deleted_at IS NULL
  RETURNING *
  `

	if err = pgxscan.Get(
		ctx,
		tx,
		&onboardedUser,
		q,
		user.Role,
		user.EducationLevel,
		user.Status,
		user.UpdatedAt,
		user.Id,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return User{}, ErrUserDoesNotExist
		}

		log.Err(err).Msg("Failed to update user for onboarding")
		return
	}

	return onboardedUser, nil
}
