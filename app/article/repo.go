package article

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
	ErrArticleCategoryNameExists   = errors.New("Article category with that name exists")
	ErrArticleCategoryDoesNotExist = errors.New("Article category does not exist")
)

func findArticleCategories(ctx context.Context, tx pgx.Tx, search string, limit uint) (categories []*ArticleCategory, err error) {
	q := "SELECT * FROM article_categories WHERE name ILIKE '%' || $1 || '%' AND deleted_at IS NULL LIMIT $2"

	categories = make([]*ArticleCategory, limit)
	if err = pgxscan.Select(ctx, tx, &categories, q, search, limit); err != nil {
		log.Err(err).Msg("Failed to find article categories")
		return
	}

	return categories, nil
}

func findArticleCategoryById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (category ArticleCategory, err error) {
	q := "SELECT * FROM article_categories WHERE id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &category, q, id); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return category, ErrArticleCategoryDoesNotExist
		}

		log.Err(err).Msg("Failed to find article category by id")
		return category, err
	}

	return category, nil
}

func saveArticleCategory(ctx context.Context, tx pgx.Tx, category ArticleCategory) (err error) {
	q := "INSERT INTO article_categories(id, name) VALUES($1, $2)"

	_, err = tx.Exec(ctx, q, category.Id, category.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrArticleCategoryNameExists
			}
		}

		log.Err(err).Msg("Failed to create article category")
		return err
	}

	return nil
}

func deleteArticleCategoryById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (err error) {
	q := "UPDATE article_categories SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	_, err = tx.Exec(ctx, q, id)
	if err != nil {
		log.Err(err).Msg("Failed to delete article category")
		return err
	}

	return nil
}

func updateArticleCategoryById(ctx context.Context, tx pgx.Tx, id ulid.ULID, name string) (category ArticleCategory, err error) {
	q := "UPDATE article_categories SET name = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL RETURNING *"

	err = pgxscan.Get(ctx, tx, &category, q, id, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return category, ErrArticleCategoryNameExists
			}
		}

		if err.Error() == "scanning one: no rows in result set" {
			return category, ErrArticleCategoryDoesNotExist
		}

		log.Err(err).Msg("Failed to update article category")
		return category, err
	}

	return category, nil
}
