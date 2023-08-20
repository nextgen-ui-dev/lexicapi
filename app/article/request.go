package article

import "gopkg.in/guregu/null.v4"

type createArticleCategoryReq struct {
	Name string `json:"name"`
}

type updateArticleCategoryReq struct {
	Name string `json:"name"`
}

type createArticleReq struct {
	CategoryId      string      `json:"category_id"`
	Title           string      `json:"title"`
	ThumbnailUrl    null.String `json:"thumbnail_url"`
	OriginalUrl     string      `json:"original_url"`
	Source          string      `json:"source"`
	Author          null.String `json:"author"`
	IsPublished     null.Bool   `json:"is_published"`
	OriginalContent string      `json:"original_content"`
}

type updateArticleReq struct {
	CategoryId   string      `json:"category_id"`
	Title        string      `json:"title"`
	ThumbnailUrl null.String `json:"thumbnail_url"`
	OriginalUrl  string      `json:"original_url"`
	Source       string      `json:"source"`
	Author       null.String `json:"author"`
	IsPublished  null.Bool   `json:"is_published"`
}

type ArticlePaginationDirection string

const (
	NEXT     ArticlePaginationDirection = "next"
	PREVIOUS                            = "previous"
)

type createArticleTextReq struct {
	Content    string `json:"content"`
	Difficulty string `json:"difficulty"`
	IsAdapted  bool   `json:"is_adapted"`
}

type updateArticleTextReq struct {
	Content    string `json:"content"`
	Difficulty string `json:"difficulty"`
	IsAdapted  bool   `json:"is_adapted"`
}

type generateOpenAIArticleTextReq struct {
	Content    string `json:"content"`
	Difficulty string `json:"difficulty"`
	IsAdapted  bool   `json:"is_adapted"`
}

type regenerateOpenAIArticleTextReq struct {
	Content    string `json:"content"`
	Difficulty string `json:"difficulty"`
	IsAdapted  bool   `json:"is_adapted"`
}
