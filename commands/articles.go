package commands

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/andreyvit/mongobulk"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

func GetArticles(session *mgo.Session, siteCodes []string, gannettSearchAPIKey string) {
	var startTime time.Time = time.Now()

	var articleWait sync.WaitGroup
	var totalArticles int = 0
	articleChannel := make(chan *m.SearchArticle, len(siteCodes)*100)

	// Fetch each markets' articles in parallel
	log.Printf("Fetching articles for all sites ...")
	for _, code := range siteCodes {
		articleWait.Add(1)
		go func(code string) {
			defer articleWait.Done()
			articles := api.GetArticlesByDay(code, time.Now(), gannettSearchAPIKey)

			for _, article := range articles {
				articleChannel <- article
			}
		}(code)
	}
	articleWait.Wait()
	close(articleChannel)
	log.Printf("...Done fetching articles")

	coll := session.DB("").C("ToScrape")
	bulk := mongobulk.New(coll, mongobulk.Config{})

	// Iterate over all the articles, and determine whether or not we need to
	// summarize the articles
	log.Printf("Determining which articles need to be scraped...")
	for article := range articleChannel {
		if shouldSummarizeArticle(article, session) {
			totalArticles += 1

			bulk.Upsert(bson.M{"article_id": article.AssetId}, &m.ScrapeRequest{
				ArticleID:  article.AssetId,
				ArticleURL: article.Urls.LongUrl,
			})
		}
	}

	err := bulk.Finish()
	if err != nil {
		log.Printf("ERROR: Failed to store articles to be scraped: %v", err)
	}

	log.Printf("Article processing done (%v). Total Articles Found: %d.", time.Now().Sub(startTime), totalArticles)
}
