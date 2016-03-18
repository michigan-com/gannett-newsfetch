package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/gannett-newsfetch/config"
	fetch "github.com/michigan-com/gannett-newsfetch/gannettApi/fetch"
	format "github.com/michigan-com/gannett-newsfetch/gannettApi/format"
	"github.com/michigan-com/gannett-newsfetch/lib"
	"github.com/michigan-com/gannett-newsfetch/model"
	parse "github.com/michigan-com/gannett-newsfetch/parse/body"
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
	var summarizedArticles int = 0
	articleChannel := make(chan *fetch.ArticleIn, len(apiConfig.SiteCodes)*100)
	summaryChannel := make(chan int, len(apiConfig.SiteCodes)*100)

	// Fetch each markets' articles in parallel
	log.Info("Fetching articles for all sites ...")
	for _, code := range apiConfig.SiteCodes {
		articleWait.Add(1)
		go func(code string) {
			defer articleWait.Done()
			articles := fetch.GetArticlesByDay(code, time.Now())

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
	log.Info("Scraping articles, if necessary...")
	for article := range articleChannel {
		articleWait.Add(1)
		totalArticles += 1
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
	log.Info("...Done scraping articles")

	// Grab the body text for the articles that need summarization
	toSummarize := make([]interface{}, 0, len(summaryChannel))
	for articleId := range summaryChannel {
		summarizedArticles += 1
		articleWait.Add(1)
		go func(articleId int) {
			defer articleWait.Done()
			articleCol := session.DB("").C("Article")
			articleContent := fetch.GetArticleContent(articleId)
			body := parse.ParseArticleBodyHtml(articleContent.FullText)

			query := bson.M{"article_id": articleId}
			update := bson.M{"$set": bson.M{"body": body}}
			articleCol.Update(query, update)

			toSummarize = append(toSummarize, bson.M{"article_id": articleId})
			toSummarize = append(toSummarize, bson.M{"article_id": articleId})
		}(articleId)
	}
	articleWait.Wait()

	// Save the articles we're going to summarize, and run the summarizer
	if len(toSummarize) > 0 {
		log.Info("Summarizing articles...")
		bulk := session.DB("").C("ToSummarize").Bulk()
		bulk.Upsert(toSummarize...)
		_, err := bulk.Run()
		if err != nil {
			log.Info(err)
		}
		ProcessSummaries()
		log.Info("...Done summarizing articles")
	} else {
		log.Info("Hey look at that, no new articles to summarize")
	}

	log.Infof(`

	Article processing done (%v)

		Total Articles Processed:	%d
		Total Articles Summarized:	%d
	`, time.Now().Sub(startTime), totalArticles, summarizedArticles)
}
