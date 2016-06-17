package model

type ScrapeRequest struct {
	ArticleID  int    `bson:"article_id"`
	ArticleURL string `bson:"article_url"`
}
