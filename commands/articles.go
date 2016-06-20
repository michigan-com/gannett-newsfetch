package commands

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

func GetArticles(session *mgo.Session, siteCodes []string, gannettAPIKey string) {
	var startTime time.Time = time.Now()

	var articleWait sync.WaitGroup
	var totalArticles int = 0
	articleChannel := make(chan *m.SearchArticle, len(siteCodes)*100)

	// Fetch each markets' articles in parallel
	log.Info("Fetching articles for all sites ...")
	for _, code := range siteCodes {
		articleWait.Add(1)
		go func(code string) {
			defer articleWait.Done()
			articles := api.GetArticlesByDay(code, time.Now(), gannettAPIKey)

			for _, article := range articles {
				articleChannel <- article
			}
		}(code)
	}
	articleWait.Wait()
	close(articleChannel)
	log.Info("...Done fetching articles")

	if session == nil {
		log.Print("No Mongo Uri specified, no articles will be saved")
		return
	}

	bulk := session.DB("").C("ToScrape").Bulk()

	// Iterate over all the articles, and determine whether or not we need to
	// summarize the articles
	log.Info("Determining which articles need to be scraped...")
	for article := range articleChannel {

		if shouldSummarizeArticle(article, session) {
			totalArticles += 1

			bulk.Upsert(bson.M{"article_id": article.AssetId}, &m.ScrapeRequest{
				ArticleID:  article.AssetId,
				ArticleURL: article.Urls.LongUrl,
			})
		}
	}

	_, err := bulk.Run()
	if err != nil {
		log.Errorf("Failed to store articles to be scraped: %v", err)
	}
	log.Info("...Done")

	log.Infof(`

	Article processing done (%v)

		Total Articles Found:	%d
	`, time.Now().Sub(startTime), totalArticles)
}
