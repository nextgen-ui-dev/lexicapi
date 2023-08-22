package friend

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrFriendDoesNotExist        = errors.New("friend does not exist")
	ErrFriendRequestAlreadyExist = errors.New("friend request already exists")
	ErrAlreadyFriends            = errors.New("requester and requestee are already friends")
)

func findPendingOrFriendedFriendByUserIds(ctx context.Context, tx pgx.Tx, firstUserId ulid.ULID, secondUserId ulid.ULID) (friend Friend, err error) {
	q := `
	SELECT *
	FROM friends
	WHERE
	  ((requester_id = $1 AND requestee_id = $2) OR (requester_id = $2 AND requestee_id = $1)) AND
	  status != 'rejected' AND
	  deleted_at IS NULL
	`

	if err = pgxscan.Get(ctx, tx, &friend, q, firstUserId, secondUserId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Friend{}, ErrFriendDoesNotExist
		}

		log.Err(err).Msg("Failed to find friend by user ids")
		return Friend{}, err
	}

	return friend, nil
}

func createFriend(ctx context.Context, tx pgx.Tx, friend Friend) (newFriend Friend, err error) {
	q := `INSERT INTO friends(id, requester_id, requestee_id, status, created_at) VALUES
	($1, $2, $3, $4, $5)
	RETURNING *
	`

	if err = pgxscan.Get(
		ctx,
		tx,
		&newFriend,
		q,
		friend.Id,
		friend.RequesterId,
		friend.RequesteeId,
		friend.Status,
		friend.CreatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return Friend{}, ErrFriendRequestAlreadyExist
			}

			log.Err(err).Msg("Failed to create friend")
			return
		}

		log.Err(err).Msg("Failed to create friend")
		return
	}

	return newFriend, nil
}
