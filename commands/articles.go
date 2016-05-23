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
	m "github.com/michigan-com/gannett-newsfetch/model"
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
	summaryChannel := make(chan *m.Article, len(apiConfig.SiteCodes)*100)

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
	log.Info("Saving articles, determining summary necessity...")
	for article := range articleChannel {
		articleWait.Add(1)
		totalArticles += 1
		go func(article *fetch.ArticleIn) {
			mongoArticle := format.FormatSearchArticleForSaving(article)
			shouldSummarize := m.ShouldSummarizeArticle(mongoArticle, session)
			mongoArticle.Save(session)

			if shouldSummarize {
				summaryChannel <- mongoArticle
			}

			articleWait.Done()
		}(article)
	}
	articleWait.Wait()
	close(summaryChannel)
	log.Info("...Done scraping articles")

	// Grab the body text for the articles that need summarization
	log.Info("Grabbing article bodies from new/updated articles...")
	toSummarize := make([]interface{}, 0, len(summaryChannel))
	for article := range summaryChannel {
		articleWait.Add(1)
		go func(article *m.Article) {
			defer articleWait.Done()
			articleCol := session.DB("").C("Article")
			articleId := article.ArticleId
			articleContent := fetch.GetArticleBody(article.Url)
			body := parse.ParseArticleBodyHtml(articleContent.FullText)
			storyHighlights := articleContent.StoryHighlights

			if body == "" {
				log.Infof("No body text for article %v, summary will be skipped", articleId)
				return
			}

			summarizedArticles += 1

			query := bson.M{"article_id": articleId}
			update := bson.M{"$set": bson.M{"body": body, "storyHighlights": storyHighlights}}
			articleCol.Update(query, update)

			toSummarize = append(toSummarize, query)
			toSummarize = append(toSummarize, query)
		}(article)
	}
	articleWait.Wait()
	log.Info("...Done grabbing article bodies")

	// Now, look for articles that show up in chartbeat but not in the search/v4 api

	// Save the articles we're going to summarize, and run the summarizer
	if len(toSummarize) > 0 {
		log.Info("Summarizing articles...")
		_, err := ProcessSummaries(toSummarize, session)
		if err != nil {
			log.Errorf("Failed to process summaries: %v", err)
		}
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
