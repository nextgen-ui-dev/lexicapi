package article

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type ArticleTextDifficultyPreset string

const (
	ADVANCED     ArticleTextDifficultyPreset = "ADVANCED"
	INTERMEDIATE                             = "INTERMEDIATE"
	EASY                                     = "EASY"
)

type ArticleText struct {
	Id         ulid.ULID `json:"id"`
	ArticleId  ulid.ULID `json:"article_id"`
	Content    string    `json:"content"`
	Difficulty string    `json:"difficulty"`
	IsAdapted  bool      `json:"is_adapted"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  null.Time `json:"updated_at"`
	DeletedAt  null.Time `json:"deleted_at"`
}

func NewArticleText(
	articleIdStr string,
	content string,
	difficulty string,
	isAdapted bool,
) (ArticleText, map[string]error) {
	errs := make(map[string]error)

	articleId, err := validateArticleId(articleIdStr)
	if err != nil {
		errs["article_id"] = err
	}
	if err = validateArticleTextContent(content); err != nil {
		errs["content"] = err
	}
	if err = validateArticleTextDifficulty(difficulty); err != nil {
		errs["difficulty"] = err
	}
	if len(errs) != 0 {
		return ArticleText{}, errs
	}

	id := ulid.Make()

	return ArticleText{
		Id:         id,
		ArticleId:  articleId,
		Content:    content,
		Difficulty: difficulty,
		IsAdapted:  isAdapted,
		CreatedAt:  time.Now(),
	}, nil
}
