package article

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type Article struct {
	Id           ulid.ULID   `json:"id"`
	CategoryId   ulid.ULID   `json:"category_id"`
	Title        string      `json:"title"`
	ThumbnailUrl null.String `json:"thumbnail_url"`
	OriginalUrl  string      `json:"original_url"`
	Source       string      `json:"source"`
	Author       null.String `json:"author"`
	IsPublished  bool        `json:"is_published"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    null.Time   `json:"updated_at"`
	DeletedAt    null.Time   `json:"deleted_at"`
}

func NewArticle(
	categoryIdStr string,
	title string,
	thumbnailUrl null.String,
	originalUrl string,
	source string,
	author null.String,
	isPublished null.Bool,
) (Article, map[string]error) {
	errs := make(map[string]error)

	categoryId, err := validateArticleCategoryId(categoryIdStr)
	if err != nil {
		errs["category_id"] = err
	}
	if err = validateArticleTitle(title); err != nil {
		errs["title"] = err
	}
	if err = validateArticleThumbnailUrl(thumbnailUrl.String); err != nil {
		errs["thumbnail_url"] = err
	}
	if err = validateArticleOriginalUrl(originalUrl); err != nil {
		errs["original_url"] = err
	}
	if err = validateArticleSource(source); err != nil {
		errs["source"] = err
	}
	if err = validateArticleAuthor(author.String); err != nil {
		errs["author"] = err
	}

	if len(errs) != 0 {
		return Article{}, errs
	}

	id := ulid.Make()

	return Article{
		Id:           id,
		CategoryId:   categoryId,
		Title:        title,
		ThumbnailUrl: thumbnailUrl,
		OriginalUrl:  originalUrl,
		Source:       source,
		Author:       author,
		IsPublished:  isPublished.Valid && isPublished.Bool,
		CreatedAt:    time.Now(),
	}, nil
}

func (a *Article) Update(
	categoryIdStr null.String,
	title null.String,
	thumbnailUrl null.String,
	originalUrl null.String,
	source null.String,
	author null.String,
	isPublished null.Bool,
) map[string]error {
	errs := make(map[string]error)

	if categoryIdStr.Valid {
		categoryId, err := validateArticleCategoryId(categoryIdStr.String)
		if err != nil {
			errs["category_id"] = err
		}
		a.CategoryId = categoryId
	}

	if title.Valid {
		if err := validateArticleTitle(title.String); err != nil {
			errs["title"] = err
		}
		a.Title = title.String
	}

	if thumbnailUrl.Valid {
		if err := validateArticleThumbnailUrl(thumbnailUrl.String); err != nil {
			errs["thumbnail_url"] = err
		}
		a.ThumbnailUrl = thumbnailUrl
	}

	if originalUrl.Valid {
		if err := validateArticleOriginalUrl(originalUrl.String); err != nil {
			errs["original_url"] = err
		}
		a.OriginalUrl = originalUrl.String
	}

	if source.Valid {
		if err := validateArticleSource(source.String); err != nil {
			errs["source"] = err
		}
		a.Source = source.String
	}

	if author.Valid {
		if err := validateArticleAuthor(author.String); err != nil {
			errs["author"] = err
		}
		a.Author = author
	}

	if isPublished.Valid {
		a.IsPublished = isPublished.Bool
	}

	if len(errs) != 0 {
		return errs
	}

	a.UpdatedAt = null.NewTime(time.Now(), true)

	return nil
}

func (a *Article) Delete() {
	if !a.DeletedAt.Valid {
		a.DeletedAt = null.NewTime(time.Now(), true)
	}
}
