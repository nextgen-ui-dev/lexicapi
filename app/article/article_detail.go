package article

type ArticleDetail struct {
	Article
	Texts map[string]ArticleText `json:"texts"`
}
