package model

import "gopkg.in/mgo.v2/bson"

type BreakingNewsSnapshot struct {
	Articles []*BreakingNewsArticle `bson:"articles" json:"articles"`
}

type BreakingNewsArticle struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ArticleId   int           `bson:"article_id" json:"article_id"`
	Headline    string        `bson:"headline" json:"headline`
	Subheadline string        `bson:"subheadline" json:"subheadline"`
}
