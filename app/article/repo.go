package article

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrArticleCategoryNameExists   = errors.New("Article category with that name exists")
	ErrArticleCategoryDoesNotExist = errors.New("Article category does not exist")
	ErrArticleDoesNotExist         = errors.New("Article does not exist")
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

func findArticles(
	ctx context.Context, tx pgx.Tx,
	query string, categoryId ulid.ULID, pageSize uint, direction ArticlePaginationDirection,
	cursor ulid.ULID,
) (articles Articles, err error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var categoryFilter string
	if categoryId != (ulid.ULID{}) {
		categoryFilter = "category_id = $2"
	} else {
		categoryFilter = "category_id != $2"
	}
	rowNumber := `WITH rows AS (
    SELECT ROW_NUMBER() OVER (ORDER BY id DESC) row, id 
    FROM articles 
    WHERE deleted_at IS NULL AND title ILIKE '%' || $1 || '%' AND ` + categoryFilter + `
    ORDER BY id DESC
    )`
	sBuilder := psql.Select("rows.row", "a.*").From("articles a").InnerJoin("rows ON rows.id = a.id")
	sBuilder = sBuilder.Where("title ILIKE '%' || ? || '%'").Where("deleted_at IS NULL")

	if categoryId != (ulid.ULID{}) {
		sBuilder = sBuilder.Where(sq.Eq{"category_id": categoryId})
	} else {
		sBuilder = sBuilder.Where(sq.NotEq{"category_id": categoryId})
	}

	switch direction {
	case NEXT:
		if cursor != (ulid.ULID{}) {
			sBuilder = sBuilder.Where(sq.LtOrEq{"a.id": cursor})
		} else {
			sBuilder = sBuilder.Where(sq.NotEq{"a.id": cursor})
		}
		sBuilder = sBuilder.OrderBy("a.id DESC").Limit(uint64(pageSize) + 1)
	case PREVIOUS:
		if cursor != (ulid.ULID{}) {
			sBuilder = sBuilder.Where(sq.GtOrEq{"a.id": cursor})
		} else {
			sBuilder = sBuilder.Where(sq.NotEq{"a.id": cursor})
		}
		sBuilder = sBuilder.OrderBy("a.id ASC").Limit(uint64(pageSize) + 1)
	}

	q, _, err := sBuilder.ToSql()
	if err != nil {
		log.Err(err).Msg("Failed to find articles")
		return
	}

	articleQuery := rowNumber + " " + q
	var listOfArticles []*ArticleWithRowNumber
	if err = pgxscan.Select(ctx, tx, &listOfArticles, articleQuery, query, categoryId, cursor); err != nil {
		log.Err(err).Msg("Failed to get articles")
		return articles, err
	}

	totalQuery := strings.ReplaceAll(q, "rows.row, a.*", "COUNT(*) total")
	totalQuery = strings.ReplaceAll(totalQuery, "INNER JOIN rows ON rows.id = a.id", "")
	totalQuery = strings.ReplaceAll(totalQuery, fmt.Sprintf("AND a.id >= $3 ORDER BY a.id ASC LIMIT %d", pageSize+1), "")
	totalQuery = strings.ReplaceAll(totalQuery, fmt.Sprintf("AND a.id <= $3 ORDER BY a.id DESC LIMIT %d", pageSize+1), "")
	totalQuery = strings.ReplaceAll(totalQuery, fmt.Sprintf("AND a.id <> $3 ORDER BY a.id ASC LIMIT %d", pageSize+1), "")
	totalQuery = strings.ReplaceAll(totalQuery, fmt.Sprintf("AND a.id <> $3 ORDER BY a.id DESC LIMIT %d", pageSize+1), "")
	var totalResult struct {
		Total uint
	}
	if err = pgxscan.Get(ctx, tx, &totalResult, totalQuery, query, categoryId); err != nil {
		log.Err(err).Msg("Failed to get articles")
		return articles, err
	}

	if direction == PREVIOUS {
		for i := 0; i < len(listOfArticles)/2; i++ {
			j := len(listOfArticles) - i - 1
			listOfArticles[i], listOfArticles[j] = listOfArticles[j], listOfArticles[i]
		}
	}

	var newCursor null.String
	firstRow := uint(0)
	lastRow := uint(0)
	if len(listOfArticles) > int(pageSize) {
		if direction == PREVIOUS {
			newCursor = null.StringFrom(listOfArticles[0].Id.String())
		} else {

			newCursor = null.StringFrom(listOfArticles[len(listOfArticles)-1].Id.String())
		}
		firstRow = listOfArticles[0].Row
		lastRow = listOfArticles[len(listOfArticles)-2].Row
		listOfArticles = listOfArticles[:len(listOfArticles)-1]
	} else if len(listOfArticles) > 0 {
		if direction == PREVIOUS {
			newCursor = null.StringFrom(listOfArticles[0].Id.String())
		} else {

			newCursor = null.StringFrom(listOfArticles[len(listOfArticles)-1].Id.String())
		}
		firstRow = listOfArticles[0].Row
		lastRow = listOfArticles[len(listOfArticles)-1].Row
	}

	if newCursor.String == cursor.String() {
		firstRow, lastRow = 0, 0
		newCursor = null.NewString("", false)
		listOfArticles = make([]*ArticleWithRowNumber, 0)
	}

	articles = Articles{
		ArticlesMetadata: ArticlesMetadata{
			Cursor:   newCursor,
			Total:    totalResult.Total,
			FirstRow: firstRow,
			LastRow:  lastRow,
		},
		Articles: listOfArticles,
	}

	return articles, nil
}

func saveArticle(ctx context.Context, tx pgx.Tx, article Article) (Article, error) {
	if _, err := findArticleCategoryById(ctx, tx, article.CategoryId); err != nil {
		return article, err
	}

	q := `
  INSERT INTO articles(id, category_id, title, thumbnail_url, original_url, source, author, is_published, created_at) VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  RETURNING *
  `

	var newArticle Article
	if err := pgxscan.Get(
		ctx,
		tx,
		&newArticle,
		q,
		article.Id,
		article.CategoryId,
		article.Title,
		article.ThumbnailUrl,
		article.OriginalUrl,
		article.Source,
		article.Author,
		article.IsPublished,
		article.CreatedAt,
	); err != nil {
		log.Err(err).Msg("Failed to save article")
		return newArticle, err
	}

	return newArticle, nil
}

func saveArticleText(ctx context.Context, tx pgx.Tx, text ArticleText) (ArticleText, error) {
	if _, err := findArticleById(ctx, tx, text.ArticleId); err != nil {
		return text, err
	}

	q := `INSERT INTO article_texts(id, article_id, content, difficulty, is_adapted, created_at) VALUES
  ($1, $2, $3, $4, $5, $6)
  ON CONFLICT(id)
  DO UPDATE SET article_id = $2, content = $3, difficulty = $4, is_adapted = $5, updated_at = $6
  RETURNING *
  `

	var newText ArticleText
	if err := pgxscan.Get(
		ctx,
		tx,
		&newText,
		q,
		text.Id,
		text.ArticleId,
		text.Content,
		text.Difficulty,
		text.IsAdapted,
		text.CreatedAt,
	); err != nil {
		log.Err(err).Msg("Failed to save article text")
		return text, err
	}

	return newText, nil
}

func findArticleById(ctx context.Context, tx pgx.Tx, id ulid.ULID) (article Article, err error) {
	q := "SELECT * FROM articles WHERE id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Get(ctx, tx, &article, q, id); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return article, ErrArticleDoesNotExist
		}

		log.Err(err).Msg("Failed to find article by id")
		return article, err
	}

	return article, nil
}

func findArticleTextsByArticleId(ctx context.Context, tx pgx.Tx, articleId ulid.ULID) (texts []*ArticleText, err error) {
	if _, err := findArticleById(ctx, tx, articleId); err != nil {
		return texts, err
	}

	q := "SELECT * FROM article_texts WHERE article_id = $1 AND deleted_at IS NULL"

	if err = pgxscan.Select(ctx, tx, &texts, q, articleId); err != nil {
		log.Err(err).Msg("Failed to find article texts")
		return
	}

	return texts, nil
}

func updateArticleById(ctx context.Context, tx pgx.Tx, article Article) (updatedArticle Article, err error) {
	if _, err := findArticleCategoryById(ctx, tx, article.CategoryId); err != nil {
		return article, err
	}

	q := `UPDATE articles
  SET category_id = $1, title = $2, thumbnail_url = $3, original_url = $4, 
  source = $5, author = $6, is_published = $7, updated_at = $8
  WHERE id = $9 AND deleted_at IS NULL
  RETURNING *
  `

	err = pgxscan.Get(
		ctx,
		tx,
		&updatedArticle,
		q,
		article.CategoryId,
		article.Title,
		article.ThumbnailUrl,
		article.OriginalUrl,
		article.Source,
		article.Author,
		article.IsPublished,
		article.UpdatedAt,
		article.Id,
	)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return updatedArticle, ErrArticleDoesNotExist
		}

		log.Err(err).Msg("Failed to update article category")
		return updatedArticle, err
	}

	return updatedArticle, nil
}

func deleteArticle(ctx context.Context, tx pgx.Tx, article Article) (err error) {
	q := "UPDATE articles SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"

	_, err = tx.Exec(ctx, q, article.DeletedAt, article.Id)
	if err != nil {
		log.Err(err).Msg("Failed to delete article")
		return err
	}

	return nil
}

func deleteArticleTextsByArticle(ctx context.Context, tx pgx.Tx, article Article) (err error) {
	q := "UPDATE article_texts SET deleted_at = $1 WHERE article_id = $2 AND deleted_at IS NULL"

	_, err = tx.Exec(ctx, q, article.DeletedAt, article.Id)
	if err != nil {
		log.Err(err).Msg("Failed to delete article texts")
		return err
	}

	return nil
}
