package article

import (
	"context"

	"github.com/rs/zerolog/log"
)

func createCollection(ctx context.Context, creatorIdStr string, body createCollectionReq) (collection Collection, errs map[string]error, err error) {
	collection, errs = NewCollection(creatorIdStr, body.Name, body.Visibility)
	if errs != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create collection")
		return
	}

	defer tx.Rollback(ctx)

	collection, err = insertCollection(ctx, tx, collection)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create collection")
		return
	}

	return collection, nil, nil
}

func updateCollection(ctx context.Context, collectionIdStr, creatorIdStr string, body updateCollectionReq) (collection Collection, errs map[string]error, err error) {
	collectionId, err := validateCollectionId(collectionIdStr)
	if err != nil {
		return
	}

	creatorId, err := validateCollectionCreatorId(creatorIdStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to update collection")
		return
	}

	defer tx.Rollback(ctx)

	collection, err = findCollectionById(ctx, tx, collectionId)
	if err != nil {
		return
	}

	errs, err = collection.Update(creatorId, body.Name, body.Visibility)
	if err != nil || errs != nil {
		return
	}

	collection, err = updateCollectionEntity(ctx, tx, collection)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to update collection")
		return
	}

	return collection, nil, nil
}