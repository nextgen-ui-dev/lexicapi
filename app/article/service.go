package article

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
)

func getArticleCategories(ctx context.Context, query string, limit uint) (categories []*ArticleCategory, err error) {
	query = strings.TrimSpace(query)

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get article categories")
		return
	}

	defer tx.Rollback(ctx)

	categories, err = findArticleCategories(ctx, tx, query, limit)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get article categories")
		return
	}

	return categories, nil
}

func getArticleCategoryById(ctx context.Context, idStr string) (category ArticleCategory, err error) {
	id, err := validateArticleCategoryId(idStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get article category by id")
		return
	}

	defer tx.Rollback(ctx)

	category, err = findArticleCategoryById(ctx, tx, id)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get article category by id")
		return
	}

	return category, nil
}

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

func deleteArticleCategory(ctx context.Context, idStr string) (err error) {
	id, err := validateArticleCategoryId(idStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to delete article category")
		return
	}

	defer tx.Rollback(ctx)

	err = deleteArticleCategoryById(ctx, tx, id)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to delete article category")
		return
	}

	return nil
}

func updateArticleCategory(ctx context.Context, idStr, name string) (category ArticleCategory, err error) {
	id, err := validateArticleCategoryId(idStr)
	if err != nil {
		return
	}
	if err = validateArticleCategoryName(name); err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to update article category")
		return
	}

	defer tx.Rollback(ctx)

	category, err = updateArticleCategoryById(ctx, tx, id, name)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to update article category")
		return
	}

	return category, nil
}
