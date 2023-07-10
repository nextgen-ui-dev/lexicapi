package article

import (
	"gopkg.in/guregu/null.v4"
)

type ArticlesMetadata struct {
	Cursor   null.String `json:"cursor"`
	Total    uint        `json:"total"`
	FirstRow uint        `json:"first_row"`
	LastRow  uint        `json:"last_row"`
}

type Articles struct {
	ArticlesMetadata
	Articles []*ArticleWithRowNumber `json:"articles"`
}

type ArticleViewModel struct {
	Article
	Teaser       string `json:"teaser"`
	CategoryName string `json:"category_name"`
}

type ArticleDetail struct {
	Article
	CategoryName string                 `json:"category_name"`
	Texts        map[string]ArticleText `json:"texts"`
}

type ArticleWithRowNumber struct {
	ArticleViewModel
	Row uint `json:"row"`
}
