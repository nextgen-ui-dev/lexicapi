package article

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

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