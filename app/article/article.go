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
	isPublished bool,
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
		IsPublished:  isPublished,
		CreatedAt:    time.Now(),
	}, nil
}
