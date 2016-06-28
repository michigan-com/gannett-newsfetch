package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/mgo.v2/bson"

	newsfetch "github.com/michigan-com/gannett-newsfetch"
	c "github.com/michigan-com/gannett-newsfetch/commands"
	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

var config newsfetch.Config

func init() {
	var err error
	config, err = newsfetch.ParseConfig()
	if err != nil {
		panic(err)
	}

	config.MongoURI = "mongodb://localhost:27017/gannett-newsfetch-test"
}

func TestIntegration(t *testing.T) {
	jsonFiles := []string{
		"./testData/expectedArticleWithVideo.json",
		// "./testData/expectedArticleNoPhoto.json",
	}
	session, err := newsfetch.SetupMongoSession(config.MongoURI)
	if err != nil {
		t.Fatalf("Error connecting to Mongo: %v", err)
	}
	defer session.Close()

	for _, jsonFile := range jsonFiles {
		testArticle := getTestArticle(jsonFile)
		if testArticle == nil {
			t.Fatalf("Test article json reading failed")
		}

		testArticleId := testArticle.ArticleId
		toScrapeCol := session.DB("").C("ToScrape")
		articleCol := session.DB("").C("Article")
		toScrapeCol.Insert(bson.M{"article_id": testArticleId})

		c.ScrapeAndSummarize(session, nil, nil, 0, 0, config.MongoURI, config.SummaryVEnv, config.GannettAssetAPIKey)

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
			t.Fatalf("Article %d failed article comparison: %s", testArticleId, errString)
		}
	}

	session.DB("").DropDatabase()
}

func TestBreakingNewsIntegration(t *testing.T) {
	session, err := newsfetch.SetupMongoSession(config.MongoURI)
	if err != nil {
		t.Fatalf("Error connecting to Mongo: %v", err)
	}
	defer session.Close()

	testBreakingNewsUrl := "http://www.freep.com/story/news/local/michigan/2016/06/06/insanity-defense-kalamazoo-shootings/85516404/"
	testBreakingArticle := &m.SearchArticle{
		AssetId:    123123,
		Urls:       m.Urls{LongUrl: testBreakingNewsUrl},
		Headline:   "Test test test",
		PromoBrief: "test test test test test",
	}
	// testArticleId := lib.GetArticleId(testBreakingNewsUrl)
	breakingChannel := make(chan *m.SearchArticle, 1)
	breakingChannel <- testBreakingArticle
	close(breakingChannel)

	// First, make sure thats we don't save this article as a breaking news alert because
	// we haven't summarized the article yet
	breakingArticles := c.SaveBreakingArticles(breakingChannel, session)
	if len(breakingArticles) != 0 {
		t.Fatalf("we should not be saving breaking news articles that havent been scraped yet")
	}

	// Run the scraping process, and summarize the necessary article
	c.ScrapeAndSummarize(session, nil, nil, 0, 0, config.MongoURI, config.SummaryVEnv, config.GannettAssetAPIKey)

	// Now, we should get one breaking news alert with this newly scraped article
	breakingChannel = make(chan *m.SearchArticle, 1)
	breakingChannel <- testBreakingArticle
	close(breakingChannel)
	breakingArticles = c.SaveBreakingArticles(breakingChannel, session)
	if len(breakingArticles) != 1 {
		t.Fatalf("Should have one breaking article now")
	}

	session.DB("").DropDatabase()
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

	if videoCheck, errStr := compareVideos(articleOne.Video, articleTwo.Video); !videoCheck {
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
	if photoOne == nil && photoTwo != nil {
		return false, "ArticleOne doesnt have a photo, but ArticleTwo does"
	} else if photoOne != nil && photoTwo == nil {
		return false, "ArticleOne has a photo, ArticleTwo does not"
	}

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
	if photoInfoCheck, errStr := comparePhotoInfo(photoOne.Thumbnail, photoTwo.Thumbnail); !photoInfoCheck {
		return false, errStr
	}

	if photoInfoCheck, errStr := comparePhotoInfo(photoOne.Small, photoTwo.Small); !photoInfoCheck {
		return false, errStr
	}

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

func compareVideos(videoOne, videoTwo *m.AssetVideo) (bool, string) {
	if videoOne == nil && videoTwo != nil {
		return false, "ArticleOne doesnt have a video, but ArticleTwo does"
	} else if videoOne != nil && videoTwo == nil {
		return false, "ArticleTwo has a video, but ArticleTwo does not"
	}

	if videoOne.Thumbnail != videoTwo.Thumbnail {
		return false, "Video thumbnails dont match"
	}

	if videoOne.VideoStill != videoTwo.VideoStill {
		return false, "VideoStills dont match"
	}

	if videoOne.Length != videoTwo.Length {
		return false, "Video length doesnt match"
	}

	if renditionCheck, errStr := compareVideoRenditions(videoOne.Renditions, videoTwo.Renditions); !renditionCheck {
		return false, errStr
	}

	return true, ""
}

func compareVideoRenditions(renditionsOne, renditionsTwo []m.VideoRendition) (bool, string) {
	if len(renditionsOne) != len(renditionsTwo) {
		return false, "Length of renditions dont match"
	}

	for i, rendition := range renditionsOne {
		otherRendition := renditionsTwo[i]

		if rendition.EncodingRate != otherRendition.EncodingRate {
			return false, fmt.Sprintf("Rendition %d has mismatched encoding rates", i)
		}

		if rendition.Height != otherRendition.Height {
			return false, fmt.Sprintf("rendition %d heights dont match", i)
		}

		if rendition.Width != otherRendition.Width {
			return false, fmt.Sprintf("rendition %d widths dont match", i)
		}

		if rendition.Size != otherRendition.Size {
			return false, fmt.Sprintf("rendition %d sizes dont match", i)
		}

		if rendition.Url != otherRendition.Url {
			return false, fmt.Sprintf("rendition %d urls dont match", i)
		}

		if rendition.Duration != otherRendition.Duration {
			return false, fmt.Sprintf("rendition %d durations dont match", i)
		}
	}

	return true, ""
}

func getTestArticle(jsonFilePath string) *m.Article {
	file, err := ioutil.ReadFile(jsonFilePath)
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
