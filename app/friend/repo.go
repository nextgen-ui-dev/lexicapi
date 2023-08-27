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

func findFriendById(ctx context.Context, tx pgx.Tx, friendId ulid.ULID) (friend Friend, err error) {
	q := "SELECT * FROM friends WHERE id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &friend, q, friendId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Friend{}, ErrFriendDoesNotExist
		}

		log.Err(err).Msg("Failed to find friend by id")
		return
	}

	return friend, nil
}

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

func updateFriend(ctx context.Context, tx pgx.Tx, friend Friend) (Friend, error) {
	q := `
	UPDATE friends
	SET status = $2, deleted_at = $3
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING *`
	
	if err := pgxscan.Get(ctx, tx, &friend, q, friend.Id, friend.Status, friend.DeletedAt); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return friend, ErrFriendDoesNotExist
		}

		log.Err(err).Msg("Failed to update friend")
		return friend, err
	}

	return friend, nil
}

func findFriendsByUserId(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (friends []*FriendDetail, err error) {
	q := `
	SELECT fr.*
	FROM (
	  SELECT fr1.*, u.name, u.email, u.image_url
	  FROM (
	    SELECT fr.*
	    FROM friends fr
	    WHERE EXISTS (
		  SELECT u.id
		  FROM users u
		  WHERE u.id = fr.requester_id
		) AND
		  fr.requester_id = $1 AND
		  fr.status = 'friended' AND
		  fr.deleted_at IS NULL
	  ) fr1
	  INNER JOIN users u
	  ON u.id = fr1.requestee_id
	  UNION
	  SELECT fr2.*, u.name, u.email, u.image_url
	  FROM (
	    SELECT fr.*
	    FROM friends fr
		WHERE EXISTS (
		  SELECT u.id
		  FROM users u
		  WHERE u.id = fr.requestee_id
		) AND
		  fr.requestee_id = $1 AND
		  fr.status = 'friended' AND
		  fr.deleted_at IS NULL
	  ) fr2
	  INNER JOIN users u
	  ON u.id = fr2.requester_id
	) fr
	ORDER BY fr.name ASC
	`

	friends = []*FriendDetail{}
	if err = pgxscan.Select(ctx, tx, &friends, q, userId); err != nil {
		log.Err(err).Msg("Failed to find friends by user id")
		return
	}

	return friends, nil
}

func findSentFriendRequestsByUserId(ctx context.Context, tx pgx.Tx, userId ulid.ULID) (friends []*FriendDetail, err error) {
	q := `
	SELECT fr.*, u.name, u.email, u.image_url
	FROM (
	  SELECT fr.*
	  FROM friends fr
	  WHERE EXISTS (
	    SELECT u.id
	    FROM users u
	    WHERE u.id = fr.requester_id
	  ) AND
	    fr.requester_id = $1 AND
	    fr.status = 'pending' AND
	    fr.deleted_at IS NULL
	) fr
	INNER JOIN users u
	ON u.id = fr.requestee_id
	ORDER BY u.name ASC
	`

	friends = []*FriendDetail{}
	if err = pgxscan.Select(ctx, tx, &friends, q, userId); err != nil {
		log.Err(err).Msg("Failed to find sent friend requests by user id")
		return
	}

	return friends, nil
}