package article

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type ArticleCategory struct {
	Id        ulid.ULID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

func NewArticleCategory(name string) (category ArticleCategory, err error) {
	if err = validateArticleCategoryName(name); err != nil {
		return category, err
	}

	category = ArticleCategory{
		Id:        ulid.Make(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	return category, nil
}
