package article

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrCollectionDoesNotExist = errors.New("collection does not exist")
)

func findCollectionById(ctx context.Context, tx pgx.Tx, collectionId ulid.ULID) (collection Collection, err error) {
	q := `
	SELECT *
	FROM collections
	WHERE id = $1 AND deleted_at IS NULL
	`

	if err = pgxscan.Get(ctx, tx, &collection, q, collectionId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return collection, ErrCollectionDoesNotExist
		}

		log.Err(err).Msg("Failed to find collection by id")
		return collection, err
	}

	return collection, nil
}

func insertCollection(ctx context.Context, tx pgx.Tx, collection Collection) (newCollection Collection, err error) {
	q := `
	INSERT INTO collections (id, creator_id, name, visibility, created_at) VALUES
	($1, $2, $3, $4, $5)
	ON CONFLICT (id)
	DO NOTHING
	RETURNING *
	`

	if err := pgxscan.Get(
		ctx,
		tx,
		&newCollection,
		q,
		collection.Id,
		collection.CreatorId,
		collection.Name,
		collection.Visibility,
		collection.CreatedAt,
	); err != nil {
		log.Err(err).Msg("Failed to insert collection")
		return newCollection, err
	}

	return newCollection, nil
}

func updateCollectionEntity(ctx context.Context, tx pgx.Tx, collection Collection) (Collection, error) {
	q := `
	UPDATE collections
	SET name = $2, visibility = $3, updated_at = $4
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING *
	`

	if err := pgxscan.Get(
		ctx,
		tx,
		&collection,
		q,
		collection.Id,
		collection.Name,
		collection.Visibility,
		collection.UpdatedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Collection{}, ErrCollectionDoesNotExist
		}

		log.Err(err).Msg("Failed to update collection entity")
		return Collection{}, err
	}

	return collection, nil
}
