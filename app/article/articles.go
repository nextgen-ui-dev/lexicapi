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

type ArticleWithRowNumber struct {
	ArticleViewModel
	Row uint `json:"row"`
}
