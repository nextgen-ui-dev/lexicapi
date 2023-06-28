package article

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/jellydator/validation/is"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidArticleId           = validation.NewError("article:invalid_article_id", "Invalid article id")
	ErrArticleTitleEmpty          = validation.NewError("article:title_empty", "Article title can't be empty")
	ErrArticleTitleTooLong        = validation.NewError("article:title_too_long", "Article title can't be longer than 255 characters")
	ErrInvalidArticleThumbnailUrl = validation.NewError("article:invalid_thumbnail_url", "Invalid thumbnail url")
	ErrInvalidArticleOriginalUrl  = validation.NewError("article:invalid_original_url", "Invalid original url")
	ErrArticleOriginalUrlEmpty    = validation.NewError("article:original_url_empty", "Original url can't be empty")
	ErrArticleSourceEmpty         = validation.NewError("article:source_empty", "Source can't be empty")
	ErrArticleSourceToolong       = validation.NewError("article:source_too_long", "Source can't be longer than 255 characters")
	ErrArticleAuthorTooLong       = validation.NewError("article:author_too_long", "Author can't be longer than 255 characters")
)

func validateArticleId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidArticleId
	}

	return id, nil
}

func validateArticleTitle(title string) error {
	title = strings.TrimSpace(title)
	return validation.Validate(
		&title,
		validation.Required.ErrorObject(ErrArticleTitleEmpty),
		validation.Length(1, 255).ErrorObject(ErrArticleTitleTooLong),
	)
}

func validateArticleThumbnailUrl(url string) error {
	url = strings.TrimSpace(url)
	return validation.Validate(
		&url,
		validation.When(
			!validation.IsEmpty(url),
			is.URL.ErrorObject(ErrInvalidArticleThumbnailUrl),
		),
	)
}

func validateArticleOriginalUrl(url string) error {
	url = strings.TrimSpace(url)
	return validation.Validate(
		&url,
		validation.Required.ErrorObject(ErrArticleOriginalUrlEmpty),
		is.URL.ErrorObject(ErrInvalidArticleOriginalUrl),
	)
}

func validateArticleSource(source string) error {
	source = strings.TrimSpace(source)
	return validation.Validate(
		&source,
		validation.Required.ErrorObject(ErrArticleSourceEmpty),
		validation.Length(1, 255).ErrorObject(ErrArticleSourceToolong),
	)
}

func validateArticleAuthor(author string) error {
	author = strings.TrimSpace(author)
	return validation.Validate(
		&author,
		validation.Length(1, 255).ErrorObject(ErrArticleAuthorTooLong),
	)
}
