package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/mgo.v2/bson"

	c "github.com/michigan-com/gannett-newsfetch/commands"
	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

var testMongoUrl string = "mongodb://localhost:27017/gannett-newsfetch-test"

func TestIntegration(t *testing.T) {
	testArticle := getTestArticle()
	if testArticle == nil {
		t.Fatalf("Test article json reading failed")
	}

	testArticleId := testArticle.ArticleId
	session := lib.DBConnect(testMongoUrl)
	toScrapeCol := session.DB("").C("ToScrape")
	articleCol := session.DB("").C("Article")
	toScrapeCol.Insert(bson.M{"article_id": testArticleId})

	c.ScrapeAndSummarize(testMongoUrl)

	count, err := toScrapeCol.Count()
	if count != 0 {
		t.Fatalf("Should be no article IDs in the toScrape collection, there are %d", count)
	} else if err != nil {
		t.Fatalf("Failed to get toScrape Count: %v", err)
	}

	storedArticle := &m.Article{}
	err = articleCol.Find(bson.M{"article_id": testArticleId}).One(storedArticle)
	if err != nil {
		t.Fatalf("failed to get article from article collection: %v", err)
	}

	sameArticle, errString := compareArticles(testArticle, storedArticle)
	if !sameArticle {
		t.Fatalf("Articles are not equal: %s", errString)
	}
}

func compareArticles(articleOne, articleTwo *m.Article) (bool, string) {
	if articleOne.ArticleId != articleTwo.ArticleId {
		return false, fmt.Sprintf("Ids don't match: %d %d", articleOne.ArticleId, articleTwo.ArticleId)
	}

	if articleOne.Headline != articleTwo.Headline {
		return false, "Headlines dont match"
	}

	if articleOne.Subheadline != articleTwo.Subheadline {
		return false, "Subheadlines dont match"
	}

	if articleOne.Section != articleTwo.Section {
		return false, "Sections dont match"
	}

	if len(articleOne.Sections) != len(articleTwo.Sections) {
		return false, "Number of sections dont match"
	}

	for i, section := range articleOne.Sections {
		sectionTwo := articleTwo.Sections[i]
		if sectionTwo != section {
			return false, "Sections dont match"
		}
	}

	if articleOne.Source != articleTwo.Source {
		return false, "Sources dont match"
	}

	if !lib.SameTime(articleOne.Created_at, articleTwo.Created_at) {
		return false, "Created at dates dont match"
	}

	if !lib.SameTime(articleOne.Updated_at, articleTwo.Updated_at) {
		return false, "Updated at dates dont match"
	}

	if !lib.SameTime(articleOne.Timestamp, articleTwo.Timestamp) {
		return false, "Timestamps dont match"
	}

	if articleOne.Url != articleTwo.Url {
		return false, "Urls dont match"
	}

	if articleOne.ShortUrl != articleTwo.ShortUrl {
		return false, "Short urls dont match"
	}

	if photoCheck, errStr := comparePhotos(articleOne.Photo, articleTwo.Photo); !photoCheck {
		return false, errStr
	}

	if articleOne.Body != articleTwo.Body {
		return false, "Body doesnt match"
	}

	if len(articleOne.Summary) != len(articleTwo.Summary) {
		return false, "Summary lenghts dont match"
	}

	for i, summarySentence := range articleOne.Summary {
		summaryTwo := articleTwo.Summary[i]

		if summarySentence != summaryTwo {
			return false, "Summary sentences dont match"
		}
	}

	if len(articleOne.StoryHighlights) != len(articleTwo.StoryHighlights) {
		return false, "story highlights length didnt match"
	}

	for i, highlight := range articleOne.StoryHighlights {
		highlightTwo := articleTwo.StoryHighlights[i]
		if highlight != highlightTwo {
			return false, "Highlight doesnt match"
		}
	}

	return true, ""
}

func comparePhotos(photoOne, photoTwo *m.Photo) (bool, string) {
	if photoOne.Caption != photoTwo.Caption {
		return false, "Captions dont match"
	}

	if photoOne.Credit != photoTwo.Credit {
		return false, "Credits dont match"
	}

	if photoInfoCheck, errStr := comparePhotoInfo(photoOne.Full, photoTwo.Full); !photoInfoCheck {
		return false, errStr
	}

	// TODO
	// if photoInfoCheck, errStr := comparePhotoInfo(photoOne.Thumbnail, photoTwo.Thumbnail); !photoInfoCheck {
	// 	return false, errStr
	// }
	//
	// if photoInfoCheck, errStr := comparePhotoInfo(photoOne.Small, photoTwo.Small); !photoInfoCheck {
	// 	return false, errStr
	// }

	return true, ""
}

func comparePhotoInfo(photoInfoOne, photoInfoTwo m.PhotoInfo) (bool, string) {
	if photoInfoOne.Url != photoInfoTwo.Url {
		return false, "Urls dont match"
	}

	if photoInfoOne.Width != photoInfoTwo.Width {
		return false, "widthds dont match"
	}

	if photoInfoOne.Height != photoInfoTwo.Height {
		return false, "height dont match"
	}

	return true, ""
}

func getTestArticle() *m.Article {
	file, err := ioutil.ReadFile("./testData/expectedArticle.json")
	if err != nil {
		fmt.Println("failed to read json article from filesystem")
		return nil
	}

	article := &m.Article{}
	err = json.Unmarshal(file, article)
	if err != nil {
		fmt.Println("failed to unmarshall json article into memory")
		return nil
	}
	return article
}
