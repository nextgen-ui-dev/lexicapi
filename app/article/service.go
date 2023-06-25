package article

import (
	"context"

	"github.com/rs/zerolog/log"
)

func createArticleCategory(ctx context.Context, name string) (category ArticleCategory, err error) {
	category, err = NewArticleCategory(name)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create article category")
		return
	}

	defer tx.Rollback(ctx)

	err = saveArticleCategory(ctx, tx, category)
	if err != nil {
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create article category")
		return
	}

	return category, nil
}
