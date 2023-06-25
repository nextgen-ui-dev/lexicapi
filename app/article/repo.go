package article

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

var (
	ErrArticleCategoryNameExists = errors.New("Article category with that name exists")
)

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
