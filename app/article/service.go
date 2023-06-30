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

func createArticle(ctx context.Context, body createArticleReq) (articleDetail ArticleDetail, errs map[string]error, err error) {
	article, errs := NewArticle(body.CategoryId, body.Title, body.ThumbnailUrl, body.OriginalUrl, body.Source, body.Author, body.IsPublished)
	if errs != nil {
		return articleDetail, errs, nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to create article")
		return
	}

	defer tx.Rollback(ctx)

	article, err = saveArticle(ctx, tx, article)
	if err != nil {
		return
	}

	originalText, errs := NewArticleText(
		article.Id.String(),
		body.OriginalContent,
		string(ADVANCED),
		false,
	)
	if errs != nil {
		return articleDetail, errs, nil
	}
	originalText, err = saveArticleText(ctx, tx, originalText)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to create article")
		return
	}

	return ArticleDetail{Article: article, Texts: map[string]ArticleText{originalText.Difficulty: originalText}}, nil, nil
}

func getArticleById(ctx context.Context, idStr string) (articleDetail ArticleDetail, err error) {
	id, err := validateArticleId(idStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to get article by id")
		return
	}

	defer tx.Rollback(ctx)

	article, err := findArticleById(ctx, tx, id)
	if err != nil {
		return
	}

	texts, err := findArticleTextsByArticleId(ctx, tx, article.Id)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to get article by id")
		return
	}

	textMap := make(map[string]ArticleText)
	for _, text := range texts {
		textMap[text.Difficulty] = *text
	}

	return ArticleDetail{Article: article, Texts: textMap}, nil
}

func updateArticle(ctx context.Context, idStr string, body updateArticleReq) (article Article, errs map[string]error, err error) {
	id, err := validateArticleId(idStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to update article")
		return
	}

	defer tx.Rollback(ctx)

	article, err = findArticleById(ctx, tx, id)
	if err != nil {
		return
	}

	if errs = article.Update(
		body.CategoryId,
		body.Title,
		body.ThumbnailUrl,
		body.OriginalUrl,
		body.Source,
		body.Author,
		body.IsPublished,
	); errs != nil {
		return
	}

	article, err = updateArticleById(ctx, tx, article)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to update article")
		return
	}

	return article, nil, nil
}

func removeArticle(ctx context.Context, idStr string) (err error) {
	id, err := validateArticleId(idStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to remove article")
		return
	}

	defer tx.Rollback(ctx)

	article, err := findArticleById(ctx, tx, id)
	if err != nil {
		return
	}

	article.Delete()
	if err = deleteArticle(ctx, tx, article); err != nil {
		return
	}
	if err = deleteArticleTextsByArticle(ctx, tx, article); err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to remove article")
		return
	}

	return nil
}
