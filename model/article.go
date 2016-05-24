package model

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Id              bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	ArticleId       int           `bson:"article_id" json:"article_id"`
	Headline        string        `bson:"headline" json:"headline`
	Subheadline     string        `bson:"subheadline" json:"subheadline"`
	Section         string        `bson:"section" json:"section"`
	Subsection      string        `bson:"subsection" json:"subsection"`
	Sections        []string      `bson:"sections" json"sections"`
	Source          string        `bson:"source" json:"source"`
	Created_at      time.Time     `bson:"created_at" json:"created_at"`
	Updated_at      time.Time     `bson:"updated_at" json:"updated_at"`
	Timestamp       time.Time     `bson:"timestamp" json:"timestamp"`
	Url             string        `bson:"url" json:"url"`
	ShortUrl        string        `bson:"shortUrl" json:"shortUrl"`
	Photo           *Photo        `bson:"photo" json:"photo"`
	Body            string        `bson:"body" json:"body"`
	Summary         []string      `bson"summary" json:"summary"`
	StoryHighlights []string      `bson"storyHighlights" json:"storyHighlights"`
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
	Thumbnail PhotoInfo `bson:"thumbnail"` // deprecated
	Small     PhotoInfo `bson:"small"`

	Crops map[string]PhotoInfo `bson:"crops"`
}

/*
	Implement the Save() interface
*/
func (a *Article) Save(session *mgo.Session) {
	articleCol := session.DB("").C("Article")

	update := bson.M{
		"$set": bson.M{
			"headline":        a.Headline,
			"subheadline":     a.Subheadline,
			"section":         a.Section,
			"subsection":      a.Subsection,
			"source":          a.Source,
			"sections":        a.Sections,
			"updated_at":      a.Updated_at,
			"timestamp":       a.Timestamp,
			"created_at":      a.Created_at,
			"url":             a.Url,
			"photo":           a.Photo,
			"body":            a.Body,
			"storyHighlights": a.StoryHighlights,
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

func IsBlacklisted(url string) bool {
	blacklist := []string{
		"/videos/",
		"/police-blotter/",
		"/interactives/",
		"facebook.com",
		"/errors/404",
		"http://live.",
	}

	for _, item := range blacklist {
		if strings.Contains(url, item) {
			return true
		}
	}

	return false
}
