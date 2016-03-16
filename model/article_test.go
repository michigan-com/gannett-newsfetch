package model

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/gannett-newsfetch/lib"
)

var testDbString = "mongodb://127.0.0.1:27017/_gannettTest"

func TestArticleComparison(t *testing.T) {
	var inputArticle *Article = &Article{}
	var defaultArticleId int = -1
	var defaultTimestamp time.Time = time.Date(2016, 1, 1, 0, 0, 0, 0, time.Local)
	var laterTimestamp time.Time = defaultTimestamp.Add(1 * time.Hour)
	var session = lib.DBConnect(testDbString)
	var articleCol = session.DB("").C("Article")
	var shouldSummarize bool
	defer session.Close()

	inputArticle.ArticleId = defaultArticleId
	inputArticle.Created_at = defaultTimestamp

	// Tests the case where theres the article doesnt yet exist
	shouldSummarize = ShouldSummarizeArticle(inputArticle, session)
	if !shouldSummarize {
		t.Fatal("Should decide to summarize articles if the article is not yet found")
	}

	// Insert the article, and ensure that we dont need to update because the timestamp
	// is the same
	articleCol.Insert(bson.M{
		"article_id": inputArticle.ArticleId,
		"timestamp":  inputArticle.Created_at,
	})
	shouldSummarize = ShouldSummarizeArticle(inputArticle, session)
	if shouldSummarize {
		t.Fatal("Found articles with identical timestamps should not be summarized")
	}

	// Make a new timestamp for this article that
	inputArticle.Created_at = laterTimestamp
	shouldSummarize = ShouldSummarizeArticle(inputArticle, session)
	if !shouldSummarize {
		t.Fatal("Found articles with new timestamps should be re-summarized")
	}

	articleCol.RemoveAll(bson.M{"article_id": defaultArticleId})
}
