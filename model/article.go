package model

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ArticleId   int           `bson:"article_id" json:"article_id"`
	Headline    string        `bson:"headline" json:"headline`
	Subheadline string        `bson:"subheadline" json:"subheadline"`
	Section     string        `bson:"section" json:"section"`
	Subsection  string        `bson:"subsection" json:"subsection"`
	Sections    []string      `bson:"sections" json"sections"`
	Source      string        `bson:"source" json:"source"`
	Created_at  time.Time     `bson:"created_at" json:"created_at"`
	Updated_at  time.Time     `bson:"updated_at" json:"updated_at"`
	Timestamp   time.Time     `bson:"timestamp" json:"timestamp"`
	Url         string        `bson:"url" json:"url"`
	ShortUrl    string        `bson:"shortUrl" json:"shortUrl"`
	Photo       *Photo        `bson:"photo" json:"photo"`
	Body        string        `bson:"body" json:"body"`
}

type PhotoInfo struct {
	Url    string `bson:"url"`
	Width  int    `bson:"width"`
	Height int    `bson:"height"`
}

type Photo struct {
	Caption   string    `bson:"caption"`
	Credit    string    `bson:"credit"`
	Full      PhotoInfo `bson:"full"`
	Thumbnail PhotoInfo `bson:"thumbnail"`
}

/*
	Implement the Save() interface
*/
func (a *Article) Save(session *mgo.Session) {
	articleCol := session.DB("").C("Article")

	update := bson.M{
		"$set": bson.M{
			"headline":    a.Headline,
			"subheadline": a.Subheadline,
			"section":     a.Section,
			"subsection":  a.Subsection,
			"source":      a.Source,
			"sections":    a.Sections,
			"updated_at":  a.Updated_at,
			"timestamp":   a.Timestamp,
			"created_at":  a.Created_at,
			"url":         a.Url,
			"photo":       a.Photo,
		},
	}

	_, err := articleCol.Upsert(bson.M{"article_id": a.ArticleId}, update)
	if err != nil {
		log.Warningf(`

		Failed to save Article %d. Upsert failed:

			Err: %v
		`, a.ArticleId, err)
	}

	return
}

/*
	We should summarize the article under two scenarios:

		1) This article does not yet exist in the database
		2) This article exists in the database, but the timestamp has been updated

	Does a lookup based on Article.ArticleId
*/
func ShouldSummarizeArticle(article *Article, session *mgo.Session) bool {
	var storedArticle *Article = &Article{}
	collection := session.DB("").C("Article")
	err := collection.Find(bson.M{"article_id": article.ArticleId}).One(storedArticle)
	if err == mgo.ErrNotFound {
		return true
	} else if !sameDate(article.Created_at, storedArticle.Created_at) {
		return true
	}
	return false
}

func sameDate(date1, date2 time.Time) bool {
	date1 = date1.UTC()
	date2 = date2.UTC()
	return date1.Year() == date2.Year() &&
		date1.Month() == date2.Month() &&
		date1.Day() == date2.Day() &&
		date1.Hour() == date2.Hour() &&
		date1.Minute() == date2.Minute() &&
		date1.Second() == date2.Second()
}
