package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/michigan-com/gannett-newsfetch/config"
	fetch "github.com/michigan-com/gannett-newsfetch/gannettApi/fetch"
	format "github.com/michigan-com/gannett-newsfetch/gannettApi/format"
	"github.com/michigan-com/gannett-newsfetch/lib"
	"github.com/michigan-com/gannett-newsfetch/model"
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Get articles from the gannett API based on the news source",
	Run:   articleCmdRun,
}

func articleCmdRun(command *cobra.Command, args []string) {
	var envConfig, _ = config.GetEnv()
	var apiConfig, _ = config.GetApiConfig()
	var articleWait sync.WaitGroup
	articleChannel := make(chan *fetch.ArticleIn, len(apiConfig.SiteCodes)*100)
	summaryChannel := make(chan int, len(apiConfig.SiteCodes)*100)

	// Fetch each markets' articles in parallel
	for _, code := range apiConfig.SiteCodes {
		articleWait.Add(1)
		go func(code string) {
			articles := fetch.GetArticlesByDay(code, time.Now())
			log.Infof("got articles for %s", code)
			log.Info(len(articles))

			for _, article := range articles {
				articleChannel <- article
			}

			log.Infof("Done adding articles")
			articleWait.Done()
		}(code)
	}
	articleWait.Wait()
	close(articleChannel)

	if envConfig.MongoUri == "" {
		log.Print("No Mongo Uri specified, no articles will be saved")
		return
	}

	session := lib.DBConnect(envConfig.MongoUri)
	defer session.Close()

	// Iterate over all the articles, and determine whether or not we need to
	// summarize the articles
	for article := range articleChannel {
		articleWait.Add(1)
		go func(article *fetch.ArticleIn) {
			mongoArticle := format.FormatArticleForSaving(article)
			shouldSummarize := model.ShouldSummarizeArticle(mongoArticle, session)
			mongoArticle.Save(session)

			if shouldSummarize {
				summaryChannel <- mongoArticle.ArticleId
			}

			articleWait.Done()
		}(article)
	}
	articleWait.Wait()
	close(summaryChannel)
}
