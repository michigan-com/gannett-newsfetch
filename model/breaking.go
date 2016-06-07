package model

type BreakingNewsSnapshot struct {
	Articles []*BreakingNewsArticle `bson:"articles" json:"articles"`
}

type BreakingNewsArticle struct {
	ArticleId   int         `bson:"article_id" json:"article_id"`
	Headline    string      `bson:"headline" json:"headline"`
	Subheadline string      `bson:"subheadline" json:"subheadline"`
	Photo       *Photo      `bson:"photo" json:"photo"`
	Video       *AssetVideo `bson:"video" json:"video"`
}
