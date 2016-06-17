package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/brvtyclient"
	"github.com/michigan-com/gannett-newsfetch/config"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

var cleanupCommand = &cobra.Command{
	Use:   "scrape-and-summarize",
	Short: "Grab stories that we see in chartbeat but not the Gannett API",
	Run:   scrapeAndSummarizeCmd,
}

func scrapeAndSummarizeCmd(command *cobra.Command, args []string) {
	var config, _ = config.GetEnv()
	ScrapeAndSummarize(config)
}

func ScrapeAndSummarize(config EnvConfig) {
	var articleWait sync.WaitGroup
	if config.MongoUri == "" {
		log.Warning("No mongo URI specified, this command is basically useless")
		return
	}

	session := lib.DBConnect(config.MongoUri)
	toScrape := session.DB("").C("ToScrape")
	defer session.Close()

	client := brvtyclient.New(config.BrvtyURL)

	for {
		toSummarize := make([]interface{}, 0, 100)
		var requests []m.ScrapeRequest

		log.Info("Finding articles in need of scraping...")
		err := toScrape.Find(bson.M{}).Select(bson.M{"article_id": true, "article_url": true, "_id": false}).All(&requests)
		if err != nil {
			log.Errorf("Error loading ToScrape collection: %v", err)
		}

		if len(requests) > 0 {
			log.Infof("...scraping %d articles...", len(requests))
			for _, request := range requests {
				articleWait.Add(1)
				go func(request m.ScrapeRequest) {
					defer articleWait.Done()
					assetArticleContent := api.GetAssetArticleContent(request.ArticleID)

					mongoArticle := api.FormatAssetArticleForSaving(assetArticleContent)
					mongoArticle.Save(session)

					articleIdQuery := bson.M{"article_id": request.ArticleID}
					toSummarize = append(toSummarize, articleIdQuery)
					toSummarize = append(toSummarize, articleIdQuery)

					toScrape.Remove(articleIdQuery)
				}(request)
			}
			log.Infof("...Done scraping articles")
			articleWait.Wait()
		} else {
			log.Infof("...no articles in need of scraping")
		}

		if len(toSummarize) > 0 {
			log.Info("Summarizing articles...")
			_, err := ProcessSummaries(toSummarize, config.MongoUri)
			if err != nil {
				log.Errorf("Failed to process summaries: %v", err)
			}
			log.Info("...Done processing summaries")
		} else {
			log.Info("No articles to summarize.")
		}

		if loop > 0 {
			log.Infof("Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			log.Info("...and now I'm awake!")
		} else {
			break
		}
	}
}
