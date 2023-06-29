package article

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidArticleTextId         = validation.NewError("article:invalid_article_text_id", "Invalid article text id")
	ErrArticleTextContentEmpty      = validation.NewError("article:content_empty", "Content can't be empty")
	ErrArticleTextDifficultyEmpty   = validation.NewError("article:difficulty_empty", "Difficulty can't be empty")
	ErrArticleTextDifficultyTooLong = validation.NewError("article:difficulty_too_long", "Difficulty can't be longer than 25 characters")
)

func validateArticleTextId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidArticleTextId
	}

	return id, nil
}

func validateArticleTextContent(content string) error {
	content = strings.TrimSpace(content)
	return validation.Validate(
		&content,
		validation.Required.ErrorObject(ErrArticleTextContentEmpty),
	)
}

func validateArticleTextDifficulty(difficulty string) error {
	difficulty = strings.TrimSpace(difficulty)
	return validation.Validate(
		&difficulty,
		validation.Required.ErrorObject(ErrArticleTextDifficultyEmpty),
		validation.Length(1, 25).ErrorObject(ErrArticleTextDifficultyTooLong),
	)
}
