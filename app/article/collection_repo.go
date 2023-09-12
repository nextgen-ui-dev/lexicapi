package article

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func findCollectionDetail(ctx context.Context, tx pgx.Tx, collectionId, userId ulid.ULID) (detail CollectionDetail, err error) {
	collection, err := findCollectionById(ctx, tx, collectionId)
	if err != nil {
		return
	}

	//TODO: handle for shared users when implemented
	if collection.CreatorId.Compare(userId) != 0 {
		return detail, ErrCollectionDoesNotExist
	}

	q := `
	SELECT c.*, u.name creator_name, COUNT(ca.id) number_of_articles
	FROM collections c
	INNER JOIN users u
	ON u.id = c.creator_id
	INNER JOIN collection_articles ca
	ON ca.collection_id = c.id
	WHERE
	  c.id = $1 AND
	  c.creator_id = $2 AND
	  c.deleted_at IS NULL AND
	  ca.deleted_at IS NULL
	GROUP BY c.id, u.name
	ORDER BY c.created_at DESC
	`

	var metadata CollectionMetadata
	if err = pgxscan.Get(ctx, tx, &metadata, q, collectionId, userId); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return detail, ErrCollectionDoesNotExist
		}

		log.Err(err).Msg("Failed to find collection detail")
		return detail, err
	}

	q = `
	SELECT
	  a.*,
	  (CASE WHEN ac.deleted_at IS NULL THEN ac.name ELSE 'Deleted Category' END) category_name,
	  (CASE WHEN LENGTH(at.content) >= 255 THEN SUBSTRING(at.content, 1, 255) || '...' ELSE at.content END) teaser
	FROM articles a
	INNER JOIN article_categories ac
	ON a.category_id = ac.id
	INNER JOIN article_texts at
	ON a.id = at.article_id
	INNER JOIN collection_articles ca
	ON ca.article_id = a.id
	INNER JOIN collections c
	ON c.id = ca.collection_id
	WHERE
	  a.is_published = TRUE AND
	  a.deleted_at IS NULL AND
	  at.difficulty = 'ADVANCED' AND
	  at.deleted_at IS NULL AND
	  c.id = $1 AND
	  ca.deleted_at IS NULL
	`

	var articles []*ArticleViewModel
	if err = pgxscan.Select(ctx, tx, &articles, q, collectionId); err != nil {
		log.Err(err).Msg("Failed to find collection detail")
		return
	}

	detail = CollectionDetail{
		CollectionMetadata: metadata,
		Articles: articles,
	}

	return detail, nil
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

func deleteCollectionEntity(ctx context.Context, tx pgx.Tx, collection Collection) (Collection, error) {
	q := `
	UPDATE collections
	SET deleted_at = $2
	WHERE id = $1 AND deleted_at IS NULL
	RETURNING *
	`

	if err := pgxscan.Get(
		ctx,
		tx,
		&collection,
		q,
		collection.Id,
		collection.DeletedAt,
	); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return Collection{}, ErrCollectionDoesNotExist
		}

		log.Err(err).Msg("Failed to delete collection entity")
		return Collection{}, err
	}

	return collection, nil
}

func findCollectionsByCreatorId(ctx context.Context, tx pgx.Tx, creatorId ulid.ULID) (collections []*CollectionMetadata, err error) {
	q := `
	SELECT c.*, u.name creator_name, COUNT(ca.id) number_of_articles
	FROM collections c
	INNER JOIN users u
	ON u.id = c.creator_id
	INNER JOIN collection_articles ca
	ON ca.collection_id = c.id
	WHERE
	  c.creator_id = $1 AND
	  c.deleted_at IS NULL AND
	  ca.deleted_at IS NULL
	GROUP BY c.id, u.name
	ORDER BY c.created_at DESC
	`

	collections = []*CollectionMetadata{}
	if err = pgxscan.Select(ctx, tx, &collections, q, creatorId); err != nil {
		log.Err(err).Msg("Failed to find collections by creator id")
		return
	}

	return collections, nil
}

func findAddedCollectionsByArticleIdAndCreatorId(ctx context.Context, tx pgx.Tx, articleId, creatorId ulid.ULID) (collections []*Collection, err error) {
	if _, err := findArticleById(ctx, tx, articleId); err != nil {
		return collections, err
	}

	q := `
	SELECT c.*
	FROM collections c
	WHERE EXISTS (
	  SELECT ca.id
	  FROM collection_articles ca
	  WHERE EXISTS (
		SELECT a.id
		FROM articles a
		WHERE
		  a.id = $2 AND
		  a.id = ca.article_id AND
		  a.deleted_at IS NULL
	  ) AND
	    c.id = ca.collection_id AND
		ca.deleted_at IS NULL
	) AND EXISTS (
	  SELECT u.id
	  FROM users u
	  WHERE 
	    u.id = $1 AND
		u.id = c.creator_id AND
		u.deleted_at IS NULL
	) AND
	  c.deleted_at IS NULL
	ORDER BY c.created_at DESC
	`

	collections = []*Collection{}
	if err = pgxscan.Select(ctx, tx, &collections, q, creatorId, articleId); err != nil {
		log.Err(err).Msg("Failed to find added collections by article id and creator id")
		return
	}

	return collections, nil
}

func createCollectionArticles(ctx context.Context, tx pgx.Tx, articleId, creatorId ulid.ULID, collectionIds []ulid.ULID) (collections []*Collection, err error) {
	if _, err := findArticleById(ctx, tx, articleId); err != nil {
		return collections, err
	}

	q := `
	UPDATE collection_articles ca
	SET deleted_at = $3
	FROM collections c, users u
	WHERE 
	  ca.collection_id = c.id AND
	  ca.article_id = $1 AND
	  c.creator_id = u.id AND
	  u.id = $2 AND
	  ca.deleted_at IS NULL
	`

	now := time.Now()

	if _, err = tx.Exec(ctx, q, articleId, creatorId, now); err != nil {
		log.Err(err).Msg("Failed to create collection articles")
		return
	}

	if len(collectionIds) > 0 {
		q = "INSERT INTO collection_articles (id, collection_id, article_id, created_at) VALUES"

		params := []any{articleId, now}
		paramCount := 3
		for i, collectionId := range collectionIds {
			q += fmt.Sprintf("\n($%d, $%d, $1, $2)", paramCount, paramCount+1)
			if i+1 < len(collectionIds) {
				q += ","
			}

			params = append(params, ulid.Make(), collectionId)
			paramCount += 2
		}

		q += "\nON CONFLICT DO NOTHING"

		if _, err = tx.Exec(ctx, q, params...); err != nil {
			log.Err(err).Msg("Failed to create collection articles")
			return
		}
	}

	collections, err = findAddedCollectionsByArticleIdAndCreatorId(ctx, tx, articleId, creatorId)
	if err != nil {
		return
	}

	return collections, nil
}

func findPublicCollections(ctx context.Context, tx pgx.Tx) (collections []*CollectionMetadata, err error) {
	q := `
	SELECT c.*, u.name creator_name, COUNT(1) number_of_articles
	FROM collections c
	INNER JOIN users u
	ON u.id = c.creator_id
	WHERE
	  c.visibility = 'public' AND
	  c.deleted_at IS NULL
	GROUP BY c.id, u.name
	ORDER BY c.created_at DESC
	`

	collections = []*CollectionMetadata{}
	if err = pgxscan.Select(ctx, tx, &collections, q); err != nil {
		log.Err(err).Msg("Failed to find public collections")
		return
	}

	return collections, nil
}
