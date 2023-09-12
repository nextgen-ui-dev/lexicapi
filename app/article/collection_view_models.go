package article

type CollectionMetadata struct {
	Collection
	CreatorName      string `json:"creator_name"`
	NumberOfArticles uint   `json:"number_of_articles"`
}

type CollectionDetail struct {
	CollectionMetadata
	Articles []*ArticleViewModel `json:"articles"`
}