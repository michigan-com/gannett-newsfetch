package commands

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/michigan-com/gannett-newsfetch/config"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Get articles published on the current day from the Gannett API",
	Run:   articleCmdRun,
}

func articleCmdRun(command *cobra.Command, args []string) {
	var startTime time.Time = time.Now()
	var envConfig, _ = config.GetEnv()
	var apiConfig, _ = config.GetApiConfig()
	var articleWait sync.WaitGroup
	var totalArticles int = 0
	articleChannel := make(chan *api.SearchArticle, len(apiConfig.SiteCodes)*100)
	articlesToScrape := make([]interface{}, 0, len(apiConfig.SiteCodes)*100)

	// Fetch each markets' articles in parallel
	log.Info("Fetching articles for all sites ...")
	for _, code := range apiConfig.SiteCodes {
		articleWait.Add(1)
		go func(code string) {
			defer articleWait.Done()
			articles := api.GetArticlesByDay(code, time.Now())

			for _, article := range articles {
				articleChannel <- article
			}
		}(code)
	}
	articleWait.Wait()
	close(articleChannel)
	log.Info("...Done fetching articles")

	if envConfig.MongoUri == "" {
		log.Print("No Mongo Uri specified, no articles will be saved")
		return
	}

	session := lib.DBConnect(envConfig.MongoUri)
	defer session.Close()

	// Iterate over all the articles, and determine whether or not we need to
	// summarize the articles
	log.Info("Determining which articles need to be scraped...")
	for article := range articleChannel {

		if shouldSummarizeArticle(article, session) {
			totalArticles += 1

			articleIdQuery := bson.M{"article_id": article.AssetId}
			articlesToScrape = append(articlesToScrape, articleIdQuery)
			articlesToScrape = append(articlesToScrape, articleIdQuery)
		}
	}

	bulk := session.DB("").C("ToScrape").Bulk()
	bulk.Upsert(articlesToScrape...)
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

/*
	We should summarize the article under two scenarios:

		1) This article does not yet exist in the database
		2) This article exists in the database, but the timestamp has been updated

	Does a lookup based on Article.ArticleId
*/
func shouldSummarizeArticle(article *api.SearchArticle, session *mgo.Session) bool {
	// Don't summarize if it's a blacklisted article
	if m.IsBlacklisted(article.Urls.LongUrl) {
		return false
	}

	var storedArticle *m.Article = &m.Article{}
	collection := session.DB("").C("Article")
	err := collection.Find(bson.M{"article_id": article.AssetId}).One(storedArticle)
	datePublished := lib.GannettDateStringToDate(article.DatePublished)
	if err == mgo.ErrNotFound {
		return true
	} else if !lib.SameDate(datePublished, storedArticle.Created_at) {
		return true
	} else if len(storedArticle.Summary) == 0 {
		return true
	}
	return false

}
